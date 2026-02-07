# Netsnitch

Netsnitch is a modular network scanning framework written in Go.
It is designed to explore low-level networking, concurrency, and clean
system-oriented architecture.

The core idea behind Netsnitch is to provide a flexible scanning engine
that can execute different scan types without being tightly coupled
to their internal logic.

----

## Key Features

- Modular scan architecture (easily extensible)
- Task-based execution engine
- Concurrent scan execution using goroutines
- Supports running multiple scan strategies in parallel
- Clean separation between scan engine and scan implementations
- Linux-focused networking

- ## Supported Scan Types
- 
- ARP-based host discovery
- TCP port scanning
- more comming...

----

## Extending the Framework

Netsnitch is designed to be easily extensible through self-registered scan modules.

To add a new scan type, a new package can be created inside the `scans` directory.
Each scan module follows a simple structure based on shared input, output, and task
packages.

Scan execution is handled through the task abstraction. The core engine does not
need to know scan-specific rules — it simply executes tasks via the `Task.Execute`
method.

This design allows new scan implementations to be added without modifying
the core framework logic.

------

## Usage

Example commands:

###ARP Scan:
sudo go run cmd/netsnitch/main.go arp 192.168.0.0/24

###TCP port Scan:
sudo go run cmd/netsnitch/main.go tcp --ports:1-100 192.168.0.5

###CIDR notation:
sudo go run cmd/netsnitch/main.go tcp --ports:1-100 192.168.0.0/24

Port Selection:
--ports:all           scan all ports
--ports:1-100         port range
--ports:1,2,3,4,5,6   explicit port list

without --ports flag it will scan predifine ports, settigns are getting from domain/config.go

###Parallel Scans:
sudo go run cmd/netsnitch/main.go arp 192.168.0.0/24 "&&" tcp --ports:1-100 192.168.0.1

###Target Selection:
192.168.0.5             single IP address
192.168.0.1,192.168.0.5 multiple IP addresses
192.168.0.0/24          CIDR network
