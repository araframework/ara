package ara

import (
    "testing"
    "os"
    "bytes"
    "io/ioutil"
)

func TestInitRouter(t *testing.T) {
    setup(t)
    defer teardown(t)

    router := NewRouter()

    if len(router.node.children) <= 0 {
        t.Fatal("empty route map")
    }

    abcNode := router.node.children["abc"]

    if abcNode.name != "abc" || abcNode.handler == nil {
        t.Fatal("Node: /abc")
    }

    idNode := abcNode.children["id"]
    if idNode.name != "id" || idNode.handler == nil {
        t.Fatal("Node: /abc/{id}")
    }
}

func TestBuildRoute(t *testing.T) {
    setup(t)
    defer teardown(t)
    router := NewRouter()
    err := router.buildNode("/abc XxxHandler")
    if err != nil {
        t.Fatal(err)
    }

    abcNode := router.node.children["abc"]
    if abcNode.name != "abc" {
        t.Fatal("uri error: /abc")
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
    buffer.WriteString("/         FS:static\n")
    buffer.WriteString("/abc/{id}        AbcIdHandler\n")
//    buffer.WriteString("/abc/{id}/xyz    AbcIdXyzHandler\n")
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
