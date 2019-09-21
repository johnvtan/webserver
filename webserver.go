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

    fmt.Println("Closing listener")
    ln.Close();
}

// thread per request/response
func handleConnection(conn net.Conn) {
    reader := bufio.NewReader(conn)
    req, err := http.ParseRequest(reader)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Closing connection")
    fmt.Println(req)
    conn.Close()
}
