package main

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/russross/blackfriday"
)

// Node contains all the information of a single organizational node
type Node struct {
	Name      string
	Path      string
	Title     string
	Config    NodeConfig
	Markdown  string
	Children  map[string]*Node
	DependsOn []*Node
}

// NodeConfig contains information that is represented in the toml header of a given node
type NodeConfig struct {
	Title     string
	DependsOn []string
}

// parseFile turns a toml-prefixed markdown file into a Node structure
func parseFile(filePath string, name string) (*Node, error) {
	node := Node{
		Name:     name,
		Path:     filePath,
		Children: make(map[string]*Node),
	}

	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fileString := string(fileContents)
	// Read toml header if present
	if len(fileString) > 3 && fileString[0:3] == "+++" {
		parts := strings.Split(fileString, "+++")
		headerString := parts[1]
		fileString = strings.Join(parts[2:], "+++")
		if _, err = toml.Decode(headerString, &node.Config); err != nil {
			return nil, err
		}
	}

	if len(node.Config.Title) > 0 {
		node.Title = node.Config.Title
	} else {
		node.Title = node.Name
	}
	node.Markdown = renderMarkdown(fileString)

	return &node, nil
}

// parseDirectory turns the project file into a Node, and adds any children nodes that it finds
func parseDirectory(directory string, name string) (*Node, error) {
	var projectFile string
	if len(name) == 0 {
		projectFile = "project.md"
	} else {
		projectFile = name + ".md"
	}

	projectFilePath := path.Join(directory, projectFile)
	node, err := parseFile(projectFilePath, name)

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		var child *Node
		filePath := path.Join(directory, f.Name())
		if f.Name() == projectFile {
			continue
		} else if f.IsDir() {
			child, err = parseDirectory(filePath, f.Name())
		} else if path.Ext(f.Name()) == ".md" {
			name := f.Name()[:len(f.Name())-3]
			child, err = parseFile(filePath, name)
		}
		if err != nil {
			return nil, err
		}
		if child != nil {
			node.Children[child.Name] = child
		}
	}

	return node, err
}

func renderMarkdown(markdown string) string {
	return string(blackfriday.MarkdownCommon([]byte(markdown)))
}

// ParseTopDirectory is just a shorthand to parse the top level, where the project has no implied name
func ParseTopDirectory(directory string) (*Node, error) {
	return parseDirectory(directory, "")
}
