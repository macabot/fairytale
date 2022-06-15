package fairy

import (
	"fmt"
	"strings"
)

type iFrameSize [2]int

var (
	SizeDesktop        = iFrameSize{0, 0}
	Size_iPhone_11_Pro = iFrameSize{375, 812}
)

var IFrameSizes = [...]iFrameSize{
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
	return strings.ReplaceAll(i.String(), " ", "-")
}

func iFrameSizeFromSlug(s string) (iFrameSize, error) {
	return iFrameSizeFromString(strings.ReplaceAll(s, "-", " "))
}

func mustIFrameSizeFromString(s string) iFrameSize {
	size, err := iFrameSizeFromString(s)
	if err != nil {
		panic(err)
	}
	return size
}

func iFrameSizeFromString(s string) (iFrameSize, error) {
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
	iFrameSize iFrameSize
	landscape  bool
}
