package ui

import (
	"fmt"
	"log"
    "time"
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
	// instance table column names and indicies
	COL_EC2_ID = iota
	COL_EC2_AMI
	COL_EC2_TYPE
	COL_EC2_STATE
	// COL_EC2_STATEREASON
)
const (
    // volumes table column names and indicies
	COL_EBS_ID = iota
	COL_EBS_SIZE
	COL_EBS_TYPE
	COL_EBS_IOPS
	COL_EBS_STATE
)
const (
	HELP_EC2_MAIN = `
	d               Describe instance
    r               Manually refresh list of instances
	e               Edit instance (WIP)
    ^l              Filter and List AMIs
	`
	HELP_EBS_MAIN = `
    r               Refresh list of volumes
	e               Edit volume (WIP)
	`
)

// global ui elements (TODO: perhaps i should make them local somehow)
var (
	instancesFlex             *eFlex                // container for the main EC2 page
	description               = tview.NewTextView() // instance description
	instancesTable            = NewEtable()         // instance status as in web ui
	volumesTable              = NewEtable()
	volumesFlex               *eFlex
	editVolumesGrid           *eGrid
	editInstancesGrid         *eGrid
	instanceOfferingsDropdown = tview.NewDropDown()
	instanceStatusRadioButton = NewRadioButtons([]string{"Start", "Stop", "Stop (Force)", "Hibernate", "Reboot", "Terminate"}) // all buttons are enabled by default
	EC2InstancesStateMachine  = common.NewEC2InstancesStateMachine()
)

type ec2Service struct {
	mainUI
	Model *model.EC2Model
	// logger log.Logger        // TODO

	// service specific data
	instances []ec2.Instance
	volumes   []ec2.Volume
}

func NewEC2Service(config aws.Config, app *tview.Application, rootPage *ePages, statBar *StatusBar) *ec2Service {

	return &ec2Service{
		mainUI: mainUI{
			MainApp:   app,
			RootPage:  rootPage,
			StatusBar: statBar,
		},
		Model:     model.NewEC2Model(config),
		instances: nil,
		volumes:   nil,
	}
}

func (ec2svc *ec2Service) InitView() {

	// Hacks
	instancesFlex = NewEFlex(ec2svc.RootPage)
	editInstancesGrid = NewEgrid(ec2svc.RootPage)
	volumesFlex = NewEFlex(ec2svc.RootPage)
	editVolumesGrid = NewEgrid(ec2svc.RootPage)

	ec2svc.drawElements()
	ec2svc.setCallbacks()

	// Configuration for ui elements
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
    editInstancesGrid.SetColumns(1, 0, 0, 1)
     editInstancesGrid.SetRows(1, 10, 1)
	editInstancesGrid.EAddItem(instanceStatusRadioButton, 1, 1, 1, 1, 0, 0, true)
	editInstancesGrid.EAddItem(instanceOfferingsDropdown, 1, 2, 1, 1, 0, 0, false)

	ec2svc.RootPage.EAddPage("Instances", instancesFlex, true, false)         // TODO: page names and such; resize=true, visible=false
	// ec2svc.RootPage.EAddPage("Edit Instance", editInstancesGrid, true, false) // TODO: page names and such
	ec2svc.RootPage.EAddPage("Volumes", volumesFlex, true, false)             // TODO: page names and such; resize=true, visible=false
	ec2svc.RootPage.EAddPage("Edit Volume", editVolumesGrid, true, false)     // TODO: page names and such

	ec2svc.WatchChanges()

}

// Fills ui elements with appropriate initial data
func (ec2svc *ec2Service) drawElements() {
	// Draw main instancesTable
	ec2svc.fillInstancesTable()
	ec2svc.fillVolumesTable()

    // Instance types allowed in current region
	offerings := ec2svc.Model.ListOfferings()
	opts := make([]string, len(offerings))
	for idx := 0; idx < len(offerings); idx++ {
		opts[idx] = string(offerings[idx].InstanceType)
	}
	instanceOfferingsDropdown.SetOptions(opts, nil)
}

