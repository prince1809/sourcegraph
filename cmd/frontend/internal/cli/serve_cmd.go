package cli

import "log"

// Main is the main entrypoint for the frontend server progral¥m.
func Main() error {
	log.SetFlags(0)
	log.SetPrefix("")
	println("frontend server started")
	// Connect to the database for the frontend server program.
	return nil
}
