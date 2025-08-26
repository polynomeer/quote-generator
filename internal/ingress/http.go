package ingress

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/polynomeer/quote-generator/internal/generator"
)

func RegisterHTTP(mux *http.ServeMux, gen *generator.Generator) {
	mux.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/v1/quotes", func(w http.ResponseWriter, r *http.Request) {
		syParam := r.URL.Query().Get("symbols")
		if syParam == "" {
			http.Error(w, "symbols required", http.StatusBadRequest)
			return
		}
		symbols := strings.Split(syParam, ",")
		quotes := gen.Snapshot(symbols)
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		_ = enc.Encode(quotes)
	})
}
