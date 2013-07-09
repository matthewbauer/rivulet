package main

import (
  "net/http"
)

type Authentication interface {
  NewContext(r *http.Request) Context
}
