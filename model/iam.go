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

func (mdl *IAModel) GetCurrentUserInfo() *types.User {
	resp, err := mdl.model.GetUser(context.TODO(), &iam.GetUserInput{})
	if err != nil {
		log.Println(err)
	}
	return resp.User
}
