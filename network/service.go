package network

type NetworkService interface {
	Download() error
	Upload() error
}
