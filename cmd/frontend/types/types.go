package types

import "time"

// ExternalService is a connection to an external service.
type ExternalService struct {
	ID          int64
	Kind        string
	DisplayName string
	Config      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeleteAt    *time.Time
}


