package ara

import (
    "testing"
    "os"
    "bytes"
    "io/ioutil"
    "fmt"
)

func TestInitRouter(t *testing.T) {
    setup(t)
    defer teardown(t)

    router := NewRouter()

    if len(router.routes) <= 0 {
        t.Fatal("empty route map")
    }

    fmt.Println(router.String())

    rt := router.routes[0]

    if rt.uri != "/abc" || rt.handler == nil {
        t.Fatal("route: " + rt.uri)
    }
}

func TestBuildRoute(t *testing.T) {
    rt, err := buildRoute("/abc XxxHandler")
    if err != nil {
        t.Fatal(err)
    }

    if rt.uri != "/abc" {
        t.Fatal("uri error:", rt.uri)
    }
}

func setup(t *testing.T) {
    err := os.Mkdir("conf", 0700)
    if err != nil {
        if !os.IsExist(err) {
            t.Error(err)
            return
        }
    }
    var buffer bytes.Buffer
    buffer.WriteString("#this is the 1st line comment\n")
    buffer.WriteString(" #and with space at start, end   \n")
    buffer.WriteString("/abc      AbcHandler\n")
    buffer.WriteString("/         static/\n")
    err = ioutil.WriteFile("conf/router", buffer.Bytes(), 0600)
    if err != nil {
        t.Error(err)
        return
    }
}

func teardown(t *testing.T) {
    if err := os.RemoveAll("conf"); err != nil {
        t.Error(err)
    }
}
