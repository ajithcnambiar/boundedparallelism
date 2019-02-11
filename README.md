# boundedparallelism
This code is intended to demonstrate bounded parallelism for concurrent http requests. In real world, this pattern can be applied in applications that perform large number of http requests and number of connections got to be constant.


## Pre-requisites
Install golang, configure GOPATH etc

## Steps
1. mkdir $GOPATH/src/concurrentHTTP && cd $GOPATH/src/concurrentHTTP
2. git clone https://github.com/ajithcnambiar/boundedparallelism
2. cd boundedparallelism
3. go run test.go

This will run a server listening at 6060. A curl on port 6060 will do an iteration of the test and server logs the http status.
```
curl 'http://localhost:6060/
```
The reason for test server is to allow go profiling tests using pprof https://golang.org/pkg/net/http/pprof/

## Design


## References
Code is based on bounded parallelism as explained in the blog https://blog.golang.org/pipelines



