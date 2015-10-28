package ara
import "net/http"

const (
    NODE_STATIC = iota // normal node, ex: abc in /abc
    NODE_DYNAMIC       // dynamic node, ex: {id} in /abc/{id}
)

type Node struct {
    name string
    value string  // the value if node type is dynamic, ex: value of id in /abc/{id}
    nodeType uint
    handler http.Handler // the function to handle this request
    children map[string]*Node
}

func NewNode(name string, nodeType uint, h http.Handler) *Node{
    return &Node{name: name, nodeType: nodeType, handler: h, children: make(map[string]*Node)}
}

func (node *Node)String() string {
    str := "Node:[" +
    "name:" + node.name + "," +
    "value:" + node.value + "," +
    "children:" + string(len(node.children)) + "," +
    "type:" + string(node.nodeType) +
    "]"
    return str
}