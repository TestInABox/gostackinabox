package common

import (
    "net/http"
    "io"
)

type HttpReply struct {
    Status       HttpStatusCode
    Headers      http.Header
    Trailers     http.Header
    ResponseData io.ReadCloser
    Length       int64
}
