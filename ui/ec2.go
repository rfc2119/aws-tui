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
    "github.com/rfc2119/simple-state-machine"
	"github.com/rivo/tview"
)


// TODO: *pukes*
const (
    // table column names
	COL_EC2_ID = iota
	COL_EC2_AMI
	COL_EC2_TYPE
	COL_EC2_STATE
	COL_EC2_STATEREASON
)
const (
    COL_EBS_ID = iota
    COL_EBS_SIZE
    COL_EBS_TYPE
    COL_EBS_IOPS
    COL_EBS_STATE
)
const (
	HELP_EC2_MAIN = `
	?               View this help message
	d               Describe instance
    r               Refresh list of instances
	e               Edit instance (WIP)
	TAB             Move to neighboring windows
    ^l              Filter and List AMIs
	q               Move back one page (will exit this help message)
	`
	HELP_EBS_MAIN = `
	?               View this help message
    r               Refresh list of volumes
	e               Edit volume (WIP)
	TAB             Move to neighboring windows
	q               Move back one page (will exit this help message)
	`
	HELP_EC2_EDIT_INSTANCE = `
    Space           Select Option in a radio box
    TAB              Move to neighboring windows
	q               Move back one page (will exit this help message)
    `
)

// global ui elements (TODO: perhaps i should make them local somehow)
var (
    instancesFlex *eFlex            // the main container
    description = tview.NewTextView() // instance description
    instancesTable = tview.NewTable()          // instance status as in web ui
    volumesTable = tview.NewTable()
    volumesFlex *eFlex
    editVolumesGrid *eGrid
    editInstancesGrid *eGrid
    instanceOfferingsDropdown = tview.NewDropDown()
    instanceStatusRadioButton = NewRadioButtons([]string{"Start", "Stop", "Stop (Force)", "Hibernate", "Reboot", "Terminate"}) // all buttons are enabled by default
    EC2InstancesStateMachine = common.NewEC2InstancesStateMachine()
)

type ec2Service struct {
	mainUI
	Model *model.EC2Model
	// logger log.Logger

	// service specific data
	instances []ec2.Instance
    volumes []ec2.Volume
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
        volumes: nil,
	}
}

func (ec2svc *ec2Service) InitView() {

	// hacks
	instancesFlex = NewEFlex(ec2svc.RootPage)
	editInstancesGrid = NewEgrid(ec2svc.RootPage)
    volumesFlex = NewEFlex(ec2svc.RootPage)
    editVolumesGrid = NewEgrid(ec2svc.RootPage)

	ec2svc.drawElements()
	ec2svc.setCallbacks()

	// configuration for ui elements
	instancesTable.SetBorders(false)
	instancesTable.SetSelectable(true, false) // rows: true, colums: false means select only rows
	instancesTable.Select(1, 1)
	instancesTable.SetFixed(0, 2)
	volumesTable.SetBorders(false)
	volumesTable.SetSelectable(true, false) // rows: true, colums: false means select only rows
	volumesTable.Select(1, 1)
	volumesTable.SetFixed(0, 2)

	instancesFlex.HelpMessage = HELP_EC2_MAIN
    instancesFlex.SetDirection(tview.FlexColumn)
    instancesFlex.SetFullScreen(true)
    instancesFlex.EAddItem(instancesTable, 0, 2, true)
    instancesFlex.EAddItem(description, 0, 1, false)

	volumesFlex.HelpMessage = HELP_EBS_MAIN
    volumesFlex.SetDirection(tview.FlexColumn)
    volumesFlex.SetFullScreen(true)
    volumesFlex.EAddItem(volumesTable, 0, 1, true)

	instanceStatusRadioButton.SetBorder(true).SetTitle("Status")
	instanceOfferingsDropdown.SetLabel("Type")
	editInstancesGrid.HelpMessage = HELP_EC2_EDIT_INSTANCE
	editInstancesGrid.SetSize(2, 4, 0, 0) // SetSize(numRows, numColumns, rowSize, columnSize int)
	editInstancesGrid.EAddItem(instanceStatusRadioButton, 0, 0, 1, 2, 0, 0, true)
	editInstancesGrid.EAddItem(instanceOfferingsDropdown, 0, 2, 1, 2, 0, 0, false)

	ec2svc.RootPage.EAddPage("Instances", instancesFlex, true, false)         // TODO: page names and such; resize=true, visible=false
	ec2svc.RootPage.EAddPage("Edit Instance", editInstancesGrid, true, false) // TODO: page names and such
	ec2svc.RootPage.EAddPage("Volumes", volumesFlex, true, false)         // TODO: page names and such; resize=true, visible=false
	ec2svc.RootPage.EAddPage("Edit Volume", editVolumesGrid, true, false) // TODO: page names and such

	ec2svc.WatchChanges()

}

