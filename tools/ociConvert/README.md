#OCI Convert
- [From aci to oci](#how-to-try)

##Build
```
go get -d github.com/huawei-openlab/oct/tools/ociConvert
cd $GOPATH/src/github.com/huawei-openlab/oct/tools/ociConvert
make
```

##How To try
It is easy to use this tool, we provide a `demo rkt file` with an `image.json and pod_runtime.json` file,
both of them are downloaded from [AppC Example](https://github.com/appc/spec/tree/master/examples)

```
./ociConvert a2i test/image.json test/pod_runtime.json

```
