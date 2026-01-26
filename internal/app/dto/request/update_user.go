package request

type UpdateUser struct {
	ActorID  string
	UserID   string
	Email    *string
	Password *string
	Name     *string
	Role     *string
}
