package router_test

import (
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "testing"

    "github.com/TestInABox/gostackinabox/common"
    "github.com/TestInABox/gostackinabox/router"
    "github.com/TestInABox/gostackinabox/service"
    "github.com/TestInABox/gostackinabox/util"
)


func validateStatus(t *testing.T, expectedStatus int, response *http.Response) {
    strStatus := fmt.Sprintf("%d", expectedStatus)

    if response.StatusCode != expectedStatus {
        t.Errorf(
            "Response Status Code (%d) does not match submitted status code (%d)",
            response.StatusCode,
            expectedStatus,
        )
    }
    if response.Status != strStatus {
        t.Errorf(
            "Response Status (%s) does not match submitted status string (%s)",
            response.Status,
            strStatus,
        )
    }
}

func validateResponseBody(t *testing.T, expectedBody string, body io.ReadCloser) {
    // bodyData is array that is larger than the expected data read
    // so that it will actually receive data. Read() will not allocate
    // an appropriate buffer to receive the data so this must be done instead
    expectedLength := len(expectedBody)
    bodyData := make([]byte, 2*expectedLength)
    readLength, readErr := body.Read(bodyData)
    if readErr != nil {
        t.Errorf("Unexpectedly received a read error: %v", readErr)
    }
    if readLength != expectedLength {
        t.Errorf("Unexpected data amount read: %d != %d", readLength, expectedLength)
    }
    strBodyData := string(bodyData[:expectedLength])
    if strBodyData != expectedBody {
        t.Errorf("Body data doesn't match: %s != %s", strBodyData, expectedBody)
    }
}

