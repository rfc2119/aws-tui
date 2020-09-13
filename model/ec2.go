package model

import (
	"context"
	// "fmt"
	"time"
	"log"
	"rfc2119/aws-tui/common"


	// "github.com/apex/log"
	// "github.com/apex/log/handlers/cli"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// each channel is concerned with a service, and only the view to the model may use the channel. for example, for a designated ec2 worker channel, only the view responsible for ec2 may listen to the channel and consume items
// this might break in the future. sometimes, multiple benefeciaries exist for a single work. for example, when deleting an ebs volume, the ec2 console should also make use of the deletion command/action to update the affected instance. i don't know how to approach this (yet)
// TODO: it doens't make sense to export the type and have a New() function in the same time, right ?
type EC2Model struct {
	model   *ec2.Client
	// watchers []watcher		// any watcher should register itself here; a watcher is a routine that polls for changes from the model
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

func (mdl *EC2Model) GetEC2Instances() []ec2.Reservation {

	req := mdl.model.DescribeInstancesRequest(&ec2.DescribeInstancesInput{})
	resp, err := req.Send(context.Background()) // the background context is never canceled
	if err != nil {                             // TODO: recover, as this get us a segfault when request fails (maybe return an empty reservation ?)
		log.Println(err)
	}
	// fmt.Printf("%T:%#v", resp, resp)
	// spew.Dump(resp.Reservations[0].Instances[0].ImageId)
	// fmt.Println(resp.Reservations)
	// spew.Dump(resp.Reservations)
	return resp.Reservations // TODO: nextToken and maxNumber if n instances is huge (use pagination)
}

// lists all instance types offered
func (mdl *EC2Model) ListOfferings() []ec2.InstanceTypeOffering{	// TODO: region, filters
	req := mdl.model.DescribeInstanceTypeOfferingsRequest(&ec2.DescribeInstanceTypeOfferingsInput{})
	resp, err := req.Send(context.TODO())
	if err != nil {			// TODO: graceful error handling
		log.Println(err)
	}
	return resp.InstanceTypeOfferings	// TODO: paginator
}

// DispatchWatchers sets the appropriate timer and calls each watcher
func (mdl *EC2Model) DispatchWatchers() {
	ticker := time.NewTicker(5 * time.Second)	// TODO: 5
	// i parametrized the goroutine only to make the channel send only
	go func(t *time.Ticker, ch chan<- common.Action, client *ec2.Client){		// dispatcher goroutine
		for {
			<-t.C
			watcher1(client, ch, true)
			// log.Println("watcher sent data")
		}
	}(ticker, mdl.Channel, mdl.model)
}

// What the duck ? DescribeInstanceStatus (the API itself) requires that an instance be in the running state
// TODO: not sure this is the best way to do this
func watcher1(client *ec2.Client, ch chan<- common.Action, describeAll bool){
	// mdl.watchers = append(mdl.watchers, watcher)

	req := client.DescribeInstanceStatusRequest(&ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: &describeAll,
	})
	resp, err := req.Send(context.TODO())
	if err != nil {
		log.Println(err)
	}
	sendMe := common.Action{Type: common.ACTION_INSTANCE_STATUS_UPDATE, Data: common.InstanceStatusesUpdate(resp.InstanceStatuses)}	// TODO: paginator
	ch <- sendMe
}
// type watcher interface{
// 	register(mdl *EC2Model)		// watcher registers self in a model instance
// 	run()				// go watch
// }
