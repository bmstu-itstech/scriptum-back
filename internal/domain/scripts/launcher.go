package scripts

type Launcher interface {
	Launch(Job) Result
}
