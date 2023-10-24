package filter

import (
	"context"
	"net/http"
)

const (
	ASC                         = "ASC"
	DESC                        = "DESC"
	FilterOptionsContextKey Str = "filter_options"
)

type Str string

type Options struct {
	Where map[string][]string
}

// The following Middleware injects filtering options into request context
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Параметры для оператора Where
		fields := make(map[string][]string)
		for k, v := range r.URL.Query() {

			if k != "sort_by" && k != "sort_order" && k != "next_person_id" {
				fields[k] = v
			}

		}

		// Если параметров для фильтрации нет
		if len(fields) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		options := Options{
			Where: fields,
		}

		// Наполним контекст запроса новой парой ключ/значение
		ctx := context.WithValue(r.Context(), FilterOptionsContextKey, options)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
