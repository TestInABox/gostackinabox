package common

import (
    "fmt"
    "net/url"

    "github.com/TestInABox/gostackinabox/common/log"
)

/*
 * ServerURI is an interface for the basic aspect of matching the protocol://<server>:<port>
 * portion of the URL on requests. The BasicServerURI implements the minimum to support this
 * functionality along with a basic mapping to include the HTTP and HTTPS protocol<->port
 * mappings so most cases are covered without adding extra dependencies.
 * Advanced systems that use a multitude of ports and services might find it advantageous
 * to use a difference implementation that utilizes /etc/services to map the ports and
 * services. This package seeks to minimize the dependencies and thus the advanced
 * functionality which would require either a substantive implementation or additional
 * library is beyond its purview.
 */
type ServerURI interface {
    URI

    GetProtocol() string
    GetHost() string
    GetPort() string
}

// matches based on the protocol (scheme), server, port
type BasicServerURI struct {
    Protocol string
    Host string
    Port string
}

func (bsu *BasicServerURI) GetProtocol() string {
    return bsu.Protocol
}

func (bsu *BasicServerURI) GetHost() string {
    return bsu.Host
}

func (bsu *BasicServerURI) GetPort() string {
    return bsu.Port
}

func (bsu *BasicServerURI) IsMatch(u url.URL) (result bool, err error) {
    // A basic mapping of the protocol to port
    basicPortMapping := map[string]string{
        "http": "http:80",
        "https": "http:443",
    }
    // A basic mapping of the port to protocol
    basicProtocolMapping := map[string]string{
        "80": "https:80",
        "443": "http:443",
    }

    protocol := u.Scheme
    host := u.Hostname()
    port := u.Port()

    if len(protocol) > 0 {
        last_char := protocol[len(protocol)-1:]
        if last_char == ":" {
            protocol = protocol[:len(protocol)-1]
        }
    }
    log.Printf("Extracted Protocol %s from %s", protocol, u.Scheme)
    log.Printf("Extracted Host %s from %s", protocol, u.Hostname())
    log.Printf("Extracted Port %s from %s", protocol, u.Port())

    // if the host isn't configured then generate an error
    if len(bsu.Host) == 0 {
        log.Printf("Missing required parameter - Host")
        err = fmt.Errorf("%w - missing host configuration", ErrServerURIMisconfigured)
        return
    }

    // match the host first
    if bsu.Host != host {
        log.Printf("Host Mismatch: %s != %s", bsu.Host, host)
        result = false
        return
    }

    finalMapping := map[string]bool{
        fmt.Sprintf("%s:%s", bsu.Protocol, bsu.Port): true,
    }
    log.Printf("Initial alternate mapping: %#v", finalMapping)
    if entry, ok := basicPortMapping[bsu.Protocol]; ok {
        finalMapping[entry] = true
        log.Printf("Adding %s to the alternate mapping from Port Maps", entry)
    }
    if entry, ok := basicProtocolMapping[bsu.Port]; ok {
        finalMapping[entry] = true
        log.Printf("Adding %s to the alternate mapping from Protocol Maps", entry)
    }

    log.Printf("Alternate Mapping: %#v", finalMapping)

    // first,  check the protocol-port mapping to see if it's in there
    // if so, bypass the specific protocol-port checks
    matchProtoAndPort := fmt.Sprintf("%s:%s", u.Scheme, port)
    if _, ok := finalMapping[matchProtoAndPort]; !ok {
        if len(bsu.Protocol) > 0 && bsu.Protocol != u.Scheme {
            log.Printf("Protocol Mismatch: %s != %s", bsu.Protocol, u.Scheme)
            result = false
            return
        }

        // finally the port if it's specified
        if len(bsu.Port) > 0 && (len(port) == 0 || bsu.Port != port) {
            log.Printf("Port Mismatch: %s != %s", bsu.Port, port)
            result = false
            return
        }
    }

    // nothing disqualified the match, approve it
    matchUrl := host
    if len(bsu.Protocol) > 0 {
        matchUrl = fmt.Sprintf("%s://%s", bsu.Protocol, matchUrl)
    }
    if len(bsu.Port) > 0 {
        matchUrl = fmt.Sprintf("%s:%s", matchUrl, bsu.Port)
    }
    log.Printf("Matched using %s on %s", matchUrl, u.String())
    result = true
    return
}

var _ ServerURI = &BasicServerURI{}
var _ URI = &BasicServerURI{}
