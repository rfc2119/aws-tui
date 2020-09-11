package ui

import (
	"fmt"
	// "log"
	"time"

	"github.com/gdamore/tcell"
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
	MainApp  *tview.Application
	RootPage *ePages
}

// as usual, root.go contains some type definitions and configs
// exported methods of names similar to the original ui elements are prefixed with the vowel 'E' (capital E) for no reason. similarily, 'e' prefixes the custom ui elements defined

// =================================
// ePages definition and methods
type ePages struct {
	*tview.Pages
	// helpMsg string
	previousPage     tview.Primitive
	previousPageName string
}

func NewEPages() *ePages {
	return &ePages{
		Pages:            tview.NewPages(),
		previousPage:     nil,
		previousPageName: "",
	}
}

// func (p *ePages) EShowPage(name string) *ePages{
// 	p.previousPageName, p.previousPage = p.GetFrontPage()
// 	p.ShowPage(name)
//
// }
// GetFrontPage() returns the last added page that is visible, that's why we need the if visible condition
func (p *ePages) EAddPage(name string, item tview.Primitive, resize, visible bool) *ePages {
	if visible {
		p.previousPageName, p.previousPage = p.GetFrontPage()
	}
	p.AddPage(name, item, resize, visible)
	return p

}

// use only if not adding new pages (switching to pages already created)
// TODO: often the current page is a dialog box or some confirmation message. in that case, use recordPreviousPage to not record the transition to the dialog box. this is a defeciency in the current architecture
func (p *ePages) ESwitchToPage(name string, recordPreviousPage bool) *ePages {
	if recordPreviousPage {
		p.previousPageName, p.previousPage = p.GetFrontPage()
	}
	p.SwitchToPage(name)
	return p

}

// use either EAddPage or ESwitchToPage; do not use both in one call
func (p *ePages) EAddAndSwitchToPage(name string, item tview.Primitive, resize bool) *ePages {
	p.EAddPage(name, item, resize, true)
	p.SwitchToPage(name)
	return p

}
func (p *ePages) ESwitchToPreviousPage(recordPreviousPage bool) *ePages {
	return p.ESwitchToPage(p.previousPageName, recordPreviousPage)
}

func (p *ePages) DisplayHelpMessage(msg string) *ePages {

	helpPage := tview.NewTextView()
	helpPage.SetBackgroundColor(tcell.ColorBlue).SetTitle("HALP").SetTitleAlign(tview.AlignCenter).SetBorder(true)
	helpPage.SetText(msg)
	helpPage.SetDoneFunc(func(key tcell.Key) { // TODO: uhhh we can't just save the previous page as it gets destroyed; what do ? what about the recordPreviousPage boolean ?
		// p.RemovePage("help")
		// p.SwitchToPage(p.previousPageName)
		p.ESwitchToPreviousPage(false)
	})

	return p.EAddAndSwitchToPage("help", helpPage, true) // "help" page gets overriden each time; resizable
}
func (p *ePages) GetPreviousPage() (string, tview.Primitive) {
	return p.previousPageName, p.previousPage
}

// ==================================
// eGrid definition and methods
type eGrid struct {
	*tview.Grid
	Members              []tview.Primitive // TODO: KeyCtrlW; equivalent to the unexported member 'items' in tview.Grid
	CurrentMemberInFocus int                // index of the current member that has focus
	HelpMessage          string
	parent               *ePages // parent is used to display help message and navigate back to previous page (TODO: maybe the grid can do this itself ?)
}

