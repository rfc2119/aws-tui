package ui

import (
	"fmt"
    "reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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

type inputCapturer struct {
    // setKeyToFunc(p tview.Primitive, keyToFunc map[tcell.Key]func()){
    keyToFunc map[tcell.Key]func()
}

func (i *inputCapturer) UpdateKeyToFunc(keyToF map[tcell.Key]func()){
    for k, v := range keyToF {
        i.keyToFunc[k] = v
    }
}
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
		Pages:     tview.NewPages(),
		pageStack: []string{},
		HelpMessage:          "NO HELP MESSAGE (maybe submit a pull request ?)",
        inputCapturer: &inputCapturer{ keyToFunc: make(map[tcell.Key]func()) },
    }
    // p.inputCapturer.UpdateKeyToFunc(map[tcell.Key]func(){
    //     tcell.Key('?'): func(){ p.DisplayHelpMessage(p.HelpMessage)},
    //     tcell.Key('q'): func(){ p.ESwitchToPreviousPage() },
    // })
    p.setKeyToFunc()
	return &p
}


func (p *ePages) setKeyToFunc(){        // TODO: see repeated method on other types
        p.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
            uKey := event.Key()
            if event.Rune() != 0 {
                uKey = tcell.Key(event.Rune())
            }
            for k, f := range p.keyToFunc{
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

// use to go forward one page. do not use it if you intend not to go back to the page (for confirmation boxes for example). instead, use the normal tview.SwitchToPage or tview.AddAndSwitchToPage
func (p *ePages) ESwitchToPage(name string) *ePages {
	currentPageName, _ := p.GetFrontPage()
	p.pageStack = append(p.pageStack, currentPageName)
	p.SwitchToPage(name)
	return p

}

func (p *ePages) EAddAndSwitchToPage(name string, item tview.Primitive, resize bool) *ePages {
	p.EAddPage(name, item, resize, false) // visible=false as GetFrontPage() gets the last visible page
	return p.ESwitchToPage(name)

}

// use to move backward one page
func (p *ePages) ESwitchToPreviousPage() *ePages {
	if len(p.pageStack) > 0 {
		p.SwitchToPage(p.pageStack[len(p.pageStack)-1])
		// p.pageStack[len(p.pageStack) - 1] = nil		// TODO
		p.pageStack = p.pageStack[:len(p.pageStack)-1]
	}
	return p
}

func (p *ePages) DisplayHelpMessage(msg string) *ePages {

	helpPage := tview.NewTextView()
	helpPage.SetBackgroundColor(tcell.ColorBlue).SetTitle("HALP ME").SetTitleAlign(tview.AlignCenter).SetBorder(true)
	helpPage.SetText(msg)
	helpPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// p.RemovePage("help")
		if event.Rune() == 'q' {
			p.ESwitchToPreviousPage()
		}
		return event
	})

	return p.EAddAndSwitchToPage("help", helpPage, true) // "help" page gets overriden each time; resizable=true
}

func (p *ePages) GetPreviousPageName() string {
	return p.pageStack[len(p.pageStack)-1]
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
    *inputCapturer
	Members              []tview.Primitive // equivalent to the unexported member 'items' in tview.Grid
	CurrentMemberInFocus int               // index of the current member that has focus
	HelpMessage          string
	parent               *ePages // parent is used to display help message and navigate back to previous page (TODO: maybe the flex can do this itself ?)
}