// Set function callbacks for different ui elements
func (ec2svc *ec2Service) setCallbacks() {

	// instancesTable
	instancesTableCallbacks := map[tcell.Key]func(){
		tcell.Key('d'): func() {
			row, _ := instancesTable.GetSelection() // TODO: multi selection
			description.SetText(fmt.Sprintf("%v", ec2svc.instances[row-1]))
		},
		tcell.Key('e'): func() {
			// Configuring the state radio button
			row, _ := instancesTable.GetSelection() // TODO: multi selection
			editInstanceStatusRadioButton(instancesTable.GetCell(row, COL_EC2_STATE).Text)
			// TODO: configure the "instance type" drop down
			editInstancesGrid.SetTitle(instancesTable.GetCell(row, COL_EC2_ID).Text)
            ec2svc.showGenericModal(editInstancesGrid, 40, 40)
		},
		tcell.Key('r'): func() {
			ec2svc.StatusBar.SetText("refreshing instances list")
			ec2svc.fillInstancesTable()
		},
	}
	instancesTable.UpdateKeyToFunc(instancesTableCallbacks)

	// instancesFlex container for EC2 instances table
	instancesFlexCallBacks := map[tcell.Key]func(){
		tcell.KeyCtrlL: func() { ec2svc.chooseAMIFilters() },
	}
	instancesFlex.SetShiftFocusFunc(ec2svc.MainApp)
	instancesFlex.UpdateKeyToFunc(instancesFlexCallBacks)

	editInstancesGrid.SetShiftFocusFunc(ec2svc.MainApp)

	// Radio button for instance status
	instanceStatusRadioButtonCallBacks := map[tcell.Key]func(){
		tcell.Key(' '): func() {
			currOpt := instanceStatusRadioButton.GetCurrentOptionName()
			msg := fmt.Sprintf("%s instance ?", currOpt)
			ec2svc.showConfirmationBox(msg, func() {
				// ec2svc.StatusBar.SetText(fmt.Sprintf("%#v", test))
				ec2svc.StatusBar.SetText(fmt.Sprintf("%sing instance", currOpt))
				row, _ := instancesTable.GetSelection() // TODO: multi selection
				instanceIds := []string{instancesTable.GetCell(row, COL_EC2_ID).Text}
				switch strings.ToLower(currOpt) { // TODO: do something w/ return value
				case "start": // TODO: magic names
					ec2svc.Model.StartEC2Instance(instanceIds)
				case "stop":
					ec2svc.Model.StopEC2Instance(instanceIds, false, false)
				case "hibernate":
					ec2svc.Model.StopEC2Instance(instanceIds, false, true) // hibernate=true
				case "stop (force)":
					ec2svc.Model.StopEC2Instance(instanceIds, true, false) // force=true
				case "reboot":
					ec2svc.Model.RebootEC2Instance(instanceIds)
				case "terminate":
					ec2svc.Model.TerminateEC2Instance(instanceIds)
				}
			})
		},
	}
	instanceStatusRadioButton.UpdateKeyToFunc(instanceStatusRadioButtonCallBacks)

	// Dropdown for instance types
    // TODO: This is by no means near ready. See https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-resize.html for all conditions
	instanceOfferingsDropdown.SetSelectedFunc(func(text string, index int) {
        msg := fmt.Sprintf("Change instance type to %s ?", text)
		ec2svc.showConfirmationBox(msg, func(){
            row, _ := instancesTable.GetSelection()
            if instancesTable.GetCell(row, COL_EC2_STATE).Text != "stopped" {
                ec2svc.StatusBar.SetText("Cannot change instance type: instance is not in the stopped state")
                return
            }
            ec2svc.Model.ChangeInstanceType(instancesTable.GetCell(row, COL_EC2_ID).Text, text)
        })
	})

	volumesTableCallBacks := map[tcell.Key]func(){
		tcell.Key('e'): func() { ec2svc.editVolumes() },
		tcell.Key('r'): func() { ec2svc.fillVolumesTable() },
	}
	volumesTable.UpdateKeyToFunc(volumesTableCallBacks)

    // The flex container holding the volumes table
	volumesFlex.SetShiftFocusFunc(ec2svc.MainApp)
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
	grid.SetBorders(true).SetTitle("HAAAAAALP")
	dropDownVolumeType.SetOptions([]string{"Magnetic (standard)", "General Purpose SSD (gp2)", "Provisioned IOPS SSD (io1)", "Provisioned IOPS SSD (io2)"}, nil)
	inputFieldVolumeIops.SetText(volumesTable.GetCell(row, COL_EBS_IOPS).Text)
	inputFieldVolumeSize.SetText(volumesTable.GetCell(row, COL_EBS_SIZE).Text)
	// dropDownVolumeType.SetIndex()       // TODO
	// grid.SetSize(2, 2, 10, 20) // numRows, numCols, rowSize, colSize
	grid.SetRows(3, 3)
	grid.SetColumns(1, 0, 0, 1)
	grid.EAddItem(dropDownVolumeType, 0, 1, 1, 1, 0, 0, true) // row, col, rowSpan, colSpan, minGridHeight, minGridWidth, focus
	grid.EAddItem(inputFieldVolumeSize, 1, 1, 1, 1, 0, 0, false)
	grid.EAddItem(radioButtonVolumeStatus, 0, 2, 2, 2, 0, 0, false)
	grid.SetShiftFocusFunc(ec2svc.MainApp)
	ec2svc.showGenericModal(grid, 50, 10)
}

