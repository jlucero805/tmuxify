package main

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
)

//go:embed tmux.tmpl
var embeddedTemplate []byte

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

func readConfigFile(name string) (string, error) {
	data, err := os.ReadFile(filepath.Join(DATA.ConfigDir, name))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func parseSession(sessionString string) (Session, error) {
	var session Session
	_, err := toml.Decode(
		sessionString,
		&session,
	)
	if err != nil {
		return Session{}, err
	}

	indexSessionWins(&session)

	return session, nil
}

func indexSessionWins(session *Session) {
	for i := 0; i < len(session.Win); i++ {
		session.Win[i].Index = i + 1
	}
}

func writeDataFile(data, name string) error {
	dataFilePath := filepath.Join(DATA.DataDir, name)

	err := os.MkdirAll(filepath.Dir(dataFilePath), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(dataFilePath, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}

func renderTemplate(session Session) (string, error) {
	tmpl, err := template.New("tmuxify").Funcs(template.FuncMap{
		"isFirst": func(index int) bool {
			return index == 0
		},
	}).Parse(string(embeddedTemplate))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, session)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func tmuxify(entry os.DirEntry) {
	data, err := readConfigFile(entry.Name())
	if err != nil {
		panic(err)
	}

	session, err := parseSession(data)
	if err != nil {
		panic(err)
	}

	templateString, err := renderTemplate(session)
	if err != nil {
		panic(err)
	}

	err = writeDataFile(templateString, entry.Name())
	if err != nil {
		panic(err)
	}
}

func tmuxifyDirEntries(dirEntries []os.DirEntry) {
	for _, entry := range dirEntries {
		tmuxify(entry)
	}
}

func getTmuxifyConfigDirEntries() ([]os.DirEntry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	DATA.ConfigDir = filepath.Join(homeDir, CONFIG_DIR_SUFFIX)
	DATA.DataDir = filepath.Join(homeDir, DATA_DIR_SUFFIX)

	entries, err := os.ReadDir(DATA.ConfigDir)
	if err != err {
		return nil, err
	}

	return entries, nil
}

func main() {
	entries, err := getTmuxifyConfigDirEntries()
	if err != nil {
		panic(err)
	}

	tmuxifyDirEntries(entries)
}
