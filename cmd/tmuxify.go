package main

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
)

type Window struct {
	Root  string
	Nvim  bool
	Index int
}

type Session struct {
	Root string
	Name string
	Win  []Window
}

const CONFIG_DIR_SUFFIX = ".config/tmuxify"

const DATA_DIR_SUFFIX = ".local/share/tmuxify"

type Data struct {
	ConfigDir string
	DataDir   string
}

var DATA = &Data{}

func tmuxify(entry os.DirEntry) {
	entryFilepath := filepath.Join(DATA.DataDir, entry.Name())
	data, err := os.ReadFile(filepath.Join(DATA.ConfigDir, entry.Name()))
	if err != nil {
		panic(err)
	}

	var session Session
	_, err = toml.Decode(
		string(data),
		&session,
	)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(session.Win); i++ {
		session.Win[i].Index = i + 1
	}

	tmpl, err := template.New("tmux.tmpl").Funcs(template.FuncMap{
		"isFirst": func(index int) bool {
			return index == 0
		},
	}).ParseFiles("tmux.tmpl")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, session)
	templateString := buf.String()

	err = os.MkdirAll(filepath.Dir(entryFilepath), 0755)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(entryFilepath, []byte(templateString), 0644)
	if err != nil {
		panic(err)
	}
}

func tmuxifyDirEntries(dirEntries []os.DirEntry) {
	for _, entry := range dirEntries {
		tmuxify(entry)
	}
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DATA.ConfigDir = filepath.Join(homeDir, CONFIG_DIR_SUFFIX)
	DATA.DataDir = filepath.Join(homeDir, DATA_DIR_SUFFIX)

	entries, err := os.ReadDir(DATA.ConfigDir)
	if err != err {
		panic(err)
	}

	tmuxifyDirEntries(entries)
}