// fills ui elements with appropriate initial data
func (ec2svc *ec2Service) drawElements() {
	// draw main instancesTable
	ec2svc.fillInstancesTable()
	ec2svc.fillVolumesTable()

	offerings := ec2svc.Model.ListOfferings()
	opts := make([]string, len(offerings))
	for idx := 0; idx < len(offerings); idx++ {
		opts[idx] = string(offerings[idx].InstanceType)
	}
	instanceOfferingsDropdown.SetOptions(opts, nil)
}

// set function callbacks for different ui elements
func (ec2svc *ec2Service) setCallbacks() {

	// main instancesTable
	instancesTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			row, _ := instancesTable.GetSelection()
		switch event.Rune() {
		case 'd':
			description.SetText(fmt.Sprintf("%v", ec2svc.instances[row-1]))
		case 'e':
            // TODO: custom config for edit table here
            editInstancesGrid.SetTitle(instancesTable.GetCell(row, COL_EC2_TYPE).Text)
			ec2svc.RootPage.ESwitchToPage("Edit Instance") // TODO: page names and such
		case 'r':
			ec2svc.StatusBar.SetText("refreshing instances list")
			ec2svc.fillInstancesTable()
		}

		return event
	})

    // TODO: unify grids
	// main instancesFlex
	instancesFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {

		case tcell.KeyTab:
			if len(instancesFlex.Members) > 0 {
				instancesFlex.CurrentMemberInFocus++
				if instancesFlex.CurrentMemberInFocus == len(instancesFlex.Members) { //  instancesFlex.CurrentMemberInFocus %= len(instancesFlex.Members)
					instancesFlex.CurrentMemberInFocus = 0
				}
                ec2svc.MainApp.SetFocus(instancesFlex.Members[instancesFlex.CurrentMemberInFocus])

			}
		case tcell.KeyCtrlL:
			// build modal and let user choose AMI filters
			ec2svc.chooseAMIFilters()

		case tcell.KeyRune:
			switch event.Rune() {
			case '?':
				instancesFlex.DisplayHelp()
			case 'q':
				ec2svc.RootPage.ESwitchToPreviousPage()
				ec2svc.StatusBar.SetText("exit ec2")
			}
		}
		return event
	})

	// edit grid (TODO: copy pasta from above)
	editInstancesGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// ec2svc.StatusBar.SetText(fmt.Sprintf("grid: %v. radio button: %v. dropdown: %v", editInstancesGrid.HasFocus(), instanceStatusRadioButton.HasFocus(), instanceOfferingsDropdown.HasFocus())) // TODO
		switch event.Key() {
		case tcell.KeyTab:
			if len(editInstancesGrid.Members) > 0 {
				editInstancesGrid.CurrentMemberInFocus++
				if editInstancesGrid.CurrentMemberInFocus == len(editInstancesGrid.Members) { //  editInstancesGrid.CurrentMemberInFocus %= len(editInstancesGrid.Members)
					editInstancesGrid.CurrentMemberInFocus = 0
				}
                ec2svc.MainApp.SetFocus(instancesFlex.Members[instancesFlex.CurrentMemberInFocus])
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case '?':
				editInstancesGrid.DisplayHelp()
			case 'q':
				ec2svc.RootPage.ESwitchToPreviousPage()
			}
		}
		return event
	})

	// radio button
	instanceStatusRadioButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter || event.Rune() == ' ' {
            currOpt := instanceStatusRadioButton.GetCurrentOptionName()
            msg := fmt.Sprintf("%s instance ?", currOpt)
            ec2svc.showConfirmationBox(msg, func(){
            // ec2svc.StatusBar.SetText(fmt.Sprintf("%#v", test))
                ec2svc.StatusBar.SetText(fmt.Sprintf("%sing instance", currOpt))
                row, _ := instancesTable.GetSelection()     // TODO: multi selection
                instanceIds := []string{instancesTable.GetCell(row, COL_EC2_ID).Text}
                switch strings.ToLower(currOpt){    // TODO: do something w/ return value
                case "start":                       // TODO: magic names
                    ec2svc.Model.StartEC2Instance(instanceIds)
                case "stop":
                    ec2svc.Model.StopEC2Instance(instanceIds, false, false)
                case "hibernate":
                    ec2svc.Model.StopEC2Instance(instanceIds, false, true)  // hibernate=true
                case "stop (force)":
                    ec2svc.Model.StopEC2Instance(instanceIds, true, false)  // force=true
                case "reboot":
                    ec2svc.Model.RebootEC2Instance(instanceIds)
                case "terminate":
                    ec2svc.Model.TerminateEC2Instance(instanceIds)
                }
            })
		}
		return event
	})

	// dropdown instance offering
	instanceOfferingsDropdown.SetSelectedFunc(func(text string, index int) {
        // choice := ec2svc.showConfirmationBox(fmt.Sprintf("Change instance type to %s ?", text))
        // if choice == "Ok" {
        //     // TODO
        // }
	})

	// volumes table
	volumesTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'e':
			// ec2svc.RootPage.ESwitchToPage("Edit Volume") // TODO: page names and such
			ec2svc.StatusBar.SetText("edit volume")
            ec2svc.editVolumes()
		case 'r':
			ec2svc.StatusBar.SetText("refreshing volumes list")
			ec2svc.fillVolumesTable()
		}

		return event
	})
	volumesFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {

		case tcell.KeyTab:
			if len(volumesFlex.Members) > 0 {
				volumesFlex.CurrentMemberInFocus++
				if volumesFlex.CurrentMemberInFocus == len(volumesFlex.Members) { //  volumesFlex.CurrentMemberInFocus %= len(volumesFlex.Members)
					volumesFlex.CurrentMemberInFocus = 0
				}
                ec2svc.MainApp.SetFocus(volumesFlex.Members[volumesFlex.CurrentMemberInFocus])

			}
		case tcell.KeyRune:
			switch event.Rune() {
			case '?':
				volumesFlex.DisplayHelp()
			case 'q':
				ec2svc.RootPage.ESwitchToPreviousPage()
				ec2svc.StatusBar.SetText("exit ebs")
			}
		}
		return event
	})
}

