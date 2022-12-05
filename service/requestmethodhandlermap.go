package service

import (
    "fmt"
)

// map[request method]service
type RequestMethodHandlerMap map[string]Service

func (rmhm RequestMethodHandlerMap) AddHandler(method string, service Service) (err error) {
    // the `method` isn't validated as it doesn't really matter, and leaving it
    // unvalidated means we have a greater breadth of call support

    if service == nil {
        err = fmt.Errorf("%w: Cannot register a nil service", ErrInvalidService)
        return
    }

    _, ok := rmhm[method]
    if ok {
        err = fmt.Errorf("%w: Method %s", ErrRequestHandlerAlreadyRegister, method)
        return
    }

    rmhm[method] = service
    return
}
