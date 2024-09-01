package erk

import (
  "fmt"
  "log/slog"
  "net/http"
)

// HTTPError is an error which contains extra information which can be processes by a panic http handler
type HTTPError interface {
  error
  HTTP() (reason string, code int)
}

type staticHTTPError struct {
  error
  code int
}

func (e staticHTTPError) HTTP() (reason string, code int) {
  return e.Error(), e.code
}

// NotFound wraps err so that it implements HTTPError if it is non-nil.
func NotFound(err error) error {
  if err == nil {
    return nil
  }

  return staticHTTPError{err, http.StatusNotFound}
}

// BadRequest wraps err so that it implements HTTPError if it is non-nil.
func BadRequest(err error) error {
  if err == nil {
    return nil
  }

  return staticHTTPError{err, http.StatusBadRequest}
}

// Forbidden wraps err so that it implements HTTPError if it is non-nil.
func Forbidden(err error) error {
  if err == nil {
    return nil
  }

  return staticHTTPError{err, http.StatusForbidden}
}

// Unauthorized wraps err so that it implements HTTPError if it is non-nil.
func Unauthorized(err error) error {
  if err == nil {
    return nil
  }

  return staticHTTPError{err, http.StatusUnauthorized}
}

// PanicHandler wraps inner to handle panics. Panicked values should implement HTTPError.
func PanicHandler(inner http.Handler) http.Handler {
  return panicHandler{inner}
}

type panicHandler struct {
  inner http.Handler
}

func (h panicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  defer func() {
    if err := recover(); err != nil {
      if httpErr, isHttpErr := err.(HTTPError); isHttpErr {
        reason, code := httpErr.HTTP()
        if code == http.StatusInternalServerError {
          slog.ErrorContext(r.Context(), "internal server error", "error", err)
        }
        http.Error(w, reason, code)
      } else {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        slog.WarnContext(r.Context(), "recovered from panic with non HTTPError value", "type", fmt.Sprintf("%T", err))
        slog.ErrorContext(r.Context(), "internal server error", "error", err)
      }
    }
  }()

  h.inner.ServeHTTP(w, r)
}
