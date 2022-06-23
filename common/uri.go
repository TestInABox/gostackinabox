package common

import (
    "errors"
    "net/url"
)

type URI interface {
    IsMatch(url.URL) (bool, error)
}

var (
    ErrServerURIMisconfigured error = errors.New("Misconfigured ServerURI")
    ErrPathURIMisconfigured error = errors.New("Misconfigured PathURI")
)
