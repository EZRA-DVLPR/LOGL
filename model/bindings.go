package model

import (
	"log"

	"fyne.io/fyne/v2/data/binding"
)

// define struct that will contain the bindings
type AppModel struct {
	MaxProcesses binding.Int
	Progress     binding.Float
}

// create global instance of struct for binding
var GlobalModel = &AppModel{
	MaxProcesses: binding.NewInt(),
	Progress:     binding.NewFloat(),
}

// set initial values
func init() {
	GlobalModel.MaxProcesses.Set(1)
	GlobalModel.Progress.Set(0)
}

// INFO: All below functions are just for convenience on managing the bindings

func GetMaxProcesses() (int, error) {
	return GlobalModel.MaxProcesses.Get()
}

func SetMaxProcesses(val int) error {
	return GlobalModel.MaxProcesses.Set(val)
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
