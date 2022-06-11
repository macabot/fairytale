package fairy

type Node interface {
	Name() string
	Children() []Node
	Tale() *Tale
	IsOpen() bool
	SetIsOpen(bool)
	Path() []int
	Parent() Node
	SetParent(Node)
}

var _ Node = &Branch{}

type Branch struct {
	name     string
	children []Node
	tale     *Tale
	isOpen   bool
	parent   Node
}

func (b Branch) Name() string           { return b.name }
func (b Branch) Children() []Node       { return b.children }
func (b Branch) Tale() *Tale            { return b.tale }
func (b Branch) IsOpen() bool           { return b.isOpen }
func (b *Branch) SetIsOpen(isOpen bool) { b.isOpen = isOpen }
func (b *Branch) Path() []int           { return getNodePath(b) }
func (b Branch) Parent() Node           { return b.parent }
func (b *Branch) SetParent(parent Node) { b.parent = parent }

func NewTree(children ...Node) Node {
	tree := &Branch{children: children, isOpen: true}
	for _, child := range children {
		child.SetParent(tree)
	}
	return tree
}

func NewBranch(name string, children ...Node) Node {
	branch := &Branch{name: name, children: children}
	for _, child := range children {
		child.SetParent(branch)
	}
	return branch
}

func getNodePath(node Node) []int {
	nodes := []Node{node}
	for node.Parent() != nil {
		nodes = append(nodes, node.Parent())
		node = node.Parent()
	}
	path := make([]int, len(nodes)-1)
	i := 0
	for j := len(nodes) - 1; j > 0; j-- {
		parent := nodes[j]
		child := nodes[j-1]
		for pathIndex, parentChild := range parent.Children() {
			if parentChild == child {
				path[i] = pathIndex
				break
			}
		}
		i++
	}
	return path
}
