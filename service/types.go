package service

import (
    "errors"
    "net/url"

    "github.com/TestInABox/gostackinabox/common"
)

type Service interface {
    IsSubService() bool
    GetName() string
    GetMatcher() common.URI
    GetHandler(u url.URL) (common.HttpHandler, error)
    RegisterHandler(subHandler Service) error
    RegisterMethodHandler(method common.HttpVerb, handler common.HttpHandler) error
    Init(name string, matcher common.URI) error
}

// map[service name]service
type ServiceHandlerMap map[string]Service


var (
    ErrInvalidService error = errors.New("Service handler is invalid")
    ErrInvalidServiceRegex error = errors.New("Invalid service regex")
    ErrRequestHandlerInvalid error = errors.New("Service: Method Handler Invalid")
    ErrRequestHandlerAlreadyRegister error = errors.New("Service: Method Already Registered")
    ErrServiceHandlerAlreadyRegister error = errors.New("Service: Handler Already Registered")
    ErrNoHandlerFunc error = errors.New("No handler func")

    ErrNotImplemented error = errors.New("Not implemented")
)
