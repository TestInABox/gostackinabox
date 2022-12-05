package client

import (
    "net/http"

    "github.com/TestInABox/gostackinabox/router"
)

var DefaultClient = &http.Client{
    Transport: &router.Router{},
}
