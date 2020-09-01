package model

import (
	"context"
	"fmt"
	"rfc2119/aws-tui/common"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/aws"
)

// TODO: it doens't make sense to export the type and have a New() function in the same time
type EC2Model struct {
	model *ec2.Client
	Channel chan common.Action	// channel from model to view
	Name    string // use the convenient map to assign the correct name
}
func NewEC2Model(config aws.Config) *EC2Model{
	return &EC2Model{
		model: ec2.New(config),
		Name: common.ServiceNames[common.SERVICE_EC2],
		Channel: make(chan common.Action),
	}
}

func (svc *EC2Model) GetEC2Instances() []ec2.Reservation {

	req := svc.model.DescribeInstancesRequest(&ec2.DescribeInstancesInput{})
	resp, err := req.Send(context.Background()) // the background context is never canceled
	if err != nil {		// TODO: recover, as this get us a segfault when request fails (maybe return an empty reservation ?)
		fmt.Println(err)
	}
	// fmt.Printf("%T:%#v", resp, resp)
	// spew.Dump(resp.Reservations[0].Instances[0].ImageId)
	// fmt.Println(resp.Reservations)
	// spew.Dump(resp.Reservations)
	return resp.Reservations // TODO: nextToken and maxNumber if n instances is huge
}
