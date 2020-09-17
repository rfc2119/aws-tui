package ui

import (
	"fmt"
	"log"
	"rfc2119/aws-tui/common"
	"rfc2119/aws-tui/model"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2" // TODO: should probably remove this
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	COL_ID = iota
	COL_AMI
	COL_TYPE
	COL_STATE
	COL_STATEREASON
)

const (
	HELP_EC2_MAIN = `
	?               View this help message
	d               Describe instance
    r               Refresh list of instances
	e               Edit instance (WIP)
	^w              Move to neighboring windows
    ^l              Filter and List AMIs
	q               Move back one page (will exit this help message)
	`
	HELP_EC2_EDIT_INSTANCE = `
    Space           Select Option in a radio box
	q               Move back one page (will exit this help message)
    `
)

// global ui elements (TODO: perhaps i should make them local somehow)
var grid *eGrid                       // the main container
var description = tview.NewTextView() // instance description
var table = tview.NewTable()          // instance status as in web ui
var gridEdit *eGrid
var instanceOfferingsDropdown = tview.NewDropDown()
var instanceStatusRadioButton = NewRadioButtons([]string{"Start", "Stop", "Hibernate", "Reboot", "Terminate"}) // all buttons are enabled by default

// TODO: it doesn't make sense to export the type and have a New() function in the same time
type ec2Service struct {
	mainUI
	Model *model.EC2Model
	// logger log.Logger

	// service specific data
	instances []ec2.Instance
}

// config: the aws client config that will create the service (the underlying model)
func NewEC2Service(config aws.Config, app *tview.Application, rootPage *ePages, statBar *StatusBar) *ec2Service {

	// var components []viewComponent
	// for _, elm := range elements {
	// 	viewComponent := viewComponent{
	// 		ID:      fmt.Sprintf("%p", elm),
	// 		Service: ServiceNames[SERVICE_EC2],
	// 		Element: elm,
	// 	}

	// 	components = append(components, viewComponent)
	// }
	return &ec2Service{
		mainUI: mainUI{
			MainApp:   app,
			RootPage:  rootPage,
			StatusBar: statBar,
		},
		Model:        model.NewEC2Model(config),
		instances: nil,
	}
}

func (ec2svc *ec2Service) InitView() {

	ec2svc.StatusBar.SetText("starting ec2 service")
	// hacks
	grid = NewEgrid(ec2svc.RootPage)
	gridEdit = NewEgrid(ec2svc.RootPage)

	ec2svc.drawElements()
	ec2svc.setCallbacks()

	// configuration for ui elements
	ec2svc.StatusBar.SetText("Status")

	table.SetBorders(false)
    table.SetBorderAttributes(tcell.AttrBold)
	table.SetSelectable(true, false) // rows: true, colums: false means select only rows
	table.Select(1, 1)
	table.SetFixed(0, 3)

	grid.HelpMessage = HELP_EC2_MAIN
	grid.SetRows(-3, -1, 2)
	grid.EAddItem(table, 0, 0, 30, 1, 0, 0, true)
	grid.EAddItem(description, 30, 0, 10, 1, 0, 0, false)
	// grid.EAddItem(StatusBar, 40, 0, 1, 1, 0, 0, false) // AddItem(p Primitive, row, column, rowSpan, colSpan, minGridHeight, minGridWidth int, focus bool)

	instanceStatusRadioButton.SetBorder(true).SetTitle("Status")
	instanceOfferingsDropdown.SetLabel("Type")
	gridEdit.HelpMessage = HELP_EC2_EDIT_INSTANCE
	gridEdit.SetSize(2, 4, 0, 0) // SetSize(numRows, numColumns, rowSize, columnSize int)
	gridEdit.EAddItem(instanceStatusRadioButton, 0, 0, 1, 2, 0, 0, true)
	gridEdit.EAddItem(instanceOfferingsDropdown, 0, 2, 1, 2, 0, 0, false)

	ec2svc.RootPage.EAddPage("Instances", grid, true, false)         // TODO: page names and such; resize=true, visible=false
	ec2svc.RootPage.EAddPage("Edit Instance", gridEdit, true, false) // TODO: page names and such

	ec2svc.WatchChanges()

}

// TODO: // Convert *string to string value: str = aws.StringValue(strPtr)
// fills ui elements with appropriate initial data
func (ec2svc *ec2Service) drawElements() {
	// draw main table
	ec2svc.fillMainTable()

	// TODO: draw edit grid
	offerings := ec2svc.Model.ListOfferings()
	opts := make([]string, len(offerings))
	for idx := 0; idx < len(offerings); idx++ {
		opts[idx] = string(offerings[idx].InstanceType) // TODO: do not forget pagination
	}
	instanceOfferingsDropdown.SetOptions(opts, nil)
}

