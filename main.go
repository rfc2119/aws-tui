package main

import (
	// "context"        // TODO
	// "runtime/debug"
	"fmt"
	"github.com/rfc2119/aws-tui/common"
	"github.com/rfc2119/aws-tui/ui"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/rivo/tview"
)

const (
	MAIN_HELP_MSG = `
Welcome to the unofficial AWS Terminal Interface. This is a very much work-in-progress and I appreciate your feedback, issues, code improvements, ... etc. Please submit them at https://github.com/rfc2119/aws-tui

Common keys found across all windows:

	TAB             Move to neighboring windows
	?               View help messages (if available)
	q               Move back one page (will exit this help message)
    Space           Select Option in a radio box/tree view (except in a confirmation box)
    hjkl		Movement keys

There's likely a help page for every window, so please use '?' when in doubt. Use Ctrl-C to exit the application.
`
)

var (
	// Filled by goreleaser at build time
	version = ""
	commit  = ""
	date    = ""
	// branch  = ""
)

func main() {

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	config, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Application, root element and status bar
	app := tview.NewApplication()
	pages := ui.NewEPages()
	statusBar := ui.NewStatusBar()

	// Services
	ec2svc := ui.NewEC2Service(config, app, pages, statusBar)
	ec2svc.InitView() // TODO: call only when user selects the service
	iamsvc := ui.NewIAMService(config, app, pages, statusBar)

	// UI elements
	mainContainer := tview.NewFlex() // Flex container for the status bar and application pages/window
	frontPage := ui.NewEFlex(pages)  // The front page which holds the info and tree view
	info := tview.NewTextView()	// The side bar informational view
	tree := tview.NewTreeView()	// The menu holding all available services

	// Filling the tree with initial values
	rootNode := tview.NewTreeNode("Services")
	for service, desc := range common.AWServicesDescriptions {
		if desc.Available {
			nodeLevel1 := tview.NewTreeNode(desc.Name)
			for _, subItemName := range common.ServiceChildrenNames[service] {
				nodeLevel2 := tview.NewTreeNode(subItemName)
				nodeLevel1.AddChild(nodeLevel2)
				// _tmpChild.SetExpanded(false)
			}
			nodeLevel1.Collapse()
			rootNode.AddChild(nodeLevel1)
		}
	}

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		children := node.GetChildren()
		if len(children) == 0 && pages.HasPage(node.GetText()) { // go to page
			pages.ESwitchToPage(node.GetText()) // TODO: unify page names

		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	// Filling the info box with initial values
	currentIAMUser := iamsvc.Model.GetCurrentUserInfo()
	fmt.Fprintf(info,
		`
    IAM User name: %7s
    IAM User arn:  %20s
    Region:        %7s

    Build Version: %s
    Build Commit:  %s
    Build Date:    %s
    SDK Name:      %7s
    SDK Version:   %-7s
    `, *currentIAMUser.UserName, *currentIAMUser.Arn, config.Region, version, commit, date, aws.SDKName, aws.SDKVersion)

	// UI config
	tree.SetRoot(rootNode)
	tree.SetCurrentNode(rootNode)

	frontPage.HelpMessage = MAIN_HELP_MSG
	frontPage.SetDirection(tview.FlexColumn)
	frontPage.AddItem(tree, 0, 3, true)
	frontPage.AddItem(info, 0, 2, false)

	mainContainer.SetDirection(tview.FlexRow).SetFullScreen(true)
	mainContainer.AddItem(pages, 0, 107, true)    // AddItem(item Primitive, fixedSize, proportion int, focus bool)
	mainContainer.AddItem(statusBar, 0, 1, false) // 107:1 seems fair ?

	pages.EAddPage("Services", frontPage, true, true) // (name string, item tview.Primitive, resize, visible bool)
	statusBar.SetText("Welcome to the terminal interface for AWS. Type '?' to get help")
	if err := app.SetRoot(mainContainer, true).SetFocus(mainContainer).Run(); err != nil {
		panic(err)
	}
}
