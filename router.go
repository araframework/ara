package ara

import (
    "fmt"
    "bufio"
    "os"
    "path"
    "strings"
    "errors"
)

const (
    GET = "GET"
    POST = "POST"
    PUT = "PUT"
    DELETE = "DELETE"
)

type route struct {
    method  string // http method
    uri     string // the request uri
    handler string // the function to handle this request
}

var routeMap map[string]route

func initRouter() {
    alog.Debug("router:init()")

    routeMap = make(map[string]route)

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

        rt, err := buildRoute(line)
        if err != nil {
            continue
        }

        routeMap[rt.method + rt.uri] = rt
    }
}

// validate the route configured in router file
func buildRoute(line string) (rt route, err error) {
    routeItems := strings.SplitN(line, " ", 3)
    for _, item := range routeItems {
        if strings.TrimSpace(item) == "" {
            alog.Debug("empty item for line" + line)
            alog.Debug(item)
            err = errors.New("invalid route: " + line)
            return
        }
    }

    rt = route{strings.TrimSpace(routeItems[0]), strings.TrimSpace(routeItems[1]), strings.TrimSpace(routeItems[2])}
    return
}