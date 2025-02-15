package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Temp struct {
	row int
	col int
}

// creates the DBRender that will be display the DB
func createDBRender() (dbRender *widget.Table) {
	// create Lists
	// nameList, mainList, mainPlusList, compList := makeLists()
	dbTable := makeLists()

	return dbTable
}

func makeLists() (dbTable *widget.Table) {
	// obtain data to insert into table
	rows, cols := 20000, 4

	// populate table with info
	dbTable = widget.NewTableWithHeaders(
		// dims
		func() (int, int) { return rows, cols },
		// create an empty cell
		func() fyne.CanvasObject { return widget.NewLabel("") },
		// populate table with content
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(fmt.Sprintf("%d", id.Row+1))
		},
	)

	// name of each header
	headers := []string{"Game Name", "Main Story", "Main + Sides", "Completionist"}

	// make the headers bold
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
			obj.(*widget.Label).SetText(headers[id.Col])
		} else {
			// for row index, start at 1:rows
			obj.(*widget.Label).SetText(fmt.Sprintf("%d", id.Row+1))
		}
	}

	// expand column widths
	// TODO: Save these settings before close
	// Adjust size of columns to be spread evenly when the horizontal spacing is larger than the min
	dbTable.SetColumnWidth(0, 150)
	dbTable.SetColumnWidth(1, 200)
	dbTable.SetColumnWidth(2, 300)
	dbTable.SetColumnWidth(3, 400)
	return
}
