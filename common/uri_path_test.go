package common_test

import (
    "errors"
    "net/url"
    "regexp"
    "testing"

    "github.com/TestInABox/gostackinabox/common"
)

func Test_Common_PathURI(t *testing.T) {
    type TestScenarioExpectParameters struct {
        result bool
        err    error
    }

    type TestScenarioParameters struct {
        regexpValue string
        checkURL url.URL

        beginFn func(t *testing.T, pu *common.PathURI)
        expectFn func(t *testing.T, tsep TestScenarioExpectParameters)
    }

    type TestScenario struct {
        name string
        setup func(t *testing.T) TestScenarioParameters
    }

    var TestScenarios = []TestScenario{
        {
            name: "invalid regex",
            setup: func(t *testing.T) TestScenarioParameters {
                return TestScenarioParameters{
                    checkURL:  url.URL{},
                    beginFn: func(t *testing.T, pu *common.PathURI) {
                        pu.Path = nil
                    },
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.result {
                            t.Errorf("Unexpectedly got result value of %t", tsep.result)
                        }
                        if !errors.Is(tsep.err, common.ErrPathURIMisconfigured) {
                            t.Errorf(
                                "Unexpected error result: %#v != %#v",
                                tsep.err,
                                common.ErrPathURIMisconfigured,
                            )
                        }
                    },
                }
            },
        },
        {
            name: "valid regex invalid value",
            setup: func(t *testing.T) TestScenarioParameters {
                return TestScenarioParameters{
                    checkURL:  url.URL{},
                    beginFn: func(t *testing.T, pu *common.PathURI) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.result {
                            t.Errorf("Unexpectedly got result value of %t", tsep.result)
                        }
                        if !errors.Is(tsep.err, common.ErrPathURIMisconfigured) {
                            t.Errorf(
                                "Unexpected error result: %#v != %#v",
                                tsep.err,
                                common.ErrPathURIMisconfigured,
                            )
                        }
                    },
                }
            },
        },
        {
            name: "valid",
            setup: func(t *testing.T) TestScenarioParameters {
                uv, _ := url.Parse("http://example.com/hello")
                return TestScenarioParameters{
                    checkURL:  *uv,
                    regexpValue: "^/hello",
                    beginFn: func(t *testing.T, pu *common.PathURI) {},
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

                regex := regexp.MustCompile(parameters.regexpValue)
                pathURI := &common.PathURI{
                    Path: regex,
                }

                parameters.beginFn(t, pathURI)

                if pathURI.Path != nil {
                    regexValue := pathURI.Path.String()
                    t.Logf("Regex(length: %d): \"%s\"", len(regexValue), regexValue)
                } else {
                    t.Logf("Empty regex")
                }

                result, err := pathURI.IsMatch(parameters.checkURL)
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
