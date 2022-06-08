package fairy

import (
	"fmt"
	"strings"
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
	return strings.ReplaceAll(i.String(), " ", "-")
}

func iFrameSizeFromSlug(s string) (IFrameSize, error) {
	return iFrameSizeFromString(strings.ReplaceAll(s, "-", " "))
}

func mustIFrameSizeFromString(s string) IFrameSize {
	size, err := iFrameSizeFromString(s)
	if err != nil {
		panic(err)
	}
	return size
}

func iFrameSizeFromString(s string) (IFrameSize, error) {
	switch s {
	case "Desktop":
		return SizeDesktop, nil
	case "iPhone 11 Pro":
		return Size_iPhone_11_Pro, nil
	default:
		return [2]int{}, fmt.Errorf("cannot convert '%s' to IFrameSize", s)
	}
}

type AdminSettings struct {
	iFrameSize IFrameSize
	landscape  bool
}
