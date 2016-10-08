package backend

// Options stores our various ways to connect to backends.
// Redis for example.
type Options struct {
}

// Backend defines our various methods for connecting to data stores.
type Backend interface {
	Add(string, string) (bool, error)
	RM(string, string) (bool, error)
	RMAll(string) (bool, error)
	Get(string) ([]string, error)
	GetAll() ([]string, error)
	GetCount() (int, error)
	Init() error
	Close()
}
