package ui

import (
	"fmt"

	"rfc2119/aws-tui/model"
	"rfc2119/aws-tui/common"
	"github.com/aws/aws-sdk-go-v2/aws"
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
	x		Delete instance ?
	^w		Move to neighboring windows
	ESC		Move back one page
	`
)

// local ui elements
var grid = NewEgrid()                 // the main container
var description = tview.NewTextView() // instance description
var table = tview.NewTable()          // instance status as in web ui
// var flex = tview.NewFlex()
var statusBar = tview.NewTextView()
var colNames = []string{"ID", "AMI", "Type", "State", "StateReason"} // TODO

// TODO: it doens't make sense to export the type and have a New() function in the same time
type ec2Service struct {
	service
	Model *model.EC2Model
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

	reservations := ec2svc.Model.GetEC2Instances() // directly invokes a method on the model

	for halpIdx := 0; halpIdx < len(colNames); halpIdx++ {
		table.SetCell(0, halpIdx,
			tview.NewTableCell(colNames[halpIdx]).SetAlign(tview.AlignCenter).SetSelectable(false))
	}
	for rowIdx, reservation := range reservations {
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
	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			ec2svc.RootPage.ESwitchToPage("Services") // TODO: page names and such
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}

	})
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'd' {
			row, _ := table.GetSelection()
			description.SetText(fmt.Sprintf("%v", reservations[row-1].Instances[0]))
		}

		return event
	})
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlW {
			statusBar.SetText("moving to another item") // TODO
			if len(grid.Members) > 0 {
				grid.CurrentMemberInFocus++
				grid.CurrentMemberInFocus %= len(grid.Members)
				ec2svc.MainApp.SetFocus(*grid.Members[grid.CurrentMemberInFocus]) // * hmmm
			}
		}else if event.Rune() == '?' {
			ec2svc.RootPage.DisplayHelpMessage(HELP_EC2_MAIN)
		}
		return event
	})

	// ui config
	statusBar.SetText("Status")
	table.SetBorders(false)
	table.SetSelectable(true, false) // rows: true, colums: false means select only rows
	table.Select(1, 1)
	table.SetFixed(0, 3)
	grid.SetRows(-3, -1, 2)
	grid.EAddItem(table, 0, 0, 30, 1, 0, 0, true)
	grid.EAddItem(description, 30, 0, 10, 1, 0, 0, false)
	grid.EAddItem(statusBar, 40, 0, 1, 1, 0, 0, false)
	// AddItem(item Primitive, fixedSize, proportion int, focus bool)
	ec2svc.RootPage.AddPage("Instances", ec2svc.GetMainElement(), true, false) // TODO: page names and such
}

func (svc *ec2Service) GetMainElement() tview.Primitive {
	return grid
	// return flex
}

// TODO: this should be changed
// returns the row index of given instance ids in the table
func getInstanceLocationInTable(ids []*string) []int{
    rowIndices := make([]int, len(ids))
    for rowIdx, id := range ids{
            rowIndices[rowIdx] =  -1
        for instanceIdx := 1; instanceIdx < table.GetRowCount(); instanceIdx++ {
            if instanceId := table.GetCell(instanceIdx, COL_ID).Text; instanceId == *id{
                rowIndices[rowIdx] =  instanceIdx
                break
            }
        }
    }
    return rowIndices
}
func  handleActionInstanceStatusUpdate(action common.Action){
    // TODO: should we re-check the inteface type ?
    val := action.Data.(common.InstancesStatusUpdate)
    for idx, status := range val.NewStatuses{
        rowLocation := getInstanceLocationInTable([]*string{status.InstanceId})[0]
        cell := table.GetCell(rowLocation, COL_STATE)
        oldStatus := string(status.InstanceState.Name)      // state v status
        if oldStatus != cell.Text{     // if new state
            // cell.Text = status.InstanceStatusName
            table.SetCell(rowLocation, COL_STATE, tview.NewTableCell(oldStatus).SetBackgroundColor(tcell.ColorGreen))
            fmt.Printf("changed idx %d", idx)
        }


    }

}
func (svc *ec2Service) WatchChanges() {

    go svc.Model.DispatchRoutines()
    go func(){                      // watches changes and call the action handler
        for {
        action := <-svc.Model.Channel
        switch action.Data.(type) {
        case common.InstancesStatusUpdate:
        case common.InstanceStatusUpdate:
            go handleActionInstanceStatusUpdate(action)
        }
    }
}()
}
