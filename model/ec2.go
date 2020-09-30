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

// Each channel is concerned with a service, and only the view to the model may use the channel. for example, for a designated ec2 worker channel, only the view responsible for ec2 may listen to the channel and consume items
// this might break in the future. sometimes, multiple benefeciaries exist for a single work. for example, when deleting an ebs volume, the ec2 console should also make use of the deletion command/action to update the affected instance. i don't know how to approach this (yet)
type EC2Model struct {
	model *ec2.Client
	Channel chan common.Action // Unbuffered channel from model to view (see above)
	Name    string             // Use the convenient map to assign the correct name
	// logger	log.Logger

}

func NewEC2Model(config aws.Config) *EC2Model {
	return &EC2Model{
		model:   ec2.New(config),
		Name:    common.ServiceNames[common.SERVICE_EC2],
		Channel: make(chan common.Action),
	}
}

func (mdl *EC2Model) StartEC2Instances(instanceIds []string) ([]ec2.InstanceStateChange, error){

	req := mdl.model.StartInstancesRequest(&ec2.StartInstancesInput{
        InstanceIds: instanceIds,
    })
	resp, err := req.Send(context.TODO())
    if err != nil {
        return nil, err
    }
    return resp.StartingInstances, err
}

func (mdl *EC2Model) StopEC2Instances(instanceIds []string, force, hibernate bool) ([]ec2.InstanceStateChange, error) {
	req := mdl.model.StopInstancesRequest(&ec2.StopInstancesInput{
        InstanceIds: instanceIds,
        Hibernate: aws.Bool(hibernate),
        Force: aws.Bool(force),
    })
	resp, err := req.Send(context.TODO())
    if err != nil {
        return nil, err
    }
    return resp.StoppingInstances, err

}
func (mdl *EC2Model) RebootEC2Instances(instanceIds []string) error {
	req := mdl.model.RebootInstancesRequest(&ec2.RebootInstancesInput{
        InstanceIds: instanceIds,
    })
	_, err := req.Send(context.TODO())
    return err
}
func (mdl *EC2Model) TerminateEC2Instances(instanceIds []string) ([]ec2.InstanceStateChange, error){

	req := mdl.model.TerminateInstancesRequest(&ec2.TerminateInstancesInput{
        InstanceIds: instanceIds,
    })
	resp, err := req.Send(context.TODO())
    if err != nil {
        return nil, err
    }
    return resp.TerminatingInstances, err
}
func (mdl *EC2Model) GetEC2Instances() ([]ec2.Instance, error) {

    req := mdl.model.DescribeInstancesRequest(&ec2.DescribeInstancesInput{})
    paginator := ec2.NewDescribeInstancesPaginator(req)
    var instances []ec2.Instance
    for paginator.Next(context.TODO()){
        for _, reservation := range paginator.CurrentPage().Reservations {
            instances = append(instances, reservation.Instances...)
        }
    }
    if paginator.Err() != nil {
        return nil, paginator.Err()
    }
	return instances, paginator.Err()
}

// Lists all instance types offered in the default region
func (mdl *EC2Model) ListOfferings() ([]ec2.InstanceTypeOffering, error) { // TODO: region, filters
	req := mdl.model.DescribeInstanceTypeOfferingsRequest(&ec2.DescribeInstanceTypeOfferingsInput{})
	resp, err := req.Send(context.TODO())
    if err != nil {
        return nil, err
    }
	return resp.InstanceTypeOfferings, err
}

// Lists AMIs offered
func (mdl *EC2Model) ListAMIs(filterMap map[string]string) ([]ec2.Image, error) {
	var filters []ec2.Filter
	for filterName, filterValue := range filterMap {
		// if filterValue != ""
			filters = append(filters, ec2.Filter{Name: aws.String(filterName), Values: strings.Split(filterValue, ",")})
		// }
	}
	req := mdl.model.DescribeImagesRequest(&ec2.DescribeImagesInput{Filters: filters})
	resp, err := req.Send(context.TODO())

    if err != nil {
        return nil, err
    }
	return resp.Images, err
}

