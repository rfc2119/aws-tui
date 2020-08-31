package ui

import (
    "fmt"

    "rfc2119/aws-tui/services"

)
func InitEC2View(){

	grid := ui.NewEgrid()
	description := tview.NewTextView()
	colNames := []string{"ID", "AMI", "Type", "State", "StateReason"}

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
			// table.SetCell(rowIdx, 1, instanceIdCell)
			// table.SetCell(rowIdx, 2, instanceAMICell)
			// table.SetCell(rowIdx, 3, instanceTypeCell)
			table.SetCell(rowIdx+1, colIdx, cell)
		}
	}
	table.SetBorders(false)
	table.SetSelectable(true, false) // rows: true, colums: false means select only rows
	table.Select(1, 1)
	table.SetFixed(0, 3)
	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			pages.SwitchToPage("EC2")
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
            fmt.Println("move to another item")     // TODO
            // for _, member := range grid.members{
            //     fmt.Printf("%#v", (*member).HasFocus())
                // box, isBox := member.(tview.Box)
                // if isBox {
                // if !box.HasFocus(){
                //     app.SetFocus(member)
                // }
            // }
            // }
		}
		return event
	})
}
