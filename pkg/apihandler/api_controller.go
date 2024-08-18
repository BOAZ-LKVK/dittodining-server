package apihandler

type APIController interface {
	Pattern() string
	Handlers() []*APIHandler
}
