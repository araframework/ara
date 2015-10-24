package ara
import "net/http"

type Controller struct {}

func (c *Controller) NotFound(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("My controller not found......."))
}

func (c *Controller) getAttribute(key string) string {
    return "" //TODO
}