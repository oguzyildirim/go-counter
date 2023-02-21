package rest

import (
	"net/http"
)

// CounterService
type CounterService interface {
	Create() error
	Find() (string, error)
}

// MiddlewareService
type MiddlewareService interface {
	Handle(next http.Handler) http.Handler
}

// CounterHandler
type CounterHandler struct {
	svc         CounterService
	ratelimiter MiddlewareService
}

// NewCounterHandler
func NewCounterHandler(svc CounterService, ratelimiter MiddlewareService) *CounterHandler {
	return &CounterHandler{
		svc:         svc,
		ratelimiter: ratelimiter,
	}
}

// Register connects the handlers to the router
func (c *CounterHandler) Register(r *http.ServeMux) {
	r.Handle("/count", c.ratelimiter.Handle(c.find()))
	r.Handle("/", c.ratelimiter.Handle(c.create()))
}

// FindCounterResponse defines the response returned back after finding count
type FindCounterResponse struct {
	Count string `json:"count"`
}

func (c *CounterHandler) find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val, err := c.svc.Find()
		if err != nil {
			renderErrorResponse(r.Context(), w, "create failed", err)
			return
		}
		renderResponse(w, &FindCounterResponse{Count: val}, http.StatusOK)
	}
}

func (c *CounterHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := c.svc.Create()
		if err != nil {
			renderErrorResponse(r.Context(), w, "create failed", err)
			return
		}
		renderResponse(w, struct{}{}, http.StatusCreated)
	}
}
