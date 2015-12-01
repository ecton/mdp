package main

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/op/go-logging"
)

func gatherDependencies(node *Node, nodeMap map[string]*Node) error {
	if len(node.Name) > 0 {
		logging.MustGetLogger("mdp").Infof("Adding %v to map", node.Name)
		nodeMap[node.Name] = node
	}
	for _, child := range node.Children {
		if err := gatherDependencies(child, nodeMap); err != nil {
			return err
		}
	}
	return nil
}

func associateDependencies(node *Node, nodeMap map[string]*Node) error {
	if node.Config.DependsOn != nil {
		node.DependsOn = make([]*Node, len(node.Config.DependsOn))
		for index, dependency := range node.Config.DependsOn {
			node.DependsOn[index] = nodeMap[dependency]
			if node.DependsOn[index] == nil {
				return fmt.Errorf("Unknown dependency on %v: %v", node.Name, dependency)
			}
		}
	}

	for _, child := range node.Children {
		if err := associateDependencies(child, nodeMap); err != nil {
			return err
		}
	}
	return nil
}

func populateDependencies(mainNode *Node) error {
	allDependencies := make(map[string]*Node)
	if err := gatherDependencies(mainNode, allDependencies); err != nil {
		return err
	}
	logging.MustGetLogger("mdp").Infof("%v", len(allDependencies))
	if err := associateDependencies(mainNode, allDependencies); err != nil {
		return err
	}
	return nil
}

// RenderNode renders the index.html for a project which includes the project
// summary and the graph
func RenderNode(destinationPath string, node *Node) error {
	if err := populateDependencies(node); err != nil {
		return err
	}

	tmplData, err := Asset("templates/basic.html")
	if err != nil {
		return err
	}
	tmpl, err := template.New("node").Parse(string(tmplData))
	if err != nil {
		return err
	}

	file, err := os.Create(path.Join(destinationPath, "index.html"))
	if err != nil {
		return err
	}
	if err = tmpl.Execute(file, node); err != nil {
		return err
	}
	err = file.Close()

	return err
}
