package util_test

import (
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "testing"

    "github.com/TestInABox/gostackinabox/util"
)


func TestUtilStringToResponseBody(t *testing.T) {
    type TestScenarioExpectParameters struct {
        readCloser io.ReadCloser
    }

    type TestScenarioParameters struct {
        value string
        beginFn func(t *testing.T)
        expectFn func(t *testing.T, tsep TestScenarioExpectParameters)
    }

    type TestScenario struct {
        name string
        setup func(t *testing.T) TestScenarioParameters
    }

    var TestScenarios = []TestScenario{
        {
            name: "null string",
            setup: func(t *testing.T) TestScenarioParameters {
                var emptyString string
                var expectedReadCount int // = 0
                readBufferSize := 100 // should be larger than the expectedReadCount
                return TestScenarioParameters{
                    value: emptyString,
                    beginFn: func(t *testing.T) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.readCloser == nil {
                            t.Error("Failed to get a read closer")
                        } else {
                            dataBuf := make([]byte, readBufferSize)
                            count, readErr := tsep.readCloser.Read(dataBuf)
                            if count != expectedReadCount {
                                t.Errorf("Read %d bytes when %d were expected", count, expectedReadCount)
                            }
                            if !errors.Is(readErr, io.EOF) {
                                t.Errorf("Received unexpected error: %v instead of %v", readErr, io.EOF)
                            }
                            resultValue := string(dataBuf[:count])
                            compareResult := strings.Compare(resultValue, emptyString)

                            t.Logf("Transfer Result: %d - %v", count, readErr)
                            t.Logf("Original(%d): %s", len(emptyString), emptyString)
                            t.Logf("Transfer(%d): %s", len(resultValue), resultValue)
                            if compareResult != 0 {
                                t.Errorf("Read \"%s\" when expecting \"%s\" - %d", resultValue, emptyString, compareResult)
                            }
                        }
                    },
                }
            },
        },
        {
            name: "valid string",
            setup: func(t *testing.T) TestScenarioParameters {
                var stringValue string = "alpha beta delta gamma epsilon"
                var expectedReadCount int = len(stringValue)
                readBufferSize := 100 // should be larger than the expectedReadCount
                return TestScenarioParameters{
                    value: stringValue,
                    beginFn: func(t *testing.T) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.readCloser == nil {
                            t.Error("Failed to get a read closer")
                        } else {
                            dataBuf := make([]byte, readBufferSize)
                            count, readErr := tsep.readCloser.Read(dataBuf)
                            if count != expectedReadCount {
                                t.Errorf("Read %d bytes when %d were expected", count, expectedReadCount)
                            }
                            if readErr != nil {
                                t.Errorf("Received unexpected error: %v", readErr)
                            }
                            resultValue := string(dataBuf[:count])
                            //if resultValue != stringValue {
                            compareResult := strings.Compare(resultValue, stringValue)

                            t.Logf("Transfer Result: %d - %v", count, readErr)
                            t.Logf("Original(%d): %s", len(stringValue), stringValue)
                            t.Logf("Transfer(%d): %s", len(resultValue), resultValue)
                            if compareResult != 0 {
                                t.Errorf("Read \"%s\" when expecting \"%s\"", resultValue, stringValue)
                            }
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
                result := util.StringToResponseBody(parameters.value)
                parameters.expectFn(
                    t,
                    TestScenarioExpectParameters{
                        readCloser: result,
                    },
                )
            },
        )
    }
}

func TestUtilGetUrlBaseResource(t *testing.T) {
    type TestScenarioExpectParameters struct {
        BaseUrl string
    }

    type TestScenarioParameters struct {
        inputUrl *url.URL
        beginFn func(t *testing.T)
        expectFn func(t *testing.T, tsep TestScenarioExpectParameters)
    }

    type TestScenario struct {
        name string
        setup func(t *testing.T) TestScenarioParameters
    }

    var TestScenarios = []TestScenario{
        {
            name: "invalid URL",
            setup: func(t *testing.T) TestScenarioParameters {
                return TestScenarioParameters {
                    inputUrl: nil,
                    beginFn: func(t *testing.T) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.BaseUrl != "" {
                            t.Errorf("Received an unexpected URL value: \"%s\"", tsep.BaseUrl)
                        }
                    },
                }
            },
        },
        {
            name: "valid URL",
            setup: func(t *testing.T) TestScenarioParameters {
                baseUrl := "https://foobar"
                fullUrl := fmt.Sprintf("%s/deadbeef/v1", baseUrl)
                urlGen, urlGenErr := url.Parse(fullUrl)
                if urlGenErr != nil {
                    t.Errorf("Failed to generate the URL object: %v", urlGenErr)
                }
                return TestScenarioParameters {
                    inputUrl: urlGen,
                    beginFn: func(t *testing.T) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.BaseUrl != baseUrl {
                            t.Errorf("Received \"%s\" when \"%s\" was expected", tsep.BaseUrl, baseUrl)
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

                result := util.GetUrlBaseResource(parameters.inputUrl)
                parameters.expectFn(
                    t,
                    TestScenarioExpectParameters{
                        BaseUrl: result,
                    },
                )
            },
        )
    }
}

func TestUtilEqualHeaders(t *testing.T) {
    type TestScenarioExpectParameters struct {
        result bool
        err    error
    }

    type TestScenarioParameters struct {
        leftHeader  *http.Header
        rightHeader *http.Header
        beginFn     func(t *testing.T)
        expectFn    func(t *testing.T, tsep TestScenarioExpectParameters)
    }

    type TestScenario struct {
        name string
        setup func(t *testing.T) TestScenarioParameters
    }

    var TestScenarios = []TestScenario{
        {
            name: "nil left",
            setup: func(t *testing.T) TestScenarioParameters {
                rightHeader := make(http.Header)
                return TestScenarioParameters{
                    leftHeader: nil,
                    rightHeader: &rightHeader,
                    beginFn: func(t *testing.T) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.err == nil {
                            t.Errorf("Failed to receive an expected error")
                        }
                        if tsep.result != false {
                            t.Errorf("Unexpected received a match")
                        }
                    },
                }
            },
        },
        {
            name: "nil right",
            setup: func(t *testing.T) TestScenarioParameters {
                leftHeader := make(http.Header)
                return TestScenarioParameters{
                    leftHeader: &leftHeader,
                    rightHeader: nil,
                    beginFn: func(t *testing.T) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.err == nil {
                            t.Errorf("Failed to receive an expected error")
                        }
                        if tsep.result != false {
                            t.Errorf("Unexpected received a match")
                        }
                    },
                }
            },
        },
        {
            name: "lengths don't match",
            setup: func(t *testing.T) TestScenarioParameters {
                leftHeader := make(http.Header)
                rightHeader := make(http.Header)
                rightHeader["foo"] = []string{
                    "bar",
                    "d34d",
                    "b33f",
                }
                return TestScenarioParameters{
                    leftHeader: &leftHeader,
                    rightHeader: &rightHeader,
                    beginFn: func(t *testing.T) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.err == nil {
                            t.Errorf("Failed to receive an expected error")
                        }
                        if tsep.result != false {
                            t.Errorf("Unexpected received a match")
                        }
                    },
                }
            },
        },
        {
            name: "missing keys",
            setup: func(t *testing.T) TestScenarioParameters {
                leftHeader := make(http.Header)
                leftHeader["bar"] = []string{
                    "foo",
                    "d34d",
                    "b33f",
                }
                rightHeader := make(http.Header)
                rightHeader["foo"] = []string{
                    "bar",
                    "d34d",
                    "b33f",
                }
                return TestScenarioParameters{
                    leftHeader: &leftHeader,
                    rightHeader: &rightHeader,
                    beginFn: func(t *testing.T) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.err == nil {
                            t.Errorf("Failed to receive an expected error")
                        }
                        if tsep.result != false {
                            t.Errorf("Unexpected received a match")
                        }
                    },
                }
            },
        },
        {
            name: "key value length match failure",
            setup: func(t *testing.T) TestScenarioParameters {
                leftHeader := make(http.Header)
                leftHeader["foo"] = []string{
                    "foo",
                    "d34d",
                    "b33f",
                }
                rightHeader := make(http.Header)
                rightHeader["foo"] = []string{
                    "bar",
                    "d34d",
                }
                return TestScenarioParameters{
                    leftHeader: &leftHeader,
                    rightHeader: &rightHeader,
                    beginFn: func(t *testing.T) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.err == nil {
                            t.Errorf("Failed to receive an expected error")
                        }
                        if tsep.result != false {
                            t.Errorf("Unexpected received a match")
                        }
                    },
                }
            },
        },
        {
            name: "key values match failure",
            setup: func(t *testing.T) TestScenarioParameters {
                leftHeader := make(http.Header)
                leftHeader["foo"] = []string{
                    "foo",
                    "d34d",
                    "b33f",
                }
                rightHeader := make(http.Header)
                rightHeader["foo"] = []string{
                    "foo",
                    "d34d",
                    "sortie",
                }
                return TestScenarioParameters{
                    leftHeader: &leftHeader,
                    rightHeader: &rightHeader,
                    beginFn: func(t *testing.T) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.err == nil {
                            t.Errorf("Failed to receive an expected error")
                        }
                        if tsep.result != false {
                            t.Errorf("Unexpected received a match")
                        }
                    },
                }
            },
        },
        {
            name: "success",
            setup: func(t *testing.T) TestScenarioParameters {
                leftHeader := make(http.Header)
                leftHeader["foo"] = []string{
                    "foo",
                    "d34d",
                    "b33f",
                }
                rightHeader := make(http.Header)
                rightHeader["foo"] = []string{
                    "foo",
                    "d34d",
                    "b33f",
                }
                return TestScenarioParameters{
                    leftHeader: &leftHeader,
                    rightHeader: &rightHeader,
                    beginFn: func(t *testing.T) {},
                    expectFn: func(t *testing.T, tsep TestScenarioExpectParameters) {
                        if tsep.err != nil {
                            t.Errorf("Unexpected received an error: %v", tsep.err)
                        }
                        if tsep.result != true {
                            t.Errorf("Unexpected received a match failure")
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

                result, err := util.EqualHeaders(
                    parameters.leftHeader,
                    parameters.rightHeader,
                )
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
