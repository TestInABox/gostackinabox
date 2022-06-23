package router

import (
    "errors"
)

var (
    ErrResponseBuildingInternalError error = errors.New("Service Router: Internal Error while building response")
    ErrServiceHandlerAlreadyRegister error = errors.New("Service Router: Already Registered")
    ErrInvalidRequest error = errors.New("Service Router: Invalid Request")
)
