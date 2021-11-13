package ui

import (
	"github.com/rfc2119/aws-tui/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rivo/tview"
)

var (
	tblUsers = NewEtable()
	flexUsers *eFlex                // Container for the main page
	formUsers *tview.Form
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
		Model: model.NewIAModel(config),
	}

}

func (iamsvc *iamService) InitView() {

	var (
		inputFieldUserName         = tview.NewInputField().SetLabel("User name")
		checkBoxProgAccess         = tview.NewCheckbox().SetLabel("Programatic access")
		checkBoxConsoleAccess      = tview.NewCheckbox().SetLabel("Console access")
		// radioButtonConsolePassword = NewRadioButtons()
		checkBoxRequirePassReset   = tview.NewCheckbox().SetLabel("Require Password reset at next login")
	)
	// TODO: call back on checking the "console access" box to hide/show other components
	idxCheckBoxConsolePassword := formUsers.GetFormItemIndex("Console access")
	if idxCheckBoxConsolePassword == -1 {
		formUsers.AddFormItem(checkBoxRequirePassReset)
		// form.AddFormItem(radioButtonConsolePassword)
	}
	formUsers.AddFormItem(inputFieldUserName)
	formUsers.AddFormItem(checkBoxProgAccess)
	formUsers.AddFormItem(checkBoxConsoleAccess)

	flexUsers = NewEFlex(iamsvc.RootPage)
}

// GetAccessKeyLastUsed https://docs.aws.amazon.com/IAM/latest/APIReference/API_GetAccessKeyLastUsed.html
// CreateLoginProfile https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateLoginProfile.html
// CreateAccessKey https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateAccessKey.html
// CreateUser https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateUser.html
