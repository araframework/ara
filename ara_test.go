package ara

import (
    "testing"
    "net/http"
    "io"
)

type MyController struct {
    Controller
}

func TestStart(t *testing.T) {
    setup(t)
    defer teardown(t)
    // Start()
}

func (c *MyController) AbcHandler(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "AbcHandler")
}
