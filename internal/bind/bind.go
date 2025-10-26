// internal/bind/bind.go
package bind

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
	qdec     = func() *schema.Decoder {
		d := schema.NewDecoder()
		d.IgnoreUnknownKeys(true)
		d.ZeroEmpty(true)
		return d
	}()
)

// FieldError is a clean error for clients.
type FieldError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message,omitempty"`
}

func toFieldErrors(err error) []FieldError {
	var ferrs []FieldError
	if err == nil {
		return ferrs
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			ferrs = append(ferrs, FieldError{
				Field: fe.Field(),
				Tag:   fe.Tag(),
				Param: fe.Param(),
			})
		}
		return ferrs
	}
	// generic
	return []FieldError{{Field: "", Tag: "invalid", Message: err.Error()}}
}

// JSON decodes a JSON body into T and validates it.
func JSON[T any](r *http.Request) (T, []FieldError, error) {
	var dst T
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&dst); err != nil {
		return dst, []FieldError{{Tag: "json", Message: err.Error()}}, err
	}
	if err := validate.Struct(dst); err != nil {
		return dst, toFieldErrors(err), err
	}
	return dst, nil, nil
}

// Query decodes ?query params into T (use `form:"..."` tags) and validates it.
func Query[T any](r *http.Request) (T, []FieldError, error) {
	var dst T
	if err := qdec.Decode(&dst, r.URL.Query()); err != nil {
		return dst, []FieldError{{Tag: "query", Message: err.Error()}}, err
	}
	if err := validate.Struct(dst); err != nil {
		return dst, toFieldErrors(err), err
	}
	return dst, nil, nil
}

// Path fills struct fields from route params (use `path:"id"` tags) and validates it.
// Works with chi (chi.URLParam), gin (c.Param), etc. Here’s a chi-friendly version.
type pathGetter func(*http.Request, string) string

func PathWith[T any](r *http.Request, get pathGetter) (T, []FieldError, error) {
	var dst T
	v := reflect.ValueOf(&dst).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("path")
		if tag == "" {
			continue
		}
		if v.Field(i).Kind() == reflect.String {
			v.Field(i).SetString(get(r, tag))
		}
	}
	if err := validate.Struct(dst); err != nil {
		return dst, toFieldErrors(err), err
	}
	return dst, nil, nil
}

// Simple helpers to map common libraries without importing them here.
// For chi:
func ChiPath[T any](r *http.Request, chiURLParam func(*http.Request, string) string) (T, []FieldError, error) {
	return PathWith[T](r, chiURLParam)
}

// Helper to format joined error messages if you want a single string.
func JoinMessages(ferrs []FieldError) string {
	if len(ferrs) == 0 {
		return ""
	}
	var parts []string
	for _, e := range ferrs {
		if e.Field != "" {
			parts = append(parts, e.Field+":"+e.Tag)
		} else {
			parts = append(parts, e.Tag)
		}
	}
	return strings.Join(parts, ", ")
}

// Expose validator if you need custom tags in main.
func Validator() *validator.Validate { return validate }
