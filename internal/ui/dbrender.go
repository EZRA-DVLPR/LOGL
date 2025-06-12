package ui

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
	"github.com/EZRA-DVLPR/GameList/model"
)

var prevWidth float32

// makes the table and reflects changes based on values of bindings
func createDBRender(availableThemes map[string]ColorTheme) (dbRender *widget.Table) {
	// if db exists then get the data
	log.Println("Checking existence of local DB")
	if dbhandler.CheckDBExists() {
		log.Println("DB exists. Obtaining data with stored defaults")
		// no initial search query so use ""
		dbData.Set(dbhandler.SortDB())
	} else {
		log.Println("No DB found!")
		dbhandler.CreateDB()
	}
	data, _ := dbData.Get()

	// make the table with size of data. default 1
	numRows := 1
	if len(data) != 0 {
		numRows = len(data)
	}

	st, _ := model.GetSelectedTheme()
	currTheme := availableThemes[st]

	log.Println("Creating the table template")
	// populate table with info
	dbRender = widget.NewTableWithHeaders(
		// table dims
		func() (int, int) { return numRows, 4 },
		// create empty cells with dflt bg color and empty text
		func() fyne.CanvasObject {
			bg := canvas.NewRectangle(hexToColor(currTheme.Background))
			label := widget.NewLabel("")
			return container.NewStack(bg, label)
		},
		// populate table with content
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			stack := obj.(*fyne.Container)
			bg := stack.Objects[0].(*canvas.Rectangle)
			label := stack.Objects[1].(*widget.Label)

			// if there is data in DB then display it o/w display "No Data"
			if len(data) > 1 {
				if id.Col == 0 {
					// if game name is too long, then truncate and append "...". o/w display entire game name
					if len(data[id.Row][0]) < 48 {
						label.SetText(data[id.Row][0])
					} else {
						label.SetText(data[id.Row][0][:45] + "...")
					}
				} else {
					// display the time data
					label.SetText(fmt.Sprintf("%v", data[id.Row][id.Col]))
				}
			} else {
				label.SetText("No Data")
			}

			// no selected row, so apply background and alt background
			if id.Row%2 == 0 {
				bg.FillColor = hexToColor(currTheme.Background)
			} else {
				bg.FillColor = hexToColor(currTheme.AltBackground)
			}
		},
	)

	// highlight the row of the cell clicked if its not a divider -> dividers have negative position values
	dbRender.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Col >= 0 {
			model.SetSelectedRow(id.Row)
		}
	}

	// set up the headers
	log.Println("Setting up table headers")
	width := w.Content().Size().Width
	dbRender = headerSetup(dbRender, width, availableThemes)

	// refresh DBRender when the theme changes
	model.AddSelectedThemeListener(func(string) {
		log.Println("Selected Theme changed. Adjusting Table")
		UpdateDBData()
		updateTableColors(dbRender, availableThemes)
		dbRender.Refresh()
	})

	// change contents of dbData binding when sort order changes
	model.AddSortOrderListener(func(val bool) {
		log.Println("Sort Order changed. Adjusting Table")
		width := w.Content().Size().Width
		UpdateDBData()
		dbRender = updateTable(dbRender, width, availableThemes)
		dbRender.Refresh()
	})

	// change contents of dbData binding when sort category changes
	model.AddSortCategoryListener(func(val string) {
		log.Println("Sort Category changed. Adjusting Table")
		width := w.Content().Size().Width
		UpdateDBData()
		dbRender = updateTable(dbRender, width, availableThemes)
		dbRender.Refresh()
	})

	// change contents of dbData binding when search text changes
	model.AddSearchTextListener(func(val string) {
		log.Println("Search Text changed. Adjusting Table")
		width := w.Content().Size().Width
		UpdateDBData()
		dbRender = updateTable(dbRender, width, availableThemes)
		dbRender.Refresh()
	})

	// selectedRow changes
	model.AddSelectedRowListener(func(int) {
		log.Println("Selected Row changed. Adjusting Table")
		selRow, _ := model.GetSelectedRow()
		width := w.Content().Size().Width
		dbRender.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
			// get the label from the stack
			stack := obj.(*fyne.Container)
			bg := stack.Objects[0].(*canvas.Rectangle)
			label := stack.Objects[1].(*widget.Label)
			label.SetText("No Data")

			// highlighting
			if id.Row == selRow {
				bg.FillColor = hexToColor(currTheme.HoverColor)
			} else {
				if id.Row%2 == 0 {
					bg.FillColor = hexToColor(currTheme.Background)
				} else {
					bg.FillColor = hexToColor(currTheme.AltBackground)
				}
			}
		}
		dbRender = updateTable(dbRender, width, availableThemes)

		// scroll to the new location selected
		var selCell widget.TableCellID
		selCell.Row = selRow
		selCell.Col = 0
		dbRender.ScrollTo(selCell)
		dbRender.ScrollToLeading()
	})

	go fixTableSize(dbRender)

	return
}

// sets dbData with given opts
func UpdateDBData() {
	dbData.Set(dbhandler.SortDB())
	sr, _ := model.GetSelectedRow()
	// ensures updates occur when there is no selected row
	if sr != -1 {
		model.SetSelectedRow(-1)
	} else {
		model.SetSelectedRow(-2)
	}
}

