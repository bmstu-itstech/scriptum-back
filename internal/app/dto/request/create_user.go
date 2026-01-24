package request

type CreateUser struct {
	ActorID  string
	Email    string
	Password string
	Name     string
	Role     string
}
