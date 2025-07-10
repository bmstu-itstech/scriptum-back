package scripts

type Notifier interface {
	Notify(Result) error
}
