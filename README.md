## interface{} based skiplist

You can insert kinds of values into list, as long as you implements 
```golang
type Node interface {
	Compare(interface{}) int // compare node, -1 less 0 equal, 1 greater
}
```