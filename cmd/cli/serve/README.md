# Directory File Server CLI Tool

This tool allows you to serve the contents of any directory over HTTP, making it accessible to any device on your local network.

## Usage

Build the tool:

```bash
make build
```

Run the server, specifying the directory to serve (default is current directory):

```bash
./bin/serve --dir /path/to/directory --port 8080
```

- `--dir`: Directory to serve (default: current directory)
- `--port`: Port to listen on (default: 8080)

## Example

```bash
./bin/serve --dir ./public --port 9000
```

Then, on any device on your local network, open a browser and go to:

```text
http://<your-local-ip>:9000
```

The tool will print the correct URL when started.

## Notes

- Works on any OS (Linux, WSL2, Ubuntu dev containers, etc.)
- No dependencies beyond the Go standard library
- Useful for quick file sharing or static site hosting on a LAN
