## k9s

* type `App` in *internal/ui/app.go* (which is the main ui screen) implements the interface `?` in ?
* *internal/{model,view,dao,client}/types.go*

## tcell
* aliases for popular keys (https://github.com/gdamore/tcell/blob/v1.3.0/key.go#L456)
```go
const (
    KeyBackspace  = KeyBS
    KeyTab        = KeyTAB
    KeyEsc        = KeyESC
    KeyEscape     = KeyESC
    KeyEnter      = KeyCR
    KeyBackspace2 = KeyDEL
)
```

# docui
* changing between panels/items is done with app.SetFocus(); gotta see how it's done in tview.Form
