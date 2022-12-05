package common_test

import (
    "net/http"
    "net/url"
    "testing"

    "github.com/TestInABox/gostackinabox/common"
)

func Test_Common_HttpCall(t *testing.T) {
    hc := &common.HttpCall{
        Method: "test",
        Url: &url.URL{},
        Headers: http.Header{},
        Request: &http.Request{},
    }
    if hc.Method != "test" {
        t.Errorf("Failed to set method")
    }
}
