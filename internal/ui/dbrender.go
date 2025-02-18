package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
)

// makes the table and reflects changes based on bindings
// PERF: make my own widget (EZRATableWidget) that has the following features:
//  1. clicking cell highlights row of cells
//  2. get column widths for each column
//  3. set size of column based on size of window
func createDBRender(sortType binding.String, opt binding.Bool) (dbRender *widget.Table) {
	var data [][]string

	// given the bindings create the table with the new set of data
	sortingType, _ := sortType.Get()
	sortingOpt, _ := opt.Get()
	if sortingOpt {
		data = dbhandler.SortDB(sortingType, "ASC")
	} else {
		data = dbhandler.SortDB(sortingType, "DESC")
	}

	// make the table with size of data. default 1 -- Display "No data"
	numRows := 1
	if len(data) != 0 {
		numRows = len(data)
	}

	// populate table with info
	dbRender = widget.NewTableWithHeaders(
		// table dims
		func() (int, int) { return numRows, 4 },
		// create an empty cell
		func() fyne.CanvasObject { return widget.NewLabel("") },
		// populate table with content
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			// if there is data in DB then display it
			// o/w display "No Data"
			if numRows > 1 {
				if id.Col == 0 {
					// if game name is too long, then truncate and append "...". o/w display entire game name
					if len(data[id.Row][0]) < 48 {
						obj.(*widget.Label).SetText(data[id.Row][0])
					} else {
						obj.(*widget.Label).SetText(data[id.Row][0][:45] + "...")
					}
				} else {
					// display the time data
					obj.(*widget.Label).SetText(fmt.Sprintf("%v", data[id.Row][id.Col]))
				}
			} else {
				obj.(*widget.Label).SetText("No Data")
			}
		},
	)
	dbRender = headerSetup(dbRender)

	// listener to update the contents of the table when value of sorting opt changes
	opt.AddListener(binding.NewDataListener(func() {
		dbRender = updateTable(opt, sortType, data, dbRender)
		dbRender.Refresh()
	}))

	// listener to update the contents of the table when value of sorting sortType changes
	sortType.AddListener(binding.NewDataListener(func() {
		dbRender = updateTable(opt, sortType, data, dbRender)
		dbRender.Refresh()
	}))
	return
}

// given bindings, data, and table will update the contents of the given table
func updateTable(opt binding.Bool, sortBy binding.String, data [][]string, dbRender *widget.Table) *widget.Table {
	sortingType, _ := sortBy.Get()
	sortingOpt, _ := opt.Get()
	if sortingOpt {
		data = dbhandler.SortDB(sortingType, "ASC")
	} else {
		data = dbhandler.SortDB(sortingType, "DESC")
	}
	numRows := 1
	if len(data) != 0 {
		numRows = len(data)
	}

	dbRender.Length = func() (int, int) { return numRows, 4 }
	dbRender.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
		// if there is data in DB then display it
		// o/w display "No Data"
		if numRows > 1 {
			if id.Col == 0 {
				// if game name is too long, then truncate and append "...". o/w display entire game name
				if len(data[id.Row][0]) < 48 {
					obj.(*widget.Label).SetText(data[id.Row][0])
				} else {
					obj.(*widget.Label).SetText(data[id.Row][0][:45] + "...")
				}
			} else {
				obj.(*widget.Label).SetText(fmt.Sprintf("%v", data[id.Row][id.Col]))
			}
		} else {
			obj.(*widget.Label).SetText("No Data")
		}
	}
	dbRender = headerSetup(dbRender)
	dbRender.ScrollToLeading()
	dbRender.ScrollToTop()
	return dbRender
}

func headerSetup(dbTable *widget.Table) *widget.Table {
	// name of each column header
	headers := []string{"Game Name", "Main Story", "Main + Sides", "Completionist"}

	// setup for creating the headers
	dbTable.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabelWithStyle(
			// add placeholder for at least 6 characters, making 6 digit numbers display nicely
			"______",
			// set text to be centered
			fyne.TextAlignCenter,
			// set text to be bold
			fyne.TextStyle{Bold: true},
		)
	}

	// make headers display content
	dbTable.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		if id.Col >= 0 && id.Col < len(headers) {
			// display headers defined prev
			obj.(*widget.Label).SetText(headers[id.Col])
		} else {
			// for row index, start at 1:rows
			obj.(*widget.Label).SetText(fmt.Sprintf("%d", id.Row+1))
		}
	}

	// set column widths
	dbTable.SetColumnWidth(0, 400)
	dbTable.SetColumnWidth(1, 200)
	dbTable.SetColumnWidth(2, 200)
	dbTable.SetColumnWidth(3, 200)
	return dbTable
}
