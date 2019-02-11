# boundedparallelism
This code is intended to demonstrate bounded parallelism for operations to be performed concurrently in batches. In real world, this pattern can be applied in applications that perform large number of http requests and number of connections got to be constant.


## Pre-requisites
Install golang, configure GOPATH etc

## Steps
1. cd $GOPATH/src
2. git clone https://github.com/ajithcnambiar/boundedparallelism
3. cd boundedparallelism
4. go run test.go

This will run a server listening at 6060. A curl on port 6060 will do an iteration of the test. In the test, input numbers from 0 to 9 are pushed to Stage 1 Channel. The stage has 5 functions that reads input from Stage 1 Channel. Stage 2 performs opertion on each input(add 20 to each input), and writes result into set of Channels. Stage 3 consumes these Channels and merge the result into a single Channel. Check for server logs for details.
```
curl 'http://localhost:6060/
```
The reason for test server is to allow go profiling tests using pprof https://golang.org/pkg/net/http/pprof/

## Design


## References
Code is based on bounded parallelism as explained in the blog https://blog.golang.org/pipelines