// TODO: could this be a generic filter box ?
// Pops up a box to filter list of AMIs. Filters are defined in file *common/ec2*
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
		// Save current filter value if existed before or if there's new text
		if txt := inputField.GetText(); txt != "" || exists {
			ec2svc.StatusBar.SetText(fmt.Sprintf("prev text: %s, exists: %v, prevName: %s", previousText, exists, prevName))
			if prevName != "" { // avoid initial value of prevName
				filters[prevName] = txt
			}
		}
		// Set auto complete for the current selected text. copied from demos/inputfield
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
		inputField.SetText(filters[option]) // Restore value for selected option, or clear the field
		prevName = option
	}

	buttonCancelFunc := func() { ec2svc.RootPage.ESwitchToPreviousPage() }
	buttonSaveFunc := func() {
		ec2svc.StatusBar.SetText("Grabbing the list of AMIs")
		amis := ec2svc.Model.ListAMIs(filters)
		instancesTableAMI := tview.NewTable()

		// Drawing the instancesTable
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
	ec2svc.showGenericModal(form, 80, 10) // 80x10 seems good for my screen
}

// Shows a generic modal box (rather than a confirmation-only box) centered at screen
// Props to skanehira from the docker tui "docui" for this! code is at github.com/skanehira/docui
func (ec2svc *ec2Service) showGenericModal(p tview.Primitive, width, height int) {
	var centeredModal *eGrid
	// unfortunately you can't access grid's minumum width or height. what to do ?
	// if g, ok := p.(*eGrid); ok {    // TODO: grid inside centered grid correctly; tview.Grid
	//     centeredModal = g
	// log.Println("OUR GRID")
	// // trying a flex instead
	// centeredModal := NewEFlex(ec2svc.RootPage).SetFullScreen(false).AddItem(
	//     tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(p, width, 0, true),
	//     height, 0, true)
	// centeredModal.SetColumns(0, width, 0).
	//                 SetRows(0, height, 0)
	// } else {
	centeredModal = NewEgrid(ec2svc.RootPage)
	centeredModal.SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true) // focus=true
		// }
	currPageName := ec2svc.RootPage.GetCurrentPageName()
	ec2svc.RootPage.EAddAndSwitchToPage("centered modal", centeredModal, true) // resize=true
	ec2svc.RootPage.ShowPage(currPageName)                                     // redraw on top (bottom ?) of the box

}

// Shows a modal box with msg and switches back to previous page. This is useful for one-off usage (no nested boxes)
func (ec2svc *ec2Service) showConfirmationBox(msg string, doneFunc func()) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"Ok", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// selectedButtonLabel = aws.String(buttonLabel)
			if buttonLabel == "Ok" && doneFunc != nil {
				doneFunc()
			}
			ec2svc.RootPage.ESwitchToPreviousPage()
		})
	ec2svc.RootPage.EAddAndSwitchToPage("modal", modal, false)      // resize=false
	ec2svc.RootPage.ShowPage(ec2svc.RootPage.GetPreviousPageName()) // +1
}

// TODO: refactor
// Fills the table for EBS volumes with volume data
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
			if v := stringFromAWSVar(item); v != "" { // helper function
				cell := tview.NewTableCell(v)
				volumesTable.SetCell(rowIdx+1, colIdx, cell)
			} else {
				ec2svc.StatusBar.SetText(fmt.Sprintf("possible invalid converstion: %#v", item))
			} // TODO: message gets cleared on the spot
		}
	}
}

