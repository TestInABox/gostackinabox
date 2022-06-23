package common

type HttpStatusCode int

const (
    HttpStatus_MethodNotSupport     HttpStatusCode = 405
    HttpStatus_RouteNotHandled      HttpStatusCode = 595
    HttpStatus_ServiceError         HttpStatusCode = 596
    HttpStatus_ServiceSubRouteError HttpStatusCode = 597
)

func GetHttpStatus(code int) HttpStatusCode {
    return HttpStatusCode(code)
}
