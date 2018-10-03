# dime-a-tap

Make TLS/SSL traffic readable in the clear by using a MITM proxy.

### Getting started

To start the proxy on port 443 (HTTPS):

```
go get github.com/syncsynchalt/dime-a-tap/cmds/dime-a-tap
export PATH=$PATH:~/go/bin
dime-a-tap 443
```

Now use /etc/hosts, a captive DNS server, or similar to redirect hosts and devices to your proxy for a given hostname.

### Creating a certificate store

To avoid security warnings you'll need to create a CA and distribute it to your devices:

```
dime-a-tap ca-init /tmp/cadir
dime-a-tap -cadir /tmp/cadir 443
```

You can create `ca.key` and `ca.crt` yourself using `openssl` or similar tools if you prefer.

Install the file `/tmp/cadir/ca.crt` on your hosts or devices as a trusted CA.
