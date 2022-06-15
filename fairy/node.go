package fairy

type Node interface {
	name() string
	children() []Node
	tale() *Tale
	isOpen() bool
	setIsOpen(bool)
}

var _ Node = &Branch{}

type Branch struct {
	myName     string
	myChildren []Node
	myTale     *Tale
	myIsOpen   bool
}

func (b Branch) name() string           { return b.myName }
func (b Branch) children() []Node       { return b.myChildren }
func (b Branch) tale() *Tale            { return b.myTale }
func (b Branch) isOpen() bool           { return b.myIsOpen }
func (b *Branch) setIsOpen(isOpen bool) { b.myIsOpen = isOpen }

func NewTree(children ...Node) Node {
	return &Branch{myChildren: children, myIsOpen: true}
}

func NewBranch(name string, children ...Node) Node {
	return &Branch{myName: name, myChildren: children}
}
