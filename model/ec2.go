package model

import (
	"context"
	"strings"
	"time"

	"github.com/rfc2119/aws-tui/common"

	"github.com/aws/aws-sdk-go-v2/aws"
	// "github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Each channel is concerned with a service, and only the view to the model may use the channel. for example, for a designated ec2 worker channel, only the view responsible for ec2 may listen to the channel and consume items
// this might break in the future. sometimes, multiple benefeciaries exist for a single work. for example, when deleting an ebs volume, the ec2 console should also make use of the deletion command/action to update the affected instance. i don't know how to approach this (yet)
type EC2Model struct {
	model   *ec2.Client
	Channel chan common.Action // Unbuffered channel from model to view (see above)
	Name    string             // Use the convenient map to assign the correct name
	// logger	log.Logger

}

func NewEC2Model(config aws.Config) *EC2Model {
	return &EC2Model{
		model:   ec2.NewFromConfig(config),
		Name:    common.AWServicesDescriptions[common.ServiceEc2].Name,
		Channel: make(chan common.Action),
	}
}

func (mdl *EC2Model) StartEC2Instances(instanceIds []string) ([]types.InstanceStateChange, error) {
	// TODO: function docstring
	resp, err := mdl.model.StartInstances(context.TODO(), &ec2.StartInstancesInput{
		InstanceIds: instanceIds,
	})
	if err != nil {
		return nil, err
	}
	return resp.StartingInstances, err
}

func (mdl *EC2Model) StopEC2Instances(instanceIds []string, force, hibernate bool) ([]types.InstanceStateChange, error) {
	// TODO:  function docstring
	resp, err := mdl.model.StopInstances(context.TODO(), &ec2.StopInstancesInput{
		InstanceIds: instanceIds,
		Hibernate:   aws.Bool(hibernate),
		Force:       aws.Bool(force),
	})
	if err != nil {
		return nil, err
	}
	return resp.StoppingInstances, err
}

func (mdl *EC2Model) RebootEC2Instances(instanceIds []string) error {
	// TODO:  function docstring
	_, err := mdl.model.RebootInstances(context.TODO(), &ec2.RebootInstancesInput{
		InstanceIds: instanceIds,
	})
	return err
}

func (mdl *EC2Model) TerminateEC2Instances(instanceIds []string) ([]types.InstanceStateChange, error) {
	// TODO:  function docstring

	resp, err := mdl.model.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
		InstanceIds: instanceIds,
	})
	if err != nil {
		return nil, err
	}
	return resp.TerminatingInstances, err
}

func (mdl *EC2Model) GetEC2Instances() ([]types.Instance, error) {
	// TODO:  function docstring
	// TODO:
	var (
		instances []types.Instance
	)
	// Version 2 of the SDK addresses the above issues, and provides context per page.
	// It also enables easy iteration over API results that span multiple page. Here’s
	// a quick example of paginator usage with v2 SDK
	// req, err := mdl.model.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	paginator := ec2.NewDescribeInstancesPaginator(mdl.model, &ec2.DescribeInstancesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		for _, reservation := range page.Reservations {
			instances = append(instances, reservation.Instances...)
		}
	}
	return instances, nil
}

// Lists all instance types offered in the default region
func (mdl *EC2Model) ListOfferings() ([]types.InstanceTypeOffering, error) { // TODO: region, filters
	// TODO:  function docstring
	resp, err := mdl.model.DescribeInstanceTypeOfferings(context.TODO(), &ec2.DescribeInstanceTypeOfferingsInput{})
	if err != nil {
		return nil, err
	}
	return resp.InstanceTypeOfferings, err
}

// Lists AMIs offered
func (mdl *EC2Model) ListAMIs(filterMap map[string]string) ([]types.Image, error) {
	// TODO:  function docstring
	var filters []types.Filter
	for filterName, filterValue := range filterMap {
		// if filterValue != ""
		filters = append(filters, types.Filter{
			Name:   aws.String(filterName),
			Values: strings.Split(filterValue, ","),
		})
	}
	resp, err := mdl.model.DescribeImages(context.TODO(), &ec2.DescribeImagesInput{Filters: filters})
	if err != nil {
		return nil, err
	}
	return resp.Images, err
}

// Changes instance type (or resize) to instType. For the restrictions on resizing an instance, see https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-resize.html
func (mdl *EC2Model) ChangeInstanceType(instId, instType string) error {
	// TODO:  function docstring
	_, err := mdl.model.ModifyInstanceAttribute(context.TODO(), &ec2.ModifyInstanceAttributeInput{
		InstanceId:   aws.String(instId),
		InstanceType: &types.AttributeValue{Value: aws.String(instType)},
	})
	// return resp.ModifyInstanceAttributeOutput         // an empty struct is returned
	return err
}
func (mdl *EC2Model) ListVolumes() ([]types.Volume, error) {
	// TODO:  function docstring
	var (
		volumes []types.Volume
	)
	paginator := ec2.NewDescribeVolumesPaginator(mdl.model, &ec2.DescribeVolumesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		volumes = append(volumes, page.Volumes...)
	}

	return volumes, nil
}
func (mdl *EC2Model) AttachVolume(volId, instId, dev string) (ec2.AttachVolumeOutput, error) {
	// TODO:  function docstring
	resp, err := mdl.model.AttachVolume(context.TODO(), &ec2.AttachVolumeInput{
		Device:     aws.String(dev),
		InstanceId: aws.String(instId),
		VolumeId:   aws.String(volId),
	})
	if err != nil {
		return ec2.AttachVolumeOutput{}, err
	}
	return *resp, nil
}

