package cmd

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

	"github.com/Phantas0s/testomatic/internal/config"
	"github.com/Phantas0s/watcher"
	"github.com/gen2brain/beeep"
)

const (
	current = "current"
	dir     = "dir"
	all     = "all"
)

var (
	conf        config.YamlConf
	file        = flag.String("config", ".testomatic.yml", "The config file")
	showWatched = flag.Bool("show", false, "Show files watched")
)

func Run() {
	flag.Parse()
	w := watcher.New()

	data, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}

	if err := conf.Parse(data); err != nil {
		log.Fatal(err)
	}

	// file or directory writting (so two in total)
	w.SetMaxEvents(2)
	w.FilterOps(watcher.Write)

	w.IgnoreHiddenFiles(conf.Watch.IgnoreHidden)
	w.FilterFiles(conf.Watch.Regex)
	if err := w.AddRecursive(conf.Watch.Root); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	if *showWatched {
		for path, f := range w.WatchedFiles() {
			fmt.Printf("%s: %s\n", path, f.Name())
		}
	}

	go func() {
		for {
			select {
			case event := <-w.Event:
				if !event.IsDir() {
					result := fireCmd(event)
					fmt.Println(result)

					if !conf.Notification.Disable {
						notify(conf, result)
					}
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

func fireCmd(event watcher.Event) string {
	if conf.Command.IgnorePath {
		return execCmd(conf.Command.Bin, conf.Command.Options)
	}

	path := event.Path

	if filepath.IsAbs(event.Path) && conf.Command.Abs != true {
		path = CreateRelative(event.Path, conf.Watch.Root, conf.Command.Scope)
	}

	options := append(conf.Command.Options, path)
	return execCmd(conf.Command.Bin, options)
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

// CreateRelative path from absolute path
// A point "." means that the path begins by the current directory
func CreateRelative(path, confPath, scope string) string {
	newpath := make([]string, 2)
	// Get the current directory where testomatic runs
	if strings.Index(confPath, ".") != -1 {
		confPath = "."
		currentDir, _ := filepath.Abs(".")
		split := strings.SplitAfter(currentDir, "/")
		currentDir = split[len(split)-1]
		newpath = strings.SplitAfter(path, currentDir)
	} else {
		newpath = strings.SplitAfter(path, confPath)
	}

	if scope == dir {
		return confPath + filepath.Dir(newpath[1])
	}

	if scope == current {
		return confPath + newpath[1]
	}

	if scope == all {
		return confPath
	}

	return ""
}

func notify(conf config.YamlConf, result string) {
	mess := ""

	if conf.Notification.DisplayResult {
		mess = result
	}

	if match, _ := regexp.MatchString(conf.Notification.RegexFailure, result); match {
		beeep.Notify("Failure!", mess, conf.Notification.ImgFailure)
	}

	if match, _ := regexp.MatchString(conf.Notification.RegexSuccess, result); match {
		beeep.Alert("Success!", mess, conf.Notification.ImgSuccess)
	}
}
