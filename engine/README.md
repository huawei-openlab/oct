# Testing Cases And Framework For Open Container Specifications


The oc-testing project aims to promote the Open Container Initiative by providing a universal testing framewrk for all the container projects. The testing includes: oci specification, container function and container performance.

## `oci.conf` Example

```
runmode = dev

listenmode = https
httpscertfile = cert/containerops/containerops.crt
httpskeyfile = cert/containerops/containerops.key


[log]
filepath = log/containerops-log

```

## `runtime.conf` Example

```
runmode = dev

listenmode = https
httpscertfile = cert/containerops/containerops.crt
httpskeyfile = cert/containerops/containerops.key

[log]
filepath = log/containerops-log
```

