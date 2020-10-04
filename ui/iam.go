package ui

import (

	"github.com/rfc2119/aws-tui/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rivo/tview"
)



type iamService struct {
	mainUI
	Model *model.IAModel
	// logger log.Logger

}
func NewIAMService(config aws.Config, app *tview.Application, rootPage *ePages, statBar *StatusBar) *iamService {

	return &iamService{
		mainUI: mainUI{
			MainApp:   app,
			RootPage:  rootPage,
			StatusBar: statBar,
		},
		Model:        model.NewIAModel(config),
	}
}
