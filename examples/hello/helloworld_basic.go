package hello

import (
    "github.com/TestInABox/gostackinabox/common"
    "github.com/TestInABox/gostackinabox/common/log"
    "github.com/TestInABox/gostackinabox/service"
    "github.com/TestInABox/gostackinabox/util"
)


/*
 * HelloWorldBasicService is the most basic handler implementation where
 * everything is served via the singular method assigned to FuncHandler
 * on the object.
 */
type HelloWorldBasicService struct {
    service.ServiceHandler
}

func NewHelloWorldBasicService() (s service.Service, err error) {
     hwb := &HelloWorldBasicService{}
     err = hwb.Init(
        "helloWorldBasic",
        &common.BasicServerURI{
            Protocol: "https",
            Host: "hello.world",
            Port: "",
        },
    )
     if err != nil {
        hwb.FuncHandler = hwb.HelloWorldHandler
     }
     s = hwb
     return
}

func (hwb *HelloWorldBasicService) HelloWorldHandler(request *common.HttpCall) (result *common.HttpReply, err error) {
    log.Printf("Hello World Basic Service FuncHandler")
    msg := "hello world!"
    expectedLength := int64(len(msg))
    result = &common.HttpReply{
        Status: common.GetHttpStatus(200),
        ResponseData: util.StringToResponseBody(msg),
        Length: expectedLength,
    }
    return
}

var _ service.Service = &HelloWorldBasicService{}
