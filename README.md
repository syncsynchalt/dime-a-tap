# dime-a-tap

MITM proxy to make TLS/SSL traffic readable in the clear.

Unencrypted traffic is sent over loopback to make it easily tcpdumpable.

### Getting started

To start the proxy on port 443 (HTTPS):

```
$ go get github.com/syncsynchalt/dime-a-tap/cmds/dime-a-tap
$ export PATH=$PATH:~/go/bin
$ dime-a-tap 443
```

Use /etc/hosts, captive DNS, or similar to redirect hosts and devices to your proxy for a given hostname.

### Creating a certificate store

To avoid security warnings you'll want to create a CA and distribute it to your devices:

```
$ dime-a-tap ca-init /tmp/cadir
$ dime-a-tap -cadir /tmp/cadir 443
```

Install the certificate in `/tmp/cadir/ca.crt` as a trusted CA on your hosts or devices.

### Capturing the unencrypted data

To capture intercepted data, there are two options.

Use `-capturedir {dir}` to write the unencrypted client (.c) and server (.s) conversation to files in that dir. Example:
```
$ mkdir /tmp/captures
$ dime-a-tap -capturedir /tmp/captures 443 &
(send traffic through the tap)
$ ls /tmp/captures
  total 56
  -rw-r--r--  1 user  wheel   75 Oct  4 12:45 192.168.69.42:52981.20181004124516.667781.c
  -rw-r--r--  1 user  wheel  756 Oct  4 12:45 192.168.69.42:52981.20181004124516.733675.s
  -rw-r--r--  1 user  wheel    0 Oct  4 12:45 192.168.69.42:52981.20181004124516.735306.c
  -rw-r--r--  1 user  wheel   75 Oct  4 12:45 192.168.69.42:52989.20181004124551.808247.c
  -rw-r--r--  1 user  wheel  756 Oct  4 12:45 192.168.69.42:52989.20181004124551.875861.s
  -rw-r--r--  1 user  wheel    0 Oct  4 12:45 192.168.69.42:52989.20181004124551.877488.c
  -rw-r--r--  1 user  wheel   75 Oct  4 12:46 192.168.69.42:52992.20181004124609.494528.c
  -rw-r--r--  1 user  wheel  297 Oct  4 12:46 192.168.69.42:52992.20181004124609.554621.s
  -rw-r--r--  1 user  wheel  459 Oct  4 12:46 192.168.69.42:52992.20181004124609.555327.s
  -rw-r--r--  1 user  wheel    0 Oct  4 12:46 192.168.69.42:52992.20181004124609.556733.c
```

Or use `tcpdump` on localhost:4430 to create a pcap file suitable for use with wireshark.  Example:
```
$ dime-a-tap 443 &
$ tcpdump -i lo0 -s 0 -w capture.pcap port 4430
(send traffic through the tap)
```
