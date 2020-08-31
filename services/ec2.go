package services

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/client"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type ec2Service struct {
	Model   *ec2.EC2
	View    []viewComponent
	Channel chan Action
	Name    string // use the convenient map to assign the correct name
}

// config: the aws client config that will create the service (the underlying model)
// elm: at least one tview element that will act as the view for the service
// ch: the channel between the model and the view
func NewEC2Service(config *client.Config, elements []*tview.Primitive, ch chan Action) {

	// config, err := external.LoadDefaultAWSConfig()
	// if err != nil {
	// 	panic("unable to load SDK config, " + err.Error())
	// }
	var components []viewComponent
	for _, elm := range elements {
		viewComponent := &viewComponent{
			ID:      fmt.Sprintf("%p", elm),
			Service: ServiceNames[EC2],
			Element: elm,
		}

		components = append(components, viewComponent)
	}
	model := ec2.New(config)
	return &ec2Service{
		Model:   ec2.New(config),
		Name:    ServiceNames[EC2],
		Channel: ch,
		View:    components,
	}
}
func (svc *ec2Service) GetEC2Instances() []ec2.Reservation {

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	req := svc.Model.DescribeInstancesRequest(&ec2.DescribeInstancesInput{})
	resp, err := req.Send(context.Background()) // the background context is never canceled
	if err != nil {
		fmt.Println("error")
	}
	// fmt.Printf("%T:%#v", resp, resp)
	// spew.Dump(resp.Reservations[0].Instances[0].ImageId)
	// fmt.Println(resp.Reservations)
	// spew.Dump(resp.Reservations)
	return resp.Reservations // TODO: nextToken and maxNumber if n instances is huge
}
