package model

import (
	"context"
	"log"
    "strings"
	"rfc2119/aws-tui/common"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// each channel is concerned with a service, and only the view to the model may use the channel. for example, for a designated ec2 worker channel, only the view responsible for ec2 may listen to the channel and consume items
// this might break in the future. sometimes, multiple benefeciaries exist for a single work. for example, when deleting an ebs volume, the ec2 console should also make use of the deletion command/action to update the affected instance. i don't know how to approach this (yet)
type EC2Model struct {
	model *ec2.Client
	Channel chan common.Action // channel from model to view (see above)
	Name    string             // use the convenient map to assign the correct name
	// logger	log.Logger

}

func NewEC2Model(config aws.Config) *EC2Model {
	return &EC2Model{
		model:   ec2.New(config),
		Name:    common.ServiceNames[common.SERVICE_EC2],
		Channel: make(chan common.Action),
	}
}

func (mdl *EC2Model) StartEC2Instance(instanceIds []string) []ec2.InstanceStateChange{

	req := mdl.model.StartInstancesRequest(&ec2.StartInstancesInput{
        InstanceIds: instanceIds,
    })
	resp, err := req.Send(context.TODO())
    printAWSError(err)      // TODO: graceful error handling
	// return resp.InstanceTypeOfferings // TODO: paginator
    return resp.StartingInstances
}

func (mdl *EC2Model) StopEC2Instance(instanceIds []string, force, hibernate bool)[]ec2.InstanceStateChange {
	req := mdl.model.StopInstancesRequest(&ec2.StopInstancesInput{
        InstanceIds: instanceIds,
        Hibernate: aws.Bool(hibernate),
        Force: aws.Bool(force),
    })
	resp, err := req.Send(context.TODO())
    printAWSError(err)      // TODO: graceful error handling
	// return resp.InstanceTypeOfferings // TODO: paginator
    return resp.StoppingInstances

}
func (mdl *EC2Model) RebootEC2Instance(instanceIds []string){
	req := mdl.model.RebootInstancesRequest(&ec2.RebootInstancesInput{
        InstanceIds: instanceIds,
    })
	_, err := req.Send(context.TODO())
    printAWSError(err)      // TODO: graceful error handling
}
func (mdl *EC2Model) TerminateEC2Instance(instanceIds []string) []ec2.InstanceStateChange {

	req := mdl.model.TerminateInstancesRequest(&ec2.TerminateInstancesInput{
        InstanceIds: instanceIds,
    })
	resp, err := req.Send(context.TODO())
    printAWSError(err)      // TODO: graceful error handling
    return resp.TerminatingInstances
}
func (mdl *EC2Model) GetEC2Instances() []ec2.Instance {

    req := mdl.model.DescribeInstancesRequest(&ec2.DescribeInstancesInput{})
    paginator := ec2.NewDescribeInstancesPaginator(req)
    var instances []ec2.Instance
    for paginator.Next(context.TODO()){
        for _, reservation := range paginator.CurrentPage().Reservations {
            instances = append(instances, reservation.Instances...)
        }
    }

    printAWSError(paginator.Err())      // TODO: graceful error handling
	return instances
}

// lists all instance types offered
func (mdl *EC2Model) ListOfferings() []ec2.InstanceTypeOffering { // TODO: region, filters
	req := mdl.model.DescribeInstanceTypeOfferingsRequest(&ec2.DescribeInstanceTypeOfferingsInput{})
	resp, err := req.Send(context.TODO())
	if err != nil { // TODO: graceful error handling
		log.Println(err)
	}
	return resp.InstanceTypeOfferings // TODO: paginator
}

// lists AMIs offered
func (mdl *EC2Model) ListAMIs(filterMap map[string]string) []ec2.Image {
	// TODO: assert length
	var filters []ec2.Filter
	for filterName, filterValue := range filterMap {
		// if filterValue != ""
			filters = append(filters, ec2.Filter{Name: aws.String(filterName), Values: strings.Split(filterValue, ",")})
		// }
	}
	req := mdl.model.DescribeImagesRequest(&ec2.DescribeImagesInput{Filters: filters})
	resp, err := req.Send(context.TODO())
    printAWSError(err)      // TODO: graceful error handling

	return resp.Images
}

func (mdl *EC2Model) ListVolumes() []ec2.Volume {
	var (
        volumes []ec2.Volume
    )
	req := mdl.model.DescribeVolumesRequest(&ec2.DescribeVolumesInput{})
    paginator := ec2.NewDescribeVolumesPaginator(req)
    for paginator.Next(context.TODO()){
        volumes = append(volumes, paginator.CurrentPage().Volumes...)
    }

    printAWSError(paginator.Err())      // TODO: graceful error handling
	return volumes
}
// TODO: i don't understand why this is better than a simple print
func printAWSError(err error) error {
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            default:
                log.Println(aerr.Error())
            }
        } else {
        // Print the error, cast err to awserr.Error to get the Code and
        // Message from an error.
        log.Println(err.Error())
        }
    }
    return err
}
// DispatchWatchers sets the appropriate timer and calls each watcher
func (mdl *EC2Model) DispatchWatchers() {
	ticker := time.NewTicker(5 * time.Second) // TODO: 5
	// i parametrized the goroutine only to make the channel send only
	go func(t *time.Ticker, ch chan<- common.Action, client *ec2.Client) { // dispatcher goroutine
		for {
			<-t.C
			watcher1(client, ch, true)
			// log.Println("watcher sent data")
		}
	}(ticker, mdl.Channel, mdl.model)
}

// TODO: not sure this is the best way to do this watcher/listner combo
func watcher1(client *ec2.Client, ch chan<- common.Action, describeAll bool) {
	// mdl.watchers = append(mdl.watchers, watcher)

	req := client.DescribeInstanceStatusRequest(&ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: &describeAll,
	})
	resp, err := req.Send(context.TODO())
    printAWSError(err)      // TODO: graceful error handling
	sendMe := common.Action{Type: common.ACTION_INSTANCE_STATUS_UPDATE, Data: common.InstanceStatusesUpdate(resp.InstanceStatuses)} // TODO: paginator
	ch <- sendMe
}
