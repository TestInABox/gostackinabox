package common_test

import (
    "testing"

    "github.com/TestInABox/gostackinabox/common"
)

func Test_Common_HttpStatus(t *testing.T) {
    status := map[common.HttpStatusCode]int{}
    status[common.HttpStatus_MethodNotSupport] = 405
    status[common.HttpStatus_RouteNotHandled] = 595
    status[common.HttpStatus_ServiceError] = 596
    status[common.HttpStatus_ServiceSubRouteError] = 597

    for v, k := range status {
        if int(v) != k {
            t.Errorf("Key Value %v != %v", v, k)
        }
    }
}
