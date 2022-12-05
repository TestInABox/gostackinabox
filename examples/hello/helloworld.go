package hello

import (
    "net/http"

    "github.com/TestInABox/gostackinabox/common"
    "github.com/TestInABox/gostackinabox/common/log"
    "github.com/TestInABox/gostackinabox/service"
    "github.com/TestInABox/gostackinabox/util"
)


/*
 * HelloWorldService is a little more complex than HelloWorldBasic
 * where it shows how to allow the single level to handle all the
 * different HTTP verbs through mapping specific methods to any
 * given HTTP verb.
 *
 * Do note that the system is not limited to just the HTTP verbs
 * that are defined by the IETF HTTP related RFCs.
 */
type HelloWorldService struct {
    service.ServiceHandler
}

func NewHelloWorldService() (s service.Service, err error) {
     log.Printf("Creating a Hello World Service")
     hw := &HelloWorldService{}
     err = hw.Init(
        "helloWorld",
        &common.BasicServerURI{
            Protocol: "https",
            Host: "hello.world",
            Port: "",
        },
    )
     if err == nil {
        err = hw.RegisterMethodHandler(common.HttpVerb_Get, common.HttpHandler(hw.MyGetHandler))
        if err == nil {
            err = hw.RegisterMethodHandler(common.HttpVerb_Head, common.HttpHandler(hw.MyHeadHandler))
            if err != nil {
                log.Printf("Failed to register HEAD handler: %#v", err)
            }
        }else {
            log.Printf("Failed to registered GET handler: %#v", err)
        }
     } else {
        log.Printf("Failed to initialize service: %#v", err)
     }
     s = hw
     return
}

func (hw *HelloWorldService) MyGetHandler(request *common.HttpCall) (result *common.HttpReply, err error) {
    log.Printf("Hello World Service GET Handler")
    msg := "hello world!"
    expectedLength := int64(len(msg))
    result = &common.HttpReply{
        Status: common.GetHttpStatus(200),
        ResponseData: util.StringToResponseBody(msg),
        Length: expectedLength,
    }
    return
}

func (hw *HelloWorldService) MyHeadHandler(request *common.HttpCall) (result *common.HttpReply, err error) {
    log.Printf("Hello World Service HEAD Handler")
    msg := "hello world!"
    result = &common.HttpReply{
        Status: common.GetHttpStatus(200),
        Headers: make(http.Header),
    }
    result.Headers["x-msg"] = []string{msg}
    return
}

var _ service.Service = &HelloWorldService{}
