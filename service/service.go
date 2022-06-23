package service

import (
    "fmt"
    "net/url"

    "github.com/TestInABox/gostackinabox/common"
    "github.com/TestInABox/gostackinabox/common/log"
    "github.com/TestInABox/gostackinabox/util"
)

/*
    In the Python Stack-In-A-Box implementation a sub-service
    was just another service; however, Python also had 3rd Party
    tools for doing the initial intercept at the URL Domain Level.

    The service handler here acts like the 3rd Party tools in Python
    to capture the domain level (e.g the `https://example.com` portion
    of `https://example.com/some/uri/path/object'.

    Some basic services may wish to implement their entire stack using
    just the ServiceHandler as a base. More advanced services will
    want to combine it with a series of registered methods and Service
    instances to offload and simplify the handling of complex URI paths.
*/

// a service must match at the URL Domain Level
type ServiceHandler struct {
    isSubService bool
    Name string
    Matcher common.URI
    // handle the local verbs (e.g GET/POST/OPTION/etc on /)
    FuncHandler common.HttpHandler
    MethodMap   common.HttpHandlerMap
    // handle sub-routes (e.g  GET/POST/OPTION/etc on /<object>)
    //SubServices ServiceMethodHandlerMap
    SubServices ServiceHandlerMap
}

func (sh *ServiceHandler) Init(name string, matcher common.URI) (err error) {
    sh.Name = name
    sh.Matcher = matcher
    sh.MethodMap = make(common.HttpHandlerMap)
    sh.SubServices = make(ServiceHandlerMap)
    sh.FuncHandler =  sh.DefaultFuncHandler

    switch sh.Matcher.(type) {
    case common.ServerURI:
        // only recognize the ServerURI matcher define the root state
        sh.isSubService = false
    default:
        // all other matchers be subservices
        sh.isSubService = true
    }

    return
}

func (sh *ServiceHandler) IsSubService() bool {
    return sh.isSubService
}

func (sh *ServiceHandler) GetName() string {
    return sh.Name
}

func (sh *ServiceHandler) GetMatcher() common.URI {
    return sh.Matcher
}

func (sh *ServiceHandler) MethodHandler(request *common.HttpCall) (result *common.HttpReply, err error) {
    for httpVerb, httpVerbHandler := range sh.MethodMap {
        if request.Method == common.HttpVerb(httpVerb) {
            result, err = httpVerbHandler(request)
            return
        }
    }

    msg := fmt.Sprintf("%s on %s is unhandled", request.Method, request.Url.String())
    expectedLength := int64(len(msg))
    result = &common.HttpReply{
        Status: common.GetHttpStatus(405),
        ResponseData: util.StringToResponseBody(msg),
        Length: expectedLength,
    }
    return
}

func (sh *ServiceHandler) DefaultFuncHandler(request *common.HttpCall) (result *common.HttpReply, err error) {
    msg := "Unhandled"
    expectedLength := int64(len(msg))
    result = &common.HttpReply{
        Status: common.GetHttpStatus(500),
        ResponseData: util.StringToResponseBody(msg),
        Length: expectedLength,
    }
    return
}

func (sh *ServiceHandler) GetHandler(requestUrl url.URL) (result common.HttpHandler, err error) {
    log.Printf("Checking if any handlers respond to %s", requestUrl.String())
    // first is there any sub service that handles the route
    for serviceName, serviceHandler := range sh.SubServices {
        // see if this service handles the URL
        matcher := serviceHandler.GetMatcher()
        matchResult, matchErr := matcher.IsMatch(requestUrl)
        if matchErr != nil {
            // there's a problem with the matcher, test needs to be fixed
            log.Printf("Retreiving matcher for service %s had error: %#v", serviceName, matchErr)
            err = fmt.Errorf("Service %s generated an error: %w", serviceName, matchErr)
            return
        }
        log.Printf("%s ? %#v : %t", serviceName, requestUrl, matchResult)
        if matchResult {
            // get the handler for the service
            handler, handlerErr := serviceHandler.GetHandler(requestUrl)
            if handlerErr != nil {
                log.Printf("Retreiving handler for service %s had error: %#v", serviceName, handlerErr)
                err = handlerErr
                return
            }

            log.Printf("Service %s supports URI %s using handler %v", serviceName, requestUrl.String(), handler)
            result = handler
            return
        }
    }

    log.Printf("No Subservices handling the URL %s", requestUrl.String())

    log.Printf("Checking for method handlers (count: %d)", len(sh.MethodMap))
    if len(sh.MethodMap) > 0 {
        log.Printf("Using the Method Handler to handle the URL %s", requestUrl.String())
        return sh.MethodHandler, nil
    }

    log.Printf("Using primary handler to handle the URL %s", requestUrl.String())
    result = sh.FuncHandler
    return
}

