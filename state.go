package fairytale

import (
	"net/url"
	"strings"

	"github.com/macabot/fairytale/internal/console"
	"github.com/macabot/hypp"
)

type AdminSettings struct {
	IFrameSize IFrameSize
	Rotation   Rotation
}

type State[S hypp.State] struct {
	hypp.EmptyState
	tree             Node[S]
	current          []int
	settings         AdminSettings
	assets           []*hypp.VNode
	selectedPanelTab int
}

func NewState[S hypp.State](tree Node[S]) *State[S] {
	return &State[S]{tree: tree}
}

func (s State[S]) Tree() Node[S]                       { return s.tree }
func (s State[S]) Current() []int                      { return s.current }
func (s *State[S]) SetCurrent(current []int)           { s.current = current }
func (s State[S]) Settings() AdminSettings             { return s.settings }
func (s *State[S]) SetSettings(settings AdminSettings) { s.settings = settings }
func (s State[S]) Assets() []*hypp.VNode               { return s.assets }
func (s *State[S]) SetAssets(assets []*hypp.VNode)     { s.assets = assets }

func (s State[S]) SelectedPanelTab() int        { return s.selectedPanelTab }
func (s *State[S]) SetSelectedPanelTab(tab int) { s.selectedPanelTab = tab }

func (s State[S]) GetTale(path []int) *Tale[S] {
	node := s.tree
	for _, i := range path {
		node = node.Children()[i]
	}
	return node.Tale()
}

func (s State[S]) CurrentTale() *Tale[S] {
	return s.GetTale(s.current)
}

func (s State[S]) TalePaths() [][]int {
	var talePaths [][]int
	var walk func(Node[S], []int)
	walk = func(node Node[S], path []int) {
		if node.Tale() != nil {
			talePaths = append(talePaths, path)
		}
		for i, child := range node.Children() {
			childPath := make([]int, len(path)+1)
			copy(childPath, path)
			childPath[len(childPath)-1] = i
			walk(child, childPath)
		}
	}
	walk(s.tree, nil)
	return talePaths
}

func (s State[S]) Clone() *State[S] {
	return &s
}

func (s State[S]) ToURL(forceCurrent []int) *url.URL {
	current := s.current
	if forceCurrent != nil {
		current = forceCurrent
	}
	slugs := make([]string, len(current))
	node := s.tree
	for i, pathI := range current {
		node = node.Children()[pathI]
		slugs[i] = node.Slug()
	}
	path := "/" + strings.Join(slugs, "/")

	return &url.URL{
		Fragment: path,
	}
}

func (s *State[S]) UpdateCurrentFromURL(u *url.URL) {
	path := u.Fragment
	if path == "" {
		path = "/"
	}
	if path == "/" {
		return
	}

	slugs := strings.Split(path, "/")
	node := s.tree
	// Skip first slug which is an empty string.
	current := make([]int, len(slugs)-1)
	found := false
	for i := 1; i < len(slugs); i++ {
		found = false
		children := node.Children()
		for pathI, child := range children {
			if child.Slug() == slugs[i] {
				current[i-1] = pathI
				node = child
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	if !found {
		console.Warn("Could not find tale for 'path' in URL fragment.")
	} else {
		s.current = current
		node := s.tree
		for _, pathI := range current {
			node = node.Children()[pathI]
			node.SetIsOpen(true)
		}
	}
}
