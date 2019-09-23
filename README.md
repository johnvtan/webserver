# Description
A really simple HTTP web server written in Go without using the net/http package. The goal of this
project is to gain a better understanding of how HTTP works.

This currently only supports `GET` requests, and returns a 501 error for any other request.

The server is hard coded to run on port 8081.

# Running
```golang
go run webserver.go
```

Then, you can see the output by either typing in `localhost:8081` into your browser or doing:
```
curl -i -X GET localhost:8081
```