// set function callbacks for different ui elements
func (ec2svc *ec2Service) setCallbacks() {

	// main table
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'd':
			row, _ := table.GetSelection()
			description.SetText(fmt.Sprintf("%v", ec2svc.instances[row-1]))
		case 'e':
			ec2svc.RootPage.ESwitchToPage("Edit Instance") // TODO: page names and such
		case 'r':
			ec2svc.StatusBar.SetText("refreshing instances list")
			ec2svc.fillMainTable()
		}

		return event
	})

    // TODO: unify grids
	// main grid
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {

		case tcell.KeyCtrlW:
			if len(grid.Members) > 0 {
				grid.CurrentMemberInFocus++
				if grid.CurrentMemberInFocus == len(grid.Members) { //  grid.CurrentMemberInFocus %= len(grid.Members)
					grid.CurrentMemberInFocus = 0
				}
				for { // a HACK to not focus on non-focusable items
					nextMemberToFocus := grid.Members[grid.CurrentMemberInFocus]
					ec2svc.MainApp.SetFocus(nextMemberToFocus)
					if !nextMemberToFocus.GetFocusable().HasFocus() { // item didn't get focus despite giving it. cycle to the next member
						grid.CurrentMemberInFocus++
						if grid.CurrentMemberInFocus == len(grid.Members) { //  grid.CurrentMemberInFocus %= len(grid.Members)
							grid.CurrentMemberInFocus = 0
						}
					} else {
						break
					}
				}
			}
		case tcell.KeyCtrlL:
			// build modal and let user choose filters
			ec2svc.chooseAMIFilters()

		case tcell.KeyRune:
			switch event.Rune() {
            case 'e':
                ec2svc.StatusBar.SetText("E FROM GRID")
			case '?':
				grid.DisplayHelp()
			case 'q':
				// ec2svc.RootPage.ESwitchToPage("Services") // TODO: page names and such
				ec2svc.RootPage.ESwitchToPreviousPage()
				ec2svc.StatusBar.SetText("exit ec2")
			}
		}
		return event
	})

	// edit grid (TODO: copy pasta from above)
	gridEdit.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ec2svc.StatusBar.SetText(fmt.Sprintf("grid: %v. radio button: %v. dropdown: %v", gridEdit.HasFocus(), instanceStatusRadioButton.HasFocus(), instanceOfferingsDropdown.HasFocus())) // TODO
		switch event.Key() {
		case tcell.KeyCtrlW:
			if len(gridEdit.Members) > 0 {
				gridEdit.CurrentMemberInFocus++
				if gridEdit.CurrentMemberInFocus == len(gridEdit.Members) { //  gridEdit.CurrentMemberInFocus %= len(gridEdit.Members)
					gridEdit.CurrentMemberInFocus = 0
				}
				for { // a HACK to not focus on non-focusable items
					nextMemberToFocus := gridEdit.Members[gridEdit.CurrentMemberInFocus]
					ec2svc.MainApp.SetFocus(nextMemberToFocus)
					if !nextMemberToFocus.GetFocusable().HasFocus() { // item didn't get focus despite giving it. cycle to the next member
						gridEdit.CurrentMemberInFocus++
						if gridEdit.CurrentMemberInFocus == len(gridEdit.Members) { //  gridEdit.CurrentMemberInFocus %= len(gridEdit.Members)
							gridEdit.CurrentMemberInFocus = 0
						}
					} else {
						break
					}
				}
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case '?':
				gridEdit.DisplayHelp()
			case 'q':
				ec2svc.RootPage.ESwitchToPreviousPage()
			}
		}
		return event
	})

	// radio button
	instanceStatusRadioButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter || event.Rune() == ' ' {
			ec2svc.showConfirmationBox(fmt.Sprintf("%s instance ?", instanceStatusRadioButton.GetCurrentOptionName()))
		}
		return event
	})

	// dropdown instance offering
	instanceOfferingsDropdown.SetSelectedFunc(func(text string, index int) {
		ec2svc.showConfirmationBox(fmt.Sprintf("Change instance type to %s ?", text))
	})

}

