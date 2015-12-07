package main

import (
	"fmt"

	"github.com/op/go-logging"
)

// Node contains all the information of a single organizational node
type Node struct {
	Name         string
	Path         string
	Title        string
	Config       NodeConfig
	Markdown     string
	Children     map[string]*Node
	DependsOn    []*Node
	DependedOnBy []*Node
}

// NodeConfig contains information that is represented in the toml header of a given node
type NodeConfig struct {
	Title     string
	DependsOn []string
}

func gatherDependencies(node *Node, nodeMap map[string]*Node) error {
	if len(node.Name) > 0 {
		nodeMap[node.Name] = node
	}
	for _, child := range node.Children {
		if err := gatherDependencies(child, nodeMap); err != nil {
			return err
		}
	}
	return nil
}

func addDependedOnByLink(node *Node, dependent *Node) {
	if node.DependedOnBy == nil {
		node.DependedOnBy = make([]*Node, 0)
	}
	node.DependedOnBy = append(node.DependedOnBy, dependent)
}

func associateDependencies(node *Node, nodeMap map[string]*Node) error {
	if node.Config.DependsOn != nil {
		node.DependsOn = make([]*Node, len(node.Config.DependsOn))
		for index, dependency := range node.Config.DependsOn {
			node.DependsOn[index] = nodeMap[dependency]
			if node.DependsOn[index] == nil {
				return fmt.Errorf("Unknown dependency on %v: %v", node.Name, dependency)
			}
			addDependedOnByLink(nodeMap[dependency], node)
		}
	}

	for _, child := range node.Children {
		if err := associateDependencies(child, nodeMap); err != nil {
			return err
		}
	}
	return nil
}

func PopulateDependencies(mainNode *Node) error {
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
