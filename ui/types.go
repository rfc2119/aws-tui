package ui

import (
	"fmt"
	// "log"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// A common type used to hold keyboard keys to functions
type inputCapturer struct {
	// setKeyToFunc(p tview.Primitive, keyToFunc map[tcell.Key]func()){
	keyToFunc map[tcell.Key]func()
}

func (i *inputCapturer) UpdateKeyToFunc(keyToF map[tcell.Key]func()) {
	for k, v := range keyToF {
		i.keyToFunc[k] = v
	}
}

type layoutContainer struct {
	*inputCapturer
	Members              []tview.Primitive // equivalent to the unexported member 'items' in tview.Grid
	CurrentMemberInFocus int               // index of the current member that has focus
}

type mainUI struct {
	MainApp   *tview.Application
	RootPage  *ePages
	StatusBar *StatusBar
}

// If only we can shift focus to grid/flex members without a tview.Application...
func (u mainUI) enableShiftingFocus(l *layoutContainer) {

	l.inputCapturer.UpdateKeyToFunc(map[tcell.Key]func(){
		tcell.KeyTab: func() {
			if len(l.Members) > 0 {
				l.CurrentMemberInFocus++
				if l.CurrentMemberInFocus >= len(l.Members) { //  l.CurrentMemberInFocus %= len(l.Members)
					l.CurrentMemberInFocus = 0
				}
				u.MainApp.SetFocus(l.Members[l.CurrentMemberInFocus])
			}
		},
		tcell.KeyBacktab: func() {
			if len(l.Members) > 0 {
				l.CurrentMemberInFocus--
				if l.CurrentMemberInFocus < 0 { //  l.CurrentMemberInFocus %= len(l.Members)
					l.CurrentMemberInFocus = len(l.Members) - 1
				}
				u.MainApp.SetFocus(l.Members[l.CurrentMemberInFocus])
			}
		},
	})
}

// Shows a modal box with msg and switches back to previous page. This is useful for one-off usage (no nested boxes)
func (u mainUI) showConfirmationBox(msg string, rememberLastPage bool, doneFunc func()) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"Ok", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Ok" && doneFunc != nil {
				go func() {
					doneFunc()       // TODO: it's a mess with nested dialogues
					u.MainApp.Draw() // The key to get this done right
				}()
			}
			u.RootPage.ESwitchToPreviousPage()
			u.RootPage.ShowPage(u.RootPage.GetPreviousPageName()) // +1
		})
	// TODO: remove this and put a fixed page name. see if anything crashes
	pageName := fmt.Sprintf("%p", &modal)
	if rememberLastPage {
		u.RootPage.EAddAndSwitchToPage(pageName, modal, false) // resize=false
		u.RootPage.ShowPage(u.RootPage.GetCurrentPageName())   // +1
	} else {
		currPageName := u.RootPage.GetCurrentPageName()
		u.RootPage.AddAndSwitchToPage(pageName, modal, false) // resize=false
		u.RootPage.ShowPage(currPageName)                     // +1
	}

}

// Shows a generic modal box (rather than a confirmation-only box) centered at screen
// Props to skanehira from the docker tui "docui" for this! code is at github.com/skanehira/docui
func (u mainUI) showGenericModal(p tview.Primitive, width, height int, rememberLastPage bool) {
	var centeredModal *eGrid
	// unfortunately you can't access grid's minumum width or height. what to do ?
	// if g, ok := p.(*eGrid); ok {    // TODO: grid inside centered grid correctly; tview.Grid
	//     centeredModal = g
	// log.Println("OUR GRID")
	// // trying a flex instead
	// centeredModal := NewEFlex(u.RootPage).SetFullScreen(false).AddItem(
	//     tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(p, width, 0, true),
	//     height, 0, true)
	// centeredModal.SetColumns(0, width, 0).
	//                 SetRows(0, height, 0)
	// } else {
	centeredModal = NewEgrid(u.RootPage)
	centeredModal.SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true) // focus=true
		// }
	if g, ok := p.(*eGrid); ok {
		centeredModal.HelpMessage = g.HelpMessage
	} // TODO: eFlex
	// TODO: remove this and put a fixed page name. see if anything crashes
	// pageName := fmt.Sprintf("%p", &centeredModal)
	pageName := "centered modal"
	if rememberLastPage {
		u.RootPage.EAddAndSwitchToPage(pageName, centeredModal, true) // resize=true
		u.RootPage.ShowPage(u.RootPage.GetPreviousPageName())         // redraw on top (bottom ?) of the box
	} else {
		currPageName := u.RootPage.GetCurrentPageName()
		u.RootPage.AddAndSwitchToPage(pageName, centeredModal, true) // resize=true
		u.RootPage.ShowPage(currPageName)                            // redraw on top (bottom ?) of the box
	}

}

