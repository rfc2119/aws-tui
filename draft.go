
	// drop down
	// dropdown := tview.NewDropDown().
	// 	SetLabel("Select an option (hit Enter): ").
	// 	SetOptions([]string{"First", "Second", "Third", "Fourth", "Fifth"}, nil)
	// if err := app.SetRoot(dropdown, true).SetFocus(dropdown).Run(); err != nil {
	// 	panic(err)
	// }

	// table
	// table := tview.NewTable().
	// 	SetBorders(true)
	// lorem := strings.Split("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
	// cols, rows := 10, 40
	// word := 0
	// for r := 0; r < rows; r++ {
	// 	for c := 0; c < cols; c++ {
	// 		color := tcell.ColorWhite
	// 		if c < 1 || r < 1 {
	// 			color = tcell.ColorYellow
	// 		}
	// 		table.SetCell(r, c,
	// 			tview.NewTableCell(lorem[word]).
	// 				SetTextColor(color).
	// 				SetAlign(tview.AlignCenter))
	// 		word = (word + 1) % len(lorem)
	// 	}
	// }
	// table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
	// 	if key == tcell.KeyEscape {
	// 		app.Stop()
	// 	}
	// 	if key == tcell.KeyEnter {
	// 		table.SetSelectable(true, true)
	// 	}
	// }).SetSelectedFunc(func(row int, column int) {
	// 	table.GetCell(row, column).SetTextColor(tcell.ColorRed)
	// 	table.SetSelectable(false, false)
	// })
	// if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
	// 	panic(err)
	// }

	// grid
	// newPrimitive := func(text string) tview.Primitive {
	// 	return tview.NewTextView().
	// 		SetTextAlign(tview.AlignCenter).
	// 		SetText(text)
	// }
	// menu := newPrimitive("Menu")
	// main := newPrimitive("Main content")
	// sideBar := newPrimitive("Side Bar")

	// grid := tview.NewGrid().
	// 	SetRows(3, 0, 3).
	// 	SetColumns(30, 0, 30).
	// 	SetBorders(true).
	// 	AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
	// 	AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	// // func (g *Grid) AddItem(p Primitive, row, column, rowSpan, colSpan, minGridHeight, minGridWidth int, focus bool) *Grid
	// 	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	// grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
	// 	AddItem(main, 1, 0, 1, 3, 0, 0, false).
	// 	AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

	// 	// Layout for screens wider than 100 cells.
	// grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
	// 	AddItem(main, 1, 1, 1, 1, 0, 100, false).
	// 	AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	// if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
	// 	panic(err)
	// }

	// pages
	// pages := tview.NewPages()
	// const pageCount = 5
	// for page := 0; page < pageCount; page++ {
	// 	func(page int) {
	// 		pages.AddPage(fmt.Sprintf("page-%d", page),
	// 			tview.NewModal().
	// 				SetText(fmt.Sprintf("This is page %d. Choose where to go next.", page+1)).
	// 				AddButtons([]string{"Next", "Quit"}).
	// 				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
	// 					if buttonIndex == 0 {
	// 						pages.SwitchToPage(fmt.Sprintf("page-%d", (page+1)%pageCount))
	// 					} else {
	// 						app.Stop()
	// 					}
	// 				}),
	// 			false,
	// 			page == 0)
	// 	}(page)
	// }
	// if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
	// 	panic(err)
	// }

	// type screenItem struct{
	//     item *tview.Box      // the superclass
	//     onSelected *func
	//     onDone *func
	// }
