package fairytale

import (
	"fmt"

	"github.com/gosimple/slug"
)

type IFrameSize [2]int

var (
	SizeDesktop        = IFrameSize{0, 0}
	Size_iPhone_11_Pro = IFrameSize{375, 812}
)

var IFrameSizes = [...]IFrameSize{
	SizeDesktop,
	Size_iPhone_11_Pro,
}

func (i *IFrameSize) Swap() {
	i[0], i[1] = i[1], i[0]
}

func (i IFrameSize) Equal(other IFrameSize) bool {
	return i[0] == other[0] && i[1] == other[1]
}

func (i IFrameSize) String() string {
	switch i {
	case SizeDesktop:
		return "Desktop"
	case Size_iPhone_11_Pro:
		return "iPhone 11 Pro"
	default:
		panic(fmt.Errorf("unknown IFrameSize: %d, %d", i[0], i[1]))
	}
}

func (i IFrameSize) Slug() string {
	return slug.Make(i.String())
}

type Rotation int

const (
	Portrait Rotation = iota
	Landscape
)

var Rotations = [...]Rotation{
	Portrait,
	Landscape,
}

func (r Rotation) String() string {
	switch r {
	case Portrait:
		return "Portrait"
	case Landscape:
		return "Landscape"
	default:
		panic(fmt.Errorf("unknown rotation: %d", r))
	}
}

func (r Rotation) Slug() string {
	return slug.Make(r.String())
}

func rotationFromSlug(s string) (Rotation, error) {
	for _, rotation := range Rotations {
		if rotation.Slug() == s {
			return rotation, nil
		}
	}
	return -1, fmt.Errorf("cannot create rotation from slug '%s'", s)
}
