package main

import (
    "net"
    "fmt"
    "bufio"
    "log"
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
        fmt.Println("got a connectoin")
        handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) error {
    reader := bufio.NewReader(conn)

    for {
        message, _ := reader.ReadString('\n')
        fmt.Print("Server got: ", string(message))
        new_msg := "HTTP/1.1 200 OK \r\n"
        conn.Write([]byte(new_msg + "\r\n"))
    }
    return nil
}