func (ec2svc *ec2Service) editVolumes() {

    grid := NewEgrid(ec2svc.RootPage)
    dropDownVolumeType := tview.NewDropDown()
    inputFieldVolumeSize := tview.NewInputField()
    inputFieldVolumeIops := tview.NewInputField()
    radioButtonVolumeStatus := NewRadioButtons([]string{"Attach", "Detach", "Force Detach", "Delete"})

    row, _ := volumesTable.GetSelection()
    inputFieldVolumeIops.SetLabel("IOPS")
    inputFieldVolumeSize.SetLabel("Size (GiB)")
    dropDownVolumeType.SetLabel("Type")
    grid.SetBorders(true).SetTitle(volumesTable.GetCell(row, COL_EBS_ID).Text)
    dropDownVolumeType.SetOptions([]string{"Magnetic (standard)", "General Purpose SSD (gp2)", "Provisioned IOPS SSD (io1)", "Provisioned IOPS SSD (io2)"}, nil)
    inputFieldVolumeIops.SetText(volumesTable.GetCell(row, COL_EBS_IOPS).Text)
    inputFieldVolumeSize.SetText(volumesTable.GetCell(row, COL_EBS_SIZE).Text)
    // dropDownVolumeType.SetIndex()       // TODO
    grid.SetSize(50, 20, 0, 0)
    grid.EAddItem(dropDownVolumeType, 0, 0, 1, 2, 0, 0, true)   // row, col, rowSpan, colSpan, minGridHeight, minGridWidth, focus
    grid.EAddItem(inputFieldVolumeSize, 1, 0, 1, 2, 0, 0, false)
    grid.EAddItem(inputFieldVolumeSize, 2, 0, 1, 2, 0, 0, false)
    grid.EAddItem(radioButtonVolumeStatus, 0, 2, 1, 2, 0, 0, false)
	// editInstancesGrid.EAddItem(instanceStatusRadioButton, 0, 0, 1, 2, 0, 0, true)
    ec2svc.showGenericModal(volumesTable.GetCell(row, COL_EBS_ID).Text, grid, 50, 20)
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
		instancesTableAMI := tview.NewTable()

		// drawing the instancesTable
		colNames := []string{"ID", "State", "Arch", "Creation Date", "Name", "Owner ID"} // TODO
		for firstColIdx := 0; firstColIdx < len(colNames); firstColIdx++ {
			instancesTableAMI.SetCell(0, firstColIdx,
				tview.NewTableCell(colNames[firstColIdx]).SetAlign(tview.AlignCenter).SetSelectable(false).SetAttributes(tcell.AttrBold))
		}
		for rowIdx, ami := range amis {
			// ownerCell := tview.NewTableCell(*ami.ImageOwnerAlias)    // TODO:
            items := []interface{}{ami.ImageId, ami.State, ami.Architecture, ami.CreationDate, ami.Name, ami.OwnerId}
			for colIdx, item := range items {
                cell := tview.NewTableCell(stringFromAWSVar(item))
				instancesTableAMI.SetCell(rowIdx+1, colIdx, cell)
			}
		}
		instancesTableAMI.SetBorders(true)
		instancesTableAMI.SetSelectable(true, false) // rows: true, colums: false means select only rows
		instancesTableAMI.Select(1, 1)
		instancesTableAMI.SetFixed(1, 1)
		instancesTableAMI.SetDoneFunc(func(key tcell.Key) {
			ec2svc.RootPage.ESwitchToPreviousPage()
		})
		ec2svc.RootPage.AddAndSwitchToPage("AMIs", instancesTableAMI, true)
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
    ec2svc.RootPage.ShowPage(currPageName)   // redraw on top (bottom ?) of the box
}
var test = make([]byte, 5)
// shows a modal box with msg and switches back to previous page. this is useful for ont-time usage (no nested boxes)
func (ec2svc *ec2Service) showConfirmationBox(msg string, doneFunc func()) {
    // var selectedButtonLabel string
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"Ok", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// selectedButtonLabel = aws.String(buttonLabel)
            if buttonLabel == "Ok" { doneFunc() }
			ec2svc.RootPage.ESwitchToPreviousPage()
		})
	ec2svc.RootPage.EAddAndSwitchToPage("modal", modal, false) // resize=false
    ec2svc.RootPage.ShowPage(ec2svc.RootPage.GetPreviousPageName())      // +1
            // ec2svc.StatusBar.SetText(fmt.Sprintf(" %#v", selectedButtonLabel))
	// return &selectedButtonLabel
}


