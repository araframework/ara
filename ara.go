package ara

import (
    "fmt"
    "os"

    "github.com/araframework/aralog"
    "net/http"
)

var alog *aralog.Logger

func init() {
    fmt.Println("ara:init()")
    // TODO log filename/path should be configured in log configuration file: aralog.toml
    logger, err := aralog.NewFileLogger("applog.log", aralog.Llongfile | aralog.Ltime)
    if err != nil {
        fmt.Println("[ERROR]new logger error: ", err)
        fmt.Println("I will exit now, sorry")
        os.Exit(1)
    }

    alog = logger
}

func Start(router *Router) {
    fmt.Println("ara:Start()")

    err := http.ListenAndServe(":8600", router)
    if err != nil {
        fmt.Println("ListenAndServe:", err)
    }

}
