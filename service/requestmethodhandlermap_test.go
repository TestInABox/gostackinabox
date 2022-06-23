package service_test

import (
    "errors"
    "testing"

    "github.com/TestInABox/gostackinabox/service"
)

func TestRequestMethodHandlerMap(t *testing.T) {
    type TestScenarioExpectParameters struct {
        err error
    }

    type TestScenarioParameters struct {
        value string
        svc service.Service
        beginFn func(t *testing.T, mapper *service.RequestMethodHandlerMap)
        expectFn func(t *testing.T, tsep TestScenarioExpectParameters)
    }

    type TestScenario struct {
        name string
        setup func(t *testing.T) TestScenarioParameters
    }

    var TestScenarios = []TestScenario{
        {
            name: "Invalid service",
            setup: func(t *testing.T) TestScenarioParameters {
                return TestScenarioParameters{
                    value: "no-matter-the-method",
                    svc: nil,
                    beginFn: func(t *testing.T, mapper *service.RequestMethodHandlerMap) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if !errors.Is(tsep.err, service.ErrInvalidService) {
                            t.Errorf("Unexpected error received; %v != %v", tsep.err, service.ErrInvalidService)
                        }
                    },
                }
            },
        },
        {
            name: "Repeated service addition",
            setup: func(t *testing.T) TestScenarioParameters {
                methodValue := "bar"
                return TestScenarioParameters{
                    value: methodValue,
                    svc: &service.ServiceHandler{},
                    beginFn: func(t *testing.T, mapper *service.RequestMethodHandlerMap) {
                        (*mapper)[methodValue] = &service.ServiceHandler{}
                    },
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if !errors.Is(tsep.err, service.ErrRequestHandlerAlreadyRegister) {
                            t.Errorf("Unexpected error received; %v != %v", tsep.err, service.ErrRequestHandlerAlreadyRegister)
                        }
                    },
                }
            },
        },
        {
            name: "Success",
            setup: func(t *testing.T) TestScenarioParameters {
                methodValue := "bar"
                return TestScenarioParameters{
                    value: methodValue,
                    svc: &service.ServiceHandler{},
                    beginFn: func(t *testing.T, mapper *service.RequestMethodHandlerMap) {
                    },
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.err != nil {
                            t.Errorf("Unexpected error received; %v", tsep.err)
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
                mapper := service.RequestMethodHandlerMap{}

                parameters := scenario.setup(t)

                parameters.beginFn(t, &mapper)

                results := TestScenarioExpectParameters{
                    err: mapper.AddHandler(parameters.value, parameters.svc),
                }
                parameters.expectFn(t, results)
            },
        )
    }
}