// TODO: refactor
func (ec2svc *ec2Service) fillVolumesTable() {
	colNames := []string{"ID", "Size (GiB)", "Type", "IOPS", "State"} // TODO
	ec2svc.volumes = ec2svc.Model.ListVolumes()
	for firstColIdx := 0; firstColIdx < len(colNames); firstColIdx++ {
		volumesTable.SetCell(0, firstColIdx,
			tview.NewTableCell(colNames[firstColIdx]).SetAlign(tview.AlignCenter).SetSelectable(false).SetAttributes(tcell.AttrBold))
	}
	for rowIdx, volume := range ec2svc.volumes {
        items := []interface{}{volume.VolumeId, volume.Size, volume.VolumeType, volume.Iops, volume.State}
		for colIdx, item := range items {
            if v := stringFromAWSVar(item); v != "" {      // helper function
                cell := tview.NewTableCell(v)
                volumesTable.SetCell(rowIdx+1, colIdx, cell)
            }else { ec2svc.StatusBar.SetText(fmt.Sprintf("possible invalid converstion: %#v", item)) }  // TODO: message gets cleared on the spot
		}
	}
}

func (ec2svc *ec2Service) fillInstancesTable() {

	colNames := []string{"ID", "AMI", "Type", "State", "StateReason"} // TODO
	ec2svc.instances = ec2svc.Model.GetEC2Instances()              // directly invokes a method on the model
	for firstColIdx := 0; firstColIdx < len(colNames); firstColIdx++ {
		instancesTable.SetCell(0, firstColIdx,
			tview.NewTableCell(colNames[firstColIdx]).SetAlign(tview.AlignCenter).SetSelectable(false).SetAttributes(tcell.AttrBold))
	}
	for rowIdx, instance := range ec2svc.instances {
		items := []interface{}{instance.InstanceId, instance.ImageId, instance.InstanceType, instance.State.Name, instance.StateReason.Message}
		for colIdx, item := range items {
            cell := tview.NewTableCell(stringFromAWSVar(item))
			instancesTable.SetCell(rowIdx+1, colIdx, cell)
		}
	}
}


