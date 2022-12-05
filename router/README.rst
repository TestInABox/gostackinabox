******
Router
******

This implements the Golang net.http.RoundTripper interface in order to
achieve the HTTP Request Interception, and then redirecting to
various services implementing the desired APIs.

Where the Python version of StackInABox allowed for different interceptor
libraries the Golang version provides its own, in part due to the fact that
there are no interceptor libraries in Golang like there are in Python; but
also due to the ease of which Golang enables the interception using a
standardized interface provided directly by Golang.
