package value

import "fmt"

type ImageTag string

func NewImageTag(prefix string, id BoxID) ImageTag {
	return ImageTag(fmt.Sprintf("%s:%s", prefix, id))
}
