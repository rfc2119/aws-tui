# arch
MVC architecture. the controller is merged into the ui.

each aws service has a type definition in the `services` namespace. for example, the ec2 service is defined as follows:
```go

type ec2Service struct {
	Model   *ec2.Client		// the backend component (aws client) which fetches the data from aws
	Service
}
```
where `Service` is an embedded type defined as:
```go

	View    []viewComponent		// all the components participating in viewing the data
	Channel chan Action		// the model to view channel
	Name    string 			// use the convenient map to assign the correct name
```
the `View` is a slice of all the ui components particiapting in the service (TODO: not sure yet). the first element should always be a `tview.Application` and the second should be the root element (here we use a `tview.Pages`). both elements are needed mainly for switching purposes (switching pages, switching focus, ... etc.)

## view
the view (ui) has direct access to the model namespace (i.e code access). it requests data from the model and display it

## model
a go channel is established between the model and the view. through this channel, the model is allowed to send periodic updates to the view to update items accordingly. the exchange protocol is specified manually through type definitions. the general type to be sent over the channel is:
```go
type Action struct {
	Type int
	Data       interface{}
}
``` 
the `Data` differs in each message. see *services/types.go* for definitions.

the view must explicitly call the `watch` method of a service with a given channel to receive updates from the model on the channel