// TODO: generalize services as a structure
// type service struct {
// 	*mainUI
// 	*aws.Client
// }

// TODO: generalize
// uses SetInputCapture on primitive
// func (i *inputCapturer) setKeyToFunc(ifce interface{}){
//     switch b := ifce.(type){
//     case *eFlex, *eGrid, *eTable, *radioButton:
//         b.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
//             uKey := event.Key()
//             if event.Rune() != 0 {
//                 uKey = tcell.Key(event.Rune())
//             }
//             for k, f := range i.keyToFunc{
//                 if k == uKey {
//                     f()
//                     break
//                 }
//             }
//             return event
//         })
//     }
// }

// Exported methods of names similar to the original ui elements (from tview package) are prefixed with the vowel 'E' (capital E) for no reason. Similarily, 'e' prefixes the custom ui elements defined
// =================================
// ePages definition and methods
type ePages struct {
	*tview.Pages
	*inputCapturer
	HelpMessage string
	pageStack   []string // used for moving backwards one page at a time
}

func NewEPages() *ePages {
	p := ePages{
		Pages:         tview.NewPages(),
		pageStack:     []string{},
		HelpMessage:   "NO HELP MESSAGE (maybe submit a pull request ?)",
		inputCapturer: &inputCapturer{keyToFunc: make(map[tcell.Key]func())},
	}
	p.inputCapturer.UpdateKeyToFunc(map[tcell.Key]func(){
		//     tcell.Key('?'): func(){ p.DisplayHelpMessage(p.HelpMessage)},
		tcell.Key('q'): func() { p.ESwitchToPreviousPage() },
	})
	p.setKeyToFunc()
	return &p
}

func (p *ePages) setKeyToFunc() { // TODO: see repeated method on other types
	p.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		uKey := event.Key()
		if event.Rune() != 0 {
			uKey = tcell.Key(event.Rune())
		}
		for k, f := range p.keyToFunc {
			if k == uKey {
				f()
				break
			}
		}
		return event
	})
}

// same as AddPage
func (p *ePages) EAddPage(name string, item tview.Primitive, resize, visible bool) *ePages {
	p.AddPage(name, item, resize, visible)
	return p

}

// Use this to go forward one page. Do not use it if you intend not to go back (confirmation boxes for example). Instead, use the normal tview.SwitchToPage or tview.AddAndSwitchToPage
func (p *ePages) ESwitchToPage(name string) *ePages {
	currentPageName := p.GetCurrentPageName()
	if p.GetPreviousPageName() != currentPageName {
		p.pageStack = append(p.pageStack, currentPageName)
	}
	p.SwitchToPage(name)
	return p

}

func (p *ePages) EAddAndSwitchToPage(name string, item tview.Primitive, resize bool) *ePages {
	p.EAddPage(name, item, resize, false) // visible=false as GetFrontPage() gets the last visible page
	return p.ESwitchToPage(name)

}

// Use this to move backward one page
func (p *ePages) ESwitchToPreviousPage() *ePages {
	if len(p.pageStack) > 0 {
		p.SwitchToPage(p.pageStack[len(p.pageStack)-1])
		// p.pageStack[len(p.pageStack) - 1] = nil		// TODO
		p.pageStack = p.pageStack[:len(p.pageStack)-1]
	}
	return p
}

// Displays the help message given. Note that there should be no nested help messages
func (p *ePages) DisplayHelpMessage(msg string) *ePages {
	helpPage := tview.NewTextView()
	helpPage.SetTitle("Help").SetTitleAlign(tview.AlignCenter).SetBorder(true)
	helpPage.SetText(msg)
	return p.EAddAndSwitchToPage("help", helpPage, true) // "help" page gets overriden each time; resizable=true
}

func (p *ePages) GetPreviousPageName() string {
	if len(p.pageStack) > 0 {
		return p.pageStack[len(p.pageStack)-1]
	}
	return "" // Invalid page name
}

func (p *ePages) GetCurrentPageName() string {
	currentPageName, _ := p.GetFrontPage()
	return currentPageName
}

// TODO: this is a copy pasta from eGrid
// =================================
// eFlex definition and methods
type eFlex struct {
	*tview.Flex
	*layoutContainer
	HelpMessage string
	parent      *ePages // parent is used to display help message and navigate back to previous page (TODO: maybe the flex can do this itself ?)
}

