package main

// This script generate markdown formatted output containing descriptions of
// the current database schema, obtained from postgres. The correct PGHOST,
// PGPORT, PGUSER etc. env variable must be set to run this script.
//
// First CLI argument is an optional filename to write the output to.

func main() {
	const dbname = "schemadoc-gen-temp"
}
