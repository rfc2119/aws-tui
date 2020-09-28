package ui

import (
	"fmt"
	"log"
	"rfc2119/aws-tui/common"
	"rfc2119/aws-tui/model"
	"strings"
	"time"

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
	COL_EBS_VOL_INSTANCE_ID = iota
	COL_EBS_VOL_STATE
	COL_EBS_VOL_DEVICE_NAME
	COL_EBS_VOL_DATA_ATTACHED
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
    HELP_EBS_EDIT_VOL =`
# Attach/Detach to/from EC2 instances
You can attach instances only available in the same AZ as the volume. You can create io1 volumes can be attached to multiple instances (see https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volumes-multi.html).

If your volume state becomes *detaching*, it is likely that you need the "Force Detach" option. As per the documentation:
    > Use this option only as a last resort to detach a volume from a failed instance, or if you are detaching a volume with the intention of deleting it. The instance doesn't get an opportunity to flush file system caches or file system metadata. If you use this option, you must perform the file system check and repair procedures
    
# Device names
On linux Devices, names /dev/sd{f-p} are valid "mount points". The kernel may rename these internally to something like /dev/xvd{f-p}. For HVM instances, /dev/sda1 or /dev/xvda is reserved for root devices.See https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/device_naming.html for a complete list of valid names and considerations

# Type, size and IOPS (if applicable)
`

)

// global ui elements (TODO: perhaps i should make them local somehow)
var (
	// Instances page
	instancesTable            = NewEtable()         // Instance table as in web UI
	instancesFlex             *eFlex                // Container for the main page
	description               = tview.NewTextView() // Instance description
	editInstancesGrid         *eGrid                // "Edit Instance" grid
	instanceOfferingsDropdown = tview.NewDropDown() // Component in "Edit Instance"
	instanceStatusRadioButton = NewRadioButtons([]string{
		"Start", "Stop", "Stop (Force)", "Hibernate", "Reboot", "Terminate",
	}) // all buttons are enabled by default. Component in "Edit Instance"
	instancesTableAMI        = NewEtable() // Table used in AMI Filtering
	EC2InstancesStateMachine = common.NewEC2InstancesStateMachine()

	// Volumes page
	volumesTable            = NewEtable()           // Main table for volumes
	volumesFlex             *eFlex                  // Container for main page
	gridEditVolume          *eGrid                  // "Edit volumes" grid
	dropDownVolumeType      = tview.NewDropDown()   // Component in "Edit volumes"
	inputFieldVolumeSize    = tview.NewInputField() // Component in "Edit volumes"
	inputFieldVolumeIops    = tview.NewInputField() // Component in "Edit volumes"
	radioButtonVolumeStatus = NewRadioButtons([]string{
		"Attach", "Detach", "Force Detach", "Delete",
	}) // Component in "Edit volumes"
	tableEditVolume        = NewEtable()
	EBSVolumesStateMachine = common.NewEBSVolumeStateMachine()
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
	gridEditVolume = NewEgrid(ec2svc.RootPage)

	ec2svc.drawElements()
	ec2svc.setCallbacks()

	// Configuration for ui elements
	instancesTable.SetBorders(false)
	instancesTable.SetSelectable(true, false) // rows: true, colums: false means select only rows
	instancesTable.Select(1, 1)
	instancesTable.SetFixed(0, 2)

	instancesFlex.HelpMessage = HELP_EC2_MAIN
	instancesFlex.SetDirection(tview.FlexColumn)
	instancesFlex.SetFullScreen(true)
	instancesFlex.EAddItem(instancesTable, 0, 2, true)
	instancesFlex.EAddItem(description, 0, 1, false)

	instanceStatusRadioButton.SetBorder(true).SetTitle("Status")
	instanceOfferingsDropdown.SetLabel("Type")
	editInstancesGrid.SetColumns(1, 0, 0, 1)
	editInstancesGrid.SetRows(1, 10, 1)
	editInstancesGrid.EAddItem(instanceStatusRadioButton, 1, 1, 1, 1, 0, 0, true)
	editInstancesGrid.EAddItem(instanceOfferingsDropdown, 1, 2, 1, 1, 0, 0, false)

	instancesTableAMI.SetBorders(true)
	instancesTableAMI.SetSelectable(true, false) // rows: true, colums: false means select only rows
	instancesTableAMI.Select(1, 1)
	instancesTableAMI.SetFixed(1, 1)

	volumesTable.SetBorders(false)
	volumesTable.SetSelectable(true, false) // rows: true, colums: false means select only rows
	volumesTable.Select(1, 1)
	volumesTable.SetFixed(0, 2)

	volumesFlex.HelpMessage = HELP_EBS_MAIN
	volumesFlex.SetDirection(tview.FlexColumn)
	volumesFlex.SetFullScreen(true)
	volumesFlex.EAddItem(volumesTable, 0, 1, true)

	inputFieldVolumeIops.SetLabel("IOPS")
	inputFieldVolumeSize.SetLabel("Size (GiB)")
	dropDownVolumeType.SetLabel("Type")
	dropDownVolumeType.SetOptions([]string{"Magnetic (standard)", "General Purpose SSD (gp2)", "Provisioned IOPS SSD (io1)", "Provisioned IOPS SSD (io2)"}, nil)
	gridEditVolume.SetBorders(true).SetTitle("Test") // Not working :(

	tableEditVolume.SetBorders(false)
	tableEditVolume.SetSelectable(true, false)
	tableEditVolume.Select(1, 1)
	tableEditVolume.SetFixed(0, 2)

    gridEditVolume.HelpMessage = HELP_EBS_EDIT_VOL
	gridEditVolume.SetRows(3, 3, 0)
	gridEditVolume.SetColumns(10, 0, 20)
	gridEditVolume.EAddItem(dropDownVolumeType, 0, 0, 1, 2, 0, 0, true) // row, col, rowSpan, colSpan, minGridHeight, minGridWidth, focus
	gridEditVolume.EAddItem(inputFieldVolumeSize, 1, 0, 1, 2, 0, 0, false)
	gridEditVolume.EAddItem(radioButtonVolumeStatus, 0, 2, 2, 1, 0, 0, false) // TODO: iops
	gridEditVolume.EAddItem(tableEditVolume, 2, 0, 1, 3, 20, 40, false)
	gridEditVolume.SetShiftFocusFunc(ec2svc.MainApp)

	ec2svc.RootPage.EAddPage("Instances", instancesFlex, true, false) // TODO: page names and such; resize=true, visible=false
	ec2svc.RootPage.EAddPage("Volumes", volumesFlex, true, false)     // TODO: page names and such; resize=true, visible=false

	ec2svc.WatchChanges()

}