func NewEFlex(parentPages *ePages) *eFlex {
	f := eFlex{
		Flex: tview.NewFlex(),
		layoutContainer: &layoutContainer{
			Members:              []tview.Primitive{},
			CurrentMemberInFocus: 0,
			inputCapturer:        &inputCapturer{keyToFunc: make(map[tcell.Key]func())},
		},
		HelpMessage: "NO HELP MESSAGE (maybe submit a pull request ?)",
		parent:      parentPages,
	}
	f.inputCapturer.UpdateKeyToFunc(map[tcell.Key]func(){
		tcell.Key('?'): func() { f.DisplayHelp() },
	})
	f.setKeyToFunc()
	return &f
}
func (f *eFlex) EAddItem(p tview.Primitive, fixedSize, proportion int, focus bool) *eFlex {
	f.AddItem(p, fixedSize, proportion, focus)
	f.Members = append(f.Members, p)
	return f
}
func (f *eFlex) DisplayHelp() {
	f.parent.DisplayHelpMessage(f.HelpMessage)
}

func (f *eFlex) setKeyToFunc() { // TODO: see repeated method on other types
	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		uKey := event.Key()
		if event.Rune() != 0 {
			uKey = tcell.Key(event.Rune())
		}
		for k, f := range f.keyToFunc {
			if k == uKey {
				f()
				break
			}
		}
		return event
	})
}

// ==================================
// eGrid definition and methods
type eGrid struct {
	*tview.Grid
	*layoutContainer
	HelpMessage string
	parent      *ePages // parent is used to display help message and navigate back to previous page (TODO: maybe the grid can do this itself ?)
}

func NewEgrid(parentPages *ePages) *eGrid {
	g := eGrid{
		Grid: tview.NewGrid(),
		layoutContainer: &layoutContainer{
			Members:              []tview.Primitive{},
			CurrentMemberInFocus: 0,
			inputCapturer:        &inputCapturer{keyToFunc: make(map[tcell.Key]func())},
		},
		HelpMessage: "NO HELP MESSAGE (maybe submit a pull request ?)",
		parent:      parentPages,
	}
	g.inputCapturer.UpdateKeyToFunc(map[tcell.Key]func(){
		tcell.Key('?'): func() { g.DisplayHelp() },
	})
	g.setKeyToFunc()
	return &g
}

// Wrapper function around tview.Grid.AddItem
func (g *eGrid) EAddItem(p tview.Primitive, row, column, rowSpan, colSpan, minGridHeight, minGridWidth int, focus bool) *eGrid {

	g.AddItem(p, row, column, rowSpan, colSpan, minGridHeight, minGridWidth, focus)
	g.Members = append(g.Members, p)
	return g
}

func (g *eGrid) DisplayHelp() {
	g.parent.DisplayHelpMessage(g.HelpMessage)
}

func (g *eGrid) setKeyToFunc() { // TODO: see repeated method on other types
	g.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		uKey := event.Key()
		if event.Rune() != 0 {
			uKey = tcell.Key(event.Rune())
		}
		for k, f := range g.keyToFunc {
			if k == uKey {
				f()
				break
			}
		}
		return event
	})
}

// eTable definition and methods
type eTable struct {
	*tview.Table
	*inputCapturer
}

func NewEtable() *eTable {
	t := eTable{
		Table:         tview.NewTable(),
		inputCapturer: &inputCapturer{keyToFunc: make(map[tcell.Key]func())},
	}
	t.setKeyToFunc()
	return &t
}
func (t *eTable) setKeyToFunc() { // TODO: see repeated method on other types
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		uKey := event.Key()
		if event.Rune() != 0 {
			uKey = tcell.Key(event.Rune())
		}
		for k, f := range t.keyToFunc {
			if k == uKey {
				f()
				break
			}
		}
		return event
	})
}

// =============================
// radio button primitive. copied from the demo https://github.com/rivo/tview/blob/master/demos/primitive
// RadioButtons implements a simple primitive for radio button selections.
type radioButtonOption struct {
	name    string
	enabled bool
}
type RadioButtons struct {
	*tview.Box
	*inputCapturer
	options       []radioButtonOption
	currentOption int // index of current selected option
}

// NewRadioButtons returns a new radio button primitive.
func NewRadioButtons(optionNames []string) *RadioButtons {
	options := make([]radioButtonOption, len(optionNames))
	for idx, name := range optionNames {
		options[idx] = radioButtonOption{name, true} // default: all enabled
	}
	r := RadioButtons{
		Box:           tview.NewBox(),
		options:       options,
		inputCapturer: &inputCapturer{keyToFunc: make(map[tcell.Key]func())},
	}
	r.setKeyToFunc()
	return &r
}

