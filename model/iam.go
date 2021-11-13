package model

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/rfc2119/aws-tui/common"
)

type IAModel struct {
	model   *iam.Client
	Channel chan common.Action // channel from model to view (see above)
	Name    string             // use the convenient map to assign the correct name

}

func NewIAModel(config aws.Config) *IAModel {
	return &IAModel{
		model:   iam.NewFromConfig(config),
		Name:    common.AWServicesDescriptions[common.ServiceIdentityAndAccessManagement].Name,
		Channel: make(chan common.Action),
	}
}

func (mdl *IAModel) GetCurrentUserInfo() (*types.User, error) {
	// ValidationError: Must specify userName when calling with non-User credentials
	// Assume the principal is an IAM user
	resp, err := mdl.model.GetUser(context.TODO(), &iam.GetUserInput{})
	if err != nil {
		log.Printf("principal is not an IAM user: %s", err)
		// pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sts#Client.GetCallerIdentity
		// TODO: Assume principal is an IAM role
		// TODO: get more information about the principal
		todoString := aws.String("TODO")
		return &types.User{UserName: todoString, Arn: todoString}, err
	}
	return resp.User, err
}
