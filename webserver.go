package main

import (
    "net"
    "fmt"
    "bufio"
    "log"
    "webserver/http"
)

func main() {
    fmt.Println("Starting Server...")
    ln, err := net.Listen("tcp", ":8081")
    defer ln.Close()
    if err != nil {
        log.Fatal(err)
    }

   // loop forever sending back the msg
    for {
        // just keep waiting for connections and handling them in a separate thread
        conn, err := ln.Accept()
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println("got a connection")
        go handleConnection(conn)
    }
}

// thread per request/response
func handleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    req, err := http.ParseRequest(reader)
    if err != nil {
        fmt.Println(err)
        return
    }

    code := "200"
    switch req.Method {
    case "GET":
        if req.URI == "/" {
            req.URI = "/hello.html"
        } else if req.URI != "/hello.html" {
            code = "404"
            req.URI = "/404.html"
        }
        conn.Write([]byte(http.ServeAddress(req.URI, code)))
    default:
        conn.Write([]byte(http.Serve501()))
    }
}
