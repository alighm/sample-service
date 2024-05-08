package log

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/alighm/sample-service/ctxutil"
	"github.com/alighm/sample-service/util"
)

func APILog(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := Logger(ctxutil.WithRequestID(r.Context(),
			strfmt.UUID(r.Header.Get("X-Request-ID")))).WithField("pkg", "api.openapi")

		start := time.Now()
		rw := NewLoggingResponseWriter(w)

		// Set Correlation ID [X-Request-ID header]
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = util.NewUUID()
		}
		if !util.IsUUID(requestID) {
			log.WithFields(map[string]interface{}{
				"api.method":    r.Method,
				"api.uri":       r.RequestURI,
				"CorrelationId": requestID,
			}).Errorf("received invalid X-Request-ID: %v", requestID)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"msg": "X-Request-ID header is invalid. Must be a UUID"})
			return
		}
		ctx = ctxutil.WithRequestID(ctx, strfmt.UUID(requestID))

		log.WithFields(map[string]interface{}{
			"api.method":    r.Method,
			"api.uri":       r.RequestURI,
			"CorrelationId": requestID,
		}).Infof("** start api call for method name: %s", name)

		req := r.Clone(ctx)
		inner.ServeHTTP(rw, req)

		log.WithFields(map[string]interface{}{
			"api.method":    r.Method,
			"api.uri":       r.RequestURI,
			"api.status":    rw.statusCode,
			"CorrelationId": requestID,
			"api.duration":  time.Since(start),
		}).Infof("**** end api call for method name: %s", name)
	})
}
