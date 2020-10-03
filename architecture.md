# Architecture
A typical Model-View architecture, where the controller logic is merged into the View (referred here to as UI). The UI has direct code access to the model, while the model can talk to the UI component through a go channel. Only one channel is established for each service. The protocol which defines the format of messages that go through the channel is defined manually. 

A service is group of a model component, main view components, and service-specific data. Each aws service has a type defined in the corresponding *.go* file under the `ui` namespace. For example, the EC2 service is defined as follows:
```go
type ec2Service struct {
	mainUI
	Model *model.EC2Model           // The underlying entity that fetches data from AWS

	// service specific data
	instances []ec2.Instance
	volumes   []ec2.Volume
}
```
where `mainUI` is an embedded type defined as:
```go

type mainUI struct {
	// View    []viewComponent
	MainApp   *tview.Application    // The main application instance
	RootPage  *ePages               // The main container for switching between pages
	StatusBar *StatusBar            // A handy status bar
}
```
In addition to `mainUI`, file *types.go* in the `ui` namespace defines custom types based on the `tview` package. In general, these types register common key bindings and extra methods to help group common behavior. 

## View
The view (UI) has direct access to the model namespace (i.e code access). it requests data from the model and display it. In addition to one-time requests, it listens to data changes on a go channel (the service channel). To activate listeners, invoke the `WatchChanges()` method for each service

## Model
A go channel is established between the model and the view. Through this channel, the model is allowed to send periodic updates to the view to update items accordingly. The exchange protocol is specified manually through type definitions. The general type to be sent over the channel is:
```go
type Action struct {
	Type int
	Data       interface{}
}
``` 
The `Data` differs in each message and is defined manually for each action. After receiving a message, the UI component also acts appropriately. See *common/types.go* for definitions of example actions.

TODO
* defining types for the `Data` field is deprecated temporarily
* explain DispatchWatchers() and view listeners