// TODO: could this be a generic filter box ?
func (ec2svc *ec2Service) chooseAMIFilters() {
	var (
		filterNames []string // need it for the dropdown
	)
	filters := make(map[string]string)
	form := tview.NewForm()
	inputField := tview.NewInputField()
	filterValuesAutoComplete := make([][]string, len(common.AMIFilters))

	for idx, filterIdx := range common.AMIFilters {
		filterNames = append(filterNames, common.FilterNames[filterIdx][0])
		filterValuesAutoComplete[idx] = common.FilterNames[filterIdx][1:]
	}
	prevName := ""
	dropDownSelectedFunc := func(option string, optionIndex int) {
		previousText, exists := filters[prevName]
		// save current filter value if existed before or if there's new text
		if txt := inputField.GetText(); txt != "" || exists {
			ec2svc.StatusBar.SetText(fmt.Sprintf("prev text: %s, exists: %v, prevName: %s", previousText, exists, prevName))
			if prevName != "" { // avoid initial value of prevName
				filters[prevName] = txt
			}
		}
		// set auto complete for the current selected text. copied from demos/inputfield
		inputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
			if len(currentText) == 0 {
				return
			}
			ec2svc.StatusBar.SetText(fmt.Sprintf("%s", filterValuesAutoComplete[optionIndex]))
			for _, word := range filterValuesAutoComplete[optionIndex] {
				if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
					entries = append(entries, word)
				}
			}
			if len(entries) < 1 {
				entries = nil
			}
			return
		})
		inputField.SetText(filters[option]) // restore value for  selected option, or clear the field
		prevName = option
	}

	buttonCancelFunc := func() { ec2svc.RootPage.ESwitchToPreviousPage() }
	buttonSaveFunc := func() {
		ec2svc.StatusBar.SetText("Grabbing the list of AMIs")
		amis := ec2svc.Model.ListAMIs(filters)
		tableAMI := tview.NewTable()

		// drawing the table
		colNames := []string{"ID", "State", "Arch", "Creation Date", "Name", "Owner ID"} // TODO
		for halpIdx := 0; halpIdx < len(colNames); halpIdx++ {
			tableAMI.SetCell(0, halpIdx,
				tview.NewTableCell(colNames[halpIdx]).SetAlign(tview.AlignCenter).SetSelectable(false))
		}
		for rowIdx, ami := range amis {
			idCell := tview.NewTableCell(*ami.ImageId)
			stateCell := tview.NewTableCell(string(ami.State))
			archCell := tview.NewTableCell(string(ami.Architecture))
			dateCell := tview.NewTableCell(*ami.CreationDate)
			// ownerCell := tview.NewTableCell(*ami.ImageOwnerAlias)    // TODO:
			nameCell := tview.NewTableCell(*ami.Name)
			ownerCell := tview.NewTableCell(*ami.OwnerId)
			cells := []*tview.TableCell{idCell, stateCell, archCell, dateCell, nameCell, ownerCell}
			for colIdx, cell := range cells {
				tableAMI.SetCell(rowIdx+1, colIdx, cell)
			}
		}
		tableAMI.SetBorders(true)
        tableAMI.SetBorderAttributes(tcell.AttrBold)
		tableAMI.SetSelectable(true, false) // rows: true, colums: false means select only rows
		tableAMI.Select(1, 1)
		tableAMI.SetFixed(1, 1)
		tableAMI.SetDoneFunc(func(key tcell.Key) {
			ec2svc.RootPage.ESwitchToPreviousPage()
		})
		ec2svc.RootPage.AddAndSwitchToPage("AMIs", tableAMI, true)
	}

	inputField = tview.NewInputField().SetLabel("Filter Value")
	form.AddDropDown("Filter Name", filterNames, 0, dropDownSelectedFunc)
	form.AddButton("Save", buttonSaveFunc)
	form.AddButton("Cancel", buttonCancelFunc)
	form.AddFormItem(inputField)
	form.SetTitle("Filter AMIs").SetBorder(true)
    ec2svc.showGenericModal("Filter AMIs", form, 80, 10)    // 80x10 seems good for my screen
}

// shows a generic modal box (rather than a confirmation-only box) centered at screen
// props to skanehira from the docker tui "docui" for this! code is at github.com/skanehira/docui
func (ec2svc *ec2Service) showGenericModal(title string, p tview.Primitive, width, height int) {
    centeredModal := tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
    currPageName := ec2svc.RootPage.GetCurrentPageName()
	ec2svc.RootPage.EAddAndSwitchToPage(title, centeredModal, true)    // resize=true
    ec2svc.RootPage.ShowPage(currPageName)   // redraw on top of the box
}

