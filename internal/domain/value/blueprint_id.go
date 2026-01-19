package value

type BlueprintID string

const BlueprintIDLength = 8

func NewBlueprintID() BlueprintID {
	return BlueprintID(NewShortUUID(BlueprintIDLength))
}
