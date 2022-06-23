package router

import (
    "fmt"
    "net/http"

    "github.com/TestInABox/gostackinabox/common"
    "github.com/TestInABox/gostackinabox/common/log"
    "github.com/TestInABox/gostackinabox/service"
    "github.com/TestInABox/gostackinabox/util"
)

type Router struct {
    ProtoMajor int
    ProtoMinor int
    RequestHandlers service.ServiceHandlerMap // TODO: Update
    DisableCompression bool
}

func New() *Router {
    return &Router{
        ProtoMajor: 1,
        ProtoMinor: 1,
        RequestHandlers: make(service.ServiceHandlerMap),
        DisableCompression: true, // compression isn't handled by this transport router
    }
}

// service name is the scheme + host and optionally the port portion of a URL
func (irt *Router) RegisterService(service string, handler service.Service) (err error) {
    log.Printf("Attempting to register service %s with handler %v", service, handler)
    if existing, ok := irt.RequestHandlers[service]; ok {
        log.Printf("Service %s already registered using handler %v", service, existing)
        err = fmt.Errorf("%w: Service %s already registered", ErrServiceHandlerAlreadyRegister, service)
        return
    }
    log.Printf("Accepting registration of service %s using handler %v", service, handler)

    irt.RequestHandlers[service] = handler
    return
}

func (irt *Router) BuildResponse(reply *common.HttpReply, request *http.Request) (response *http.Response, err error) {
    if reply == nil || request == nil {
        log.Printf("Recieved invalid parameter: Reply: %v, Request: %#v", reply, request)
        err = fmt.Errorf(
            "%w: Invalid service response or systems error (reply: %#v, request: %#v",
            ErrResponseBuildingInternalError,
            reply,
            request,
        )
        return
    }

    log.Printf("Building Reply: Status: %d, Data Length: %d", reply.Status, reply.Length)

    intStatus := int(reply.Status)
    response = &http.Response{
        Status: fmt.Sprintf("%d", intStatus),
        StatusCode: intStatus,
        Proto: fmt.Sprintf("HTTP/%d.%d", irt.ProtoMajor, irt.ProtoMinor),
        ProtoMajor: irt.ProtoMajor,
        ProtoMinor: irt.ProtoMinor,
        Header: reply.Headers,
        Body: reply.ResponseData,
        ContentLength: reply.Length,
        Trailer: reply.Trailers,
        Request: request,
        Uncompressed: false, // this doesn't do anything with compression
        TLS: nil,  // this doesn't implement TLS at all; it's all unencrypted
    }
    return
}

func (irt *Router) ServiceRouter(request *http.Request) (response *http.Response, err error) {
    // 1. build a response object
    // 2. call into a handler method providing the request and response

    // Reservered Errors:
    //  405: Method Not Supported
    //      URI handled but HTTP Method is not
    //  595: Route Not Handled
    //      http://127.0.0.1:80/uri is completely unhandled
    //  596: Service Error
    //      Service handling the request had an error
    //  597: Service doesn't handle a sub-route
    //      Service handles part of but not the entire route
    requestUrl := request.URL
    if requestUrl == nil {
        log.Printf("Request is invalid because the URL attribute is nil: %#v", request)
        err = fmt.Errorf("%w: Request has a nil URL", ErrInvalidRequest)
        return
    }
    // Normally it's an error to set these; however we need to set them for the stack to work properly
    if len(request.RequestURI) == 0 {
        request.RequestURI = "/"
    }
    if len(request.URL.Path) == 0 {
        request.URL.Path = "/"
        request.URL.RawPath = "/"
    }


    log.Printf("Attempting to handle request: Method: %s RequestURI: \"%s\"", request.Method, request.RequestURI)
    // is there a handler for the URI?
    for serviceName, serviceHandler := range irt.RequestHandlers {
        log.Printf("Attempting to match Service %s against URL \"%s\"", serviceName, request.RequestURI)
        // see if this service handles the URL
        matcher := serviceHandler.GetMatcher()
        matchResult, matchErr := matcher.IsMatch(*requestUrl)
        if matchErr != nil {
            // there's a problem with the matcher, test infrastructure needs to be fixed
            log.Printf("Matcher for Service %s generated an error. Please fix the fixture: %#v", serviceName, matchErr)
            err = fmt.Errorf("Service %s generated an error: %w", serviceName, matchErr)
            return
        }
        if matchResult {
            log.Printf("Service %s handles URI %s", serviceName, request.RequestURI)
            // get the handler for the service
            handler, handlerErr := serviceHandler.GetHandler(*requestUrl)
            if handlerErr != nil {
                log.Printf("Service %s generated an error when retrieving the handler: %#v", serviceName, handlerErr)
                err = handlerErr
                return
            }

            log.Printf("Running handler for Service %s on URI %s", serviceName, request.RequestURI)
            // attempt to let the registered service handle it
            reply, err := handler(
                &common.HttpCall{
                    Method: common.HttpVerb(request.Method),
                    Url: request.URL,
                    Headers: request.Header,
                    Request: request,
                },
            )

            // service had an error
            if err != nil {
                log.Printf("Service %s generated an error while handling the request: %#v", serviceName, err)
                msg := fmt.Sprintf(
                    "gostackinabox: service handling request had an error - %#v",
                    err,
                )
                return irt.BuildResponse(
                    &common.HttpReply{
                        Status: common.HttpStatus_ServiceError,
                        ResponseData: util.StringToResponseBody(msg),
                        Length: int64(len(msg)),
                    },
                    request,
                )
            }

            log.Printf("Service %s generated a successful response", serviceName)
            // send back the reply from the service
            return irt.BuildResponse(reply, request)
        }
    }

    // return 595
    msg := fmt.Sprintf(
        "gostackinabox: no service to handle URL '%s'",
        request.URL.String(),
    )
    return irt.BuildResponse(
        &common.HttpReply{
            Status: common.HttpStatus_RouteNotHandled,
            ResponseData: util.StringToResponseBody(msg),
            Length: int64(len(msg)),
        },
        request,
    )
}

func (irt *Router) RoundTrip(request *http.Request) (response *http.Response, err error) {
    log.Printf("Request Intercepted: %#v", request)
    log.Printf("Request URL: %#v", request.URL)
    log.Printf("Request URL: %#v", request.RequestURI)
    response, err = irt.ServiceRouter(request)
    log.Printf("Response Returned: %#v", response)
    log.Printf("Error Returned: %#v", err)
    return
}

// validate the interface
var _ http.RoundTripper = &Router{}
