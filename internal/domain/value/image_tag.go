package value

import "fmt"

type ImageTag string

func NewImageTag(prefix string, id BlueprintID) ImageTag {
	return ImageTag(fmt.Sprintf("%s:%s", prefix, id))
}
