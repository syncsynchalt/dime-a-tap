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
