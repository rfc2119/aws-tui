package ui

import (
	"fmt"
	"log"
	"rfc2119/aws-tui/common"
	"rfc2119/aws-tui/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"		// TODO: should probably remove this
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
	?		View this help message
	d		Describe instance
	x		Delete instance ? (TODO)
	e		Edit instance	(TODO)
	^w		Move to neighboring windows (TODO: unfocusable status bar)
	ESC		Move back one page (will exit this help message)
	`
)

// global ui elements (TODO: perhaps i should make them local somehow)
var grid *eGrid                       // the main container
var description = tview.NewTextView() // instance description
var table = tview.NewTable()          // instance status as in web ui
var statusBar = NewStatusBar()        // TODO: create a new unfocusable type
var gridEdit *eGrid
var instanceStatusRadioButton = NewRadioButtons([]string{"Start", "Stop", "Hibernate", "Reboot", "Terminate"})


// TODO: it doesn't make sense to export the type and have a New() function in the same time
type ec2Service struct {
	service
	Model *model.EC2Model
	// logger log.Logger

	// service specific data
	reservations []ec2.Reservation		// TODO: should be []ec2.Instance
}

// config: the aws client config that will create the service (the underlying model)
func NewEC2Service(config aws.Config, app *tview.Application, rootPage *ePages) *ec2Service {

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
		service: service{
			MainApp:  app,
			RootPage: rootPage,
		},
		Model: model.NewEC2Model(config),
	}
}

func (ec2svc *ec2Service) InitView() {

	// hacks
	grid = NewEgrid(ec2svc.RootPage)
	gridEdit = NewEgrid(ec2svc.RootPage)

	ec2svc.drawElements()
	ec2svc.setCallbacks()

	// configuration for ui elements
	statusBar.SetText("Status")

	table.SetBorders(false)
	table.SetSelectable(true, false) // rows: true, colums: false means select only rows
	table.Select(1, 1)
	table.SetFixed(0, 3)

	grid.HelpMessage = HELP_EC2_MAIN
	grid.SetRows(-3, -1, 2)
	grid.EAddItem(table, 0, 0, 30, 1, 0, 0, true)
	grid.EAddItem(description, 30, 0, 10, 1, 0, 0, false)
	grid.EAddItem(statusBar, 40, 0, 1, 1, 0, 0, false) // AddItem(p Primitive, row, column, rowSpan, colSpan, minGridHeight, minGridWidth int, focus bool)

	instanceStatusRadioButton.SetBorder(true).SetTitle("HALP")
	gridEdit.SetSize(2, 4, 10, 10) // SetSize(numRows, numColumns, rowSize, columnSize int)
	gridEdit.EAddItem(instanceStatusRadioButton, 0, 0, 1, 2, 0, 0, true)

	ec2svc.RootPage.EAddPage("Instances", grid, true, false)         // TODO: page names and such; resize=true, visible=false
	ec2svc.RootPage.EAddPage("Edit Instance", gridEdit, true, false) // TODO: page names and such

	ec2svc.WatchChanges()

}

// fills ui elements with appropriate initial data
func (ec2svc *ec2Service) drawElements() {
	// draw main table
	colNames := []string{"ID", "AMI", "Type", "State", "StateReason"} // TODO
	ec2svc.reservations = ec2svc.Model.GetEC2Instances()                    // directly invokes a method on the model
	for halpIdx := 0; halpIdx < len(colNames); halpIdx++ {
		table.SetCell(0, halpIdx,
			tview.NewTableCell(colNames[halpIdx]).SetAlign(tview.AlignCenter).SetSelectable(false))
	}
	for rowIdx, reservation := range ec2svc.reservations {
		instanceIdCell := tview.NewTableCell(*reservation.Instances[0].InstanceId)
		instanceAMICell := tview.NewTableCell(*reservation.Instances[0].ImageId)
		instanceTypeCell := tview.NewTableCell(string(reservation.Instances[0].InstanceType))
		instanceStateCell := tview.NewTableCell(string(reservation.Instances[0].State.Name))
		instanceStateReasonCell := tview.NewTableCell(*reservation.Instances[0].StateReason.Message)
		cells := []*tview.TableCell{instanceIdCell, instanceAMICell, instanceTypeCell, instanceStateCell, instanceStateReasonCell}
		for colIdx, cell := range cells {
			table.SetCell(rowIdx+1, colIdx, cell)
		}
	}


	// TODO: draw edit grid
}

// set function callbacks for different ui elements
func (ec2svc *ec2Service) setCallbacks() {

	// main table
	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			ec2svc.RootPage.ESwitchToPage("Services", false) // TODO: page names and such
		}
		// if key == tcell.KeyEnter {
		// 	table.SetSelectable(true, true)
		// }

	})
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'd' :
			row, _ := table.GetSelection()
			description.SetText(fmt.Sprintf("%v", ec2svc.reservations[row-1].Instances[0]))
		case 'e':
			ec2svc.RootPage.ESwitchToPage("Edit Instance", true) // TODO: page names and such

		}

		return event
	})

	// main grid
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlW:
			statusBar.SetText("moving to another item; statusbar focus: " + fmt.Sprintf("%s", statusBar.HasFocus())) // TODO
			if len(grid.Members) > 0 {
				grid.CurrentMemberInFocus++
				if grid.CurrentMemberInFocus == len(grid.Members) { //  grid.CurrentMemberInFocus %= len(grid.Members)
					grid.CurrentMemberInFocus = 0
				}
				for { // a HACK to not focus on non-focusable items
					nextMemberToFocus := grid.Members[grid.CurrentMemberInFocus]
					ec2svc.MainApp.SetFocus(nextMemberToFocus)
					if !nextMemberToFocus.GetFocusable().HasFocus() {          // item didn't get focus despite giving it. cycle to the next member
						grid.CurrentMemberInFocus++
						if grid.CurrentMemberInFocus == len(grid.Members) { //  grid.CurrentMemberInFocus %= len(grid.Members)
							grid.CurrentMemberInFocus = 0
						}
					}else { break }
				}
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case '?':
				ec2svc.RootPage.DisplayHelpMessage(grid.HelpMessage)
			}
		}
		return event
	})

	// edit grid
	gridEdit.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			ec2svc.RootPage.ESwitchToPage("Instances", false) // TODO: page names and such
		}
		return event
	})
}

//
func (svc *ec2Service) GetMainElement() tview.Primitive {
	return grid
	// return flex
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
