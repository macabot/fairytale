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
	tree             Node
	current          []int
	settings         AdminSettings
	assets           []*hypp.VNode
	taleEvents       []TaleEvent
	selectedPanelTab int
}

func NewState(tree Node) *State {
	return &State{tree: tree}
}

func (s State) Tree() Node                          { return s.tree }
func (s State) Current() []int                      { return s.current }
func (s *State) SetCurrent(current []int)           { s.current = current }
func (s State) Settings() AdminSettings             { return s.settings }
func (s *State) SetSettings(settings AdminSettings) { s.settings = settings }
func (s State) Assets() []*hypp.VNode               { return s.assets }
func (s *State) SetAssets(assets []*hypp.VNode)     { s.assets = assets }
func (s State) TaleEvents() []TaleEvent             { return s.taleEvents }
func (s *State) SetTaleEvents(events []TaleEvent)   { s.taleEvents = events }
func (s State) SelectedPanelTab() int               { return s.selectedPanelTab }
func (s *State) SetSelectedPanelTab(tab int)        { s.selectedPanelTab = tab }

func (s State) GetTale(path []int) *Tale {
	node := s.tree
	for _, i := range path {
		node = node.Children()[i]
	}
	return node.Tale()
}

func (s State) CurrentTale() *Tale {
	return s.GetTale(s.current)
}

func (s State) Clone() *State {
	return &s
}

func (s State) ToURL(forceCurrent []int) *url.URL {
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

func (s *State) UpdateCurrentFromURL(u *url.URL) {
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
