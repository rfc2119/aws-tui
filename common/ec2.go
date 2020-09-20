package common

import (
    "github.com/rfc2119/simple-state-machine"
	"github.com/gdamore/tcell"
)

type stateMachine struct {
    ssm.StateMachine
    color   tcell.Color
}

func (sm *stateMachine) GetColor() tcell.Color {
    return sm.color
}

var (
    AMIFilters = []int{FILTER_ARCHITECTURE, FILTER_OWNER_ALIAS, FILTER_NAME, FILTER_PLATFORM, FILTER_ROOT_DEVICE_TYPE, FILTER_STATE}

    pendingState = ssm.State{Name: "pending"}
    runningState = ssm.State{Name: "running"}
    stoppedState = ssm.State{Name: "stopped"}
    stoppingState = ssm.State{Name: "stopping"}
    rebootingState = ssm.State{Name: "rebooting"}
    shuttingDownState = ssm.State{Name: "shutting-down"}
    terminatedState = ssm.State{Name: "terminated"}
)

func NewEC2InstancesStateMachine() *stateMachine {
    // state machine for ec2 instance lifecycle (could be done with a switch case statement)
    // reference: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html

    // triggers (see ui/ec2 for trigger names)  TODO: unify trigger names
    emptyTrigger := ssm.Trigger{Key: ""}
    startTrigger := ssm.Trigger{Key: "Start"}
    stopTrigger := ssm.Trigger{Key: "Stop"}
    hibernateTrigger := ssm.Trigger{Key: "Hibernate"}
    terminateTrigger := ssm.Trigger{Key: "Terminate"}
    rebootTrigger := ssm.Trigger{Key: "Reboot"}

    // the state machine itself (initially in the "pending" state) and configs
    EC2LifeCycle := stateMachine{
        StateMachine: *ssm.NewStateMachine(pendingState),
        color: tcell.ColorDefault,
    }
    runningConfig := EC2LifeCycle.Configure(runningState)
    stoppedConfig := EC2LifeCycle.Configure(stoppedState)
    pendingConfig := EC2LifeCycle.Configure(pendingState)
    shuttingDownConfig := EC2LifeCycle.Configure(shuttingDownState)
    stoppingConfig := EC2LifeCycle.Configure(stoppingState)
    rebootingConfig := EC2LifeCycle.Configure(rebootingState)

    // configuring the running state
    runningConfig.OnEnter(func() {EC2LifeCycle.color = tcell.ColorGreen})
    runningConfig.Permit(terminateTrigger, shuttingDownState)
    runningConfig.Permit(rebootTrigger, rebootingState)
    runningConfig.Permit(stopTrigger, stoppingState)
    runningConfig.Permit(hibernateTrigger, stoppingState)

    // configuring the running state
    stoppedConfig.OnEnter(func() {EC2LifeCycle.color = tcell.ColorRed})
    stoppedConfig.Permit(startTrigger, pendingState)
    stoppedConfig.Permit(terminateTrigger, terminatedState)

    // configuring states with an empty trigger (intermediate states)
    pendingConfig.OnEnter(func() {EC2LifeCycle.color = tcell.ColorYellow})
    shuttingDownConfig.OnEnter(func() {EC2LifeCycle.color = tcell.ColorYellow})
    stoppingConfig.OnEnter(func() {EC2LifeCycle.color = tcell.ColorYellow})
    rebootingConfig.OnEnter(func() {EC2LifeCycle.color = tcell.ColorBlue})
    pendingConfig.Permit(emptyTrigger, runningState)
    shuttingDownConfig.Permit(emptyTrigger, terminatedState)
    stoppingConfig.Permit(emptyTrigger, stoppedState)
    rebootingConfig.Permit(emptyTrigger, runningState)

    return &EC2LifeCycle
}
