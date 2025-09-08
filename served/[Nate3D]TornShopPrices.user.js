// ==UserScript==
// @name         Torn PDA — Inventory Prices (TT-style + aligned + vendor tooltip)
// @description  TT-style prices with column alignment. Equipment shows price in dots row. Unit is green/red vs vendor; difference shown in tooltip. Includes sortable flyout control.
// @version      1.8.0
// @author       you
// @match        https://www.torn.com/item.php*
// @match        https://www.torn.com/*items.php*
// @run-at       document-idle
// ==/UserScript==

(() => {
    const LS_KEY = "TPDA_API_KEY";
    const CACHE_KEY = "TPDA_ITEMS_CACHE_V2";
    const CACHE_TTL_MS = 24 * 60 * 60 * 1000;

    // NEW: persist flyout state
    const FLYOUT_STATE_KEY = "TPDA_SORT_FLYOUT_STATE"; // values: "expanded" | "collapsed"

    const $ = (s, r = document) => r.querySelector(s);
    const $$ = (s, r = document) => Array.from(r.querySelectorAll(s));
    const fmtMoney = (n) => "$" + Math.round(n).toLocaleString("en-US");
    const fmtPct = (v) => `${v.toFixed(1)}%`; // one decimal is plenty

    // ---------- Styles
    const style = document.createElement("style");
    style.textContent = `
    .tpda-green { color:#39b54a !important; }
    .tpda-red   { color:#ff5a5a !important; }

    .title-wrap .title .name-wrap{
      display:flex; gap:.25rem; align-items:center; flex-wrap:nowrap; min-width:0;
    }
    .title-wrap .title .name-wrap .name,
    .title-wrap .title .name-wrap .qty{ flex:0 0 auto; }
    .title-wrap .title .name-wrap .tt-item-price{ margin-left:auto; flex:0 0 auto; }

    .tt-item-price, .tpda-price {
      white-space:nowrap;
      font-weight:600;
      font-variant-numeric: tabular-nums;
      display:grid;
      grid-auto-flow: column;
      align-items:center;
      column-gap: 8px;
      padding-right: 6px;
    }
    .tpda-qty  { color:#fff; }
    .tpda-total{ color:#39b54a; }

    .bonuses-wrap > li.tt-item-price{ display:inline-block; margin-left:8px; padding-right:6px; }
    .bonuses-wrap > li.tt-item-price.fl{ float:none; }

    /* keep legacy simple one-line green if ever used */
    .title-wrap .title .name-wrap .tt-item-price:not(.tpda-price){ color:#39b54a !important; }

    /* Sort Control Styles */
    .tpda-sort-control {
      position: fixed;
      top: 30vh;              /* Moved to ~30% down the screen */
      right: 20px;
      z-index: 10000;
      background: #1a1a1a;
      border: 1px solid #444;
      border-radius: 8px;
      padding: 12px;
      box-shadow: 0 4px 12px rgba(0,0,0,0.5);
      font-family: inherit;
      transition: transform 0.3s ease;
      min-width: 220px;
    }
    .tpda-sort-control.collapsed {
      transform: translateX(calc(100% - 40px));
    }
    .tpda-sort-control.collapsed .tpda-sort-toggle {
      position: absolute;
      left: -30px;
      top: 12px;
      z-index: 1;
    }
    .tpda-sort-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 10px;
      color: #fff;
      font-weight: bold;
      font-size: 14px;
    }
    .tpda-sort-toggle {
      background: #333;
      border: 1px solid #555;
      color: #fff;
      padding: 4px 8px;
      border-radius: 4px;
      cursor: pointer;
      font-size: 12px;
    }
    .tpda-sort-toggle:hover { background: #444; }
    .tpda-sort-content { display: flex; flex-direction: column; gap: 8px; }
    .tpda-sort-control.collapsed .tpda-sort-content { display: none; }
    .tpda-sort-select {
      background: #222;
      border: 1px solid #444;
      color: #fff;
      padding: 6px 8px;
      border-radius: 4px;
      font-size: 13px;
    }
    .tpda-sort-button {
      background: #39b54a;
      border: none;
      color: #fff;
      padding: 8px 12px;
      border-radius: 4px;
      cursor: pointer;
      font-weight: bold;
      font-size: 13px;
    }
    .tpda-sort-button:hover { background: #2ea03d; }
    .tpda-sort-button:disabled { background: #555; cursor: not-allowed; }

    /* NEW: visually mark modifier/rarity-dependent prices (equipment) */
    .tpda-variant .tpda-unit { color: #cfa8ff !important; } /* light purple */
  `;
    document.head.appendChild(style);

    // ---------- API key
    async function getApiKey() {
        let key = localStorage.getItem(LS_KEY) || "";
        if (!key) {
            key = prompt(
                "Torn PDA — Inventory Prices:\nEnter your Torn API key (stored locally).",
                ""
            );
            if (key) localStorage.setItem(LS_KEY, key.trim());
        }
        key = (key || "").trim();
        return key || null;
    }

    // ---------- Cache + fetch
    function readCache() {
        try {
            const raw = localStorage.getItem(CACHE_KEY);
            if (!raw) return null;
            const obj = JSON.parse(raw);
            if (!obj?.updated || !obj?.items) return null;
            if (Date.now() - obj.updated > CACHE_TTL_MS) return null;
            return obj.items;
        } catch { return null; }
    }
    function writeCache(items) {
        try {
            localStorage.setItem(CACHE_KEY, JSON.stringify({ updated: Date.now(), items }));
        } catch { }
    }

    async function getAllItemsMap() {
        const cached = readCache();
        if (cached) return cached;

        const key = await getApiKey();
        if (!key) return {};

        const url = `https://api.torn.com/torn/?selections=items&key=${encodeURIComponent(key)}`;
        const resp = await fetch(url, { credentials: "omit" });
        if (!resp.ok) return {};
        const data = await resp.json();
        if (data?.error) return {};

        const map = {};
        for (const [id, info] of Object.entries(data.items || {})) {
            map[id] = {
                name: info.name,
                market_value: Number(info.market_value) || 0,
                buy_price: Number(info.buy_price) || 0, // vendor price
            };
        }
        writeCache(map);
        return map;
    }

    // Nicely formatted vendor/market tooltip
    function makeVendorTooltip({ unit, vendor }) {
        const marketStr = fmtMoney(unit);
        if (vendor > 0) {
            const vendorStr = fmtMoney(vendor);
            const diff = unit - vendor;
            const diffStr = fmtMoney(Math.abs(diff));
            const pct = (diff / vendor) * 100;
            const diffSign = diff >= 0 ? "+" : "-";
            const pctSign = pct >= 0 ? "+" : "-";
            return `Vendor: ${vendorStr}\nMarket: ${marketStr}\nΔ: ${diffSign}${diffStr} (${pctSign}${Math.abs(pct).toFixed(1)}%)`;
        }
        return `Vendor: N/A\nMarket: ${marketStr}`;
    }

    function buildAlignedPriceElems({ unit, qty, total, vendor }) {
        const wrap = document.createElement('span');
        wrap.className = 'tt-item-price tpda-price';

        const unitSpan = document.createElement('span');
        unitSpan.className = (vendor > 0 && unit - vendor < 0) ? 'tpda-unit tpda-red' : 'tpda-unit tpda-green';
        unitSpan.textContent = fmtMoney(unit);

        const tip = makeVendorTooltip({ unit, vendor });
        unitSpan.setAttribute('title', tip);
        unitSpan.setAttribute('data-title', tip);

        const qtySpan = document.createElement('span');
        qtySpan.className = 'tpda-qty';
        qtySpan.textContent = `x${qty}`;

        const totalSpan = document.createElement('span');
        totalSpan.className = 'tpda-total';
        totalSpan.textContent = `= ${fmtMoney(total)}`;

        wrap.appendChild(unitSpan);
        wrap.appendChild(qtySpan);
        wrap.appendChild(totalSpan);
        return wrap;
    }

    function buildEquipmentPriceLi({ unit, vendor }) {
        const li = document.createElement('li');
        li.className = 'tt-item-price fl';

        const span = document.createElement('span');
        // NEW: mark equipment as variant-priced -> purple
        span.className = 'tpda-price tpda-variant';

        const unitSpan = document.createElement('span');
        // underlying class still set; purple override comes from .tpda-variant .tpda-unit
        unitSpan.className = (vendor > 0 && unit - vendor < 0) ? 'tpda-unit tpda-red' : 'tpda-unit tpda-green';
        unitSpan.textContent = fmtMoney(unit);

        const tip = makeVendorTooltip({ unit, vendor });
        unitSpan.setAttribute('title', tip);
        unitSpan.setAttribute('data-title', tip);

        span.appendChild(unitSpan);
        li.appendChild(span);
        return li;
    }

    // ---------- Injectors
    function injectOneLine(li, item) {
        const nameWrap = li.querySelector(".title-wrap .name-wrap");
        if (!nameWrap) return;

        const qtyAttr = li.getAttribute("data-qty");
        const qty = Number(qtyAttr || "1") || 1;
        const unit = item.market_value;
        const total = unit * qty;

        nameWrap.querySelector(":scope > .tt-item-price")?.remove();
        const price = buildAlignedPriceElems({
            unit,
            qty,
            total,
            vendor: item.buy_price,
        });
        nameWrap.appendChild(price);
    }

    function injectEquipment(li, item) {
        const bonusesUl = li.querySelector(".cont-wrap .bonuses .bonuses-wrap");
        if (!bonusesUl) return;

        li.querySelector(".title-wrap .name-wrap > .tt-item-price")?.remove();

        let liPrice = bonusesUl.querySelector(":scope > li.tt-item-price");
        const fresh = buildEquipmentPriceLi({
            unit: item.market_value,
            vendor: item.buy_price,
        });
        if (!liPrice) bonusesUl.appendChild(fresh);
        else liPrice.replaceWith(fresh);
    }

    function decorateLI(li, itemsMap) {
        if (!li || li.nodeType !== 1) return;
        const id = li.getAttribute("data-item");
        if (!id) return;
        const item = itemsMap[id];
        if (!item || !item.market_value) return;

        const isEquipment = !!li.querySelector(".cont-wrap .bonuses .bonuses-wrap");
        if (isEquipment) injectEquipment(li, item);
        else injectOneLine(li, item);
    }

    function processAll(itemsMap) {
        $$(".itemsList li[data-item]").forEach((li) => decorateLI(li, itemsMap));
    }

    function watch(itemsMap) {
        const root = $("#category-wrap") || document.body;
        const mo = new MutationObserver((muts) => {
            for (const m of muts) {
                if (m.type === "childList") {
                    m.addedNodes.forEach((n) => {
                        if (!(n instanceof HTMLElement)) return;
                        if (n.matches?.("li[data-item]")) decorateLI(n, itemsMap);
                        n.querySelectorAll?.("li[data-item]").forEach((li) =>
                            decorateLI(li, itemsMap)
                        );
                    });
                } else if (m.type === "attributes" && m.target instanceof HTMLElement) {
                    if (m.target.matches("li[data-item]")) decorateLI(m.target, itemsMap);
                }
            }
        });
        mo.observe(root, {
            childList: true,
            subtree: true,
            attributes: true,
            attributeFilter: ["data-item", "data-qty"],
        });
    }

    function addKeyResetButton() {
        const header = $(".items-wrap .title-black");
        if (!header || header.querySelector(".tpda-key-btn")) return;
        const btn = document.createElement("button");
        btn.className = "tpda-key-btn";
        btn.textContent = "API key";
        btn.title = "Set/Update Torn API key for price fetch";
        Object.assign(btn.style, {
            marginLeft: "8px",
            padding: "2px 6px",
            borderRadius: "4px",
            border: "1px solid #444",
            background: "#222",
            color: "#ccc",
            cursor: "pointer",
        });
        btn.addEventListener("click", () => {
            const key = prompt(
                "Enter Torn API key:",
                localStorage.getItem(LS_KEY) || ""
            );
            if (key !== null) {
                localStorage.setItem(LS_KEY, (key || "").trim());
                localStorage.removeItem(CACHE_KEY);
                location.reload();
            }
        });
        header.appendChild(btn);
    }

    function createSortControl(itemsMap) {
        if (document.querySelector('.tpda-sort-control')) return;

        const control = document.createElement('div');
        control.className = 'tpda-sort-control';

        const header = document.createElement('div');
        header.className = 'tpda-sort-header';
        header.innerHTML = `
            <span>Sort Items</span>
            <button class="tpda-sort-toggle">⟨</button>
        `;

        const content = document.createElement('div');
        content.className = 'tpda-sort-content';

        const select = document.createElement('select');
        select.className = 'tpda-sort-select';
        select.innerHTML = `
            <option value="none">No Sort</option>
            <option value="unit-low">Unit Price: Low → High</option>
            <option value="unit-high">Unit Price: High → Low</option>
            <option value="total-low">Total Cost: Low → High</option>
            <option value="total-high">Total Cost: High → Low</option>
        `;

        const button = document.createElement('button');
        button.className = 'tpda-sort-button';
        button.textContent = 'Sort Now';

        content.appendChild(select);
        content.appendChild(button);
        control.appendChild(header);
        control.appendChild(content);

        // NEW: restore persisted collapsed/expanded state
        const toggle = header.querySelector('.tpda-sort-toggle');
        const savedState = localStorage.getItem(FLYOUT_STATE_KEY);
        if (savedState === "collapsed") {
            control.classList.add('collapsed');
            toggle.textContent = '⟩';
        } else {
            toggle.textContent = '⟨';
        }

        // Toggle functionality + persist
        toggle.addEventListener('click', () => {
            control.classList.toggle('collapsed');
            const collapsed = control.classList.contains('collapsed');
            toggle.textContent = collapsed ? '⟩' : '⟨';
            localStorage.setItem(FLYOUT_STATE_KEY, collapsed ? "collapsed" : "expanded");
        });

        // Sort functionality
        button.addEventListener('click', () => {
            const sortType = select.value;
            if (sortType === 'none') {
                restoreOriginalOrder();
            } else {
                sortItems(sortType, itemsMap);
            }
        });

        document.body.appendChild(control);
    }

    function getItemData(li, itemsMap) {
        const itemId = li.getAttribute('data-item');
        const qty = Number(li.getAttribute('data-qty') || '1') || 1;
        const item = itemsMap[itemId];
        if (!item || !item.market_value) return null;

        return {
            element: li,
            itemId,
            qty,
            unitPrice: item.market_value,
            totalCost: item.market_value * qty,
            originalIndex: Array.from(li.parentNode.children).indexOf(li)
        };
    }

    function sortItems(sortType, itemsMap) {
        const itemsList = $('.itemsList:not([style*="display: none"])');
        if (!itemsList) return;

        cleanupItemsList(itemsList);

        const items = Array.from(itemsList.querySelectorAll('li[data-item]'))
            .map(li => getItemData(li, itemsMap))
            .filter(Boolean);

        if (items.length === 0) return;

        if (!itemsList.hasAttribute('data-original-order')) {
            const originalOrder = items.map(item => item.element.getAttribute('data-item')).join(',');
            itemsList.setAttribute('data-original-order', originalOrder);
        }

        items.sort((a, b) => {
            switch (sortType) {
                case 'unit-low': return a.unitPrice - b.unitPrice;
                case 'unit-high': return b.unitPrice - a.unitPrice;
                case 'total-low': return a.totalCost - b.totalCost;
                case 'total-high': return b.totalCost - a.totalCost;
                default: return 0;
            }
        });

        reorderItemsSafely(itemsList, items);
    }

    function restoreOriginalOrder() {
        const itemsList = $('.itemsList:not([style*="display: none"])');
        if (!itemsList || !itemsList.hasAttribute('data-original-order')) return;

        cleanupItemsList(itemsList);

        const originalOrder = itemsList.getAttribute('data-original-order').split(',');
        const itemElements = Array.from(itemsList.querySelectorAll('li[data-item]'));

        const elementMap = new Map();
        itemElements.forEach(el => {
            elementMap.set(el.getAttribute('data-item'), el);
        });

        const orderedItems = originalOrder.map(itemId => elementMap.get(itemId)).filter(Boolean);
        reorderItemsSafely(itemsList, orderedItems.map(element => ({ element })));
    }

    function cleanupItemsList(itemsList) {
        const expandedItems = itemsList.querySelectorAll('li.item-info-active, li.show-item-info');
        expandedItems.forEach(item => {
            if (item.classList.contains('show-item-info')) {
                item.remove();
            } else {
                item.classList.remove('act', 'item-info-active', 'thumbnail-active', 'item-active', 'opened');
            }
        });

        const allLiElements = Array.from(itemsList.querySelectorAll('li'));
        allLiElements.forEach(li => {
            if (!li.hasAttribute('data-item') &&
                !li.classList.contains('show-item-info') &&
                !li.querySelector('.item-wrap, .thumbnail-wrap')) {
                console.warn('Removing orphaned element:', li);
                li.remove();
            }
        });
    }

    function reorderItemsSafely(itemsList, items) {
        try {
            const fragment = document.createDocumentFragment();
            items.forEach(item => {
                const element = item.element;
                if (element && element.parentNode === itemsList) {
                    fragment.appendChild(element);
                }
            }); 
            itemsList.innerHTML = '';
            itemsList.appendChild(fragment);
        } catch (e) {
            console.error('Safe reordering failed, falling back to simple approach:', e);
            items.forEach(item => {
                const el = item.element;
                if (el && itemsList.contains(el)) itemsList.appendChild(el);
            });
        }
    }

    (async function main() {
        try {
            addKeyResetButton();
            const itemsMap = await getAllItemsMap();
            if (!Object.keys(itemsMap).length) return;
            processAll(itemsMap);
            watch(itemsMap);
            createSortControl(itemsMap);
        } catch (e) {
            console.error("Inventory Prices script error:", e);
        }
    })();
})();