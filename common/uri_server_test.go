package common_test

import (
    "errors"
    "net/url"
    "testing"

    "github.com/TestInABox/gostackinabox/common"
)

func Test_Common_BasicServerURI(t *testing.T) {
    type TestScenarioExpectParameters struct {
        result bool
        err    error
    }

    type TestScenarioParameters struct {
        checkURL url.URL

        beginFn func(t *testing.T, su *common.BasicServerURI)
        expectFn func(t *testing.T, tsep TestScenarioExpectParameters)
    }

    type TestScenario struct {
        name string
        setup func(t *testing.T) TestScenarioParameters
    }

    var TestScenarios = []TestScenario{
        {
            name: "invalid hostname",
            setup: func(t *testing.T) TestScenarioParameters {
                return TestScenarioParameters{
                    checkURL: url.URL{},
                    beginFn: func(t *testing.T, su *common.BasicServerURI) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.result {
                            t.Errorf("Unexpectedly got result value of %t", tsep.result)
                        }
                        if !errors.Is(tsep.err, common.ErrServerURIMisconfigured) {
                            t.Errorf(
                                "Unexpected error result: %#v != %#v",
                                tsep.err,
                                common.ErrServerURIMisconfigured,
                            )
                        }
                    },
                }
            },
        },
        {
            name: "hostname mismatch",
            setup: func(t *testing.T) TestScenarioParameters {
                return TestScenarioParameters{
                    checkURL: url.URL{
                        Host: "example.com",
                    },
                    beginFn: func(t *testing.T, su *common.BasicServerURI) {
                        su.Host = "com.example"
                    },
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.result {
                            t.Errorf("Unexpectedly got result value of %t", tsep.result)
                        }
                        if tsep.err != nil {
                            t.Errorf(
                                "Unexpected error result: %#v",
                                tsep.err,
                            )
                        }
                    },
                }
            },
        },
        {
            name: "protocol mismatch",
            setup: func(t *testing.T) TestScenarioParameters {
                return TestScenarioParameters{
                    checkURL: url.URL{
                        Scheme: "https",
                        Host: "example.com",
                    },
                    beginFn: func(t *testing.T, su *common.BasicServerURI) {
                        su.Host = "example.com"
                        su.Protocol = "http"
                    },
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.result {
                            t.Errorf("Unexpectedly got result value of %t", tsep.result)
                        }
                        if tsep.err != nil {
                            t.Errorf(
                                "Unexpected error result: %#v",
                                tsep.err,
                            )
                        }
                    },
                }
            },
        },
        {
            name: "protocol mismatch with colon",
            setup: func(t *testing.T) TestScenarioParameters {
                return TestScenarioParameters{
                    checkURL: url.URL{
                        Scheme: "https:",
                        Host: "example.com",
                    },
                    beginFn: func(t *testing.T, su *common.BasicServerURI) {
                        su.Host = "example.com"
                        su.Protocol = "http"
                    },
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.result {
                            t.Errorf("Unexpectedly got result value of %t", tsep.result)
                        }
                        if tsep.err != nil {
                            t.Errorf(
                                "Unexpected error result: %#v",
                                tsep.err,
                            )
                        }
                    },
                }
            },
        },
        {
            name: "port mismatch",
            setup: func(t *testing.T) TestScenarioParameters {
                return TestScenarioParameters{
                    checkURL: url.URL{
                        Scheme: "http",
                        Host: "example.com:443",
                    },
                    beginFn: func(t *testing.T, su *common.BasicServerURI) {
                        su.Host = "example.com"
                        su.Protocol = "http"
                        su.Port = "80"
                    },
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.result {
                            t.Errorf("Unexpectedly got result value of %t", tsep.result)
                        }
                        if tsep.err != nil {
                            t.Errorf(
                                "Unexpected error result: %#v",
                                tsep.err,
                            )
                        }
                    },
                }
            },
        },
        {
            name: "success",
            setup: func(t *testing.T) TestScenarioParameters {
                return TestScenarioParameters{
                    checkURL: url.URL{
                        Scheme: "http",
                        Host: "example.com:443",
                    },
                    beginFn: func(t *testing.T, su *common.BasicServerURI) {
                        su.Host = "example.com"
                        su.Protocol = "http"
                        su.Port = "443"
                    },
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if !tsep.result {
                            t.Errorf("Unexpectedly got result value of %t", tsep.result)
                        }
                        if tsep.err != nil {
                            t.Errorf(
                                "Unexpected error result: %#v",
                                tsep.err,
                            )
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

                serverURI := &common.BasicServerURI{}

                parameters.beginFn(t, serverURI)

                result, err := serverURI.IsMatch(parameters.checkURL)
                parameters.expectFn(
                    t,
                    TestScenarioExpectParameters{
                        result: result,
                        err: err,
                    },
                )

            },
        )
    }
}
