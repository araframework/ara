package ara

import (
    "fmt"
    "bufio"
    "os"
    "path"
)

func initRouter() {
    alog.Debug("router:init()")

    // the router conf file named: router
    confPath := path.Join("conf", "router")
    file, err := os.Open(confPath)
    if err != nil {
        alog.Debug(err)
        os.Exit(1)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fmt.Println(scanner.Text()) // Println will add back the final '\n'
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }
}