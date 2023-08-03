package app

// App struct
type App struct {
	repository Repository
}

// New constructor for App
func New(repository Repository) *App {
	return &App{
		repository: repository,
	}
}
