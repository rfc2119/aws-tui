package common

import (

"github.com/aws/aws-sdk-go-v2/service/ec2"
// "github.com/aws/aws-sdk-go-v2/aws"
)

// v2 of the aws sdk is used

const (
	ACTION_INSTANCE_STATUS_UPDATE = iota // iota is a counter starting from zero
    ACTION_INSTANCES_STATUS_UPDATE
	ACTION_INSTANCE_HALP_ME

	// defining the services themselves as numeric constants
	// i don't know how is this useful
	SERVICE_EC2 = iota
	SERVICE_EBS
)

// convenient map *shrugs*
var ServiceNames = map[int]string{
	SERVICE_EC2: "ec2",
	SERVICE_EBS: "ebs",
}

// actions are defined for each service. right now, this is a way to define the behavior that the view should follow. in some sense, there ought to be some generalization of this instead of defining everything manually, but we'll see how things goes

// this is the generic action structure. the data field is an interface. this permits all other structures (defined next) to be passed onto the channel. on the receiving side of the work channel, each receiver should assert the type of the data field and act accordingly
type Action struct {
	Type int
	Data interface{}
}

// these are the manually defined data structures that any first party
// should expect when receiving/sending an action. these structures are the "data" field in the action
// notice the similarity in using a name similar to the action, but in camel case
type InstanceStatusUpdate ec2.InstanceStatus

// TODO: is type(InstanceStatusesUpdate) different from type([]ec2.InstanceStatus)
type InstanceStatusesUpdate []ec2.InstanceStatus
