package main

import (
	"os"
	"path"
	"text/template"
)

// RenderNode renders the index.html for a project which includes the project
// summary and the graph
func renderNode(destinationPath string, node *Node) error {
	tmplData, err := Asset("templates/tree.html")
	if err != nil {
		return err
	}
	tmpl, err := template.New("node").Parse(string(tmplData))
	if err != nil {
		return err
	}

	fileName := node.Name
	if len(fileName) == 0 {
		fileName = "index"
	}
	file, err := os.Create(path.Join(destinationPath, fileName+".html"))
	if err != nil {
		return err
	}
	if err = tmpl.Execute(file, node); err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	for _, child := range node.Children {
		err = renderNode(destinationPath, child)
		if err != nil {
			return err
		}
	}

	return nil
}

func RenderTopNode(destinationPath string, node *Node) error {
	if err := PopulateDependencies(node); err != nil {
		return err
	}

	return renderNode(destinationPath, node)
}
