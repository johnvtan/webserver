package http

import (
    "bufio"
    "fmt"
    "io"
    "io/ioutil"
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

type HttpResponse struct {
    Code string
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

        // headerFields[1] may contain colons
        headerFields := strings.SplitN(headerLine, ":", 2)
        if len(headerFields) != 2 {
            return request, fmt.Errorf("http.ParseRequest: Bad header line: '%s'", headerLine)
        }

        // fieldName is case insensitive
        fieldName := strings.ToLower(strings.TrimSpace(headerFields[0]))
        fieldVal := strings.TrimSpace(headerFields[1])
        request.Headers[fieldName] = fieldVal
    }

    // according to rfc 2616, the presence of a message-body is signaled by the inclusion of a
    // content-length or transfer-encoding header field. For now, only support content-length.
    if lengthStr, ok := request.Headers["content-length"]; ok {
        bodyLen, err := strconv.Atoi(lengthStr)
        if err != nil {
            return request, err
        }
        if bodyLen < 0 {
            return request, fmt.Errorf("http.ParseRequest: Invalid content length: '%d'", bodyLen)
        }
        bodyBuf := make([]byte, bodyLen)
        if _, err := io.ReadFull(reader, bodyBuf); err != nil {
            return request, err
        }
        request.Body = string(bodyBuf)
    }

    return request, nil
}

func responseToString(response HttpResponse) string {
    var ret strings.Builder
    ret.WriteString("HTTP/1.1 ")
    switch response.Code {
    case "200":
        ret.WriteString("200 OK")
    case "404":
        ret.WriteString("404 NOT FOUND")
    case "501":
        ret.WriteString("501 NOT IMPLEMENTED")
    default:
        panic("http.responseToString: Got unrecognized HTTP code")
    }
    ret.WriteString("\r\n")

    for key, val := range response.Headers {
        ret.WriteString(key)
        ret.WriteString(": ")
        ret.WriteString(val)
        ret.WriteString("\r\n")
    }
    ret.WriteString("\r\n")
    ret.WriteString(response.Body)
    ret.WriteString("\r\n")
    return ret.String()
}

func ServeAddress(uri string, code string) string {
    var response HttpResponse
    response.Code = code

    // slice off leading / character
    body, err := ioutil.ReadFile(uri[1:])
    if err != nil {
        panic("http.ServeAddress: File not found")
    }

    response.Body = string(body)

    response.Headers = make(map[string]string)
    response.Headers["Content-Type"] = "text/html; charset=UTF-8"
    return responseToString(response)
}

func Serve501() string {
    var response HttpResponse
    response.Code = "501"
    response.Headers = make(map[string]string)
    return responseToString(response)
}
