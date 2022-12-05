package service_test

import (
    "errors"
    "fmt"
    "net/url"
    //"regexp"
    "testing"

    "github.com/TestInABox/gostackinabox/common"
    "github.com/TestInABox/gostackinabox/service"
)

func TestService(t *testing.T) {
    t.Run(
        "init",
        func(t *testing.T) {
            handler := service.ServiceHandler{}
            if handler.SubServices != nil {
                t.Errorf("Unexpected initialized subservices: %v", handler.SubServices)
            }
            result := handler.Init("init", &common.BasicServerURI{})
            if handler.SubServices == nil {
                t.Errorf("Unexpectedly did not initialize the subservices: %v", handler.SubServices)
            }
            if result != nil {
                t.Errorf("Unexpectedly received an error: %#v", result)
            }
        },
    )
    t.Run(
        "GetMatcher",
        func(t *testing.T) {
            type TestScenarioExpectParameters struct {
                matcher common.URI
            }

            type TestScenarioParameters struct {
                beginFn func(t *testing.T, handler *service.ServiceHandler)
                expectFn func(t *testing.T, tsep TestScenarioExpectParameters)
            }

            type TestScenario struct {
                name string
                setup func(t *testing.T) TestScenarioParameters
            }

            var TestScenarios = []TestScenario{
                {
                    name: "invalid matcher",
                    setup: func(t *testing.T) TestScenarioParameters {
                        return TestScenarioParameters{
                            beginFn: func(t *testing.T, handler *service.ServiceHandler) {
                                handler.Matcher = nil
                            },
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if tsep.matcher != nil {
                                    t.Errorf("Unexpectly received a matcher: %#v", tsep.matcher)
                                }
                            },
                        }
                    },
                },
                {
                    name: "success",
                    setup: func(t *testing.T) TestScenarioParameters {
                        theMatcher := common.BasicServerURI{
                            Protocol: "https",
                            Host: "foo.bar",
                            Port: "443",
                        }
                        return TestScenarioParameters{
                            beginFn: func(t *testing.T, handler *service.ServiceHandler) {
                                handler.Matcher = &theMatcher
                            },
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if tsep.matcher.(*common.BasicServerURI).Protocol != theMatcher.Protocol ||
                                   tsep.matcher.(*common.BasicServerURI).Host != theMatcher.Host ||
                                   tsep.matcher.(*common.BasicServerURI).Port != theMatcher.Port {
                                    t.Errorf("Unexpectly received a matcher: %#v", tsep.matcher)
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
                        handler := service.ServiceHandler{}
                        handler.Init(
                            scenario.name,
                            &common.BasicServerURI{},
                        )

                        parameters := scenario.setup(t)

                        parameters.beginFn(t, &handler)

                        results := TestScenarioExpectParameters{
                            matcher: handler.GetMatcher(),
                        }
                        parameters.expectFn(t, results)
                    },
                )
            }
        },
    )
    t.Run(
        "RegisterHandler",
        func(t *testing.T) {
            type TestScenarioExpectParameters struct {
                err error
            }

            type TestScenarioParameters struct {
                svc service.Service
                beginFn func(t *testing.T, handler *service.ServiceHandler)
                expectFn func(t *testing.T, tsep TestScenarioExpectParameters)
            }

            type TestScenario struct {
                name string
                setup func(t *testing.T) TestScenarioParameters
            }

            var TestScenarios = []TestScenario{
                {
                    name: "Invalid Service Instance",
                    setup: func(t *testing.T) TestScenarioParameters {
                        /*
                        noImplementedSvc := &service.ServiceHandler{
                            Name: "Not Implemented",
                            Matcher: &common.BasicServerURI{
                                Protocol: "http",
                                Host: "example.com",
                            },
                        }
                        */
                        return TestScenarioParameters{
                            //svc service.Service
                            beginFn: func(t *testing.T, handler *service.ServiceHandler) {},
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if !errors.Is(tsep.err, service.ErrInvalidService) {
                                    t.Errorf("Unexpected error received: %#v != %#v", tsep.err, service.ErrInvalidService)
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
                        handler := service.ServiceHandler{}
                        handler.Init(
                            scenario.name,
                            &common.BasicServerURI{
                                Protocol: "http",
                                Host: "example.com",
                            },
                        )

                        parameters := scenario.setup(t)

                        parameters.beginFn(t, &handler)

                        results := TestScenarioExpectParameters{
                            err: handler.RegisterHandler(parameters.svc),
                        }
                        parameters.expectFn(t, results)
                    },
                )
            }
        },
    )
    t.Run(
        "ValidateRegex",
        func(t *testing.T) {
            type TestScenarioExpectParameters struct {
                err error
            }

            type TestScenarioParameters struct {
                regexValue string
                isSubService bool
                beginFn func(t *testing.T, handler *service.ServiceHandler)
                expectFn func(t *testing.T, tsep TestScenarioExpectParameters)
            }

            type TestScenario struct {
                name string
                setup func(t *testing.T) TestScenarioParameters
            }

            var TestScenarios = []TestScenario{
                {
                    name: "invalid first char",
                    setup: func(t *testing.T) TestScenarioParameters {
                        return TestScenarioParameters{
                            regexValue: "foo",
                            isSubService: false,
                            beginFn: func(t *testing.T, handler *service.ServiceHandler) {},
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if !errors.Is(tsep.err, service.ErrInvalidServiceRegex) {
                                    t.Errorf("Unexpected error: %#v != %#v", tsep.err, service.ErrInvalidServiceRegex)
                                }
                            },
                        }
                    },
                },
                {
                    name: "invalid last char",
                    setup: func(t *testing.T) TestScenarioParameters {
                        return TestScenarioParameters{
                            regexValue: "^foo",
                            isSubService: false,
                            beginFn: func(t *testing.T, handler *service.ServiceHandler) {},
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if !errors.Is(tsep.err, service.ErrInvalidServiceRegex) {
                                    t.Errorf("Unexpected error: %#v != %#v", tsep.err, service.ErrInvalidServiceRegex)
                                }
                            },
                        }
                    },
                },
                {
                    name: "invalid last char for subservice",
                    setup: func(t *testing.T) TestScenarioParameters {
                        return TestScenarioParameters{
                            regexValue: "^foo$",
                            isSubService: true,
                            beginFn: func(t *testing.T, handler *service.ServiceHandler) {},
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if !errors.Is(tsep.err, service.ErrInvalidServiceRegex) {
                                    t.Errorf("Unexpected error: %#v != %#v", tsep.err, service.ErrInvalidServiceRegex)
                                }
                            },
                        }
                    },
                },
                {
                    name: "success for no-subservice",
                    setup: func(t *testing.T) TestScenarioParameters {
                        return TestScenarioParameters{
                            regexValue: "^foo$",
                            isSubService: false,
                            beginFn: func(t *testing.T, handler *service.ServiceHandler) {},
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if tsep.err != nil {
                                    t.Errorf("Unexpected error: %#v", tsep.err)
                                }
                            },
                        }
                    },
                },
                {
                    name: "success for subservice",
                    setup: func(t *testing.T) TestScenarioParameters {
                        return TestScenarioParameters{
                            regexValue: "^foo",
                            isSubService: true,
                            beginFn: func(t *testing.T, handler *service.ServiceHandler) {},
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if tsep.err != nil {
                                    t.Errorf("Unexpected error: %#v", tsep.err)
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
                        handler := service.ServiceHandler{}

                        parameters := scenario.setup(t)

                        parameters.beginFn(t, &handler)

                        results := TestScenarioExpectParameters{
                            err: handler.ValidateRegex(parameters.regexValue, parameters.isSubService),
                        }
                        parameters.expectFn(t, results)
                    },
                )
            }
        },
    )
    t.Run(
        "GetHandler",
        func(t *testing.T) {
            type TestScenarioExpectParameters struct {
                err error
                handler common.HttpHandler
            }

            type TestScenarioParameters struct {
                inputUrl url.URL
                beginFn func(t *testing.T, handler *service.ServiceHandler)
                expectFn func(t *testing.T, tsep TestScenarioExpectParameters)
            }

            type TestScenario struct {
                name string
                setup func(t *testing.T) TestScenarioParameters
            }

            var TestScenarios = []TestScenario{
                {
                    name: "nil handler",
                    setup: func(t *testing.T) TestScenarioParameters {
                        theUrl, _ := url.Parse("https://foo.bar")
                        theHandler := func(*common.HttpCall) (*common.HttpReply, error) {
                            return nil, fmt.Errorf("Not implemented")
                        }
                        return TestScenarioParameters{
                            inputUrl: *theUrl,
                            beginFn: func(t *testing.T, handler *service.ServiceHandler) {
                                handler.FuncHandler = theHandler
                            },
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if tsep.err != nil {
                                    t.Errorf("Unexpected error: %#v", tsep.err)
                                }
                                if tsep.handler == nil {
                                    t.Errorf("Did not receive a handler")
                                }
                                // would be great if it could be proven the handler above is the
                                // same one that got returned; unfortunately, that's not very easy
                                // to do in Golang.
                                /*
                                f1 := &tsep.handler
                                f2 := (*common.HttpHandler)(&theHandler)
                                if !(f1 == f2) {
                                    t.Errorf("Unexpected handler received: %#v != %#v", &tsep.handler, &theHandler)
                                }
                                */
                            },
                        }
                    },
                },
                /*
                {
                    name: "success",
                    setup: func(t *testing.T) TestScenarioParameters {
                        theUrl, _ := url.Parse("https://foo.bar")
                        theHandler := func(*common.HttpCall) (*common.HttpReply, error) {
                            return nil, fmt.Errorf("Not implemented")
                        }
                        return TestScenarioParameters{
                            inputUrl: *theUrl,
                            beginFn: func(t *testing.T, handler *service.ServiceHandler) {
                                handler.FuncHandler = theHandler
                            },
                            expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                                if tsep.handler == nil {
                                    t.Errorf("Unexpected handler returned: %#p != %#p", tsep.handler, theHandler)
                                }
                                if tsep.err != nil {
                                    t.Errorf("Unexpected error: %#v", tsep.err)
                                }
                            },
                        }
                    },
                },
                */
            }

            for _, scenario := range TestScenarios {
                t.Run(
                    scenario.name,
                    func(t *testing.T) {
                        handler := service.ServiceHandler{}

                        parameters := scenario.setup(t)

                        parameters.beginFn(t, &handler)

                        subHandler, err := handler.GetHandler(parameters.inputUrl)
                        results := TestScenarioExpectParameters{
                            handler: subHandler,
                            err: err,
                        }
                        parameters.expectFn(t, results)
                    },
                )
            }
        },
    )
}

