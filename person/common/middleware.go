package common

import (
	"net/http"
)

func CORSMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//if r.Method == http.MethodOptions {
		//	return
		//}

		next.ServeHTTP(w, r)
	})
}