// Draw draws this primitive onto the screen.
func (r *RadioButtons) Draw(screen tcell.Screen) {
	r.Box.Draw(screen)
	x, y, width, height := r.GetInnerRect()

	for index, option := range r.options { //FIXME: what if option #1 is disabled ?
		if index >= height {
			break
		}
		radioButton := "\u25ef" // Unchecked.
		if index == r.currentOption && option.enabled {
			radioButton = "\u25c9" // Checked.
		}
		format := `%s[gray] %s`
		if option.enabled {
			format = `%s[white] %s`
		}
		line := fmt.Sprintf(format, radioButton, option.name)
		tview.Print(screen, line, x, y+index, width, tview.AlignLeft, tcell.ColorYellow)
	}
}

// InputHandler returns the handler for this primitive.
func (r *RadioButtons) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch {
		case event.Key() == tcell.KeyUp, event.Rune() == 'k':
			for i := 0; i < len(r.options); i++ {
				r.currentOption--
				if r.currentOption < 0 {
					r.currentOption = len(r.options) - 1
				}
				if r.options[r.currentOption].enabled {
					break
				}
			}
		case event.Key() == tcell.KeyDown, event.Rune() == 'j':
			for i := 0; i < len(r.options); i++ {
				r.currentOption++
				if r.currentOption >= len(r.options) {
					r.currentOption = 0
				}
				if r.options[r.currentOption].enabled {
					break
				}
			}
		}
	})
}

// Return the name of the current option
func (r *RadioButtons) GetCurrentOptionName() string {
	return r.options[r.currentOption].name
}

func (r *RadioButtons) GetOptions() []string {
	opts := make([]string, len(r.options))
	for idx, opt := range r.options {
		opts[idx] = opt.name
	}
	return opts
}
func (r *RadioButtons) DisableOptionByName(name string) {
	for _, opt := range r.options {
		if opt.name == name {
			opt.enabled = false
			break
		}
	}
}

func (r *RadioButtons) DisableOptionByIdx(idx int) {
	r.options[idx].enabled = false
}

func (r *RadioButtons) EnableOptionByIdx(idx int) {
	r.options[idx].enabled = true
}
func (r *RadioButtons) setKeyToFunc() { // TODO: see repeated method on other types
	r.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		uKey := event.Key()
		if event.Rune() != 0 {
			uKey = tcell.Key(event.Rune())
		}
		for k, f := range r.keyToFunc {
			if k == uKey {
				f()
				break
			}
		}
		return event
	})
}

// ====================
// A non-focusable status bar
type StatusBar struct {
	tview.TextView
	// app *tview.Application // Needed to properly clear text. Can be omitted
	durationInSeconds int // Duration after which the status bar is  cleared
}

func NewStatusBar() *StatusBar {

	bar := StatusBar{
		TextView: *tview.NewTextView(),
		// app: app,
		durationInSeconds: 3, // TODO: Parameter
	}
	// TODO: this is a naiive way of clearing the text bar on regular intervals; no syncronization or context is used
	bar.SetChangedFunc(func() {
		// go bar.app.Draw()          // hmmmmmm
		time.Sleep(time.Duration(bar.durationInSeconds) * time.Second)
		bar.Clear() // Clear() does not trigger a changed event
	})
	// bar.SetScrollable(false) // Helps trimming the internal buffer to only the viewable area
	bar.SetBackgroundColor(tcell.ColorBlue)
	return &bar
}

// Non-focusable status bar by ignoring all key events and directing Focus() away
func (bar *StatusBar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return nil
}

func (bar *StatusBar) Focus(delegate func(p tview.Primitive)) {
	bar.Blur()
}

// helper functions
// Grabs a string representation from types returned by the model (TODO: probably refactor)
func stringFromAWSVar(awsVar interface{}) string {
	var t string
	switch v := awsVar.(type) {
	case *string:
		t = aws.StringValue(v)
	case *int:
		t = fmt.Sprint(aws.IntValue(v)) // hmmmm
	case *int64:
		// go vet being helpful as always:
		// conversion from int to string yields a string of one rune,
		// not a string of digits (did you mean fmt.Sprint(x)?)
		t = fmt.Sprint(aws.Int64Value(v)) // hmmmm
	case *time.Time:
		t = fmt.Sprint(aws.TimeValue(v))
	default:
		switch reflect.TypeOf(v).Kind() {
		case reflect.String: // should be a type derived from string ?
			t = reflect.ValueOf(v).String()
		case reflect.Int, reflect.Int64:
			t = fmt.Sprint(reflect.ValueOf(v).Int())
		default:
			t = ""
		}
	}
	return t // TODO: return error on failure
}
