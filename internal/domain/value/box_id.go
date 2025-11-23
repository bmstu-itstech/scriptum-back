package value

type BoxID string

const BoxIDLength = 8

func NewBoxID() BoxID {
	return BoxID(NewShortUUID(BoxIDLength))
}
