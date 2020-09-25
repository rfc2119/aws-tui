package common

import (
    "github.com/rfc2119/simple-state-machine"
	"github.com/gdamore/tcell"
)

type EStateMachine struct {
    ssm.StateMachine
    color   tcell.Color             // On transitioning to a new state, a new color is set
    emptyTrigger ssm.Trigger        // for state machines which have intermediate states
}

func (sm *EStateMachine) GetColor() tcell.Color { return sm.color }
func (sm *EStateMachine) GetEmptyTrigger() ssm.Trigger { return sm.emptyTrigger }


var (
    AMIFilters = []int{FILTER_ARCHITECTURE, FILTER_OWNER_ALIAS, FILTER_NAME, FILTER_PLATFORM, FILTER_ROOT_DEVICE_TYPE, FILTER_STATE}

    // states of an EC2 instance life cycle
    pendingState = ssm.State{Name: "pending"}
    runningState = ssm.State{Name: "running"}
    stoppedState = ssm.State{Name: "stopped"}
    stoppingState = ssm.State{Name: "stopping"}
    rebootingState = ssm.State{Name: "rebooting"}
    shuttingDownState = ssm.State{Name: "shutting-down"}
    terminatedState = ssm.State{Name: "terminated"}

    // states of an EBS volume life cycle
    attachedState = ssm.State{Name: "in-use"}
    detachedState = ssm.State{Name: "available"}
    deletingState = ssm.State{Name: "deleting"}     // next state is "deleted", which if reached, then the API for listVolumes will return nothing.
    creatingState = ssm.State{Name: "creating"}

)

func NewEBSVolumeStateMachine() *EStateMachine {
    // State machine for the life cycle of an EBS Volume
    // Triggers (see ui/ec2 for trigger names)  TODO: unify trigger names
    emptyTrigger := ssm.Trigger{Key: " "}       // transition to intermediate states
    attachTrigger := ssm.Trigger{Key: "Attach"}
    detachTrigger := ssm.Trigger{Key: "Detach"}
    forceDetachTrigger := ssm.Trigger{Key: "Force Detach"}
    deleteTrigger := ssm.Trigger{Key: "Delete"}

    // The state machine itself (initially in the "available" state) and configs
    EBSLifeCycle := EStateMachine {
        StateMachine: *ssm.NewStateMachine(detachedState),
        color: tcell.ColorDefault,
        emptyTrigger: emptyTrigger,
    }
    inUseConfig := EBSLifeCycle.Configure(attachedState)
    availableConfig := EBSLifeCycle.Configure(detachedState)
    deletingConfig := EBSLifeCycle.Configure(deletingState)
    creatingConfig := EBSLifeCycle.Configure(creatingState)

    inUseConfig.OnEnter(func() {EBSLifeCycle.color = tcell.ColorGreen})
    inUseConfig.Permit(detachTrigger, detachedState)
    inUseConfig.Permit(forceDetachTrigger, detachedState)

    availableConfig.OnEnter(func() {EBSLifeCycle.color = tcell.ColorBlue})
    availableConfig.Permit(attachTrigger, attachedState)
    availableConfig.Permit(deleteTrigger, deletingState)

    deletingConfig.OnEnter(func() {EBSLifeCycle.color = tcell.ColorRed})
    creatingConfig.OnEnter(func() {EBSLifeCycle.color = tcell.ColorYellow})
    creatingConfig.Permit(emptyTrigger, attachedState)

    return &EBSLifeCycle
}

func NewEC2InstancesStateMachine() *EStateMachine {
    // State machine for ec2 instance lifecycle (could be done with a switch case statement)
    // reference: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html

    // Triggers (see ui/ec2 for trigger names)  TODO: unify trigger names
    emptyTrigger := ssm.Trigger{Key: " "}       // transition to intermediate states
    startTrigger := ssm.Trigger{Key: "Start"}
    stopTrigger := ssm.Trigger{Key: "Stop"}
    stopForceTrigger := ssm.Trigger{Key: "Stop (Force)"}
    hibernateTrigger := ssm.Trigger{Key: "Hibernate"}
    terminateTrigger := ssm.Trigger{Key: "Terminate"}
    rebootTrigger := ssm.Trigger{Key: "Reboot"}

    // the state machine itself (initially in the "pending" state) and configs
    EC2LifeCycle := EStateMachine{
        StateMachine: *ssm.NewStateMachine(pendingState),
        color: tcell.ColorDefault,
        emptyTrigger: emptyTrigger,          // if there's no trigger defined, set trigger.Key to ""
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
    runningConfig.Permit(stopForceTrigger, stoppingState)
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
