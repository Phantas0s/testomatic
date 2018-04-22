package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/Phantas0s/testomatic/yamlext"
	"github.com/Phantas0s/watcher"
	"github.com/gen2brain/beeep"
	"github.com/kylelemons/go-gypsy/yaml"
)

func main() {
	w := watcher.New()

	file := flag.String("config", ".testomatic.yml", "The config file")
	config := yaml.ConfigFile(*file)

	// 2 since the event can be writting file or writting directory ... to fix
	w.SetMaxEvents(2)
	w.FilterOps(watcher.Write)
	filePath := ExtractScalar(config, "testomatic.folder")
	root := yamlext.ToMap(config.Root)

	go func() {
		for {
			select {
			case event := <-w.Event:
				if !event.IsDir() {
					result := fireCmd(config, event, filePath)
					fmt.Println(result)
					Notify(config, result)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	ext := ExtractExt(root)
	if err := w.AddSpecificFiles(filePath, ext); err != nil {
		log.Fatalln(err)
	}

	fmt.Print("Testomatic is watching the files... \n")
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func fireCmd(config *yaml.File, event watcher.Event, filepath string) string {
	cmdPath, err := config.Get("testomatic.command_path")
	if err != nil {
		fmt.Println(err)
	}

	options := ExtractOpt(yamlext.ToMap(config.Root))

	path := event.Path
	if IsAbsolute(event.Path) {
		path = CreateRelative(event.Path, filepath)
	}

	options = append(options, path)
	result := execCmd(cmdPath, options)

	return result
}

func execCmd(cmdPath string, args []string) string {
	var out bytes.Buffer
	var stderr bytes.Buffer
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

func ExtractScalar(config *yaml.File, name string) string {
	entry, err := config.Get(name)
	if err != nil {
		fmt.Println(err)
	}

	return entry
}

func ExtractOpt(root yaml.Map) []string {
	list := yamlext.ToList(yamlext.ToMap(root["testomatic"])["command_options"])
	result := make([]string, list.Len())

	for k, v := range list {
		// Delete quotes
		result[k] = strings.Replace(yamlext.ToScalar(v).String(), "'", "", -1)
		result[k] = strings.Replace(result[k], "\"", "", -1)
	}

	return result
}

func ExtractExt(root yaml.Map) []string {
	list := yamlext.ToList(yamlext.ToMap(root["testomatic"])["ext"])
	result := make([]string, list.Len())

	for k, v := range list {
		result[k] = "." + yamlext.ToScalar(v).String()
	}

	return result
}

func IsAbsolute(path string) bool {
	if path[0] == '/' {
		return true
	}

	return false
}

func CreateRelative(path string, filepath string) string {
	newpath := strings.SplitAfter(path, filepath)
	path = filepath + newpath[1]

	return path
}

// How can I test that??
func Notify(config *yaml.File, result string) {
	if match, _ := regexp.MatchString(ExtractScalar(config, "testomatic.notification.notify_success"), result); match {
		beeep.Notify("Success!", result, ExtractScalar(config, "testomatic.notification.notify_img_success"))
	} else {
		beeep.Alert("Failure!", result, ExtractScalar(config, "testomatic.notification.notify_img_failure"))
	}
}
