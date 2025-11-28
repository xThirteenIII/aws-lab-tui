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

// initialModel defines the initial state of the application.
// Defining a function to return the initial model, but could use a variable elsewhere, too.
func initialModel() model {
	return model{

		// List all the IoT Tools usable with the app.
		choices: []string{"Send IoT Jobs", "Dictionary", "Disenroll Inverter", "Upload .json to AWS S3"},

		// A map which indicates which choices are selected.
		// For now, only one choice is possible, for later states (e.g. Downloading commands from Dictionary), we can select multiples.
		// The keys refer to the indexes of the "choices" slice, above.
		selected: make(map[int]struct{}),
	}
}