// Fills the table for EC2 instances with instance data
func (ec2svc *ec2Service) fillInstancesTable() {

    colNames := []string{"ID", "AMI", "Type", "State"} // TODO: magic
	ec2svc.instances = ec2svc.Model.GetEC2Instances()  // directly invokes a method on the model
	for firstColIdx := 0; firstColIdx < len(colNames); firstColIdx++ {
		instancesTable.SetCell(0, firstColIdx,
			tview.NewTableCell(colNames[firstColIdx]).SetAlign(tview.AlignCenter).SetSelectable(false).SetAttributes(tcell.AttrBold))
	}
	for rowIdx, instance := range ec2svc.instances {
		items := []interface{}{instance.InstanceId, instance.ImageId, instance.InstanceType, instance.State.Name}
		for colIdx, item := range items {
			cell := tview.NewTableCell(stringFromAWSVar(item))
			instancesTable.SetCell(rowIdx+1, colIdx, cell)
		}
	}
}

// Dispatches goroutines to monitor changes. Assigns listeners to each action
func (svc *ec2Service) WatchChanges() {
	svc.Model.DispatchWatchers()
	go func(ch <-chan common.Action) { // listner goroutine
		for receiveMe := range ch {
			// switch receiveMe.Data.(type){
            switch receiveMe.Type {             // TODO: is Type useful anyway ?
			case common.ACTION_INSTANCE_STATUS_UPDATE:
				go listener1(receiveMe)
			default:
				log.Printf("received invalid data of type %T", receiveMe.Data)
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
		if newState != cell.Text {
			// Hop to state newState and trigger the onEnter function (to get the correct color)
			state := ssm.State{Name: newState}
			if err := EC2InstancesStateMachine.GoToState(state, true); err != nil {
				log.Println(err)
				return
			}
			colorizeRowInTable(instancesTable, rowIdx, EC2InstancesStateMachine.GetColor())
			cell.SetText(newState)  // TODO: queue draw event
            go func(){              // TODO: this is a cheap way of clearing colors
                time.Sleep(3 * time.Second)     // TODO
                colorizeRowInTable(instancesTable, rowIdx, tcell.ColorDefault)
            }()
		}
	}
}

// TODO: enum ? func (enum SummaryStatus) MarshalValue() (string, error)
// ============ helper functions
// Given an instance ID, return the row index of the instance in instancesTable t
func rowIndexFromTable(t *eTable, instanceID string) int {
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

// Colorize a row in a given instancesTable
func colorizeRowInTable(t *eTable, row int, color tcell.Color) {
	for col := 0; col < t.GetColumnCount(); col++ {
		t.GetCell(row, col).SetBackgroundColor(color)
	}
}

// Applying DFS to return all valid next triggers from current state currState
func getNextTriggersNoEmptyTriggers(currState ssm.State, emptyTriggerKey string) []ssm.Trigger {
	var (
		ret []ssm.Trigger
		// nextStates = EC2InstancesStateMachine.GetNextStates()    // TODO: states differ from next triggers
		nextTriggers = EC2InstancesStateMachine.GetNextTriggers()
	)

	for _, nextTrig := range nextTriggers {
		// EC2InstancesStateMachine.GoToState(next, false)     // triggerOnEnter=false
		if nextTrig.Key == emptyTriggerKey { // an intermediate state!
			// if EC2InstancesStateMachine.CanFire(emptyTriggerKey){
			EC2InstancesStateMachine.Fire(nextTrig.Key, nil)
			s := EC2InstancesStateMachine.State()
			ret = getNextTriggersNoEmptyTriggers(s, emptyTriggerKey)
		} else {
			ret = append(ret, nextTrig) // TODO: no better way ?
		}
	}
	return ret
}

// Edit the instance status radio button according to current state
func editInstanceStatusRadioButton(currStateText string) {
	currState := ssm.State{Name: currStateText}
	EC2InstancesStateMachine.GoToState(currState, false)
	var allowedActions []ssm.Trigger                                        // valid next actions/triggers will be returned here
	if trig := EC2InstancesStateMachine.GetEmptyTrigger(); trig.Key != "" { // empty trigger is defined
		allowedActions = getNextTriggersNoEmptyTriggers(currState, trig.Key)
	} else {
		allowedActions = EC2InstancesStateMachine.GetNextTriggers()
	}
	for idx, optName := range instanceStatusRadioButton.GetOptions() { // TODO: urgh
		enabled := false
		for _, allowedAction := range allowedActions {
			if allowedAction.Key == optName {
				instanceStatusRadioButton.EnableOptionByIdx(idx)
				enabled = true
				break
			}
		}
		if !enabled {
			instanceStatusRadioButton.DisableOptionByIdx(idx)
		}
	}
}
