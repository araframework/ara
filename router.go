package ara

import (
    "fmt"
    "bufio"
    "os"
    "path"
    "strings"
    "errors"
    "net/http"
    "reflect"
)

const (
    GET = "GET"
    POST = "POST"
    PUT = "PUT"
    DELETE = "DELETE"
)

// define route: uri pattern -> handler function name
type Route struct {
    uri     string       // the request uri
    handler http.Handler // the function to handle this request
}

// contains all routes
type Router struct {
    // Configurable Handler to be used when no route matches.
    NotFoundHandler http.Handler
    // Routes to be matched, in order.
    routes          []*Route
    // Routes by name for URL building.
    namedRoutes     map[string]*Route

    // this is for string method name reflect invoke
    ControllerValue reflect.Value
}

func (route *Route)String() string {
    str := "Route:[" +
    "uri:" + route.uri + "," +
    "]"

    return str
}

func (router * Router) String() string {
    str := "Router: ["
    for _, route := range router.routes {
        str += route.String()
        str += " "
    }
    str += "]"
    return str
}

// ------------------ init ---------------------

func NewRouter() *Router {
    router := &Router{namedRoutes: make(map[string]*Route)}
    router.initRouter()
    return router
}

func (router *Router) SetControllerValue(cValue reflect.Value) {
    router.ControllerValue = cValue
}

func (router *Router) initRouter() {
    alog.Debug("router:initRouter()")

    // the router conf file named: router
    confPath := path.Join("conf", "router")
    file, err := os.Open(confPath)
    if err != nil {
        alog.Debug(err.Error())
        os.Exit(100)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    if err := scanner.Err(); err != nil {
        alog.Debug(err.Error())
        os.Exit(101)
    }

    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(line) // Println will add back the final '\n'

        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "#") {
            continue
        }

        rt, err := router.buildRoute(line)
        if err != nil {
            continue
        }

        router.routes = append(router.routes, rt)
    }
}

// validate the route configured in router file
func (router *Router) buildRoute(line string) (rt *Route, err error) {
    routeItems := strings.SplitN(line, " ", 2)
    for _, item := range routeItems {
        if strings.TrimSpace(item) == "" {
            alog.Debug("empty item for line" + line)
            alog.Debug(item)
            err = errors.New("invalid route: " + line)
            return
        }
    }

    uri := strings.TrimSpace(routeItems[0])
    var handler http.Handler
    // if starts with "FS:", that's FileServerHandler
    handleFuncName := strings.TrimSpace(routeItems[1])
    isFs := strings.HasPrefix(handleFuncName, "FS:")
    if isFs {
        handler = http.FileServer(http.Dir(handleFuncName[3:]))
    } else {
        handler = http.HandlerFunc(router.makeHandler(handleFuncName))
    }

    // TODO replace this ugly shit whit Router.map[uri]Route
    rt = &Route{uri, handler}
    if uri == "/" {
        rootRoute = rt
    }
    return
}
///////////////////////////// not work properly yet //////////////////
func (router *Router) makeHandler(tp string) http.HandlerFunc {
    fn := func(w http.ResponseWriter, r *http.Request) {
        fmt.Println(tp)

        w.Header().Set("content-type", "application/json")

        method := router.ControllerValue.MethodByName(tp)

        in := make([]reflect.Value, 2)
        in[0] = reflect.ValueOf(w)
        in[1] = reflect.ValueOf(r)

        method.Call(in)
    }
    if fn == nil {
        return func(w http.ResponseWriter, r *http.Request) {
            fmt.Println("404")

            method := router.ControllerValue.MethodByName("NotFound")
            method.Call([]reflect.Value{})
        }
    }

    return fn
}
/////////////////////////////////////////////////////////////////////
// TODO replace this ugly shit whit Router.map[uri]Route
var rootRoute *Route

// register new path -> function
func (router *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
    rt := &Route{uri: pattern, handler: http.HandlerFunc(handler)}
    router.routes = append(router.routes, rt)
    // TODO replace this ugly shit whit Router.map[uri]Route
    if pattern == "/" {
        rootRoute = rt
    }
}

func (router *Router) Handle(pattern string, handler http.Handler) {
    rt := &Route{uri: pattern, handler: handler}
    router.routes = append(router.routes, rt)
    // TODO replace this ugly shit whit Router.map[uri]Route
    if pattern == "/" {
        rootRoute = rt
    }
}

// ---------------- runtime --------------------
// this will run in an incoming http request
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Router ServeHTTP:" + r.URL.String())
    // 1, match regular path pattern to find the handler
    // 2, call the handler function

    for _, route := range router.routes {
        if r.URL.Path == route.uri {
            alog.Debug("found handler")
            route.handler.ServeHTTP(w, r)
            return
        }
    }

    // not found any handler, run root handler "/"
    if rootRoute != nil {
        rootRoute.handler.ServeHTTP(w, r)
    }
}