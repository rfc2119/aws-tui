package ui

import (

	"github.com/rivo/tview"
)
// as usual, root.go contains some type definitions and configs
// exported methods of names similar to the original ui elements are prefixed with the vowel 'E' (capital E) for no reason. similarily, 'e' prefixes the custom ui elements defined





type ePage struct { // cPage reads as custom page
	page    *tview.Primitive
	helpMsg string
}

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
