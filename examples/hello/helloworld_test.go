package hello_test

import (
    "errors"
    "net/http"
    "io"
    "testing"

    "github.com/TestInABox/gostackinabox/examples/hello"
    "github.com/TestInABox/gostackinabox/router"
)

func Test_HelloWorldService(t *testing.T) {
    t.Run(
        "GET",
        func(t *testing.T) {
            // configure the HTTP Client
            r := router.New()
            http.DefaultClient.Transport = r

            // create a Go-Stack-In-A-Box service
            hwService, hwServiceErr := hello.NewHelloWorldService()
            if hwServiceErr != nil {
                t.Errorf("Failed to create hello World Service: %#v", hwServiceErr)
            }

            // expected result
            expectedBody := "hello world!"
            expectedBodyLength := len(expectedBody)
            expectedStatusCode := 200

            // register it with the client
            serviceUrl := "https://hello.world"
            registerErr := r.RegisterService(serviceUrl, hwService)
            if registerErr != nil {
                t.Errorf("Failed to register hello world service: %#v", registerErr)
            }

            // attempt the service call
            resp, respErr := http.Get(serviceUrl)
            if respErr != nil {
                t.Errorf("Error making HTTP Call: %#v", respErr)
            }
            if resp != nil {
                // validate the status code response
                if resp.StatusCode != expectedStatusCode {
                    t.Errorf("Unexpected status code: %d != %d", resp.StatusCode, expectedStatusCode)
                }
                // validate the body length
                if resp.ContentLength != int64(expectedBodyLength) {
                    t.Errorf("Unexpected body length: %d != %d", resp.ContentLength, expectedBodyLength)
                }

                // access the response body and validate it
                bodyData := make([]byte, 2*expectedBodyLength)
                readDataLength, readDataErr := resp.Body.Read(bodyData)
                if readDataErr != nil {
                    t.Errorf("Unexpected error reading data: %#v", readDataErr)
                }
                t.Logf("Data Length: %d, Data: %#v", readDataLength, bodyData)
                if readDataLength != expectedBodyLength {
                    t.Errorf("Unexpected data length read: %d != %d", readDataLength, expectedBodyLength)
                }

                // read gives back a byte array, to it must be converted to a string to convert
                strBodyData := string(bodyData[:readDataLength])
                if strBodyData != expectedBody {
                    t.Errorf("Unexpected body data received: \"%s\" != \"%s\"", strBodyData, expectedBody)
                }
            } else {
                t.Errorf("Unexpected nil response")
            }
        },
    )
    t.Run(
        "HEAD",
        func(t *testing.T) {
            // configure the HTTP Client
            r := router.New()
            http.DefaultClient.Transport = r

            // create a Go-Stack-In-A-Box service
            hwService, hwServiceErr := hello.NewHelloWorldService()
            if hwServiceErr != nil {
                t.Errorf("Failed to create hello World Service: %#v", hwServiceErr)
            }

            // expected result
            expectedHeaders := []string{"hello world!"}
            expectedStatusCode := 200
            // HEAD does not have a response body
            expectedBodyLength := 0

            // register it with the client
            serviceUrl := "https://hello.world"
            registerErr := r.RegisterService(serviceUrl, hwService)
            if registerErr != nil {
                t.Errorf("Failed to register hello world service: %#v", registerErr)
            }

            // attempt the service call
            resp, respErr := http.Head(serviceUrl)
            if respErr != nil {
                t.Errorf("Error making HTTP Call: %#v", respErr)
            }
            if resp != nil {
                // validate the status code response
                if resp.StatusCode != expectedStatusCode {
                    t.Errorf("Unexpected status code: %d != %d", resp.StatusCode, expectedStatusCode)
                }
                // validate the body length
                if resp.ContentLength != int64(expectedBodyLength) {
                    t.Errorf("Unexpected body length: %d != %d", resp.ContentLength, expectedBodyLength)
                }

                // access the response body and validate it
                var bodyData []byte
                readDataLength, readDataErr := resp.Body.Read(bodyData)
                if !errors.Is(readDataErr, io.EOF) {
                    t.Errorf("Unexpected error reading data: %#v", readDataErr)
                }
                if readDataLength != expectedBodyLength {
                    t.Errorf("Unexpected data length read: %d != %d", readDataLength, expectedBodyLength)
                }

                // read gives back a byte array, to it must be converted to a string to convert
                if strHeaderData, ok := resp.Header["x-msg"]; ok {
                    if len(strHeaderData) != len(expectedHeaders) {
                        t.Errorf("Header field length does not match: %d != %d", len(strHeaderData), len(expectedHeaders))
                    } else {
                        if strHeaderData[0] != expectedHeaders[0] {
                            t.Errorf("Unexpected body data received: \"%s\" != \"%s\"", strHeaderData[9], expectedHeaders[0])
                        }
                    }
                } else {
                    t.Errorf("Missing headers for `x-msg`")
                }
            } else {
                t.Errorf("Unexpected Nil response")
            }
        },
    )
}
