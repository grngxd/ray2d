package screen

type Screen interface {
	// Update is called every frame to update the screen
	Update()
	// OnEnter is called when the screen is entered
	OnEnter()
	// OnExit is called when the screen is exited
	OnExit()
}

// ScreenManager is responsible for managing the screens in the game
type ScreenManager struct {
	currentScreen Screen
}

var Manager *ScreenManager = NewManager()

// NewScreenManager creates a new ScreenManager
func NewManager() *ScreenManager {
	return &ScreenManager{
		currentScreen: &loadingScreen{},
	}
}

// SetScreen sets the current screen
func (sm *ScreenManager) SetScreen(s Screen) {
	if sm.currentScreen != nil {
		sm.currentScreen.OnExit()
	}

	sm.currentScreen = s
	sm.currentScreen.OnEnter()
}

// Update updates the current screen
func (sm *ScreenManager) Update() {
	if sm.currentScreen != nil {
		sm.currentScreen.Update()
	}
}

// GetCurrentScreen returns the current screen
func (sm *ScreenManager) GetCurrentScreen() Screen {
	return sm.currentScreen
}

// SetCurrentScreen sets the current screen
func (sm *ScreenManager) SetCurrentScreen(s Screen) {
	sm.currentScreen = s
}
