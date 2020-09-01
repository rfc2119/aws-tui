package main

import (
	// "context"
	"fmt"
	// "io/ioutil"

	"rfc2119/aws-tui/common"
	"rfc2119/aws-tui/ui"

	// "github.com/gdamore/tcell"
	"github.com/rivo/tview"
	// "github.com/davecgh/go-spew/spew"
	// "strings"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

// establishes one-way go channels between each service backend (i.e model) and the view. realistically speaking, the model is the aws sdk
// returns a slice of worker channels. each channel is concerned with a service, and only the view to the model may use the channel. for example, for a designated ec2 worker channel, only the view responsible for ec2 may listen to the channel and consume items
// this might break in the future. sometimes, multiple benefeciaries exist for a single work. for example, when deleting an ebs volume, the ec2 console should also make use of the deletion command/action to update the affected instance. i don't know how to approach this (yet)
// func establishServiceChannels() chan services.Action {
// 
// 	// TODO: available services
// 	ec2Chan := make(chan services.Action)
// 	// TODO: register a <-chan to view, chan<- to model
// 	return ec2Chan
// }

func main() {

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	config, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	fmt.Println("halp")

	// application and root element
	app := tview.NewApplication()
	pages := tview.NewPages()

	// services
	ec2svc := ui.NewEC2Service(config, app, pages)
	ec2svc.InitView()

	// service channels (to be used by the view) TODO
	// ec2Chan := ec2svc.Channel


	// main ui:
	tree := tview.NewTreeView()

	// configuring the tree:
	rootNode := tview.NewTreeNode("Services")
	var topLevelNodesNames []string
	for _, name := range common.ServiceNames{
		topLevelNodesNames = append(topLevelNodesNames, name)
	}
	levelInstances := []string{"Instances", "HALP"}
	levelEBS := []string{"Volumes", "WELP"}
	// var topLevelNodes []*tview.TreeNode
	levelOne := [][]string{levelInstances, levelEBS}

	// add levelX to top level nodes
	for idx, node := range topLevelNodesNames {
		_tmp := tview.NewTreeNode(node)
		// topLevelNodes = append(topLevelNodes, _tmp)
		for _, child := range levelOne[idx] {
			_tmpChild := tview.NewTreeNode(child)
			_tmp.AddChild(_tmpChild)
			// _tmpChild.SetExpanded(false)
		}
		_tmp.Collapse()
		rootNode.AddChild(_tmp)

	}
	tree.SetRoot(rootNode)
	tree.SetCurrentNode(rootNode)
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		children := node.GetChildren()
		if len(children) == 0 && pages.HasPage(node.GetText()) { // go to page
			pages.SwitchToPage(node.GetText())
			// tview.NewModal().SetText("children").AddButtons([]string{"ok"})

		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	// ui config
	// flex.AddItem(table, 0, 2, true)
	// flex.AddItem(description, 0, 1, true)
	pages.AddPage("Services", tree, true, true)
	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
