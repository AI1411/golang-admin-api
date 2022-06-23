//go:generate mockgen -source=uuid.go -destination=mock_model/uuid.go
package models

import "github.com/google/uuid"

type UUIDGenerator interface {
	GenerateUUID() string
}

type RandomUUIDGenerator struct{}

func (g *RandomUUIDGenerator) GenerateUUID() string {
	newUUID, _ := uuid.NewRandom()
	return newUUID.String()
}
