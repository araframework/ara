package ara

import (
    "testing"
    "net/http"
    "io"
)

type MyT T

func TestStart(t *testing.T) {
    setup(t)
    defer teardown(t)
    // Start()
}

func (t *MyT) AbcHandler(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "AbcHandler")
}
