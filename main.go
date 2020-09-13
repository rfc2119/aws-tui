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

// TODO: information about region, IAM user, sdk version used, current build version, ... etc.
// TODO: available services

func main() {

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	config, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	fmt.Println("halp")

	// application, root element and status bar
	app := tview.NewApplication()
	pages := ui.NewEPages()
	statusBar := ui.NewStatusBar()

	// services
	ec2svc := ui.NewEC2Service(config, app, pages, statusBar)
	ec2svc.InitView()

	// main ui element
	tree := tview.NewTreeView()

	// configuring the tree:
	rootNode := tview.NewTreeNode("Services")
	var topLevelNodesNames []string
	for _, name := range common.ServiceNames {
		topLevelNodesNames = append(topLevelNodesNames, name)
	}
	levelInstances := []string{"Instances"}
	levelEBS := []string{"Volumes"}
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
			pages.ESwitchToPage(node.GetText())
			// tview.NewModal().SetText("children").AddButtons([]string{"ok"})

		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	// ui config
	mainContainer := tview.NewFlex() // a flex container for the status bar and application pages/window
	mainContainer.SetDirection(tview.FlexRow).SetFullScreen(true)
	mainContainer.AddItem(pages, 0, 107, true)    //AddItem(item Primitive, fixedSize, proportion int, focus bool)
	mainContainer.AddItem(statusBar, 0, 1, false) // 107:1 seems fair ?
	pages.EAddPage("Services", tree, true, true)  // EAddPage(name string, item tview.Primitive, resize, visible bool)
    statusBar.SetText("Welcome to the terminal interface for AWS. Type '?' to get help")
	if err := app.SetRoot(mainContainer, true).SetFocus(mainContainer).Run(); err != nil {
		panic(err)
	}
}
