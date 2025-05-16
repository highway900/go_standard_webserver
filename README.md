# Simple Golang HTTP service

A set of simple HTTP services written using the standard library in go. I wrote this to show a new go team what what go services is like and how we go from simple to more complex-ish and what that does for the implementation.

Run samples
===========

install the [go compiler](https://golang.org)

1. `go run cmd/simple1/main.go`
2. `go run cmd/simple2/main.go`
3. `go run cmd/simple3/main.go`

Simple1
=======

Pretty dead simple handler and server setup. Just return a string if the request method is GET

Simple2
=======

Introduce some additional handlers and handle some JSON input and output

Simple3
=======

+ Use multiplexing to route handlers
+ Wrap handlers in some middleware functions
+ Use context to pass a request scoped requestID through the handler cycle
+ Try to gracefully handle server timeouts and interrupts / cancellations using channels
+ Use HTTPS

### Generate SSL certs for HTTPS

Create the key file for serving with TLS

```
openssl req  -new  -newkey rsa:2048  -nodes  -keyout localhost.key  -out localhost.csr
openssl  x509  -req  -days 365  -in localhost.csr  -signkey localhost.key  -out localhost.crt
```

## Additional Resources

+ [Way better article on building web services in go](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/)
  - NOTE: I have added this as this repo is now a few years old. It was nice to see not much has changed in the go eco http handler ecosystem.
+ [~~Great article on building a full service~~](https://medium.com/rungo/creating-a-simple-hello-world-http-server-in-go-31c7fd70466e)
+ [Testing HTTP handlers in GO](https://blog.questionable.services/article/testing-http-handlers-go/)
