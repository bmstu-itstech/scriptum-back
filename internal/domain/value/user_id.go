package value

const UserIDLength = 8

type UserID string

func NewUserID() UserID {
	return UserID(NewShortUUID(UserIDLength))
}