// update the contents of the given table
func updateTable(
	dbRender *widget.Table,
	width float32,
	availableThemes map[string]ColorTheme,
) *widget.Table {
	selRow, _ := model.GetSelectedRow()

	// check rows for finding dims
	data, _ := dbData.Get()
	numRows := 1
	if len(data) != 0 {
		numRows = len(data)
	}

	st, _ := model.GetSelectedTheme()
	currTheme := availableThemes[st]

	// set dims
	dbRender.Length = func() (int, int) { return numRows, 4 }
	dbRender.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
		// get the label from the stack
		stack := obj.(*fyne.Container)
		bg := stack.Objects[0].(*canvas.Rectangle)
		label := stack.Objects[1].(*widget.Label)

		// if there is data in DB then display it
		// o/w display "No Data"
		if len(data) != 0 {
			if id.Col == 0 {
				// if game name is too long, then truncate and append "...". o/w display entire game name
				if len(data[id.Row][0]) < 48 {
					label.SetText(data[id.Row][0])
				} else {
					label.SetText(data[id.Row][0][:45] + "...")
				}
			} else {
				// display the time data
				label.SetText(fmt.Sprintf("%v", data[id.Row][id.Col]))
			}
		} else {
			label.SetText("No Data")
		}

		// highlighting
		if id.Row == selRow {
			bg.FillColor = hexToColor(currTheme.HoverColor)
		} else {
			if id.Row%2 == 0 {
				bg.FillColor = hexToColor(currTheme.Background)
			} else {
				bg.FillColor = hexToColor(currTheme.AltBackground)
			}
		}
	}

	// setup headers
	dbRender = headerSetup(dbRender, width, availableThemes)

	// unselect all cells
	dbRender.UnselectAll()
	dbRender.Refresh()
	return dbRender
}

func headerSetup(
	dbTable *widget.Table,
	width float32,
	availableThemes map[string]ColorTheme,
) *widget.Table {
	st, _ := model.GetSelectedTheme()

	currTheme := availableThemes[st]

	// name of each column header
	headers := []string{"Game Name", "Main Story", "Main + Sides", "Completionist"}

	// setup for creating the headers
	dbTable.CreateHeader = func() fyne.CanvasObject {
		return container.NewStack(
			canvas.NewRectangle(hexToColor(currTheme.ScrollBarColor)),
			widget.NewLabelWithStyle(
				"------",
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true},
			),
			widget.NewButton("------", nil),
		)
	}

	// make headers display content
	dbTable.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		containerObj := obj.(*fyne.Container)
		labelBG := containerObj.Objects[0].(*canvas.Rectangle)
		label := containerObj.Objects[1].(*widget.Label)
		button := containerObj.Objects[2].(*widget.Button)

		if id.Row == -1 {
			// display column header buttons
			button.Show()
			labelBG.Hide()
			label.Hide()
			button.SetText(headers[id.Col])
			button.OnTapped = func() {
				// sortCategory gets set to whichever header was clicked
				if id.Col == 0 {
					model.SetSortCategory("name")
				} else if id.Col == 1 {
					model.SetSortCategory("main")
				} else if id.Col == 2 {
					model.SetSortCategory("mainPlus")
				} else {
					model.SetSortCategory("comp")
				}
			}
		} else {
			// display row label index, from 1:rows
			button.Hide()
			label.Show()
			labelBG.Show()
			labelBG.FillColor = hexToColor(currTheme.InputBackgroundColor)
			label.SetText(fmt.Sprintf("%d", id.Row+1))
		}
		labelBG.Refresh()
	}

	// set column widths
	dbTable.SetColumnWidth(0, 400)

	// game name has 400, and the row headers take ~70 spacing
	// all other space is to be given to the other columns
	spacing := (width - 400 - 70) / 3

	dbTable.SetColumnWidth(1, spacing)
	dbTable.SetColumnWidth(2, spacing)
	dbTable.SetColumnWidth(3, spacing)
	return dbTable
}

// check window size every 0.25 and adjust size of table col widths if it changes
func fixTableSize(dbRender *widget.Table) {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		width := w.Content().Size().Width
		if prevWidth != width {
			select {
			case <-ticker.C:
				log.Println("Window size changed! Adjusting Column widths of table")
				// update the table widths
				prevWidth = width

				// set name col width
				dbRender.SetColumnWidth(0, 400)

				// game name has 400, and the row headers take ~70 spacing
				// all other space is to be given to the other columns
				spacing := (width - 400 - 70) / 3

				dbRender.SetColumnWidth(1, spacing)
				dbRender.SetColumnWidth(2, spacing)
				dbRender.SetColumnWidth(3, spacing)

				dbRender.Refresh()
			}
		}
	}
}

func updateTableColors(
	dbRender *widget.Table,
	availableThemes map[string]ColorTheme,
) {
	selRow, _ := model.GetSelectedRow()
	st, _ := model.GetSelectedTheme()
	currTheme := availableThemes[st]
	dbRender.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
		// get the label from the stack
		stack := obj.(*fyne.Container)
		bg := stack.Objects[0].(*canvas.Rectangle)
		// highlighting
		if id.Row == selRow {
			bg.FillColor = hexToColor(currTheme.HoverColor)
		} else {
			if id.Row%2 == 0 {
				bg.FillColor = hexToColor(currTheme.Background)
			} else {
				bg.FillColor = hexToColor(currTheme.AltBackground)
			}
		}
	}

	dbRender.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		containerObj := obj.(*fyne.Container)
		labelBG := containerObj.Objects[0].(*canvas.Rectangle)
		label := containerObj.Objects[1].(*widget.Label)
		button := containerObj.Objects[2].(*widget.Button)

		if id.Row != -1 {
			// display row label index, from 1:rows
			button.Hide()
			label.Show()
			labelBG.Show()
			labelBG.FillColor = hexToColor(currTheme.InputBackgroundColor)
			label.SetText(fmt.Sprintf("%d", id.Row+1))
		}
		labelBG.Refresh()
	}
}
