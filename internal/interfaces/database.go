package interfaces

type Database interface {
	Connect() error
	Disconnect() error
}