func NewEFlex(parentPages *ePages) *eFlex {
	f := eFlex{
		Flex:                 tview.NewFlex(),
		Members:              []tview.Primitive{},
		CurrentMemberInFocus: 0,
		HelpMessage:          "NO HELP MESSAGE (maybe submit a pull request ?)",
		parent:               parentPages,
        inputCapturer: &inputCapturer{ keyToFunc: make(map[tcell.Key]func()) },
	}
    f.inputCapturer.UpdateKeyToFunc(map[tcell.Key]func(){
        tcell.Key('?'): func(){ f.DisplayHelp()},
        tcell.Key('q'): func(){ f.parent.ESwitchToPreviousPage() },
    })
    f.setKeyToFunc()
	return &f
}
func (f *eFlex) EAddItem(p tview.Primitive, fixedSize, proportion int, focus bool) *eFlex{
	f.AddItem(p, fixedSize, proportion, focus)
	f.Members = append(f.Members, p)
	return f
}
func (f *eFlex) DisplayHelp() {
	f.parent.DisplayHelpMessage(f.HelpMessage)
}
func (f *eFlex) setKeyToFunc(){        // TODO: see repeated method on other types
        f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
            uKey := event.Key()
            if event.Rune() != 0 {
                uKey = tcell.Key(event.Rune())
            }
            for k, f := range f.keyToFunc{
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
    *inputCapturer
	Members              []tview.Primitive // equivalent to the unexported member 'items' in tview.Grid
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
        inputCapturer: &inputCapturer{ keyToFunc: make(map[tcell.Key]func()) },
	}
    g.inputCapturer.UpdateKeyToFunc(map[tcell.Key]func(){
        tcell.Key('?'): func(){ g.DisplayHelp()},
        tcell.Key('q'): func(){ g.parent.ESwitchToPreviousPage() },
    })
    g.setKeyToFunc()
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

func (g *eGrid) setKeyToFunc(){        // TODO: see repeated method on other types
        g.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
            uKey := event.Key()
            if event.Rune() != 0 {
                uKey = tcell.Key(event.Rune())
            }
            for k, f := range g.keyToFunc{
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
		Table:                 tview.NewTable(),
        inputCapturer: &inputCapturer{ keyToFunc: make(map[tcell.Key]func()) },
	}
    t.setKeyToFunc()
	return &t
}
func (t *eTable) setKeyToFunc(){        // TODO: see repeated method on other types
        t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
            uKey := event.Key()
            if event.Rune() != 0 {
                uKey = tcell.Key(event.Rune())
            }
            for k, f := range t.keyToFunc{
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
type radioButtons struct {
	*tview.Box
    *inputCapturer
	options       []radioButtonOption
	currentOption int // index of current selected option
}

// NewRadioButtons returns a new radio button primitive.
func NewRadioButtons(optionNames []string) *radioButtons {
	options := make([]radioButtonOption, len(optionNames))
	for idx, name := range optionNames {
		options[idx] = radioButtonOption{name, true} // default: all enabled
	}
    r := radioButtons{
		Box:     tview.NewBox(),
		options: options,
        inputCapturer: &inputCapturer{ keyToFunc: make(map[tcell.Key]func()) },
	}
    r.setKeyToFunc()
    return &r
}

// Draw draws this primitive onto the screen.
func (r *radioButtons) Draw(screen tcell.Screen) {
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
func (r *radioButtons) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
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

// return the name of the current option
func (r *radioButtons) GetCurrentOptionName() string {
	return r.options[r.currentOption].name
}

func (r *radioButtons) GetOptions() []string {
    opts := make([]string, len(r.options))
	for idx, opt := range r.options {
        opts[idx] = opt.name
    }
    return opts
}
func (r *radioButtons) DisableOptionByName(name string) {
	for _, opt := range r.options {
		if opt.name == name {
			opt.enabled = false
			break
		}
	}
}

func (r *radioButtons) DisableOptionByIdx(idx int) {
	r.options[idx].enabled = false
}

func (r *radioButtons) EnableOptionByIdx(idx int) {
	r.options[idx].enabled = true
}
func (r *radioButtons) setKeyToFunc(){        // TODO: see repeated method on other types
        r.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
            uKey := event.Key()
            if event.Rune() != 0 {
                uKey = tcell.Key(event.Rune())
            }
            for k, f := range r.keyToFunc{
                if k == uKey {
                    f()
                    break
                }
            }
            return event
        })
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


// helper functions
// quite silly function (TODO: probably refactor)
func stringFromAWSVar(awsVar interface{}) string {
    var t string
    switch v := awsVar.(type){
    case *string:
        t = aws.StringValue(v)
    case *int:
        t = fmt.Sprint(aws.IntValue(v))      // hmmmm
    case *int64:
        // go vet being helpful as always:
        // conversion from int to string yields a string of one rune, 
        // not a string of digits (did you mean fmt.Sprint(x)?)
        t = fmt.Sprint(aws.Int64Value(v))      // hmmmm
    default:
        switch reflect.TypeOf(v).Kind(){
        case reflect.String:    // should be a type derived from string ?
            t = reflect.ValueOf(v).String()
        case reflect.Int, reflect.Int64:
            t = fmt.Sprint(reflect.ValueOf(v).Int())
        default:
            t = ""
        }
    }
    return t
}
// TODO: initially i did this to avoid writing similar functions for the same keys. i think this should be more generalized. also, is this slower than a big jump table (i.e switch/case statement) ?
// func setKeyToFunc(i interface{}){
//     // TODO: unify into one type
//     switch b := i.(type){
//     case *eFlex, *eGrid, *eTabel, *RadioButtons:
//     b.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
//         uKey := event.Key()
//         if event.Rune() != 0 {
//             uKey = tcell.Key(event.Rune())
//         }
//         for k, f := range b.keyToFunc{
//             if k == uKey {
//                 f()
//                 break
//             }
//         }
//         return event
//     })
// default:
//     fmt.Println("NOT IMPLEMENTED")
// }
// 
// }
