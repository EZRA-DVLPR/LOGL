package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
)

// makes the table and performs changes on it to send to main window
func createDBRender(sortBy string, opt string) (dbRender *widget.Table) {
	dbTable := makeTable(sortBy, opt)
	dbTable = headerSetup(dbTable)

	return dbTable
}

// PERF: make my own widget (EZRATableWidget) that has the following features:
//  1. clicking cell highlights row of cells
//  2. get column widths for each column
//  3. set size of column based on size of window
func makeTable(sortBy string, opt string) (dbTable *widget.Table) {
	// obtain data to insert into table
	dbData := dbhandler.SortDB(sortBy, opt)

	numRows := 1
	if len(dbData) != 0 {
		numRows = len(dbData)
	}

	// populate table with info
	dbTable = widget.NewTableWithHeaders(
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
					if len(dbData[id.Row][0]) < 48 {
						obj.(*widget.Label).SetText(dbData[id.Row][0])
					} else {
						obj.(*widget.Label).SetText(dbData[id.Row][0][:45] + "...")
					}
				} else if id.Col == 1 {
					obj.(*widget.Label).SetText(fmt.Sprintf("%v", dbData[id.Row][1]))
				} else if id.Col == 2 {
					obj.(*widget.Label).SetText(fmt.Sprintf("%v", dbData[id.Row][2]))
				} else {
					obj.(*widget.Label).SetText(fmt.Sprintf("%v", dbData[id.Row][3]))
				}
			} else {
				obj.(*widget.Label).SetText("No Data")
			}
		},
	)

	return
}

func headerSetup(dbTable *widget.Table) *widget.Table {
	// name of each header
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
