package http

import (
    "bufio"
    "fmt"
    "strconv"
    "strings"
)

type HttpRequest struct {
    Method string
    URI string
    VersionMajor int
    VersionMinor int
    Headers map[string]string
    Body string
}

func ParseRequest(reader *bufio.Reader) (HttpRequest, error) {
    var request HttpRequest

    // parse the first line -
    requestLine, err := reader.ReadString('\n')
    if err != nil {
        return request, err
    }

    fmt.Println(requestLine)
    requestLineFields := strings.Fields(requestLine)
    if len(requestLineFields) != 3 {
        return request, fmt.Errorf("http.ParseRequest: Malformed request line - %s", requestLine)
    }

    // parse request type - only supports get and post for now
    if requestLineFields[0] == "GET" || requestLineFields[0] == "POST" {
        request.Method = requestLineFields[0]
    } else {
        return request, fmt.Errorf("http.ParseRequest: Unimplemented request type %s", requestLineFields[0])
    }

    // Don't do anything with the URI for now
    request.URI = requestLineFields[1]

    // parse HTTP version
    httpVersion := strings.Split(requestLineFields[2], "/")
    if len(httpVersion) != 2 || httpVersion[0] != "HTTP" {
        return request, fmt.Errorf("http.ParseRequest: Bad HTTP version in request line: %s", requestLineFields[2])
    }


    versionNums := strings.Split(httpVersion[1], ".")
    if len(versionNums) != 2 {
        return request, fmt.Errorf("http.ParseRequest: Bad HTTP version number: %s", httpVersion[1])
    }

    request.VersionMajor, err = strconv.Atoi(versionNums[0])
    if err != nil {
        return request, err
    }

    request.VersionMinor, err = strconv.Atoi(versionNums[1])
    if err != nil {
        return request, err
    }

    // parse header lines
    request.Headers = make(map[string]string)
    for headerLine, err := reader.ReadString('\n'); headerLine != "\r\n"; headerLine, err = reader.ReadString('\n') {
        if err != nil {
            return request, err
        }
        headerFields := strings.SplitN(headerLine, ":", 2)
        if len(headerFields) != 2 {
            return request, fmt.Errorf("http.ParseRequest: Bad header line: '%s'", headerLine)
        }
        headerKey := strings.TrimSpace(headerFields[0])
        headerVal := strings.TrimSpace(headerFields[1])
        request.Headers[headerKey] = headerVal
        fmt.Println(headerLine)
    }

    return request, nil
}
