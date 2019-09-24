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

func fileInDir(name string) bool {
    files, _ := ioutil.ReadDir(".")
    for _, f := range files {
        if name == f.Name() {
            return true
        }
    }

    return false
}

func getUriAndCode(uri string) (string, string) {
    var retUri string
    var code string
    if uri == "/" {
        retUri = "/hello.html"
        code = "200"
    } else if fileInDir(uri[1:]) {
        retUri = uri
        code = "200"
    } else {
        retUri = "/404.html"
        code = "404"
    }
    return retUri, code
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

    switch req.Method {
    case "GET":
        uri, code := getUriAndCode(req.URI)
        conn.Write([]byte(http.ServeGet(uri, code)))
    case "HEAD":
        uri, code := getUriAndCode(req.URI)
        conn.Write([]byte(http.ServeHead(uri, code)))
    default:
        conn.Write([]byte(http.Serve501()))
    }
}
