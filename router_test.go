package ara

import (
    "testing"
    "os"
    "bytes"
    "io/ioutil"
)

func Test_initRouter(t *testing.T) {
    setup(t)
    defer teardown(t)

    initRouter()

    if len(routeMap) <= 0 {
        t.Fatal("empty route map")
    }

    rt, ok := routeMap["GET/abc"]
    if !ok {
        t.Fatal("specified route not found from map")
    }

    if rt.uri != "/abc" || rt.method != GET || rt.handler != "AbcHandler" {
        t.Fatal("route: " + rt.method + " : " + rt.uri + " : " + rt.handler)
    }
}

func Test_buildRoute(t *testing.T) {
    rt, err := buildRoute("GET /abc XxxHandler")
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
    buffer.WriteString("GET /abc      AbcHandler\n")
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
