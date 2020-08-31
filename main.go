package main

import (
	// "context"
	"fmt"
    // "io/ioutil"

    "rfc2119/aws-tui/ui"
    "rfc2119/aws-tui/services"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	// "github.com/davecgh/go-spew/spew"

	// "strings"
)

// establishes one-way go channels between each service backend (i.e model) and the view. realistically speaking, the model is the aws sdk
// returns a slice of worker channels. each channel is concerned with a service, and only the view to the model may use the channel. for example, for a designated ec2 worker channel, only the view responsible for ec2 may listen to the channel and consume items
// this might break in the future. sometimes, multiple benefeciaries exist for a single work. for example, when deleting an ebs volume, the ec2 console should also make use of the deletion command/action to update the affected instance. i don't know how to approach this (yet)
func establishServiceChannels() chan services.Action{

    // TODO: available services
    ec2Chan := make(chan services.Action)
    // TODO: register a <-chan to view, chan<- to model
    return ec2Chan
}

func main() {

	fmt.Println("halp")
	app := tview.NewApplication()
	pages := tview.NewPages()
	table := tview.NewTable()
	// flex := tview.NewFlex()
    reservations := services.GetEC2Instances()      // TODO: ec2 namespace ?
	// fmt.Println(reservations)
	// spew.Dump(reservations)
    // TODO: here was teh ui code
	tree := tview.NewTreeView()
	rootNode := tview.NewTreeNode("EC2")

	topLevelNodesNames := []string{"Instances", "EBS"}
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
	// flex.AddItem(table, 0, 2, true)
	// flex.AddItem(description, 0, 1, true)
	grid.EAddItem(table, 0, 0, 1, 1, 0, 0, true)
	grid.EAddItem(description, 1, 0, 1, 1, 0, 0, false)
	pages.AddPage("EC2", tree, true, true)
	pages.AddPage("Instances", grid, true, false)
	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}

