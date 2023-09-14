# Network Usage Prometheus Exporter

Simple Prometheus exporter for \*nix environments with support for `ip -s -s link show`.

The output should be similar to the following:

```
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN mode DEFAULT group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    RX: bytes  packets  errors  dropped overrun mcast
    1484       8        0       0       0       0
    RX errors: length   crc     frame   fifo    missed
               0        0       0       0       0
    TX: bytes  packets  errors  dropped carrier collsns
    1484       8        0       0       0       0
    TX errors: aborted  fifo   window heartbeat transns
               0        0       0       0       0
```

The command output is parsed when requested and presented on the webpage at path `/metrics`. The parsed output is cached for a default of 30 seconds, but is configurable.

Was originally designed to be run on an Ubiquiti EdgeRouter Lite ERLite-3 with a MIPS64 processor.

### Command Line Arguments

| Short | Description |
| --- | --- |
| `-t` | Cache timeout period in seconds (Default 30) |
| `-h` | Webserver host (Default empty) |
| `-p` | Webserver port (Default 15835) |
