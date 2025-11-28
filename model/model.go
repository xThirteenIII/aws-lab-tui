// package model defines the Model struct, which defines application states.
// It implements three methods which operate on the model structure.
//
// Init: a function that returns an initial command for the application to run
// Update: a function that updates incoming events and updates the model accordingly
// View: a function that renders UI based on the data in the model

package model

type model struct {
	choices  []string         // items on the tool list
	cursor   int              // which tool item the cursor is pointing at
	selected map[int]struct{} // which tool items are selected
}
