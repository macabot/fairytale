package component

import (
	"fmt"

	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func IFrame[S hypp.State](currentTale *fairytale.Tale[S], settings fairytale.AdminSettings) *hypp.VNode {
	size := settings.IFrameSize
	if settings.Orientation == fairytale.Landscape {
		size.Swap()
	}
	divProps := hypp.HProps{"class": "current-tale"}
	iFrameProps := hypp.HProps{
		"src":   "/",
		"title": taleToTitle(currentTale),
	}
	if size[0] != 0 && size[1] != 0 {
		divProps["style"] = map[string]string{
			"min-height": fmt.Sprintf("%dpx", size[1]),
		}
		iFrameProps["style"] = map[string]string{
			"width":  fmt.Sprintf("%dpx", size[0]),
			"height": fmt.Sprintf("%dpx", size[1]),
		}
	}
	return html.Div(
		divProps,
		html.Iframe(iFrameProps),
	)
}
