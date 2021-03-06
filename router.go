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

/////////////// node ////////////////////

// define route: uri pattern -> handler function name
//type Route struct {
//    uri     string       // the request uri
//    handler http.Handler // the function to handle this request
//}

// contains all routes
type Router struct {
    // Configurable Handler to be used when no route matches.
    NotFoundHandler http.Handler

    node            *Node

    // Routes to be matched, in order.
    //    routes          []*Route
    //    // Routes by name for URL building.
    //    namedRoutes     map[string]*Route

    // this is for string method name reflect invoke
//    ControllerValue reflect.Value
    controllerImpl interface{}
}

//func (router * Router) String() string {
//    str := "Router: ["
//    for _, route := range router.routes {
//        str += route.String()
//        str += " "
//    }
//    str += "]"
//    return str
//}

// ------------------ init ---------------------

func NewRouter() *Router {
    // "/" will respond all request
    node := NewNode("/", NODE_STATIC, nil)
    router := &Router{
        node: node}
    router.initRouter()
    return router
}

func (router *Router) SetController(impl interface{}) {
    router.controllerImpl = impl
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
        if strings.HasPrefix(line, "#") || line == "" {
            continue
        }

        err := router.buildNode(line)
        if err != nil {
            alog.Debug(err.Error())
            continue
        }
    }
}

// validate the route configured in router file
func (router *Router) buildNode(line string) (err error) {
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
    handleFuncName := strings.TrimSpace(routeItems[1])
    h := router.makeHttpHandler(handleFuncName)

    if uri == "/" {
        router.node.handler = h
        rootHandler = h
        return
    }

    if strings.HasSuffix(uri, "/") {
        err = errors.New("url should not ends with '/'")
        return
    }

    currNode := router.node
    uriSections := router.splitUri(uri)
    for _, section := range uriSections {
        sec, secType, err := router.getSection(section)
        if err != nil {
            return err
        }

        if currNode.children[sec] == nil {
            currNode.children[sec] = NewNode(sec, secType, nil)
        }
        currNode = currNode.children[sec]
    }

    // only the last leaf should set the handler
    currNode.handler = h

    return
}

func (router * Router) getSection(section string) (sec string, secType uint, err error) {
    if section == "" {
        err = errors.New("section is empty in url")
        return
    }

    if strings.HasPrefix(section, "{") && strings.HasSuffix(section, "}") {
        dynSectionName := strings.TrimSpace(section[1:len(section) - 1])
        if dynSectionName == "" {
            err = errors.New("section is empty in url: " + section)
            return
        }

        sec = dynSectionName
        secType = NODE_DYNAMIC
        return
    }

    sec = section
    secType = NODE_STATIC
    return
}

func (router *Router) splitUri(uri string) []string {
    sec := strings.Split(uri[1:], "/") // avoid the first slash
    return sec
}


func (router *Router) makeHttpHandler(handleFuncName string) (h http.Handler) {
    // if starts with "FS:", that's FileServerHandler

    isFs := strings.HasPrefix(handleFuncName, "FS:")
    if isFs {
        h = http.FileServer(http.Dir(handleFuncName[3:]))
    } else {
        h = http.HandlerFunc(router.makeHandler(handleFuncName))
    }
    return
}

func (router *Router) makeHandler(tp string) http.HandlerFunc {
    fn := func(w http.ResponseWriter, r *http.Request) {
        fmt.Println(tp)

        w.Header().Set("content-type", "application/json")

        controllerValue := reflect.ValueOf(router.controllerImpl)
        method := controllerValue.MethodByName(tp)

        in := make([]reflect.Value, 2)
        in[0] = reflect.ValueOf(w)
        in[1] = reflect.ValueOf(r)

        method.Call(in)
    }
    if fn == nil {
        return func(w http.ResponseWriter, r *http.Request) {
            alog.Debug("404 The method NotFound")

            controllerValue := reflect.ValueOf(router.controllerImpl)
            method := controllerValue.MethodByName("NotFound")
            method.Call([]reflect.Value{})
        }
    }

    return fn
}
/////////////////////////////////////////////////////////////////////
// TODO replace this ugly shit with Router.map[uri]Route
var rootHandler http.Handler

// register new path -> function
//func (router *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
//    rt := &Route{uri: pattern, handler: http.HandlerFunc(handler)}
//    router.routes = append(router.routes, rt)
//    // TODO replace this ugly shit with Router.map[uri]Route
//    if pattern == "/" {
//        rootRoute = rt
//    }
//}
//
//func (router *Router) Handle(pattern string, handler http.Handler) {
//    rt := &Route{uri: pattern, handler: handler}
//    router.routes = append(router.routes, rt)
//    // TODO replace this ugly shit with Router.map[uri]Route
//    if pattern == "/" {
//        rootRoute = rt
//    }
//}

// ---------------- runtime --------------------
// this will run in an incoming http request
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Router ServeHTTP:" + r.URL.String())
    // 1, match regular path pattern to find the handler
    // 2, call the handler function

    err := r.ParseForm()
    if err != nil {
        // DO NOTHING?
    }

    currNode := router.node
    sections := router.splitUri(r.URL.Path)
    for _, section := range sections {
        nextNode := currNode.children[section]

        if nextNode != nil {
            currNode = nextNode
            continue
        }

        // not found static node, try to find in dynamic way
        for key, node := range currNode.children {
            if node.nodeType == NODE_DYNAMIC {
                r.Form.Add(key, section)
                currNode = node
                break
            }
        }

        // still not found, break and let rootHandler handle it
        if currNode == nil {
            break
        }
    }

    if currNode != nil {
        currNode.handler.ServeHTTP(w, r)
    } else {
        rootHandler.ServeHTTP(w, r)
    }

    //    for _, route := range router.routes {
    //        if r.URL.Path == route.uri {
    //            alog.Debug("found handler")
    //            route.handler.ServeHTTP(w, r)
    //            return
    //        }
    //    }
    //
    //    // not found any handler, run root handler "/"
    //    if rootRoute != nil {
    //        rootRoute.handler.ServeHTTP(w, r)
    //    }
}