// Fills ui elements with appropriate initial data
func (ec2svc *ec2Service) drawElements() {
	// Draw tables
	drawFirstRowForAllTables()
	ec2svc.fillInstancesTable()
	ec2svc.fillVolumesTable()

	// Instance types allowed in current region
	var (
		offerings []ec2.InstanceTypeOffering
		err       error
	)
	if offerings, err = ec2svc.Model.ListOfferings(); err != nil {
		ec2svc.StatusBar.SetText(err.Error())
		return
	}
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
			state := ssm.State{Name: instancesTable.GetCell(row, COL_EC2_STATE).Text}
			if err := EC2InstancesStateMachine.GoToState(state, false); err != nil {
				log.Println(err)
				return
			}
			configureRadioButton(instanceStatusRadioButton, EC2InstancesStateMachine)
			// TODO: configure the "instance type" drop down
			// editInstancesGrid.SetTitle(instancesTable.GetCell(row, COL_EC2_ID).Text) //TODO
			ec2svc.showGenericModal(editInstancesGrid, 40, 80, true)
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
			ec2svc.showConfirmationBox(msg, true, func() {
				// ec2svc.StatusBar.SetText(fmt.Sprintf("%#v", test))
				row, _ := instancesTable.GetSelection() // TODO: multi selection
				instanceIds := []string{instancesTable.GetCell(row, COL_EC2_ID).Text}
				switch strings.ToLower(currOpt) { // TODO: do something w/ return value
				case "start": // TODO: magic names
					if _, err := ec2svc.Model.StartEC2Instances(instanceIds); err != nil {
						ec2svc.showConfirmationBox(err.Error(), false, nil) // TODO: spread
						ec2svc.StatusBar.SetText(err.Error())        // TODO: logging
						return
					}
				case "stop":
					if _, err := ec2svc.Model.StopEC2Instances(instanceIds, false, false); err != nil {
						// ec2svc.StatusBar.SetText(err.Error())	// TODO: logging
						return
					}
				case "hibernate":
					if _, err := ec2svc.Model.StopEC2Instances(instanceIds, false, true); err != nil { // hibernate=true
						// ec2svc.StatusBar.SetText(err.Error())	// TODO: logging
						return
					}
				case "stop (force)":
					if _, err := ec2svc.Model.StopEC2Instances(instanceIds, true, false); err != nil { // force=true
						// ec2svc.StatusBar.SetText(err.Error())	// TODO: logging
						return
					}
				case "reboot":
					if err := ec2svc.Model.RebootEC2Instances(instanceIds); err != nil {
						// ec2svc.StatusBar.SetText(err.Error())	// TODO: logging
						return
					}
				case "terminate":
					if _, err := ec2svc.Model.TerminateEC2Instances(instanceIds); err != nil {
						// ec2svc.StatusBar.SetText(err.Error())	// TODO: logging
						return
					}
				}
				ec2svc.StatusBar.SetText(fmt.Sprintf("%sing instance", currOpt))
			})
		},
	}
	instanceStatusRadioButton.UpdateKeyToFunc(instanceStatusRadioButtonCallBacks)

	// Dropdown for instance types
	// TODO: This is by no means near ready. See https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-resize.html for all conditions
	instanceOfferingsDropdown.SetSelectedFunc(func(text string, index int) {
		msg := fmt.Sprintf("Change instance type to %s ?", text)
		ec2svc.showConfirmationBox(msg, true, func() {
			row, _ := instancesTable.GetSelection()
			if instancesTable.GetCell(row, COL_EC2_STATE).Text != "stopped" {
				ec2svc.StatusBar.SetText("Cannot change instance type: instance is not in the stopped state")
				return
			}
			if err := ec2svc.Model.ChangeInstanceType(instancesTable.GetCell(row, COL_EC2_ID).Text, text); err != nil {
				ec2svc.StatusBar.SetText(err.Error())
			}
		})
	})

	volumesTableCallBacks := map[tcell.Key]func(){
		tcell.Key('r'): func() { ec2svc.fillVolumesTable() },
		tcell.Key('e'): func() {
			// Configuring the state radio button
			row, _ := volumesTable.GetSelection()
			state := ssm.State{Name: volumesTable.GetCell(row, COL_EBS_STATE).Text}
			if err := EBSVolumesStateMachine.GoToState(state, false); err != nil {
				log.Println(err)
				return
			}
			configureRadioButton(radioButtonVolumeStatus, EBSVolumesStateMachine)
			inputFieldVolumeIops.SetText(volumesTable.GetCell(row, COL_EBS_IOPS).Text) //TODO: put in main screen
			inputFieldVolumeSize.SetText(volumesTable.GetCell(row, COL_EBS_SIZE).Text)
			// dropDownVolumeType.SetIndex()       // TODO
			// TODO: generalize table drawing
			for rowIdx, info := range ec2svc.volumes[row-1].Attachments {
				items := []interface{}{info.InstanceId, info.State, info.Device, info.AttachTime}
				for colIdx, item := range items {
					cell := tview.NewTableCell(stringFromAWSVar(item))
					tableEditVolume.SetCell(rowIdx+1, colIdx, cell)
				}
			}
			tableEditVolume.SetSelectionChangedFunc(func(row, col int) {
				if row <= 1 {
					return
				}
				state := ssm.State{Name: volumesTable.GetCell(row, COL_EBS_STATE).Text}
				if err := EBSVolumesStateMachine.GoToState(state, false); err != nil {
					log.Println(err)
					return
				}
				configureRadioButton(radioButtonVolumeStatus, EBSVolumesStateMachine)
			})
			ec2svc.showGenericModal(gridEditVolume, 50, 20, true)
		},
	}
	volumesTable.UpdateKeyToFunc(volumesTableCallBacks)
	radioButtonVolumeStatusCallbacks := map[tcell.Key]func(){
		tcell.Key(' '): func() {
			currOpt := radioButtonVolumeStatus.GetCurrentOptionName()
			msg := fmt.Sprintf("%s volume ?", currOpt)
			ec2svc.showConfirmationBox(msg, true, func() {
				// ec2svc.StatusBar.SetText(fmt.Sprintf("%#v", test))
				row, _ := tableEditVolume.GetSelection()
				detachedInstanceId := tableEditVolume.GetCell(row, COL_EBS_VOL_INSTANCE_ID).Text
				detachedDeviceName := tableEditVolume.GetCell(row, COL_EBS_VOL_DEVICE_NAME).Text
				row, _ = volumesTable.GetSelection()
				volId := volumesTable.GetCell(row, COL_EBS_ID).Text
				az := aws.StringValue(ec2svc.volumes[row-1].AvailabilityZone)
				switch currOpt {
				case "Attach":
					// ec2svc.showConfirmationBox("TODO!", true, nil)
					form := tview.NewForm()
					inputFieldInstanceId := tview.NewInputField().SetLabel("Instance ID")
					inputFieldDeviceName := tview.NewInputField().SetLabel("Device Name")
                    inputFieldInstanceId.SetFieldWidth(19)  // TODO: acceptance func
                    inputFieldDeviceName.SetFieldWidth(10)
					buttonAttachFunc := func() {
						if _, err := ec2svc.Model.AttachVolume(volId, inputFieldInstanceId.GetText(), inputFieldDeviceName.GetText()); err != nil {
							ec2svc.showConfirmationBox(err.Error(), true, nil)
						}
					}
					// Set auto complete for the current selected text. copied from demos/inputfield
					inputFieldInstanceId.SetAutocompleteFunc(func(currentText string) (entries []string) {
						if len(currentText) == 0 {
							return
						}
						for _, inst := range ec2svc.instances {
							if strings.HasPrefix(*inst.InstanceId, currentText) {
								entries = append(entries, *inst.InstanceId)
							}
						}
						if len(entries) < 1 {
							entries = nil
						}
						return
					})
					form.AddButton("Attach", buttonAttachFunc)
					form.AddFormItem(inputFieldInstanceId)
					form.AddFormItem(inputFieldDeviceName)
					form.SetTitle(fmt.Sprintf("%s in %s", volId, az))
					form.SetBorder(true)
					ec2svc.showGenericModal(form, 50, 15, false)
				case "Force Detach":
					if _, err := ec2svc.Model.DetachVolume(volId, detachedInstanceId, detachedDeviceName, false); err != nil {
						ec2svc.showConfirmationBox(err.Error(), true, nil)
					}
				case "Detach":
					if _, err := ec2svc.Model.DetachVolume(volId, detachedInstanceId, detachedDeviceName, false); err != nil {
						ec2svc.showConfirmationBox(err.Error(), true, nil)
					}
				case "Delete": // TODO
					ec2svc.showConfirmationBox("TODO!", true, nil)
				}
			})

		},
	}
	radioButtonVolumeStatus.UpdateKeyToFunc(radioButtonVolumeStatusCallbacks)

	// The flex container holding the volumes table
	volumesFlex.SetShiftFocusFunc(ec2svc.MainApp)
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
		var (
			amis []ec2.Image
			err  error
		)
		if amis, err = ec2svc.Model.ListAMIs(filters); err != nil {
			ec2svc.StatusBar.SetText(err.Error()) // TODO: alert
			return
		}

		// TODO: generalize table drawing
		for rowIdx, ami := range amis {
			// ownerCell := tview.NewTableCell(*ami.ImageOwnerAlias)    // TODO:
			items := []interface{}{ami.ImageId, ami.State, ami.Architecture, ami.CreationDate, ami.Name, ami.OwnerId}
			for colIdx, item := range items {
				cell := tview.NewTableCell(stringFromAWSVar(item))
				instancesTableAMI.SetCell(rowIdx+1, colIdx, cell)
			}
		}
		ec2svc.RootPage.AddAndSwitchToPage("AMIs", instancesTableAMI, true)
	}

	inputField = tview.NewInputField().SetLabel("Filter Value")
	form.AddDropDown("Filter Name", filterNames, 0, dropDownSelectedFunc)
	form.AddButton("Save", buttonSaveFunc)
	form.AddButton("Cancel", buttonCancelFunc)
	form.AddFormItem(inputField)
	form.SetTitle("Filter AMIs").SetBorder(true)
	ec2svc.showGenericModal(form, 80, 10, true) // 80x10 seems good for my screen
}

