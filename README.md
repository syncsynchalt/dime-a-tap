# dime-a-tap

MITM proxy to make TLS/SSL traffic readable in the clear.

Unencrypted traffic is sent over loopback to make it easily tcpdumpable.

### Getting started

To start the proxy on port 443 (HTTPS):

```
go get github.com/syncsynchalt/dime-a-tap/cmds/dime-a-tap
export PATH=$PATH:~/go/bin
dime-a-tap 443
```

Now use /etc/hosts, captive DNS, or similar to redirect hosts and devices to your proxy for a given hostname.

### Creating a certificate store

To avoid security warnings you'll want to create a CA and distribute it to your devices:

```
dime-a-tap ca-init /tmp/cadir
dime-a-tap -cadir /tmp/cadir 443
```

Install the file `/tmp/cadir/ca.crt` on your hosts or devices as a trusted CA.

(You can create `ca.key` and `ca.crt` yourself using `openssl` or similar tools if you prefer to customize the CA certificate).

### Capturing the unencrypted data

To capture intercepted data, there are two options.

Use `-capturedir {dir}` to write reads and writes to files in that dir. Example:
```
mkdir /tmp/captures
dime-a-tap -capturedir /tmp/captures 443 &
(send traffic through the tap)
ls /tmp/captures
  total 56
  -rw-r--r--  1 user  wheel   75 Oct  4 12:45 127.0.0.1:52981.20181004124516.667781.c
  -rw-r--r--  1 user  wheel  756 Oct  4 12:45 127.0.0.1:52981.20181004124516.733675.s
  -rw-r--r--  1 user  wheel    0 Oct  4 12:45 127.0.0.1:52981.20181004124516.735306.c
  -rw-r--r--  1 user  wheel   75 Oct  4 12:45 127.0.0.1:52989.20181004124551.808247.c
  -rw-r--r--  1 user  wheel  756 Oct  4 12:45 127.0.0.1:52989.20181004124551.875861.s
  -rw-r--r--  1 user  wheel    0 Oct  4 12:45 127.0.0.1:52989.20181004124551.877488.c
  -rw-r--r--  1 user  wheel   75 Oct  4 12:46 127.0.0.1:52992.20181004124609.494528.c
  -rw-r--r--  1 user  wheel  297 Oct  4 12:46 127.0.0.1:52992.20181004124609.554621.s
  -rw-r--r--  1 user  wheel  459 Oct  4 12:46 127.0.0.1:52992.20181004124609.555327.s
  -rw-r--r--  1 user  wheel    0 Oct  4 12:46 127.0.0.1:52992.20181004124609.556733.c
```

Or use `tcpdump` on 127.0.0.1:4430 to create a pcap file suitable for use with wireshark etc.  Example:
```
dime-a-tap 443 &
tcpdump -i lo0 port 4430 -o capture.pcap &
(send traffic through the tap)
```
