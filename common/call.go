package common

import (
    "net/http"
    "net/url"
)

type HttpCall struct {
    // most of these fields are just
    // easy access to the request information
    Method  HttpVerb
    Url     *url.URL
    Headers http.Header
    Request *http.Request
}
