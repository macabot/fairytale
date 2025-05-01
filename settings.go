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

func MustIFrameSizeFromString(s string) IFrameSize {
	size, err := IFrameSizeFromString(s)
	if err != nil {
		panic(err)
	}
	return size
}

func IFrameSizeFromString(s string) (IFrameSize, error) {
	for _, size := range IFrameSizes {
		if size.String() == s {
			return size, nil
		}
	}
	return [2]int{}, fmt.Errorf("fairytale: cannot create iFrameSize from string '%s'", s)
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
		panic(fmt.Errorf("fairytale: unknown IFrameSize: %d, %d", i[0], i[1]))
	}
}

func (i IFrameSize) Slug() string {
	return slug.Make(i.String())
}

type Orientation int

const (
	Portrait Orientation = iota
	Landscape
)

var Orientations = [...]Orientation{
	Portrait,
	Landscape,
}

func MustOrientationFromString(s string) Orientation {
	orientation, err := OrientationFromString(s)
	if err != nil {
		panic(err)
	}
	return orientation
}

func OrientationFromString(s string) (Orientation, error) {
	for _, orientation := range Orientations {
		if orientation.String() == s {
			return orientation, nil
		}
	}
	return -1, fmt.Errorf("fairytale: cannot create orientation from string '%s'", s)
}

func (r Orientation) String() string {
	switch r {
	case Portrait:
		return "Portrait"
	case Landscape:
		return "Landscape"
	default:
		panic(fmt.Errorf("fairytale: unknown orientation: %d", r))
	}
}

func (r Orientation) Slug() string {
	return slug.Make(r.String())
}

func orientationFromSlug(s string) (Orientation, error) {
	for _, orientation := range Orientations {
		if orientation.Slug() == s {
			return orientation, nil
		}
	}
	return -1, fmt.Errorf("fairytale: cannot create orientation from slug '%s'", s)
}
