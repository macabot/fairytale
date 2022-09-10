package fairy

import (
	"net/url"
	"strings"

	"github.com/macabot/hypp"
)

type state struct {
	hypp.EmptyState
	Tree             Node
	Current          []int
	Settings         adminSettings
	Assets           []*hypp.VNode
	TaleEvents       []taleEvent
	SelectedPanelTab int
}

func (s state) getTale(path []int) *Tale {
	node := s.Tree
	for _, i := range path {
		node = node.children()[i]
	}
	return node.tale()
}

func (s state) hasTale(path []int) bool {
	node := s.Tree
	for _, i := range path {
		children := node.children()
		if i < 0 || i >= len(children) {
			return false
		}
		node = children[i]
	}
	return true
}

func (s state) currentTale() *Tale {
	return s.getTale(s.Current)
}

func (s state) clone() *state {
	return &s
}

func (s state) toURL(forceCurrent []int) *url.URL {
	current := s.Current
	if forceCurrent != nil {
		current = forceCurrent
	}
	slugs := make([]string, len(current))
	node := s.Tree
	for i, pathI := range current {
		node = node.children()[pathI]
		slugs[i] = node.slug()
	}
	path := "/" + strings.Join(slugs, "/")

	return &url.URL{
		Fragment: path,
	}
}

func (s *state) updateCurrentFromURL(u *url.URL) {
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
		children := node.children()
		for pathI, child := range children {
			if child.slug() == slugs[i] {
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
		consoleWarn("Could not find tale for 'path' in URL fragment.")
	} else {
		s.Current = current
		node := s.Tree
		for _, pathI := range current {
			node = node.children()[pathI]
			node.setIsOpen(true)
		}
	}
}
