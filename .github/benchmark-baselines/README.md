# Benchmark Baselines

This directory contains baseline benchmark results used for performance regression detection.

## Files

- `baseline-amd64.txt` - Baseline benchmarks for AMD64 architecture
- `baseline-arm64.txt` - Baseline benchmarks for ARM64 architecture
- `metadata.txt` - Metadata about when baselines were last updated

## How it works

1. When code is pushed to `main`, the benchmark workflow runs comprehensive benchmarks
2. The results are automatically stored here as new baselines
3. For pull requests, current benchmarks are compared against these baselines
4. Performance regressions are detected and reported in PR comments

## Manual baseline update

To manually update baselines:

```bash
# Run benchmarks locally
go test -bench=. -benchmem -benchtime=5s -count=3 ./... > .github/benchmark-baselines/baseline-amd64.txt

# Update metadata
echo "Updated: $(date -u +"%Y-%m-%d %H:%M:%S UTC")" > .github/benchmark-baselines/metadata.txt
echo "Commit: $(git rev-parse HEAD)" >> .github/benchmark-baselines/metadata.txt
echo "Manual update" >> .github/benchmark-baselines/metadata.txt

# Commit changes
git add .github/benchmark-baselines/
git commit -m "📊 Update benchmark baselines"
```

## Performance Thresholds

The benchmark workflow will flag:
- **Regressions**: >10% performance decrease
- **Improvements**: >10% performance increase
- **Critical**: >50% performance decrease (warnings)

## Architecture Support

Benchmarks are run on multiple architectures:
- **Linux AMD64**: Primary development and CI environment
- **Linux ARM64**: Alternative architecture for comprehensive testing

## Baseline Storage

Baselines are stored in git to ensure:
- Version control of performance expectations
- Ability to track performance changes over time
- Automatic updates as code evolves
