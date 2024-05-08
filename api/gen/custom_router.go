package openapi

type CustomRouter struct {
	Router
	rts Routes
}

func NewCustomRouter(router Router, rts Routes) *CustomRouter {
	return &CustomRouter{router, rts}
}

func (cr *CustomRouter) Routes() Routes {
	return cr.rts
}
