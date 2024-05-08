/*
 * Sample Service
 *
 * This is a sample service
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// UserAPIController binds http requests to an api service and writes the service results to the http response
type UserAPIController struct {
	service      UserAPIServicer
	errorHandler ErrorHandler
}

// UserAPIOption for how the controller is set up.
type UserAPIOption func(*UserAPIController)

// WithUserAPIErrorHandler inject ErrorHandler into controller
func WithUserAPIErrorHandler(h ErrorHandler) UserAPIOption {
	return func(c *UserAPIController) {
		c.errorHandler = h
	}
}

// NewUserAPIController creates a default api controller
func NewUserAPIController(s UserAPIServicer, opts ...UserAPIOption) Router {
	controller := &UserAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the UserAPIController
func (c *UserAPIController) Routes() Routes {
	return Routes{
		"AddUser": Route{
			strings.ToUpper("Post"),
			"/users",
			c.AddUser,
		},
		"DeleteUser": Route{
			strings.ToUpper("Delete"),
			"/users/{userId}",
			c.DeleteUser,
		},
		"GetUser": Route{
			strings.ToUpper("Get"),
			"/users/{userId}",
			c.GetUser,
		},
		"GetUsers": Route{
			strings.ToUpper("Get"),
			"/users",
			c.GetUsers,
		},
		"UpdateUser": Route{
			strings.ToUpper("Put"),
			"/users/{userId}",
			c.UpdateUser,
		},
	}
}

// AddUser - Add a new user
func (c *UserAPIController) AddUser(w http.ResponseWriter, r *http.Request) {
	xRequestIDParam := r.Header.Get("X-Request-ID")
	userParam := User{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&userParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertUserRequired(userParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertUserConstraints(userParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.AddUser(r.Context(), xRequestIDParam, userParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// DeleteUser - Deletes a user
func (c *UserAPIController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	xRequestIDParam := r.Header.Get("X-Request-ID")
	userIdParam := params["userId"]
	if userIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"userId"}, nil)
		return
	}
	result, err := c.service.DeleteUser(r.Context(), xRequestIDParam, userIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetUser - Find user by ID
func (c *UserAPIController) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	xRequestIDParam := r.Header.Get("X-Request-ID")
	userIdParam := params["userId"]
	if userIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"userId"}, nil)
		return
	}
	result, err := c.service.GetUser(r.Context(), xRequestIDParam, userIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetUsers - Get all users
func (c *UserAPIController) GetUsers(w http.ResponseWriter, r *http.Request) {
	xRequestIDParam := r.Header.Get("X-Request-ID")
	result, err := c.service.GetUsers(r.Context(), xRequestIDParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateUser - Updates a user
func (c *UserAPIController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	xRequestIDParam := r.Header.Get("X-Request-ID")
	userIdParam := params["userId"]
	if userIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"userId"}, nil)
		return
	}
	userParam := User{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&userParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertUserRequired(userParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertUserConstraints(userParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateUser(r.Context(), xRequestIDParam, userIdParam, userParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
