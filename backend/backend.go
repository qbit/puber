package backend

// Backend defines our various methods for connecting to data stores.
type Backend interface {
	Add(string, string) (bool, error)
	RM(string, string) (bool, error)
	RMAll(string) (bool, error)
	Get(string) ([]string, error)
	GetAll() (map[string][]string, error)
	GetCount() (int, error)
	Init()
}
