# Pachyderm language clients

## Go Client

The Go client is officially supported by the Pachyderm team.  It implements almost all of the functionality that is provided with the `pachctl` CLI tool, and, thus, you can easily integrated operations like `put file` into your applications.

For more info, check out the [godocs](https://godoc.org/github.com/pachyderm/pachyderm/src/client).

**Note** - A compatible version of `grpc` is needed when using the Go client.  You can deduce the compatible version from our [vendor.json](https://github.com/pachyderm/pachyderm/blob/master/src/server/vendor/vendor.json) file, where you will see something like:

```
		{
			"checksumSHA1": "mEyChIkG797MtkrJQXW8X/qZ0l0=",
			"path": "google.golang.org/grpc",
			"revision": "21f8ed309495401e6fd79b3a9fd549582aed1b4c",
			"revisionTime": "2017-01-27T15:26:01Z"
		},
```

You can then get this version via:

```
go get google.golang.org/grpc
cd $GOPATH/src/google.golang.org/grpc
git checkout 21f8ed309495401e6fd79b3a9fd549582aed1b4c
```

### Running Go Examples

The Pachyderm [godocs](https://godoc.org/github.com/pachyderm/pachyderm/src/client) reference
provides examples of how you can use the Go client API. You need to have a running Pachyderm cluster
to run these examples.

Make sure that you use your `pachd_address` in `client.NewFromAddress("<your-pachd-address>:30650")`.
For example, if you are testing on `minikube`, run
`minikube ip` to get this information.

See the [OpenCV Example in Go](https://github.com/pachyderm/pachyderm/tree/master/examples/opencv) for more
information.

## Python Client

The Python client is officially supported by the Pachyderm team. 
It implements almost all of the functionalities provided with the `pachctl` CLI tool allowing you to easily integrate operations like `create repo`, `put a file,` or `create pipeline` into your python applications.

For more info, check out our Github [repo](https://github.com/pachyderm/python-pachyderm/blob/master/README.md).

## Scala Client

Our users are currently working on a Scala client for Pachyderm. Please contact us if you are interested in helping with this or testing it out.

## Other languages

Pachyderm uses a simple [protocol buffer API](https://github.com/pachyderm/pachyderm/blob/master/src/pfs/pfs.proto). Protobufs support [a bunch of other languages](https://developers.google.com/protocol-buffers/), any of which can be used to programmatically use Pachyderm. We haven’t built clients for them yet, but it’s not too hard. It’s an easy way to contribute to Pachyderm if you’re looking to get involved.
