package notify

import (
	"log"
	"testing"
)

type fake struct {
	title   string
	message string
	img     string
}

func (f *fake) Info(title string, message string, img string) {
	f.title = title
	f.message = message
	f.img = img
}

func (f *fake) Alert(title string, message string, img string) {
	f.title = title
	f.message = message
	f.img = img
}

func assert(expected interface{}, got interface{}) {
	if expected != got {
		log.Fatalln("Expected ", expected, "got ", got)
	}
}

func TestInfo(t *testing.T) {
	notifier := new(fake)
	notifier.Info("Title", "Message", "Image")

	assert("Title", notifier.title)
	assert("Message", notifier.message)
	assert("Image", notifier.img)
}

func TestAlert(t *testing.T) {
	notifier := new(fake)
	notifier.Alert("Title", "Message", "Image")

	assert("Title", notifier.title)
	assert("Message", notifier.message)
	assert("Image", notifier.img)
}
