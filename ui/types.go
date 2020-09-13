package ui

import (
	"fmt"
	// "log"
	"time"

	// "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// type viewComponent struct {
// 	ID      string           // unique id for the component; assigned as the address of the actual ui element
// 	Service string           // which service does this component serve ? see below for defintion of services
// 	Element tview.Primitive // the ui element itself. Primitive is an interface
// }
type mainUI struct {
	// View    []viewComponent
	MainApp   *tview.Application
	RootPage  *ePages
	StatusBar *StatusBar
}

// services themselves are a way to group a model (the backend sdk) and the corresponding view. i don't know what will be the view as of this moment, but here goes nothing
// each service has a structure defined in the corresponding .go file
// a general representation of a model and view
// TODO: generalize services as a structure
// type service struct {
// 	*mainUI
// 	*aws.Client
// }

// as usual, types.go contains some type definitions and configs
// exported methods of names similar to the original ui elements (from tview package) are prefixed with the vowel 'E' (capital E) for no reason. similarily, 'e' prefixes the custom ui elements defined

// =================================
// ePages definition and methods
type ePages struct {
	*tview.Pages
	pageStack	[]string		// used for moving backwards one page at a time
}

func NewEPages() *ePages {
	return &ePages{
		Pages:            tview.NewPages(),
		pageStack:	  []string{},
	}
}

// same as AddPage
func (p *ePages) EAddPage(name string, item tview.Primitive, resize, visible bool) *ePages {
	p.AddPage(name, item, resize, visible)
	return p

}

// use to go forward one page
func (p *ePages) ESwitchToPage(name string) *ePages {
	currentPageName, _ := p.GetFrontPage()
	p.pageStack = append(p.pageStack, currentPageName)
	p.SwitchToPage(name)
	return p

}

func (p *ePages) EAddAndSwitchToPage(name string, item tview.Primitive, resize bool) *ePages {
	p.EAddPage(name, item, resize, false)		// visible=false as GetFrontPage() gets the last visible page
	return p.ESwitchToPage(name)

}

// use to move backward one page
func (p *ePages) ESwitchToPreviousPage() *ePages {
	if len(p.pageStack) > 0{
		p.SwitchToPage(p.pageStack[len(p.pageStack) - 1])
		// p.pageStack[len(p.pageStack) - 1] = nil		// TODO
		p.pageStack = p.pageStack[:len(p.pageStack) - 1]
	}
	return p
}

func (p *ePages) DisplayHelpMessage(msg string) *ePages {

	helpPage := tview.NewTextView()
	helpPage.SetBackgroundColor(tcell.ColorBlue).SetTitle("HALP ME").SetTitleAlign(tview.AlignCenter).SetBorder(true)
	helpPage.SetText(msg)
	helpPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// p.RemovePage("help")
		if event.Rune() == 'q'{
			p.ESwitchToPreviousPage()
		}
		return event
	})

	return p.EAddAndSwitchToPage("help", helpPage, true) // "help" page gets overriden each time; resizable=true
}
// ==================================
// eGrid definition and methods
type eGrid struct {
	*tview.Grid
	Members              []tview.Primitive // TODO: KeyCtrlW; equivalent to the unexported member 'items' in tview.Grid
	CurrentMemberInFocus int               // index of the current member that has focus
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
	return &g
}
func (g *eGrid) EAddItem(p tview.Primitive, row, column, rowSpan, colSpan, minGridHeight, minGridWidth int, focus bool) *eGrid {

	g.AddItem(p, row, column, rowSpan, colSpan, minGridHeight, minGridWidth, focus)
	g.Members = append(g.Members, p)
	return g
}

func (g *eGrid) DisplayHelp() {
	g.parent.DisplayHelpMessage(g.HelpMessage)
}

// =============================
// radio button primitive. copied from the demo https://github.com/rivo/tview/blob/master/demos/primitive
// RadioButtons implements a simple primitive for radio button selections.
type RadioButtons struct {
	*tview.Box
	options       []string
	currentOption int // index of current selected option
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
		switch event.Key() {
		case tcell.KeyUp:
			r.currentOption--
			if r.currentOption < 0 {
				r.currentOption = len(r.options) - 1
			}
		case tcell.KeyDown:
			r.currentOption++
			if r.currentOption >= len(r.options) {
				r.currentOption = 0
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j': // KeyDown
				r.currentOption++
				if r.currentOption >= len(r.options) {
					r.currentOption = 0
				}
			case 'k': // KeyUp
				r.currentOption--
				if r.currentOption < 0 {
					r.currentOption = len(r.options) - 1
				}
			}
		}
	})
}

func (r *RadioButtons) GetCurrentOption() string{
	return r.options[r.currentOption]
}

// ====================
// status bar
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

// non-focusable status bar by ignoring all key events and directing Focus() away
func (bar *StatusBar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return nil
}

func (bar *StatusBar) Focus(delegate func(p tview.Primitive)) {
	bar.Blur()
}
