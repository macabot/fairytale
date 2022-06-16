package fairy

import (
	"net/url"
	"strconv"
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

func (s *state) updateFromQuery(query url.Values) {
	if query.Has("path") {
		slugs := strings.Split(query.Get("path"), "/")
		node := s.Tree
		// Skip first segment which is an empty string.
		path := make([]int, len(slugs)-1)
		found := false
		for i := 1; i < len(slugs); i++ {
			found = false
			for pathI, child := range node.children() {
				if child.slug() == slugs[i] {
					path[i-1] = pathI
					node = child
					found = true
					break
				}
			}
			if !found {
				break
			}
		}
		if !found || !s.hasTale(path) {
			consoleWarn("Could not find tale for query param 'path'.")
		} else {
			s.Current = path
			node := s.Tree
			for _, i := range path {
				node = node.children()[i]
				node.setIsOpen(true)
			}
		}
	}
	if query.Has("iFrameSize") {
		if size, err := iFrameSizeFromSlug(query.Get("iFrameSize")); err != nil {
			consoleWarn("Could not parse query param 'iFrameSize'.")
		} else {
			s.Settings.iFrameSize = size
		}
	}
	if query.Has("landscape") {
		if landscape, err := strconv.ParseBool(query.Get("landscape")); err != nil {
			consoleWarn("Could not parse query param 'landscape'.")
		} else {
			s.Settings.landscape = landscape
		}
	}
}

func (s state) toQuery() url.Values {
	query := url.Values{}
	if s.Current != nil {
		slugs := make([]string, len(s.Current))
		node := s.Tree
		for i, pathI := range s.Current {
			node = node.children()[pathI]
			slugs[i] = node.slug()
		}
		query.Set("path", "/"+strings.Join(slugs, "/"))
	}
	query.Set("iFrameSize", s.Settings.iFrameSize.Slug())
	query.Set("landscape", strconv.FormatBool(s.Settings.landscape))
	return query
}