func (sh *ServiceHandler) ValidateRegex(r string, isSubService bool) (err error) {
    // Regex:
    //  1. starts with ^
    //  2. if no subservices, then it must end is $
    //  3. if there are subservices, then it may not end with a $
    firstChar := r[0:1]
    lastChar := r[len(r)-1:]
    log.Printf("r: %s, first char: %s, last char: %s, subservice: %t", r, firstChar, lastChar, isSubService)

    if firstChar != "^" {
        log.Printf("RegEx Rule Violation: First Char (%s) is not ^", firstChar)
        err = fmt.Errorf("%w: Regex must start with ^", ErrInvalidServiceRegex)
        return
    }

    if lastChar != "$" &&  !isSubService {
        log.Printf("RegEx Rule Violation (isSubService: %t): Last Char (%s) is not $ a", isSubService, firstChar)
        err = fmt.Errorf("%w: Regex pattern must end with $", ErrInvalidServiceRegex)
        return
    }

    if lastChar == "$" && isSubService {
        log.Printf("RegEx Rule Violation (isSubService: %t): Last Char (%s) is $ a", isSubService, firstChar)
        err = fmt.Errorf("%w: Subservice Regex must not end with $", ErrInvalidServiceRegex)
        return
    }

    return
}

func (sh *ServiceHandler) RegisterHandler(subHandler Service) (err error) {
    if subHandler == nil {
        err = fmt.Errorf("%w: Missing Subservice instance", ErrInvalidService)
        return
    }

    svcName := subHandler.GetName()

    if len(svcName) == 0 {
        err = fmt.Errorf("%w: Subservice must have a name", ErrInvalidService)
        return
    }

    if !subHandler.IsSubService() {
        err = fmt.Errorf("%w: Can only registere subservics", ErrInvalidService)
        return
    }

    matcher := subHandler.GetMatcher()

    switch m := matcher.(type) {
    case common.ServerURI:
        err = fmt.Errorf("%w: sub services cannot use `common.ServerURI` for their matcher", ErrInvalidService)
        return
    case *common.PathURI:
        regExErr :=  sh.ValidateRegex(m.Path.String(), true)
        if regExErr != nil {
            err = fmt.Errorf("%w: Invalid regex for subservice", regExErr)
            return
        }
    default:
        // unable to validate regex
    }

    if _, ok := sh.SubServices[svcName]; ok {
        // service is already registered
        err = fmt.Errorf("%w: Attempting to register %s multiple times.", ErrServiceHandlerAlreadyRegister, svcName)
        return
    }

    sh.SubServices[svcName] = subHandler
    return
}

func (sh *ServiceHandler) RegisterMethodHandler(method common.HttpVerb, handler common.HttpHandler) (err error) {
    if handler == nil {
        err = fmt.Errorf("%w: Missing handler method for %s", ErrRequestHandlerInvalid, method)
        return
    }

    if _, ok := sh.MethodMap[method]; ok {
        // method is already registered
        err = fmt.Errorf("%w: Attempting to register %s multiple times.", ErrRequestHandlerAlreadyRegister, method)
        return
    }
    log.Printf("Registered method %s using handler %#v", method, handler)

    sh.MethodMap[method] = handler
    return
}

var _ Service = &ServiceHandler{}
