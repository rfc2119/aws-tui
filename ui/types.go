package ui

import (

	"github.com/rivo/tview"
	// "rfc2119/aws-tui/model"
)
// type viewComponent struct {
// 	ID      string           // unique id for the component; assigned as the address of the actual ui element
// 	Service string           // which service does this component serve ? see below for defintion of services
// 	Element tview.Primitive // the ui element itself. Primitive is an interface
// }
// services themselves are a way to group a model (the backend sdk) and the corresponding view. i don't know what will be the view as of this moment, but here goes nothing
// each service has a structure defined in the corresponding .go file
// a general representation of a model and view
type service struct {
	// View    []viewComponent
	MainApp *tview.Application
	RootPage *tview.Pages

}
// as usual, root.go contains some type definitions and configs
// exported methods of names similar to the original ui elements are prefixed with the vowel 'E' (capital E) for no reason. similarily, 'e' prefixes the custom ui elements defined





// ePage definition and methods
type ePage struct {
	page    *tview.Primitive
	helpMsg string
}

// eGrid definition and methods
type eGrid struct {
	*tview.Grid
    members []*tview.Primitive      // TODO: KeyCtrlW
}

func NewEgrid() *eGrid {
	return &eGrid{
		Grid:    tview.NewGrid(),
		members: []*tview.Primitive{},
	}
}
func (g *eGrid) EAddItem(p tview.Primitive, row, column, rowSpan, colSpan, minGridHeight, minGridWidth int, focus bool) *eGrid {

	g.AddItem(p, row, column, rowSpan, colSpan, minGridHeight, minGridWidth, focus)
	g.members = append(g.members, &p)
    return g
}