// TODO: refactor
// Fills the table for EBS volumes with volume data
func (ec2svc *ec2Service) fillVolumesTable() {
	var err error
	if ec2svc.volumes, err = ec2svc.Model.ListVolumes(); err != nil {
		ec2svc.StatusBar.SetText(err.Error())
	}
	for rowIdx, volume := range ec2svc.volumes {
		items := []interface{}{volume.VolumeId, volume.Size, volume.VolumeType, volume.Iops, volume.State}
		for colIdx, item := range items {
			cell := tview.NewTableCell(stringFromAWSVar(item))
			volumesTable.SetCell(rowIdx+1, colIdx, cell)
		}
	}
}

// Fills the table for EC2 instances with instance data
func (ec2svc *ec2Service) fillInstancesTable() {
	var err error

	if ec2svc.instances, err = ec2svc.Model.GetEC2Instances(); err != nil { // directly invokes a method on the model
		ec2svc.StatusBar.SetText(err.Error())
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
			switch receiveMe.Type { // TODO: is Type useful anyway ?
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
			cell.SetText(newState) // TODO: queue draw event
			go func() {            // TODO: this is a cheap way of clearing colors
				time.Sleep(3 * time.Second) // TODO
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

// Applying DFS to return all valid next triggers from state machine sm
func getNextTriggersNoEmptyTriggers(sm *common.EStateMachine, emptyTriggerKey string) []ssm.Trigger {
	var (
		ret []ssm.Trigger
		// nextStates = EC2InstancesStateMachine.GetNextStates()    // TODO: states differ from next triggers
		nextTriggers = sm.GetNextTriggers()
	)

	for _, nextTrig := range nextTriggers {
		// EC2InstancesStateMachine.GoToState(next, false)     // triggerOnEnter=false
		if nextTrig.Key == emptyTriggerKey { // An intermediate state!
			// if EC2InstancesStateMachine.CanFire(emptyTriggerKey){
			sm.Fire(nextTrig.Key, nil) // Fire the empty trigger
			ret = getNextTriggersNoEmptyTriggers(sm, emptyTriggerKey)
		} else {
			ret = append(ret, nextTrig) // TODO: No better way ?
		}
	}
	return ret
}

// Configures a radio button (whose options are names of state actions) according to current state in a state machine sm
func configureRadioButton(rButton *RadioButtons, sm *common.EStateMachine) {
	// currState := ssm.State{Name: currStateText}
	// EC2InstancesStateMachine.GoToState(currState, false)
	var allowedActions []ssm.Trigger                  // Valid next actions/triggers will be returned here
	if trig := sm.GetEmptyTrigger(); trig.Key != "" { // Empty trigger is defined. beware that "" is not a key
		allowedActions = getNextTriggersNoEmptyTriggers(sm, trig.Key)
	} else {
		allowedActions = sm.GetNextTriggers()
	}
	for idx, optName := range rButton.GetOptions() { // TODO: urgh
		enabled := false
		for _, allowedAction := range allowedActions {
			if allowedAction.Key == optName {
				rButton.EnableOptionByIdx(idx)
				enabled = true
				break
			}
		}
		if !enabled {
			rButton.DisableOptionByIdx(idx)
		}
	}
}

// Draws the first row in all tables
func drawFirstRowForAllTables() { // Hardcoded all the way (TODO?)
	tables := []*eTable{tableEditVolume, volumesTable, instancesTable, instancesTableAMI}
	allColNames := [][]string{
		[]string{"Instance ID", "State", "Device", "Date Attached"},
		[]string{"ID", "Size (GiB)", "Type", "IOPS", "State"},
		[]string{"ID", "AMI", "Type", "State"},
		[]string{"ID", "State", "Arch", "Creation Date", "Name", "Owner ID"},
	}

	for idx := 0; idx < len(tables); idx++ {
		func(table *eTable, colNames []string) {
			for firstColIdx := 0; firstColIdx < len(colNames); firstColIdx++ {
				table.SetCell(0, firstColIdx,
					tview.NewTableCell(colNames[firstColIdx]).SetAlign(tview.AlignCenter).SetSelectable(false).SetAttributes(tcell.AttrBold))
			}
		}(tables[idx], allColNames[idx])
	}
}
