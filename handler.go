package there

import (
	"net/http"
)

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}

		httpResponse := Error(StatusInternalServerError, err)
		writeHeader(&writer, httpResponse)
		_ = httpResponse.Execute(router, request, &writer)
	}()

	method := request.Method

	httpRequest := NewHttpRequest(request)
	var httpResponse HttpResponse

	errorOut := func(err error) {
		httpResponse = Error(StatusInternalServerError, err)
		writeHeader(&writer, httpResponse)
		_ = httpResponse.Execute(router, request, &writer)
	}

	for _, middleware := range router.GlobalMiddlewares {
		httpResponse = middleware(httpRequest)
		writeHeader(&writer, httpResponse)
		if isNextMiddleware(httpResponse) {
			_ = httpResponse.Execute(router, request, &writer)
		} else {
			break
		}
	}

	var endpoint Endpoint = nil
	var middlewares = make([]Middleware, 0)

	for _, route := range router.routes {
		routeParams, ok := route.Path.Parse(request.URL.Path)
		if ok && CheckArrayContains(route.Methods, method) {
			endpoint = route.Endpoint
			middlewares = route.Middlewares
			routeParamReader := RouteParamReader(routeParams)
			httpRequest.RouteParams = &routeParamReader
			break
		}
	}

	for _, middleware := range middlewares {
		httpResponse = middleware(httpRequest)
		writeHeader(&writer, httpResponse)

		if isNextMiddleware(httpResponse) {
			_ = httpResponse.Execute(router, request, &writer)
		} else {
			break
		}
	}

	if httpResponse == nil || isNextMiddleware(httpResponse) {
		if endpoint == nil {
			endpoint = func(_ HttpRequest) HttpResponse {
				return Error(StatusNotFound, router.RouterConfiguration.RouteNotFound(request))
			}
		}
		httpResponse = endpoint(httpRequest)
		writeHeader(&writer, httpResponse)
	}

	err := httpResponse.Execute(router, request, &writer)
	if err != nil {
		errorOut(err)
		return
	}

}

func writeHeader(w *http.ResponseWriter, httpResponse HttpResponse) {
	for key, values := range httpResponse.Header().Values {
		(*w).Header().Del(key)
		for _, value := range values {
			(*w).Header().Add(key, value)
		}
	}
}

func isNextMiddleware(response HttpResponse) bool {
	switch v := response.(type) {
	case *nextMiddleware:
		return true
	case *HeaderWrapper:
		switch v.HttpResponse.(type) {
		case *nextMiddleware:
			return true
		default:
			return false
		}
	case *contextResponse:
		for {
			switch res := v.response.(type) {
			case *contextResponse:
				v = res
			case *nextMiddleware:
				return true
			default:
				return false
			}
		}
	default:
		return false
	}
}
