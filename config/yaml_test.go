package config

import (
	"io/ioutil"
	"log"
	"testing"
)

func assert(expected interface{}, got interface{}) {
	if expected != got {
		log.Fatalln("Expected ", expected, "got ", got)
	}
}

func TestParse(t *testing.T) {
	data, err := ioutil.ReadFile("./.testomatic.yml")
	if err != nil {
		log.Fatal(err)
	}

	var config YamlConf
	if err := config.Parse(data); err != nil {
		log.Fatal(err)
	}

	assert("src/Tests", config.Watch.Folder)
	assert("Test.php", config.Watch.Reg)
	assert("docker-compose", config.Command.Path)
	assert("exec", config.Command.Options[0])
	assert("-T", config.Command.Options[1])
	assert("php", config.Command.Options[2])
	assert("bin/phpunit", config.Command.Options[3])
	assert("/home/superUser/.autotest/images/success.png", config.Notification.ImgSuccess)
	assert("/home/superUser/.autotest/images/failure.png", config.Notification.ImgFailure)
	assert("OK", config.Notification.SuccessRegex)
}
