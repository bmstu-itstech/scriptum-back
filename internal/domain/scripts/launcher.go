package scripts

type Launcher interface {
	launch(Job) Result
}
