# Ara
A simple web framework that supports RESTful router rules, and text-based handler configuration(see the following example)

# Example
- Install

```
go get github.com/araframework/ara
```
- Dir structure overview

```
yourapp
    |_conf
    |   |_router     <-*required
    |_static         <-optional
    |   |_index.html
    |   |_image
    |   |_js
    |   |_fonts
    |   |_css
    |_controller.go
    |_main.go
```
- controller.go

```go
package main
import (
    "io"
    "net/http"
    "github.com/araframework/ara"
)

type Controller struct {
    ara.Controller
}

// process /hello
func (c *Controller) MyHandler(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Hello, 世界")
}

// process /hello/{id}
func (c *Controller) MyHandlerWithId(w http.ResponseWriter, r *http.Request) {
    id := r.Form.Get("id")
    w.Write([]byte(id))
}
```
- main.go

```go
package main
import (
    "github.com/araframework/ara"
)

func main() {
    controller := &controller.Controller{}
    
    router := ara.NewRouter()
    router.SetController(controller)
    
    ara.Start(router)
}
```

- conf/router

```
/             FS:static
/hello        MyHandler
/hello/{id}   MyHandlerWithId
```
- `go build` and run the executable file
- Navigate browser to `http://localhost:8600` will show the index.html
- And `http://localhost:8600/hello`will show `Hello, 世界`
- And `http://localhost:8600/hello/abc123` will show `abc123`

# TODO
- namespace support