// Changes instance type (or resize) to instType. For the restrictions on resizing an instance, see https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-resize.html
func (mdl *EC2Model) ChangeInstanceType(instId, instType string) error {
    req := mdl.model.ModifyInstanceAttributeRequest(&ec2.ModifyInstanceAttributeInput{
        InstanceId: aws.String(instId),
        InstanceType: &ec2.AttributeValue{Value: aws.String(instType)},
    })
	_, err := req.Send(context.TODO())
    // return resp.ModifyInstanceAttributeOutput         // an empty struct is returned
    return err
}
func (mdl *EC2Model) ListVolumes() ([]ec2.Volume, error) {
	var (
        volumes []ec2.Volume
    )
	req := mdl.model.DescribeVolumesRequest(&ec2.DescribeVolumesInput{})
    paginator := ec2.NewDescribeVolumesPaginator(req)
    for paginator.Next(context.TODO()){
        volumes = append(volumes, paginator.CurrentPage().Volumes...)
    }

    if paginator.Err() != nil {
        return nil, paginator.Err()
    }
	return volumes, paginator.Err()
}
func (mdl *EC2Model) AttachVolume(volId, instId, dev string) (ec2.AttachVolumeOutput, error) {
    req := mdl.model.AttachVolumeRequest(&ec2.AttachVolumeInput{
        Device: aws.String(dev),
        InstanceId: aws.String(instId),
        VolumeId: aws.String(volId),
    })
	resp, err := req.Send(context.TODO())
    if err != nil {
        return ec2.AttachVolumeOutput{}, err
    }
    return *resp.AttachVolumeOutput, nil
}

func (mdl *EC2Model) DetachVolume(volId, instId, dev string, force bool) (ec2.DetachVolumeOutput, error) {
    req := mdl.model.DetachVolumeRequest(&ec2.DetachVolumeInput{
        Device: aws.String(dev),
        Force: aws.Bool(force),
        InstanceId: aws.String(instId),
        VolumeId: aws.String(volId),
    })
	resp, err := req.Send(context.TODO())
    if err != nil {
        return ec2.DetachVolumeOutput{}, err
    }
    return *resp.DetachVolumeOutput, nil
}
func (mdl *EC2Model) ModifyVolume(iops, size int64, volType, volId string) (ec2.ModifyVolumeOutput, error) {
    input := &ec2.ModifyVolumeInput{}
    if iops != -1 { input.Iops = aws.Int64(iops) }
    if size != -1 { input.Size = aws.Int64(size) }
    if volType != "" { input.VolumeType = ec2.VolumeType(volType) }
    input.VolumeId = aws.String(volId)

    req := mdl.model.ModifyVolumeRequest(input)
	resp, err := req.Send(context.TODO())
    if err != nil {
        return ec2.ModifyVolumeOutput{}, err
    }
    return *resp.ModifyVolumeOutput, nil
}
// TODO: I don't understand yet why this is better than a simple print
func  printAWSError(err error) error {
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
	// I parametrized the goroutine only to make the channel send only
	go func(t *time.Ticker, ch chan<- common.Action, client *ec2.Client) { // dispatcher goroutine
		for {
			<-t.C
			watcher1(client, ch, true)
			// log.Println("watcher sent data")
		}
	}(ticker, mdl.Channel, mdl.model)
}

// watcher1 watches the status of all EC2 instances. A similar function in the view component "listner1" should make use of this information
func watcher1(client *ec2.Client, ch chan<- common.Action, describeAll bool) {
    var statuses []ec2.InstanceStatus
	req := client.DescribeInstanceStatusRequest(&ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: aws.Bool(describeAll),
	})
    paginator := ec2.NewDescribeInstanceStatusPaginator(req)
    for paginator.Next(context.TODO()){
        statuses = append(statuses, paginator.CurrentPage().InstanceStatuses...)
    }

    if paginator.Err() != nil {
        printAWSError(paginator.Err())      // TODO: logging and error handling
        return
    }
	action := common.Action{
        Type: common.ACTION_INSTANCES_STATUS_UPDATE,
        Data: statuses,
    }
	ch <- action
}
func watcher2(client *ec2.Client, ch chan<- common.Action) {        // TODO: filters ?

    var modifications []ec2.VolumeModification
    req := client.DescribeVolumesModificationsRequest(&ec2.DescribeVolumesModificationsInput{})
    paginator := ec2.NewDescribeVolumesModificationsPaginator(req)
    for paginator.Next(context.TODO()){
        modifications = append(modifications, paginator.CurrentPage().VolumesModifications...)
    }

    if paginator.Err() != nil {
        printAWSError(paginator.Err())      // TODO: logging and error handling
        return
    }
	action := common.Action{
        Type: common.ACTION_VOLUME_MODIFIED,
        Data: modifications,
    }
    ch<- action
}
