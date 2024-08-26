package apicontroller

type APIController interface {
	Pattern() string
	Handlers() []*APIHandler
}
