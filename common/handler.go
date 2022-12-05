package common

type HttpHandler func(*HttpCall) (*HttpReply, error)
type HttpHandlerMap map[HttpVerb]HttpHandler
