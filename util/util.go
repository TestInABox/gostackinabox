package util

import (
    "bytes"
    "fmt"
    "net/http"
    "net/url"
    "io"
    "io/ioutil"
)

func StringToResponseBody(s string) io.ReadCloser {
    b := bytes.NewBufferString(s)
    r := bytes.NewReader(b.Bytes())
    return ioutil.NopCloser(r)
}

func GetUrlBaseResource(url *url.URL) string {
    if url != nil {
        return fmt.Sprintf(
            "%s://%s",
            url.Scheme,
            url.Host,
        )
    }

    return ""
}

func EqualHeaders(l *http.Header, r *http.Header) (matches bool, err error) {
    // if the left differs from the right then return false; otherwise return true
    // NOTE: This is namely a test support function

    // make sure there's something to compare
    if l == nil || r == nil {
        err = fmt.Errorf("Either left or right is an invalid pointer: Left (%v), Right (%v)", l, r)
        matches = false
        return
    }

    if len(*l) != len(*r) {
        err = fmt.Errorf("Lengths do not match: %d ! %d\n", len(*l), len(*r))
        matches = false
        return
    }

    for lk, lv := range *l {
        rv, ok := (*r)[lk]
        if !ok {
            err = fmt.Errorf("Key %v not found on the right\n", lk)
            matches = false
            return
        }

        if len(lv) != len(rv) {
            err = fmt.Errorf("header %v  array length doesn't match: %d != %d\n", lk, len(lv), len(rv))
            matches = false
            return
        }

        for lvi, lvv := range lv {
            rvv := rv[lvi]
            if lvv != rvv {
                err = fmt.Errorf("header %v array index %d doesn't match: %v != %v\n", lk, lvi, lvv, rvv)
                matches = false
                return
            }
        }
    }

    // length matches, and each value matches so they must be equal
    matches = true
    return
}