func (mdl *EC2Model) DetachVolume(volId, instId, dev string, force bool) (ec2.DetachVolumeOutput, error) {
	// TODO:  function docstring
	resp, err := mdl.model.DetachVolume(context.TODO(), &ec2.DetachVolumeInput{
		Device:     aws.String(dev),
		Force:      aws.Bool(force),
		InstanceId: aws.String(instId),
		VolumeId:   aws.String(volId),
	})
	if err != nil {
		return ec2.DetachVolumeOutput{}, err
	}
	return *resp, nil
}

func (mdl *EC2Model) ModifyVolume(iops, size int32, volType, volId string) (ec2.ModifyVolumeOutput, error) {
	// TODO:  function docstring
	input := &ec2.ModifyVolumeInput{}
	if iops != -1 {
		input.Iops = aws.Int32(iops)
	}
	if size != -1 {
		input.Size = aws.Int32(size)
	}
	if volType != "" {
		input.VolumeType = types.VolumeType(volType)
	}
	input.VolumeId = aws.String(volId) // Required
	resp, err := mdl.model.ModifyVolume(context.TODO(), input)
	if err != nil {
		return ec2.ModifyVolumeOutput{}, err
	}
	return *resp, nil
}
func (mdl *EC2Model) DeleteVolume(volId string) (ec2.DeleteVolumeOutput, error) {
	// TODO:  function docstring
	resp, err := mdl.model.DeleteVolume(context.TODO(), &ec2.DeleteVolumeInput{
		VolumeId: aws.String(volId),
	})
	if err != nil {
		return ec2.DeleteVolumeOutput{}, err
	}
	return *resp, nil
}

func (mdl *EC2Model) CreateVolume(iops, size int32, volType, snapshotId, az string, isEncrypted, isMultiAttached bool) (ec2.CreateVolumeOutput, error) {
	// TODO:  function docstring
	// TODO: tags
	input := &ec2.CreateVolumeInput{}
	if iops != -1 {
		input.Iops = aws.Int32(iops)
	}
	if size != -1 {
		input.Size = aws.Int32(size)
	}
	if volType != "" {
		input.VolumeType = types.VolumeType(volType)
	}
	if snapshotId != "" {
		input.SnapshotId = aws.String(snapshotId)
	}
	input.MultiAttachEnabled = aws.Bool(isMultiAttached)
	input.Encrypted = aws.Bool(isEncrypted)
	input.AvailabilityZone = aws.String(az) // Required
	resp, err := mdl.model.CreateVolume(context.TODO(), input)
	if err != nil {
		return ec2.CreateVolumeOutput{}, err
	}
	return *resp, nil
}

// TODO: I don't understand yet why this is better than a simple print
// func printAWSError(err error) error {
// 	// TODO:  function docstring
// 	if err != nil {
// 		log.Println(err.Error())
// 		// if aerr, ok := err.(awserr.Error); ok {
// 		// 	switch aerr.Code() {
// 		// 	default:
// 		// 		log.Println(aerr.Error())
// 		// 	}
// 		// } else {
// 		// 	// Print the error, cast err to awserr.Error to get the Code and
// 		// 	// Message from an error.
// 		// 	log.Println(err.Error())
// 		// }
// 	}
// 	return err
// }

// DispatchWatchers sets the appropriate timer and calls each watcher
func (mdl *EC2Model) DispatchWatchers() {
	ticker := time.NewTicker(5 * time.Second) // TODO: 5
	// I parametrized the goroutine only to make the channel send only
	go func(t *time.Ticker, ch chan<- common.Action, client *ec2.Client) { // dispatcher goroutine
		for {
			<-t.C
			watcher1(client, ch, true)
			watcher2(client, ch)
			// log.Println("watcher sent data")
		}
	}(ticker, mdl.Channel, mdl.model)
}

// watcher1 watches the status of all EC2 instances. A similar function in the view component "listner1" should make use of this information
func watcher1(client *ec2.Client, ch chan<- common.Action, describeAll bool) {
	var (
		statuses []types.InstanceStatus
	)
	action := common.Action{
		Type: common.ActionInstancesStatusUpdate,
		Data: statuses,
	}
	params := &ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: aws.Bool(describeAll),
	}
	paginator := ec2.NewDescribeInstanceStatusPaginator(client, params)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			action.Type = common.ActionError
			action.Data = err
			break
		}
		statuses = append(statuses, page.InstanceStatuses...)
	}
	ch <- action
}

func watcher2(client *ec2.Client, ch chan<- common.Action) { // TODO: filters ?
	var (
		modifications []types.VolumeModification
	)
	action := common.Action{
		Type: common.ActionVolumeModified,
		Data: []types.VolumeModification{},
	}
	params := &ec2.DescribeVolumesModificationsInput{}
	paginator := ec2.NewDescribeVolumesModificationsPaginator(client, params)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			action.Type = common.ActionError
			action.Data = err
			break
		}
		modifications = append(modifications, page.VolumesModifications...)
	}
	ch <- action
}
