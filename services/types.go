package services

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/rivo/tview"
)

// v2 of the aws sdk is used

// services themselves are a way to group a model (the backend sdk) and the corresponding view. i don't know what will be the view as of this moment, but here goes nothing
// each service has a structure defined in the corresponding .go file

type viewComponent struct {
	ID      string           // unique id for the component; assigned as the address of the actual ui element
	Service string           // which service does this component serve ? see below for defintion of services
	Element *tview.Primitive // the ui element itself
}

// actions are defined for each service. right now, this is a way to define the behavior that the view should follow. in some sense, there ought to be some generalization of this instead of defining everything manually, but we'll see how things goes

const (
	// first, actions to be used

	// with the EC2 service
	ACTION_INSTANCE_STATUS_UPDATE = iota // iota is a counter starting from zero
	ACTION_INSTANCE_HALP_ME

	// ===================
	// second, defining the services themselves as numeric constants
	// i don't know how is this useful
	SERVICE_EC2
	SERVICE_EBS
	SERVICE_HALP
)

// convenient map *shrugs*
var ServiceNames = map[int]string{
	EC2:  "ec2",
	EBS:  "ebs",
	HALP: "halp",
}

// this is the generic action structure. the data field is an interface. this permits all other structures (defined next) to be passed onto the channel. on the receiving side of the work channel, each receiver should assert the type of the data field and act accordingly
type Action struct {
	Type int
	Data       interface{}
}

// these are the manually defined data structures that any first party (e.g the model)
// should expect when receiving/sending an action. these structures are the "data" field in the action
// notice the similarity in using a name similar to the action, but in camel case
type InstanceStatusUpdate struct {
	NewStatus int // TODO: an ec2 structure
	// TODO: what else do we need ?
}
