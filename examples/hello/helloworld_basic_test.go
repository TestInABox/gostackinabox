package hello_test

import (
    //"errors"
    "net/http"
    //"io"
    "testing"

    "github.com/TestInABox/gostackinabox/examples/hello"
    "github.com/TestInABox/gostackinabox/router"
)

func Test_HelloWorldBasicService(t *testing.T) {
    t.Run(
        "GET",
        func(t *testing.T) {
            // configure the HTTP Client
            r := router.New()
            http.DefaultClient.Transport = r

            // create a Go-Stack-In-A-Box service
            hwbService, hwbServiceErr := hello.NewHelloWorldBasicService()
            if hwbServiceErr != nil {
                t.Errorf("Failed to create hello World Service: %#v", hwbServiceErr)
            }

            // expected result
            expectedBody := "basic hello world!"
            expectedBodyLength := len(expectedBody)
            expectedStatusCode := 200

            // register it with the client
            serviceUrl := "https://hello.world"
            registerErr := r.RegisterService(serviceUrl, hwbService)
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
}

