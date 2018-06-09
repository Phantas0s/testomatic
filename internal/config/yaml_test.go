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
	data, err := ioutil.ReadFile("./testdata/testomatic.yml")
	if err != nil {
		log.Fatal(err)
	}

	var config YamlConf
	if err := config.Parse(data); err != nil {
		log.Fatal(err)
	}

	assert("src/Tests", config.Watch.Root)
	assert("Test.php", config.Watch.Regex)
	assert("vendor", config.Watch.Ignore[0])
	assert(true, config.Watch.IgnoreHidden)
	assert("docker-compose", config.Command.Bin)
	assert(false, config.Command.Abs)
	assert(true, config.Command.IgnorePath)
	assert("current", config.Command.Scope)
	assert("exec", config.Command.Options[0])
	assert("-T", config.Command.Options[1])
	assert("php", config.Command.Options[2])
	assert("bin/phpunit", config.Command.Options[3])
	assert("/home/superUser/.autotest/images/success.png", config.Notification.ImgSuccess)
	assert("/home/superUser/.autotest/images/failure.png", config.Notification.ImgFailure)
	assert(false, config.Notification.Disable)
	assert("ok", config.Notification.RegexSuccess)
	assert("fail", config.Notification.RegexFailure)
	assert(true, config.Notification.DisplayResult)
}
