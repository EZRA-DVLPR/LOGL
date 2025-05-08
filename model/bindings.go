package model

import (
	"log"

	"fyne.io/fyne/v2/data/binding"
)

// define struct that will contain the bindings
type AppModel struct {
	SearchText   binding.String
	SelectedRow  binding.Int
	MaxProcesses binding.Int
	Progress     binding.Float
}

// create global instance of struct for binding
var GlobalModel = &AppModel{
	SearchText:   binding.NewString(),
	SelectedRow:  binding.NewInt(),
	MaxProcesses: binding.NewInt(),
	Progress:     binding.NewFloat(),
}

// set initial values
func init() {
	GlobalModel.SearchText.Set("")
	GlobalModel.SelectedRow.Set(1)
	GlobalModel.MaxProcesses.Set(1)
	GlobalModel.Progress.Set(0)
}

// INFO: All below functions are just for convenience on managing the bindings

func GetSearchText() (string, error) {
	return GlobalModel.SearchText.Get()
}

func SetSearchText(val string) error {
	return GlobalModel.SearchText.Set(val)
}

func AddSearchTextListener(listener func(string)) binding.DataListener {
	dataListener := binding.NewDataListener(func() {
		val, _ := GlobalModel.SearchText.Get()
		listener(val)
	})
	GlobalModel.SearchText.AddListener(dataListener)
	return dataListener
}

func GetSelectedRow() (int, error) {
	return GlobalModel.SelectedRow.Get()
}

func SetSelectedRow(val int) error {
	return GlobalModel.SelectedRow.Set(val)
}

func AddSelectedRowListener(listener func(int)) binding.DataListener {
	dataListener := binding.NewDataListener(func() {
		val, _ := GlobalModel.SelectedRow.Get()
		listener(val)
	})
	GlobalModel.SelectedRow.AddListener(dataListener)
	return dataListener
}

func GetMaxProcesses() (int, error) {
	return GlobalModel.MaxProcesses.Get()
}

func SetMaxProcesses(val int) error {
	return GlobalModel.MaxProcesses.Set(val)
}

func AddMaxProcessesListener(listener func(int)) binding.DataListener {
	// create the listener
	dataListener := binding.NewDataListener(func() {
		val, _ := GlobalModel.MaxProcesses.Get()
		listener(val)
	})
	// attach the listener to the binding
	GlobalModel.MaxProcesses.AddListener(dataListener)

	return dataListener
}

func RemoveMaxProcessesListener(listener binding.DataListener) {
	GlobalModel.MaxProcesses.RemoveListener(listener)
}

func GetProgress() (float64, error) {
	return GlobalModel.Progress.Get()
}

func SetProgress(val float64) error {
	return GlobalModel.Progress.Set(val)
}

func ResetProgress() error {
	return GlobalModel.Progress.Set(0)
}

func IncrementProgress() error {
	val, err := GlobalModel.Progress.Get()
	if err != nil {
		log.Fatal("Cannot increment progress properly", err)
	}
	return GlobalModel.Progress.Set(val + 1)
}

func AddProgressListener(listener func(float64)) binding.DataListener {
	// create the listener
	dataListener := binding.NewDataListener(func() {
		val, _ := GlobalModel.Progress.Get()
		listener(val)
	})
	// attach the listener to the binding
	GlobalModel.Progress.AddListener(dataListener)

	return dataListener
}

func RemoveProgressListener(listener binding.DataListener) {
	GlobalModel.Progress.RemoveListener(listener)
}
