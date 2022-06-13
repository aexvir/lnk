package storage

import "github.com/google/uuid"

type SlugGenerator interface {
	Random() (string, error)
}

type UUIDSlugGenerator struct{}

func (us *UUIDSlugGenerator) Random() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return id.String()[:6], nil
}