func Test_Router_Router(t *testing.T) {
    t.Run(
        "New",
        func(t *testing.T) {
            irt := router.New()
            if irt == nil {
                t.Fatal("New generated a nil pointer")
            }
            if irt.RequestHandlers == nil {
                t.Errorf("New did not create the route table")
            }
            if len(irt.RequestHandlers) != 0 {
                t.Errorf("New did not generate an empty route table: %#v", irt.RequestHandlers)
            }
            expectedProtoMajor := 1
            expectedProtoMinor := 1
            expectedDisableCompression := true
            if irt.ProtoMajor != expectedProtoMajor ||
               irt.ProtoMinor != expectedProtoMinor ||
               irt.DisableCompression != expectedDisableCompression {
                t.Errorf(
                    "Unexpected parameters: ProtoMajor: %d != %d, ProtoMinor: %d != %d, DisableCompression: %t != %t",
                    irt.ProtoMajor,
                    expectedProtoMajor,
                    irt.ProtoMinor,
                    expectedProtoMinor,
                    irt.DisableCompression,
                    expectedDisableCompression,
                )
            }
        },
    )
    t.Run(
        "RegisterService",
        func(t *testing.T) {
            irt := router.New()
            if len(irt.RequestHandlers) != 0 {
                t.Errorf("Not starting with an empty route table: %#v", irt.RequestHandlers)
            }

            serviceName := "dummyService"
            serviceHandler := &service.ServiceHandler{}

            err := irt.RegisterService(serviceName, serviceHandler)
            if err != nil {
                t.Errorf("Unexpected error adding service name (%v) to the route table: %v - %#v", serviceName, err, irt.RequestHandlers)
            }

            // adding another time should generate an error
            err = irt.RegisterService(serviceName, serviceHandler)
            if !errors.Is(err, router.ErrServiceHandlerAlreadyRegister) {
                t.Errorf("Unexpected error adding a suplicate service name (%v) to the route table: %v - %#v", serviceName, err, irt.RequestHandlers)
            }
        },
    )
    t.Run(
        "BuildResponse",
        func(t *testing.T) {
            type TestScenarioExpectParameters struct {
                router   *router.Router
                response *http.Response
                err      error
            }

            type TestScenarioParameters struct {
                request *http.Request
                response *common.HttpReply
                beginFn func(t *testing.T)
                expectFn func(t *testing.T, tsep TestScenarioExpectParameters)
            }

            type TestScenario struct {
                name string
                setup func(t *testing.T) TestScenarioParameters
            }

            var TestScenarios = []TestScenario{
                {
                    name: "Nil response",
                    setup: func(t *testing.T) TestScenarioParameters {
                        return TestScenarioParameters{
                            request: nil,
                            response: nil,
                            beginFn: func(t *testing.T) {},
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if !errors.Is(tsep.err, router.ErrResponseBuildingInternalError) {
                                    t.Errorf("Unexpected error for build response: %v != %v", tsep.err, router.ErrResponseBuildingInternalError)
                                }
                                if tsep.response != nil {
                                    t.Errorf("Unexpected HttpReply received: %#v", tsep.response)
                                }
                            },
                        }
                    },
                },
                {
                    name: "Nil request",
                    setup: func(t *testing.T) TestScenarioParameters {
                        return TestScenarioParameters{
                            request: nil,
                            response: &common.HttpReply{},
                            beginFn: func(t *testing.T) {},
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if !errors.Is(tsep.err, router.ErrResponseBuildingInternalError) {
                                    t.Errorf("Unexpected error for build response: %v != %v", tsep.err, router.ErrResponseBuildingInternalError)
                                }
                                if tsep.response != nil {
                                    t.Errorf("Unexpected HttpReply received: %#v", tsep.response)
                                }
                            },
                        }
                    },
                },
                {
                    name: "valid",
                    setup: func(t *testing.T) TestScenarioParameters {
                        header := make(http.Header)
                        trailer := make(http.Header)
                        status := 200
                        msg := "the quick brown fox jumped over the hen house to escape the farmer"
                        expectedLength := int64(len(msg))
                        return TestScenarioParameters{
                            request: &http.Request{},
                            response: &common.HttpReply{
                                Status: common.HttpStatusCode(status),
                                Headers: header,
                                Trailers: trailer,
                                ResponseData: util.StringToResponseBody(msg),
                                Length: expectedLength,
                            },
                            beginFn: func(t *testing.T) {},
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if tsep.err != nil {
                                    t.Errorf("Unexpected error received: %v", tsep.err)
                                }
                                if tsep.response == nil {
                                    t.Errorf("Unexpectedly did not receive an HttpReply: %#v", tsep.response)
                                } else {
                                    strProto := fmt.Sprintf(
                                        "HTTP/%d.%d",
                                        tsep.router.ProtoMajor,
                                        tsep.router.ProtoMinor,
                                    )
                                    if tsep.response.Proto != strProto {
                                        t.Errorf(
                                            "Response Protocol (%s) does not match the expected protocol (%s)",
                                            tsep.response.Proto,
                                            strProto,
                                        )
                                    }
                                    if tsep.response.ProtoMajor != tsep.router.ProtoMajor {
                                        t.Errorf(
                                            "Response Protocol Major Version (%d) does not match the expected protocol major version (%d)",
                                            tsep.response.ProtoMajor,
                                            tsep.router.ProtoMajor,
                                        )
                                    }
                                    if tsep.response.ProtoMinor != tsep.router.ProtoMinor {
                                        t.Errorf(
                                            "Response Protocol Minor Version (%d) does not match the expected protocol minor version (%d)",
                                            tsep.response.ProtoMinor,
                                            tsep.router.ProtoMinor,
                                        )
                                    }
                                    headerMatch, headerMatchErr := util.EqualHeaders(
                                        &tsep.response.Header,
                                        &header,
                                    )
                                    if headerMatch == false {
                                        t.Errorf("%v", headerMatchErr)
                                    }
                                    trailerMatch, trailerMatchErr := util.EqualHeaders(
                                        &tsep.response.Trailer,
                                        &trailer,
                                    )
                                    if trailerMatch == false {
                                        t.Errorf("%v", trailerMatchErr)
                                    }
                                    if tsep.response.ContentLength != expectedLength {
                                        t.Errorf(
                                            "Expected Body data length doesn't match: %d != %d",
                                            tsep.response.ContentLength,
                                            expectedLength,
                                        )
                                    }
                                    if tsep.response.Uncompressed {
                                        t.Errorf("Unexpectedly found compression enabled")
                                    }
                                    if tsep.response.TLS != nil {
                                        t.Errorf("Unexpected found TLS support: %v", tsep.response.TLS)
                                    }

                                    validateStatus(t, status, tsep.response)
                                    validateResponseBody(t, msg, tsep.response.Body)
                                }
                            },
                        }
                    },
                },
            }

            for _, scenario := range TestScenarios {
                t.Run(
                    scenario.name,
                    func(t *testing.T) {
                        parameters := scenario.setup(t)

                        parameters.beginFn(t)

                        irt := router.New()
                        hr, err := irt.BuildResponse(parameters.response, parameters.request)
                        parameters.expectFn(
                            t,
                            TestScenarioExpectParameters{
                                router: irt,
                                response: hr,
                                err: err,
                            },
                        )
                    },
                )
            }
        },
    )
    t.Run(
        "RoundTrip",
        func(t *testing.T) {
            type TestScenarioExpectParameters struct {
                router   *router.Router
                response *http.Response
                err      error
            }

            type TestScenarioParameters struct {
                request *http.Request

                beginFn func(t *testing.T, irt *router.Router)
                expectFn func(t *testing.T, tsep TestScenarioExpectParameters)
            }

            type TestScenario struct {
                name string
                setup func(t *testing.T) TestScenarioParameters
            }

            var TestScenarios = []TestScenario{
                {
                    name: "invalid route",
                    setup: func(t *testing.T) TestScenarioParameters {
                        myUrl, _ := url.Parse("http://example.com/")
                        return TestScenarioParameters{
                            request: &http.Request{
                                URL: myUrl,
                            },
                            beginFn: func(t *testing.T, irt *router.Router) {
                            },
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if tsep.err != nil {
                                    t.Errorf("Unexpectedly received an error: %#v", tsep.err)
                                }
                                if tsep.response == nil {
                                    t.Errorf("Did not receive a response")
                                } else {
                                    msg := fmt.Sprintf(
                                        "gostackinabox: no service to handle URL '%s'",
                                        myUrl.String(),
                                    )
                                    validateStatus(t, int(common.HttpStatus_RouteNotHandled), tsep.response)
                                    validateResponseBody(t, msg, tsep.response.Body)
                                }
                            },
                        }
                    },
                },
                {
                    name: "invalid URL",
                    setup: func(t *testing.T) TestScenarioParameters {
                        myUrl, _ := url.Parse("http://example.com/")
                        expectedErr := errors.New("mock service error")
                        return TestScenarioParameters{
                            request: &http.Request{},
                            beginFn: func(t *testing.T, irt *router.Router) {
                                serviceName := util.GetUrlBaseResource(myUrl)
                                serviceHandler := &service.ServiceHandler{
                                    Matcher: &common.BasicServerURI{
                                        Protocol: "http",
                                        Host: "example.com",
                                    },
                                    FuncHandler: func(hc *common.HttpCall) (hr *common.HttpReply, err error) {
                                        err = expectedErr
                                        return
                                    },
                                }

                                err := irt.RegisterService(serviceName, serviceHandler)
                                if err != nil {
                                    t.Errorf("Unexpected error adding service name (%v) to the route table: %v - %#v", serviceName, err, irt.RequestHandlers)
                                }
                            },
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if !errors.Is(tsep.err, router.ErrInvalidRequest) {
                                    t.Errorf("Unexpectedly received an error: %#v != %v", tsep.err, router.ErrInvalidRequest)
                                }
                                if tsep.response != nil {
                                    t.Errorf("Unexpectedly received a response; %#v", tsep.response)
                                }
                            },
                        }
                    },
                },
                {
                    name: "service match error",
                    setup: func(t *testing.T) TestScenarioParameters {
                        myUrl, _ := url.Parse("http://example.com/")
                        expectedErr := errors.New("mock service error")
                        return TestScenarioParameters{
                            request: &http.Request{
                                URL: myUrl,
                            },
                            beginFn: func(t *testing.T, irt *router.Router) {
                                serviceName := util.GetUrlBaseResource(myUrl)
                                serviceHandler := &service.ServiceHandler{
                                    FuncHandler: func(hc *common.HttpCall) (hr *common.HttpReply, err error) {
                                        err = expectedErr
                                        return
                                    },
                                    Matcher: &common.BasicServerURI{
                                        Protocol: "http",
                                        Host: "",
                                    },
                                }

                                err := irt.RegisterService(serviceName, serviceHandler)
                                if err != nil {
                                    t.Errorf("Unexpected error adding service name (%v) to the route table: %v - %#v", serviceName, err, irt.RequestHandlers)
                                }
                            },
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if !errors.Is(tsep.err, common.ErrServerURIMisconfigured) {
                                    t.Errorf("Unexpectedly received an error: %#v != %#v", tsep.err, common.ErrServerURIMisconfigured)
                                }
                                if tsep.response != nil {
                                    t.Errorf("Unexpectedly received a response: %#v", tsep.response)
                                }
                            },
                        }
                    },
                },
                {
                    name: "service error",
                    setup: func(t *testing.T) TestScenarioParameters {
                        myUrl, _ := url.Parse("http://example.com/")
                        expectedErr := errors.New("mock service error")
                        return TestScenarioParameters{
                            request: &http.Request{
                                URL: myUrl,
                            },
                            beginFn: func(t *testing.T, irt *router.Router) {
                                serviceName := util.GetUrlBaseResource(myUrl)
                                serviceHandler := &service.ServiceHandler{
                                    Matcher: &common.BasicServerURI{
                                        Protocol: "http",
                                        Host: "example.com",
                                    },
                                    FuncHandler: func(hc *common.HttpCall) (hr *common.HttpReply, err error) {
                                        err = expectedErr
                                        return
                                    },
                                }

                                err := irt.RegisterService(serviceName, serviceHandler)
                                if err != nil {
                                    t.Errorf("Unexpected error adding service name (%v) to the route table: %v - %#v", serviceName, err, irt.RequestHandlers)
                                }
                            },
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if tsep.err != nil {
                                    t.Errorf("Unexpectedly received an error: %#v", tsep.err)
                                }
                                if tsep.response == nil {
                                    t.Errorf("Did not receive a response")
                                } else {
                                    msg := fmt.Sprintf(
                                        "gostackinabox: service handling request had an error - %#v",
                                        expectedErr,
                                    )
                                    validateStatus(t, int(common.HttpStatus_ServiceError), tsep.response)
                                    validateResponseBody(t, msg, tsep.response.Body)
                                }
                            },
                        }
                    },
                },
                {
                    name: "success",
                    setup: func(t *testing.T) TestScenarioParameters {
                        expectedStatus := 200
                        msg := "hello world"
                        myUrl, _ := url.Parse("http://example.com/")
                        return TestScenarioParameters{
                            request: &http.Request{
                                URL: myUrl,
                            },
                            beginFn: func(t *testing.T, irt *router.Router) {
                                serviceName := util.GetUrlBaseResource(myUrl)
                                serviceHandler := &service.ServiceHandler{
                                    Matcher: &common.BasicServerURI{
                                        Protocol: "http",
                                        Host: "example.com",
                                    },
                                    FuncHandler: func(hc *common.HttpCall) (hr *common.HttpReply, err error) {
                                        hr = &common.HttpReply{
                                            Status: common.HttpStatusCode(expectedStatus),
                                            ResponseData: util.StringToResponseBody(msg),
                                            Length: int64(len(msg)),
                                        }
                                        return
                                    },
                                }

                                err := irt.RegisterService(serviceName, serviceHandler)
                                if err != nil {
                                    t.Errorf("Unexpected error adding service name (%v) to the route table: %v - %#v", serviceName, err, irt.RequestHandlers)
                                }
                            },
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if tsep.err != nil {
                                    t.Errorf("Unexpectedly received an error: %#v", tsep.err)
                                }
                                if tsep.response == nil {
                                    t.Errorf("Did not receive a response")
                                } else {
                                    validateStatus(t, expectedStatus, tsep.response)
                                    validateResponseBody(t, msg, tsep.response.Body)
                                }
                            },
                        }
                    },
                },
            }

            for _, scenario := range TestScenarios {
                t.Run(
                    scenario.name,
                    func(t *testing.T) {
                        parameters := scenario.setup(t)

                        irt := router.New()

                        parameters.beginFn(t, irt)

                        req, err := irt.RoundTrip(parameters.request)
                        parameters.expectFn(
                            t,
                            TestScenarioExpectParameters{
                                router: irt,
                                response: req,
                                err: err,
                            },
                        )
                    },
                )
            }
        },
    )
}
