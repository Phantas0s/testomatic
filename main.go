package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Phantas0s/watcher"
	"github.com/kylelemons/go-gypsy/yaml"
	"log"
	"os/exec"
	"time"
)

func main() {
	w := watcher.New()

	file := flag.String("config", "config.yml", "The config file")
	config := yaml.ConfigFile(*file)

	// fireCmd(config)

	root := toYamlMap(config.Root)

	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	filePath, err := config.Get("watcher.folder")
	ext := extractExt(root)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v", ext)
	if err := w.AddSpecificFiles(filePath, ext); err != nil {
		log.Fatalln(err)
	}

	fmt.Print("File watched: \n")
	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	// Inject path of the file via command line argument

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func fireCmd(config *yaml.File) {
	cmdPath, err := config.Get("watcher.command_path")
	if err != nil {
		fmt.Println(err)
	}
	execCmd(cmdPath)
}

func execCmd(cmdPath string) {
	cmd := exec.Command(cmdPath)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%q\n", out.String())
}

func toYamlMap(node yaml.Node) yaml.Map {
	m, ok := node.(yaml.Map)
	if !ok {
		panic(fmt.Sprintf("%v is not of type map", node))
	}
	return m
}

func toYamlList(node yaml.Node) yaml.List {
	m, ok := node.(yaml.List)
	if !ok {
		panic(fmt.Sprintf("%v is not of type map", node))
	}
	return m
}

func extractExt(root yaml.Map) []string {
	list := toYamlList(toYamlMap(root["watcher"])["ext"])
	result := make([]string, list.Len())

	for k, v := range list {
		result[k] = "." + v.(yaml.Scalar).String()
	}

	return result
}