// dispatches goroutines to monitor changes; assigns listeners to each action
func (svc *ec2Service) WatchChanges() {
	svc.Model.DispatchWatchers()
	go func(ch <-chan common.Action) { // listner goroutine
        for receiveMe := range ch {
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
		rowIdx := rowIndexFromTable(instancesTable, *status.InstanceId) // TODO: check for -1
		cell := instancesTable.GetCell(rowIdx, COL_EC2_STATE)
		newState := string(status.InstanceState.Name)
		// log.Printf("old state: %s cell: %s", newState, cell.Text)
		if newState != cell.Text {
            // hop to state newState and trigger the onEnter function (to get the correct color)
            state := ssm.State{Name: newState}
            if err := EC2InstancesStateMachine.GoToState(state, true); err != nil{
                log.Println(err)
                return
            }
            colorizeRowInTable(instancesTable, rowIdx, EC2InstancesStateMachine.GetColor())  // TODO: queue draw event
			cell.SetText(newState)
		}
	}
}

// TODO: enum ? func (enum SummaryStatus) MarshalValue() (string, error)
// ============ helper functions
// given an instance ID, return the row index of the instance in instancesTable t
func rowIndexFromTable(t *tview.Table, instanceID string) int {
	idx := -1
	for rowIdx := 1; rowIdx < t.GetRowCount(); rowIdx++ { // 1 because first row is for column labels
		id := t.GetCell(rowIdx, COL_EC2_ID).Text
		if instanceID == id {
			idx = rowIdx
			break
		}
	}
	return idx
}

// colorize a row in a given instancesTable
func colorizeRowInTable(t *tview.Table, row int, color tcell.Color) {
	for col := 0; col < t.GetColumnCount(); col++ {
		t.GetCell(row, col).SetBackgroundColor(color)
	}
}

// tweak the edit grid according to each instance
func modifyEditGrid(g *eGrid, instanceIdx int) {

}
