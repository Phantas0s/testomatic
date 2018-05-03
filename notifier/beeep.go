package notifier

import "github.com/gen2brain/beeep"

type Beeep struct{}

func (b *Beeep) Info(title string, message string, img string) {
	beeep.Notify(title, message, img)
}

func (b *Beeep) Alert(title string, message string, img string) {
	beeep.Alert(title, message, img)
}
