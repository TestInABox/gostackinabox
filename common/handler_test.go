package common_test

import (
    "testing"

    "github.com/TestInABox/gostackinabox/common"
)

func Test_Common_Handler(t *testing.T) {
    fn := func(hc *common.HttpCall) (hr *common.HttpReply, err error) {
        return
    }

    key := common.HttpVerb("test")
    var hhm common.HttpHandlerMap = make(common.HttpHandlerMap)
    if len(hhm) != 0 {
        t.Errorf("Unepected non-empty map: length: %d, %#v", len(hhm), hhm)
    }
    hhm[key] = fn
    if len(hhm) != 1 {
        t.Errorf("Unepected map size: length: %d, %#v", len(hhm), hhm)
    }

    if _, ok := hhm[key]; !ok {
        t.Errorf("Unexpectedly did not find the key (%s) in the map (%v)", key, hhm)
    }
}
