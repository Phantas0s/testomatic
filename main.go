package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Phantas0s/testomatic/config"
	"github.com/Phantas0s/testomatic/notifier"
	"github.com/Phantas0s/watcher"
)

var conf config.YamlConf
var file = flag.String("config", ".testomatic.yml", "The config file")

// init is called prior to main.
func init() {
	// Change the device for logging to stdout.
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()
	w := watcher.New()

	data, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}

	if err := conf.Parse(data); err != nil {
		log.Fatal(err)
	}

	// 2 since the event can be writting file or writting directory ... to fix
	w.SetMaxEvents(2)
	w.FilterOps(watcher.Write)
	filePath := conf.Watch.Folder

	w.AllowByRegex(conf.Watch.Reg)
	if err := w.AddRecursive(filePath); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	go func() {
		for {
			select {
			case event := <-w.Event:
				if !event.IsDir() {
					result := fireCmd(conf.Command.Path, conf.Command.Options, event, filePath)
					fmt.Println(result)

					notify(conf, result)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	fmt.Print("Testomatic is watching the files... \n")
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func fireCmd(cmdPath string, options []string, event watcher.Event, confPath string) string {
	path := event.Path

	if filepath.IsAbs(event.Path) && conf.Watch.Abs != true {
		path = CreateRelative(event.Path, confPath, conf.Test)
	}

	options = append(options, path)
	result := execCmd(cmdPath, options)

	return result
}

func execCmd(cmdPath string, args []string) string {
	var out bytes.Buffer
	var stderr bytes.Buffer

	// fmt.Println(cmdPath, args)
	cmd := exec.Command(cmdPath, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &out

	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}

	return out.String()
}

func CreateRelative(
	path string,
	confPath string,
	testRealm string,
) string {
	newpath := make([]string, 2)
	if strings.Index(confPath, ".") != -1 {
		confPath = "."
		currentDir, _ := filepath.Abs(".")
		split := strings.SplitAfter(currentDir, "/")
		currentDir = split[len(split)-1]
		newpath = strings.SplitAfter(path, currentDir)
	} else {
		newpath = strings.SplitAfter(path, confPath)
	}

	if testRealm == "dir" {
		return confPath + filepath.Dir(newpath[1])
	}

	if testRealm == "current" {
		return confPath + newpath[1]
	}

	if testRealm == "all" {
		return confPath
	}

	return ""
}

func notify(conf config.YamlConf, result string) {
	notifier := new(notifier.Beeep)
	mess := ""

	if conf.Notification.DisplayResult {
		mess = result
	}

	if match, _ := regexp.MatchString(conf.Notification.SuccessRegex, result); match {
		notifier.Info("Success!", mess, conf.Notification.ImgSuccess)
	} else {
		notifier.Alert("Failure!", mess, conf.Notification.ImgFailure)
	}
}
