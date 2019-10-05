package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Phantas0s/testomatic/internal/config"
	"github.com/gen2brain/beeep"
	"github.com/radovskyb/watcher"
)

const (
	// scopes to create the relative path
	current = "current"
	dir     = "dir"
	all     = "all"
)

var (
	conf        config.YamlConf
	file        = flag.String("config", ".testomatic.yml", "The config file")
	showWatched = flag.Bool("show", false, "Show files watched")
)

func Run() error {
	flag.Parse()
	w := watcher.New()

	data, err := ioutil.ReadFile(*file)
	if err != nil {
		return err
	}

	if err := conf.Parse(data); err != nil {
		return err
	}

	// file or directory writting (so two in total)
	w.SetMaxEvents(2)
	w.FilterOps(watcher.Write)

	w.IgnoreHiddenFiles(conf.Watch.IgnoreHidden)
	w.Ignore(conf.Watch.Ignore...)

	r := regexp.MustCompile(conf.Watch.Regex)
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	if err := w.AddRecursive(conf.Watch.Root); err != nil {
		return err
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	if *showWatched {
		for path, f := range w.WatchedFiles() {
			fmt.Printf("%s: %s\n", path, f.Name())
		}
	}

	go func() error {
		for {
			select {
			case event := <-w.Event:
				if !event.IsDir() {
					result, err := fireCmd(event)
					if err != nil {
						return err
					}

					fmt.Println(*result)

					if !conf.Notification.Disable {
						notify(conf, *result)
					}
				}
			case err := <-w.Error:
				return err
			case <-w.Closed:
				return nil
			}
		}
	}()

	if err != nil {
		return err
	}

	fmt.Print("Testomatic is watching the files... \n")
	if err := w.Start(time.Millisecond * 100); err != nil {
		return err
	}

	return nil
}

func fireCmd(event watcher.Event) (*string, error) {
	if conf.Command.IgnorePath {
		result, err := execCmd(conf.Command.Bin, conf.Command.Options)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	var path *string
	path = &event.Path
	if filepath.IsAbs(event.Path) && conf.Command.Abs != true {
		var err error
		path, err = CreateRelative(event.Path, conf.Watch.Root, conf.Command.Scope)
		if err != nil {
			return nil, err
		}
	}

	options := append(conf.Command.Options, *path)
	result, err := execCmd(conf.Command.Bin, options)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func execCmd(cmdPath string, args []string) (*string, error) {
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
		fmt.Println(err)
		fmt.Println(stderr.String())
	}

	result := out.String()
	return &result, nil
}

// CreateRelative path from absolute path
// A point "." means that the path begins by the current directory
// TODO refactor...
func CreateRelative(path, confPath, scope string) (*string, error) {
	newpath := make([]string, 2)
	// Get the current directory where testomatic runs
	if strings.Index(confPath, ".") != -1 {
		confPath = "."
		currentDir, err := getProjectDir()
		if err != nil {
			return nil, err
		}

		split := strings.SplitAfter(*currentDir, "/")
		*currentDir = split[len(split)-1]
		newpath = strings.SplitAfter(path, *currentDir)
	} else {
		newpath = strings.SplitAfter(path, confPath)
	}

	if scope == dir {
		finalPath := confPath + filepath.Dir(newpath[1])
		return &finalPath, nil
	}

	if scope == current {
		finalPath := confPath + newpath[1]
		return &finalPath, nil
	}

	if scope == all {
		return &confPath, nil
	}

	return nil, nil
}
func getProjectDir() (*string, error) {
	projectDir, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}

	return &projectDir, nil
}

func notify(conf config.YamlConf, result string) {
	mess := ""

	if conf.Notification.DisplayResult {
		mess = result
	}

	if match, _ := regexp.MatchString(conf.Notification.RegexFailure, result); match {
		beeep.Notify("Failure!", mess, conf.Notification.ImgFailure)
		return
	}

	if match, _ := regexp.MatchString(conf.Notification.RegexSuccess, result); match {
		if conf.Notification.Mute {
			beeep.Notify("Success!", mess, conf.Notification.ImgSuccess)
		} else {
			beeep.Alert("Success!", mess, conf.Notification.ImgSuccess)
		}
	}
}
