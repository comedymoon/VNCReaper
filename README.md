> [!CAUTION]
> This tool was made **only — and I mean ONLY — for educational use**.  
> I’m not taking the blame for anything you do with it.  
> Use it only on machines you own or have clear permission to mess with.  
> If you break the law, that’s on you.

> [!WARNING]
> This project was partly created with help from ChatGPT.  
> AI can be wrong, so always double-check the code and know what it does before running it.

<img width="1280" height="460" alt="vncreaper" src="https://github.com/user-attachments/assets/45ae4c4c-6c08-4c1b-b018-7f58c1720810" />

---

### **VNCReaper**

# *Tool for scanning/finding open VNCs on the internet.

VNCReaper is a high-speed scanner for finding VNC and noVNC services across IP ranges.  
It’s written in Go for speed, efficiency, and the ability to handle huge amounts of targets without choking.

The goal is simple: make scanning fast, output clean, and leave room for extra modules like CVE checks and brute-force logic.

---

## Features

- Multi-threaded scanning (default: 2000 threads)
- Detects both:
  - **Classic VNC** (RFB protocol on common ports)
  - **HTTP-based noVNC** web clients
- CIDR and target list support
- Banner grab for RFB services
- noVNC detection with optional favicon hashing
- Saves results in JSON lines format
- Optional **local GUI** to view results in a browser
- Adjustable timeouts, thread count, and IP limits

---

## Installation

You’ll need **Go 1.21+**.

```bash
git clone https://github.com/YourUser/VNCReaper.git
cd VNCReaper
go build -o vncreaper ./cmd/vncreaper
```

---

## Usage

Basic CLI scan:
```bash
./vncreaper -i targets.txt -o results.json
```

GUI mode:
```bash
./vncreaper -gui -o results.json -port 7777
```
Then open `http://localhost:7777` in your browser.

---

### Options

| Flag | Description | Default |
|------|-------------|---------|
| `-i` | Input file with IPs or CIDRs | `targets.txt` |
| `-o` | Output JSON lines file | `results.json` |
| `-gui` | Run in GUI mode | `false` |
| `-port` | Port for GUI mode | `7777` |
| `-limit` | Limit total IPs scanned | `0` (unlimited) |
| `-t` | Number of concurrent threads | `2000` |
| `-timeout` | Connection timeout (ms) | `800` |
| `-http-only` | Skip TCP RFB checks, scan HTTP/noVNC only | `false` |
| `-no-favicon` | Skip favicon hashing | `false` |
| `-novnc-disabled` | Disable noVNC detection | `false` |
| `-v` | Verbose output | `false` |

---

## Output Format

Each result is saved as a single JSON object per line:

```json
{
  "ip": "192.168.1.10",
  "port": "5900",
  "protocol": "RFB",
  "banner": "RFB 003.008",
  "title": "",
  "favicon_hash": "",
  "status": "open",
  "timestamp": "2025-08-09T14:32:00Z",
  "raw": ""
}
```

---

## How it Works

1. Expands all targets from file (`IP:PORT` or CIDR ranges).
2. Scans each target/port in parallel.
3. If not in `http-only` mode:
   - Connects via TCP to check for RFB (VNC) protocol banner.
4. If noVNC detection is enabled:
   - Sends HTTP requests, looks for noVNC strings and optionally hashes favicon.
5. Logs all open services to JSON and (if GUI mode) displays them in a simple web interface.

---

## License

MIT License — use it, modify it, share it. You’re responsible for what happens.