func NewEgrid(parentPages *ePages) *eGrid {
	g := eGrid{
		Grid:                 tview.NewGrid(),
		Members:              []tview.Primitive{},
		CurrentMemberInFocus: 0,
		HelpMessage:          "NO HELP MESSAGE (maybe submit a pull request ?)",
		parent:               parentPages,
	}
// SetInputCapture installs a function which captures key events before they are forwarded to the primitive's default key event handler (implemented by InputHandler())
// g.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 	switch {
// 	case event.Key() == tcell.KeyCtrlW:
// 			// statusBar.SetText("moving to another item; statusbar focus: " + fmt.Sprintf("%s", statusBar.HasFocus())) // TODO
// 		if len(g.Members) > 0 {
// 			g.CurrentMemberInFocus++
// 			if g.CurrentMemberInFocus == len(grid.Members) { //  grid.CurrentMemberInFocus %= len(grid.Members)
// 				g.CurrentMemberInFocus = 0
// 			}
// 			// ec2svc.MainApp.SetFocus(*g.Members[grid.CurrentMemberInFocus]) // * hmmm
// 			setFocus(*g.Members[g.CurrentMemberInFocus]) // * hmmm
// 		}
// 	case event.Rune() == '?':
// 		// ec2svc.RootPage.DisplayHelpMessage(g.HelpMessage)
// 		g.parent.DisplayHelpMessage(g.HelpMessage)
// 		log.Println("FIXME") //FIXME
// 	}
// 	case event.Key() == tcell.KeyEscape:
// 		g.parent.ESwitchToPreviousPage(false)
// 
// }
return &g
}
func (g *eGrid) EAddItem(p tview.Primitive, row, column, rowSpan, colSpan, minGridHeight, minGridWidth int, focus bool) *eGrid {

	g.AddItem(p, row, column, rowSpan, colSpan, minGridHeight, minGridWidth, focus)
	g.Members = append(g.Members, p)
	return g
}


// =============================
// radio button primitive. copied from https://github.com/rivo/tview/blob/master/demos/primitive
// RadioButtons implements a simple primitive for radio button selections.
type RadioButtons struct {
	*tview.Box
	options       []string
	currentOption int
}

// NewRadioButtons returns a new radio button primitive.
func NewRadioButtons(options []string) *RadioButtons {
	return &RadioButtons{
		Box:     tview.NewBox(),
		options: options,
	}
}

// Draw draws this primitive onto the screen.
func (r *RadioButtons) Draw(screen tcell.Screen) {
	r.Box.Draw(screen)
	x, y, width, height := r.GetInnerRect()

	for index, option := range r.options {
		if index >= height {
			break
		}
		radioButton := "\u25ef" // Unchecked.
		if index == r.currentOption {
			radioButton = "\u25c9" // Checked.
		}
		line := fmt.Sprintf(`%s[white]  %s`, radioButton, option)
		tview.Print(screen, line, x, y+index, width, tview.AlignLeft, tcell.ColorYellow)
	}
}

// InputHandler returns the handler for this primitive.
func (r *RadioButtons) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		// switch event.Key() {
		switch {
		case event.Rune() == 'k':
		case event.Key() == tcell.KeyUp:
			r.currentOption--
			if r.currentOption < 0 {
				r.currentOption = 0
			}
		case event.Rune() == 'j':
		case event.Key() == tcell.KeyDown:
			r.currentOption++
			if r.currentOption >= len(r.options) {
				r.currentOption = len(r.options) - 1
			}
		}
	})
}

// ====================
// WIP status bar
// TODO: the bar still gets focused when pressing ^W
type StatusBar struct {
	*tview.TextView
	durationInSeconds int // duration after which the status bar is  cleared
}

func NewStatusBar() *StatusBar {

	bar := StatusBar{
		TextView:          tview.NewTextView(),
		durationInSeconds: 3, // TODO: parameter
	}
	// very naiive way of clearing the text bar on regular intervals; no syncronization or context is used
	bar.SetChangedFunc(func() {
		time.Sleep(time.Duration(bar.durationInSeconds) * time.Second)
		bar.Clear()
	})
	return &bar
}

// non-focusable status bar by ignoring all key events *crosses fingers*
func (bar *StatusBar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return nil
}

func (bar *StatusBar) Focus(delegate func(p tview.Primitive)){
	bar.Blur()
}
