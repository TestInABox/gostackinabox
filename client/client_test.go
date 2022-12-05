package client_test

import (
    "testing"

    "github.com/TestInABox/gostackinabox/client"
)

func Test_Client_DefaultClient(t *testing.T) {
    if client.DefaultClient == nil {
        t.Errorf("Unexpected nil client: %#v", client.DefaultClient)
    }
}
