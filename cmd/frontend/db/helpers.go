package db

// LimitOffset specifies SQL LIMIT and OFFSET counts. A pointer to it is typically embedded in
// structs that need to performs SQL queries with LIMIT and OFFSET.
type LimitOffset struct {
	Limit  int // SQL LIMIT count
	Offset int // SQL OFFSET count
}
