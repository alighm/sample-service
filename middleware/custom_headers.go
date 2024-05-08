package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/alighm/sample-service/ctxutil"
)

func CustomHeaders(versions []string, next http.HandlerFunc) http.HandlerFunc {
	if len(versions) < 1 {
		log(context.Background()).Fatalf("no versions supplied to CustomHeaders middleware")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Set Version [X-Version header]
		version := r.Header.Get("X-Version")
		if version == "" {
			version = versions[len(versions)-1]
		} else {
			versionIsValid := false
			for _, v := range versions {
				if v == version {
					versionIsValid = true
				}
			}
			if !versionIsValid {
				log(r.Context()).Errorf("received invalid X-Version: %v", version)
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"msg": fmt.Sprintf("X-Version header is invalid. Must be on of [%v]",
					strings.Join(versions, ", "))})
				return
			}
		}
		ctx = ctxutil.WithVersion(ctx, version)

		req := r.Clone(ctx)
		next.ServeHTTP(w, req)
	}
}
