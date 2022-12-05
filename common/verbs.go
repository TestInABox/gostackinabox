package common

type HttpVerb string

const (
    HttpVerb_Connect HttpVerb = "CONNECT"
    HttpVerb_Delete  HttpVerb = "DELETE"
    HttpVerb_Get     HttpVerb = "GET"
    HttpVerb_Head    HttpVerb = "HEAD"
    HttpVerb_Option  HttpVerb = "OPTION"
    HttpVerb_Patch   HttpVerb = "PATCH"
    HttpVerb_Post    HttpVerb = "POST"
    HttpVerb_Put     HttpVerb = "PUT"
    HttpVerb_Trace   HttpVerb = "TRACE"
)

func GetHttpVerb(verb string) HttpVerb {
    return HttpVerb(verb)
}
