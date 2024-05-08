/*
 * OpenAPI Go Reference Service
 *
 * This is a sample Go Reference
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"net/http"

	"github.com/alighm/sample-service/log"
)

func Logger(inner http.Handler, name string) http.Handler {
	return log.APILog(inner, name)
}
