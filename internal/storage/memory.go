package storage

import (
	"fmt"
	"sync"
	"time"
)

// Memory is an ephemeral storage implementation.
// All the links are stored in memory and indexed by their slug.
type Memory struct {
	links   map[string]*Link
	slugger SlugGenerator

	mutex sync.RWMutex
}

const maxrecursion = 5

// NewMemoryStorage instantiates a new in-memory storage.
// By default, it uses an UUIDSlugGenerator for generating slugs, but this can be
// clustomized via MemoryOptions.
func NewMemoryStorage(opts ...MemoryOption) (*Memory, error) {
	ms := Memory{
		links: make(map[string]*Link, 0),
	}

	for _, opt := range opts {
		if err := opt(&ms); err != nil {
			return nil, err
		}
	}

	if ms.slugger == nil {
		ms.slugger = &UUIDSlugGenerator{}
	}

	return &ms, nil
}

// CreateLink to the target url received as parameter.
// If the slug param is not nil, that slug will be used instead of generating
// a new random one, which allows for custom shortened links.
func (m *Memory) CreateLink(target string, slug *string) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if slug == nil || *slug == "" {
		s, err := m.genslug()
		if err != nil {
			return "", err
		}
		slug = &s
	}

	m.links[*slug] = &Link{
		Slug:      *slug,
		Target:    target,
		Histogram: make(map[string]uint64),
	}

	return *slug, nil
}

// GetLink returns the Link object associated with the specified slug.
func (m *Memory) GetLink(slug string) (link *Link, err error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	link, found := m.links[slug]
	if !found {
		err = fmt.Errorf("no link with slug %s found", slug)
		return
	}

	return link, nil
}

// DeleteLink removes a link from the database.
func (m *Memory) DeleteLink(slug string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.links, slug)

	return nil
}

// AllLinks returns all links stored in the database.
func (m *Memory) AllLinks() []*Link {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]*Link, 0, len(m.links))
	for _, link := range m.links {
		result = append(result, link)
	}

	return result
}

// GetTarget returns the full url associated with a specific slug.
// It returns an error if the slug isn't found on the database.
func (m *Memory) GetTarget(slug string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	link, found := m.links[slug]
	if !found {
		return "", fmt.Errorf("no url found for slug: %s", slug)
	}

	return link.Target, nil
}

// RegisterHit increments the hit counter for the specific slug and the current day.
// If the slug doesn't exist on the database this is noop.
func (m *Memory) RegisterHit(slug string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	bucket := now.Format("2006-01-02")

	if link := m.links[slug]; link != nil {
		link.Hits++
		link.Histogram[bucket]++
	}
}

// genslug generates a slug using the slugger function.
// If the generated slug already exists, it will keep generating slugs
// until it finds a unique one, or the maxrecursion limit is hit.
func (m *Memory) genslug() (string, error) {
	var slug string
	var err error
	iteration := 0

	for {
		if iteration > maxrecursion {
			return "", fmt.Errorf("could not generate a unique slug in %d attempts", maxrecursion)
		}

		slug, err = m.slugger.Random()
		if err != nil {
			return "", fmt.Errorf("error generating slug: %w", err)
		}

		_, dupe := m.links[slug]
		if !dupe {
			break
		}

		iteration++
	}

	return slug, err
}

type MemoryOption func(ms *Memory) error

func WithSlugGenerator(slugger SlugGenerator) MemoryOption {
	return func(ms *Memory) error {
		ms.slugger = slugger
		return nil
	}
}
