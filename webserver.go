package main

import (
    "io/ioutil"
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

    code := "404"
    switch req.Method {
    case "GET":
        if req.URI == "/" {
            req.URI = "/hello.html"
        }

        // check if the file exists in this directory
        files, err := ioutil.ReadDir(".")
        if err != nil {
            fmt.Println(err)
            return
        }
        for _, f := range files {
            if req.URI[1:] == f.Name() {
                code = "200"
                break
            }
        }

        // if we still haven't found it, set URI to 404 
        if code == "404" {
            req.URI = "/404.html"
        }
        conn.Write([]byte(http.ServeAddress(req.URI, code)))
    default:
        conn.Write([]byte(http.Serve501()))
    }
}
