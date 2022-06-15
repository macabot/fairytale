package fairy

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/macabot/hypp"
)

type state struct {
	hypp.EmptyState
	Tree             Node
	Current          []int
	Settings         AdminSettings
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
		var path []int
		if err := json.Unmarshal([]byte(query.Get("path")), &path); err != nil {
			consoleWarn("Could not parse query param 'path'.")
		} else if !s.hasTale(path) {
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
		b, err := json.Marshal(s.Current)
		if err != nil {
			panic("Could not JSON marshal State.Current")
		}
		query.Set("path", string(b))
	}
	query.Set("iFrameSize", s.Settings.iFrameSize.Slug())
	query.Set("landscape", strconv.FormatBool(s.Settings.landscape))
	return query
}
