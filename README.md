# Netsnitch

![Go](https://img.shields.io/badge/Go-1.22-blue)
![Status](https://img.shields.io/badge/status-active-orange)

**A modular high-performance network scanner built with Go.**

Netsnitch uses a scheduler and worker pool architecture to execute scan
tasks concurrently and identify services using a fingerprint engine.

## 🚀 Features

- **Concurrent scanning engine using goroutines**
- **Scheduler with worker pool design**
- **Modular scan task builders**
- **Service fingerprinting based on probes and pattern matching**
- **Configurable scanning via CLI flags**
- **Designed for extensibility and future AI analysis**

## 🔎 Supported Scan Types
  
- ARP-based host discovery
- TCP port scanning (SYN / full handshake)
- UDP port scanning
- More scan modules coming soon

## ⚙️  Installation
```bash
git clone https://github.com/Hugokoks/netsnitch
cd netsnitch
go build -o netsnitch ./cmd/netsnitch
```

## 🧩 Extending the Framework

Netsnitch is designed to be easily extensible through self-registered scan modules.

To add a new scan type, a new package can be created inside the `scans` directory.
Each scan module follows a simple structure based on shared input, output, and task
packages.

Scan execution is handled through the task abstraction. The core engine does not
need to know scan-specific rules — it simply executes tasks via the `Task.Execute`
method.

This design allows new scan implementations to be added without modifying
the core framework logic.

## 💻 Scan Usage

**ARP Scan:**

- Example: 
```bash 
sudo ./netsnitch arp 192.168.1.0/24 -t 500ms
```

- Available Flags: -r,-t

**TCP Port Scan:**

- Example:
```bash
sudo ./netsnitch tcp 192.168.1.1 -p 1-1024 -m s -o
```

- Available Flags: -m, -o, -p, -r, -t

**UDP Port Scan:**

- Example: 
```bash
sudo ./netsnitch udp 192.168.1.1 -p 53,161 -t 500ms
```

- Available Flags: -p, -r, -t

**Parallel Scans:**

```bash
sudo ./netsnitch arp 192.168.0.0/24 -t 1s "&&" tcp 192.168.0.1 -p 1-100
```

## 💻 Flags Usage

**-h:**
  
- usage:       -h
  
- description: Show help and exit.

**-m:**

- usage:       -m (f | s)

- description: Scan strategy: 'f' for full handshake, 's' for stealth SYN.

**-o:**

- usage:       -o

- description: Show only active/open ports; hide everything else.

- default:     show all

**-p:**

- usage:       -p (1-100 | 80,443)

- description: Target ports: supports ranges, lists, or single ports.

**-r:**

- usage:       -r (raw | json)

- description: Output format: table rows or machine-readable JSON.

- default:     raw

**-t:**

- usage:       -t (1s | 500ms)

- description: Network timeout: higher value means better accuracy on slow links.

- default:     1s

**Port Selection:**

- -p all           
- -p 1-100         
- -p 1,2,3,4,5,6   

**Target Selection:**

- 192.168.0.5                 
- 192.168.0.1,192.168.0.5     
- 192.168.0.0/24              

**Timeout Selection:**

- -t 200ms    
- -t 2s       
- -t 1m       

## 🏁 Roadmap

- Web UI dashboard
- AI-assisted scan analysis
- Expanded fingerprint database
- Distributed scanning
