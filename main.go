package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Phantas0s/watcher"
	"github.com/kylelemons/go-gypsy/yaml"
)

func main() {
	w := watcher.New()

	file := flag.String("config", ".testomatic.yml", "The config file")
	config := yaml.ConfigFile(*file)

	// 2 since the event can be writting file or writting directory ... to fix
	w.SetMaxEvents(2)
	w.FilterOps(watcher.Write)
	filePath := extractScalar(config, "watcher.folder")

	go func() {
		for {
			select {
			case event := <-w.Event:
				if !event.IsDir() {
					result := fireCmd(config, event, filePath)
					fmt.Println(result)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	root := toYamlMap(config.Root)
	ext := extractExt(root)
	if err := w.AddSpecificFiles(filePath, ext); err != nil {
		log.Fatalln(err)
	}

	fmt.Print("Testomatic begins: \n")
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func fireCmd(config *yaml.File, event watcher.Event, filepath string) string {
	cmdPath, err := config.Get("watcher.command_path")
	if err != nil {
		fmt.Println(err)
	}

	options := extractOpt(toYamlMap(config.Root))

	path := event.Path
	if isAbsolute(event.Path) {
		path = getRelative(event.Path, filepath)
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

func toYamlMap(node yaml.Node) yaml.Map {
	m, ok := node.(yaml.Map)
	if !ok {
		log.Fatalf("%v is not of type map", node)
	}
	return m
}

func toYamlList(node yaml.Node) yaml.List {
	m, ok := node.(yaml.List)
	if !ok {
		log.Fatalf("%v is not of type list", node)
	}
	return m
}

func toYamlScalar(node yaml.Node) yaml.Scalar {
	m, ok := node.(yaml.Scalar)
	if !ok {
		log.Fatalf("%v is not of type scalar", node)
	}
	return m
}

func extractExt(root yaml.Map) []string {
	list := toYamlList(toYamlMap(root["watcher"])["ext"])
	result := make([]string, list.Len())

	for k, v := range list {
		result[k] = "." + toYamlScalar(v).String()
	}

	return result
}

func extractOpt(root yaml.Map) []string {
	list := toYamlList(toYamlMap(root["watcher"])["command_options"])
	result := make([]string, list.Len())

	for k, v := range list {
		str := toYamlScalar(v).String()
		str = strings.Replace(str, "'", "", -1)
		result[k] = str
	}

	return result
}

func extractScalar(config *yaml.File, name string) string {
	entry, err := config.Get(name)
	if err != nil {
		fmt.Println(err)
	}

	return entry
}

func isAbsolute(path string) bool {
	if path[0] == '/' {
		return true
	}

	return false
}

func getRelative(path string, filepath string) string {
	newpath := strings.SplitAfter(path, filepath)
	path = filepath + newpath[1]

	return path
}