// shows a modal box with msg and switches back to previous page. this is useful for ont-time usage (no nested boxes). selected text is returned for convenience
func (ec2svc *ec2Service) showConfirmationBox(msg string) string {
	var selectedButtonLabel string
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"Ok", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// if buttonLabel == "Ok" {
			// 	ec2svc.StatusBar.SetText(statusMsg)
			// }
			selectedButtonLabel = buttonLabel
			ec2svc.RootPage.ESwitchToPreviousPage()
			// ec2svc.RootPage.RemovePage("modal")		// TODO: is this necessary ? this will loop over all pages
		})
	ec2svc.RootPage.EAddAndSwitchToPage("modal", modal, false) // resize=false
    ec2svc.RootPage.ShowPage(ec2svc.RootPage.GetPreviousPageName())      // +1
	return selectedButtonLabel
}

func (ec2svc *ec2Service) fillMainTable() {

	colNames := []string{"ID", "AMI", "Type", "State", "StateReason"} // TODO
	ec2svc.instances = ec2svc.Model.GetEC2Instances()              // directly invokes a method on the model
	for halpIdx := 0; halpIdx < len(colNames); halpIdx++ {
		table.SetCell(0, halpIdx,
			tview.NewTableCell(colNames[halpIdx]).SetAlign(tview.AlignCenter).SetSelectable(false))
	}
	for rowIdx, instance := range ec2svc.instances {
		instanceIdCell := tview.NewTableCell(aws.StringValue(instance.InstanceId))
		instanceAMICell := tview.NewTableCell(*instance.ImageId)
		instanceTypeCell := tview.NewTableCell(string(instance.InstanceType))
		instanceStateCell := tview.NewTableCell(string(instance.State.Name))
		instanceStateReasonCell := tview.NewTableCell(*instance.StateReason.Message)
		cells := []*tview.TableCell{instanceIdCell, instanceAMICell, instanceTypeCell, instanceStateCell, instanceStateReasonCell}
		for colIdx, cell := range cells {
			table.SetCell(rowIdx+1, colIdx, cell)
		}
	}
}

// dispatches goroutines to monitor changes; assigns listeners to each action
func (svc *ec2Service) WatchChanges() {
	svc.Model.DispatchWatchers()
	go func(ch <-chan common.Action) { // listner goroutine
		for {
			receiveMe := <-ch
			// log.Println("listener received data")
			// switch receiveMe.Data.(type){	// FIXME (see below)
			switch receiveMe.Type {
			// case common.InstanceStatusesUpdate:	// FIXME why doesn't this work ? received type is []ec2.InstanceStatus
			case common.ACTION_INSTANCE_STATUS_UPDATE:
				// log.Println("listener 1 is dispatched")
				go listener1(receiveMe)
			default:
				log.Printf("received data of type %T", receiveMe.Data)
			}
		}
	}(svc.Model.Channel)

}

// listener for watcher1
func listener1(action common.Action) {

	statuses := action.Data.(common.InstanceStatusesUpdate)
	for _, status := range statuses {
		rowIdx := rowIndexFromTable(table, *status.InstanceId) // TODO: check for -1
		cell := table.GetCell(rowIdx, COL_STATE)
		newState := string(status.InstanceState.Name)
		// log.Printf("old state: %s cell: %s", newState, cell.Text)
		if newState != cell.Text {
			// cell.SetBackgroundColor(tcell.ColorRed)		// TODO: colors for transitions; color row
			colorizeRowInTable(table, rowIdx, tcell.ColorRed)
			cell.SetText(newState)
		}

	}
}

// TODO: enum ? func (enum SummaryStatus) MarshalValue() (string, error)
// ============ helper functions
// given an instance ID, return the row index of the instance in table t
func rowIndexFromTable(t *tview.Table, instanceID string) int {
	idx := -1
	for rowIdx := 1; rowIdx < t.GetRowCount(); rowIdx++ { // 1 because first row is for column labels
		id := t.GetCell(rowIdx, COL_ID).Text
		if instanceID == id {
			idx = rowIdx
			break
		}
	}
	return idx
}

// colorize a row in a given table
func colorizeRowInTable(t *tview.Table, row int, color tcell.Color) {
	for col := 0; col < t.GetColumnCount(); col++ {
		t.GetCell(row, col).SetBackgroundColor(color)
	}
}

// tweak the edit grid according to each instance
func modifyEditGrid(g *eGrid, instanceIdx int) {

}
