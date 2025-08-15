package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type SavedSearch struct {
	Name      string    `json:"name"`
	Query     string    `json:"query"`
	Timestamp time.Time `json:"timestamp"`
}

const savedSearchFileName = "saved_search.json"

var savedSearchDir = func() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("could not determine home directory, using relative path for cache: %+v\n", err)
		return ".config"
	}
	return filepath.Join(home, ".config", "cashd")
}()

var savedSearchPath = func() string {
	return filepath.Join(savedSearchDir, savedSearchFileName)
}()

var loaded bool
var searches []SavedSearch

func LoadSavedSearches() ([]SavedSearch, error) {
	data, err := os.ReadFile(savedSearchPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []SavedSearch{}, nil
		}
		return nil, fmt.Errorf("failed to read saved searches file: %w", err)
	}

	if err := json.Unmarshal(data, &searches); err != nil {
		return nil, fmt.Errorf("failed to unmarshal saved searches: %w", err)
	}

	// Sort by timestamp in descending order (most recent first)
	sort.Slice(searches, func(i, j int) bool {
		return searches[i].Timestamp.After(searches[j].Timestamp)
	})

	loaded = true
	return searches, nil
}

func saveSavedSearches(s []SavedSearch) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal saved searches: %w", err)
	}

	if err := os.MkdirAll(savedSearchDir, 0755); err == nil {
		if err := os.WriteFile(savedSearchPath, data, 0644); err != nil {
			// Log caching error but don't fail the request
			log.Printf("Failed to write to %s: %+v", savedSearchPath, err)
		}
	}

	return nil
}

func AddOrUpdateSavedSearch(name, query string) error {
	if !loaded {
		LoadSavedSearches()
	}

	found := false
	for i := range searches {
		if searches[i].Name == name {
			searches[i].Query = query
			searches[i].Timestamp = time.Now()
			found = true
			break
		}
	}

	if !found {
		searches = append(searches, SavedSearch{
			Name:      name,
			Query:     query,
			Timestamp: time.Now(),
		})
	}

	return saveSavedSearches(searches)
}

func DeleteSavedSearch(name string) error {
	if !loaded {
		LoadSavedSearches()
	}

	var updatedSearches []SavedSearch
	for _, s := range searches {
		if s.Name != name {
			updatedSearches = append(updatedSearches, s)
		}
	}

	return saveSavedSearches(updatedSearches)
}

