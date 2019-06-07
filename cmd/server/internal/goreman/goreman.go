package goreman

import "fmt"

// Start starts up the Procfile
func Start(rpcAddr string, contents []byte) error {
	fmt.Println("Started goreman")
	return nil
}
