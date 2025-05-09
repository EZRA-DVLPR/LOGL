package ui

import (
	"fmt"
	"sync"

	"fyne.io/fyne/v2/data/binding"
)

// create the struct for handling the data -> [][]string
// PERF: Change name of MyDataBinding to something more descriptive like matrix/tablestring
type MyDataBinding struct {
	data     [][]string
	listener []binding.DataListener
	// manages goroutines to read/write properly
	lock sync.Mutex
}

// create a new data binding obj with no initial data
func NewMyDataBindingEmpty() *MyDataBinding {
	return &MyDataBinding{data: [][]string{}}
}

// create a new data binding obj with the initial data given
func NewMyDataBinding(initialData [][]string) *MyDataBinding {
	return &MyDataBinding{data: initialData}
}

// return the data from the binding, and error
func (b *MyDataBinding) Get() ([][]string, error) {
	// prevent other goroutines from writing to this binding while we read it
	b.lock.Lock()

	// release the lock so other goroutines can write again when we return
	defer b.lock.Unlock()
	return b.data, nil
}

// set updates to data and notify the listeners
func (b *MyDataBinding) Set(v any) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	// grab the new data
	newData, ok := v.([][]string)
	if !ok {
		return fmt.Errorf("Invalid type! Expected type: [][]string")
	}

	// update only if the data is different
	b.data = newData

	// notify the listeners that a change occurred
	for _, listener := range b.listener {
		listener.DataChanged()
	}
	return nil
}

// register a listener for changes
func (b *MyDataBinding) AddListener(listener binding.DataListener) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.listener = append(b.listener, listener)
}
