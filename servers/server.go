package servers


// Server is the interface of server
// abstract a interface 
type Server interface {
	// launch server on address
	Run(address string) error
}



