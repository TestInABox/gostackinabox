package common

import (
    "fmt"
    "net/url"
    "regexp"

    "github.com/TestInABox/gostackinabox/common/log"
)

// matches based on paths
type  PathURI struct {
    Method string
    Path *regexp.Regexp
}

func (pu *PathURI) IsMatch(u url.URL) (result bool, err error) {
    if pu.Path == nil || (len(pu.Path.String()) == 0) {
        log.Printf("No path regex available to match %s against", u.String())
        err = fmt.Errorf("%w, missing regular expression to match with", ErrPathURIMisconfigured)
        return
    }

    // use the regex to match
    result = pu.Path.MatchString(u.Path)
    log.Printf("Attempting to match %s against %s... match: %t", u.String(), pu.Path.String(), result)
    return
}

var _ URI = &PathURI{}
