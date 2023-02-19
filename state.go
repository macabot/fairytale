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

type TaleEvent struct {
	Key   string
	Event any
}

type State struct {
	hypp.EmptyState
	Tree             Node
	Current          []int
	Settings         AdminSettings
	Assets           []*hypp.VNode
	TaleEvents       []TaleEvent
	SelectedPanelTab int
}

func (s State) GetTale(path []int) *Tale {
	node := s.Tree
	for _, i := range path {
		node = node.Children()[i]
	}
	return node.Tale()
}

func (s State) CurrentTale() *Tale {
	return s.GetTale(s.Current)
}

func (s State) Clone() *State {
	return &s
}

func (s State) ToURL(forceCurrent []int) *url.URL {
	current := s.Current
	if forceCurrent != nil {
		current = forceCurrent
	}
	slugs := make([]string, len(current))
	node := s.Tree
	for i, pathI := range current {
		node = node.Children()[pathI]
		slugs[i] = node.Slug()
	}
	path := "/" + strings.Join(slugs, "/")

	return &url.URL{
		Fragment: path,
	}
}

func (s *State) UpdateCurrentFromURL(u *url.URL) {
	path := u.Fragment
	if path == "" {
		path = "/"
	}
	if path == "/" {
		return
	}

	slugs := strings.Split(path, "/")
	node := s.Tree
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
		s.Current = current
		node := s.Tree
		for _, pathI := range current {
			node = node.Children()[pathI]
			node.SetIsOpen(true)
		}
	}
}
