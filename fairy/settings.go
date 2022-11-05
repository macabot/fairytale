package fairy

import (
	"fmt"

	"github.com/gosimple/slug"
)

type iFrameSize [2]int

var (
	SizeDesktop        = iFrameSize{0, 0}
	Size_iPhone_11_Pro = iFrameSize{375, 812}
)

var iFrameSizes = [...]iFrameSize{
	SizeDesktop,
	Size_iPhone_11_Pro,
}

func (i *iFrameSize) Swap() {
	i[0], i[1] = i[1], i[0]
}

func (i iFrameSize) Equal(other iFrameSize) bool {
	return i[0] == other[0] && i[1] == other[1]
}

func (i iFrameSize) String() string {
	switch i {
	case SizeDesktop:
		return "Desktop"
	case Size_iPhone_11_Pro:
		return "iPhone 11 Pro"
	default:
		panic(fmt.Errorf("unknown IFrameSize: %d, %d", i[0], i[1]))
	}
}

func (i iFrameSize) Slug() string {
	return slug.Make(i.String())
}

func mustIFrameSizeFromString(s string) iFrameSize {
	size, err := iFrameSizeFromString(s)
	if err != nil {
		panic(err)
	}
	return size
}

func iFrameSizeFromString(s string) (iFrameSize, error) {
	for _, size := range iFrameSizes {
		if size.String() == s {
			return size, nil
		}
	}
	return [2]int{}, fmt.Errorf("cannot create iFrameSize from string '%s'", s)
}

type rotation int

const (
	Portrait rotation = iota
	Landscape
)

var rotations = [...]rotation{
	Portrait,
	Landscape,
}

func (r rotation) String() string {
	switch r {
	case Portrait:
		return "Portrait"
	case Landscape:
		return "Landscape"
	default:
		panic(fmt.Errorf("unknown rotation: %d", r))
	}
}

func rotationFromString(s string) (rotation, error) {
	for _, rotation := range rotations {
		if rotation.String() == s {
			return rotation, nil
		}
	}
	return -1, fmt.Errorf("cannot create rotation from string '%s'", s)
}

func mustRotationFromString(s string) rotation {
	rotation, err := rotationFromString(s)
	if err != nil {
		panic(err)
	}
	return rotation
}

func (r rotation) Slug() string {
	return slug.Make(r.String())
}

func rotationFromSlug(s string) (rotation, error) {
	for _, rotation := range rotations {
		if rotation.Slug() == s {
			return rotation, nil
		}
	}
	return -1, fmt.Errorf("cannot create rotation from slug '%s'", s)
}

type adminSettings struct {
	iFrameSize iFrameSize
	rotation   rotation
}
