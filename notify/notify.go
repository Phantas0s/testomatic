package notify

type Notifier interface {
	Info(title string, message string, img string) error
	Alert(title string, message string, img string) error
}
