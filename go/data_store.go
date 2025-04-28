/*
 * Grocery List Management API
 *
 * A simple REST API for managing a grocery shopping list, designed for student term projects.
 *
 * API version: 1.0.0
 * Contact: support@example.edu
 */

package openapi

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

// GroceryStore provides a thread-safe in-memory storage for grocery items
type GroceryStore struct {
	sync.RWMutex
	items  map[int64]*GroceryItem
	lastID int64
}

// Global instance of the grocery store to be used across the application
var Store = &GroceryStore{
	items:  make(map[int64]*GroceryItem),
	lastID: 0,
}

// AddItem adds a new grocery item to the store
func (s *GroceryStore) AddItem(item GroceryItemCreate) (*GroceryItem, error) {
	s.Lock()
	defer s.Unlock()

	// Increment the last ID
	s.lastID++

	// Create a new grocery item with the provided data
	now := time.Now()
	newItem := &GroceryItem{
		Name:      item.Name,
		Quantity:  item.Quantity,
		Unit:      item.Unit,
		Category:  item.Category,
		Notes:     item.Notes,
		Id:        s.lastID,
		CreatedAt: now,
		UpdatedAt: now,
		Purchased: false,
	}

	// Store the item
	s.items[s.lastID] = newItem

	return newItem, nil
}

// GetItemByID retrieves a grocery item by its ID
func (s *GroceryStore) GetItemByID(id int64) (*GroceryItem, error) {
	s.RLock()
	defer s.RUnlock()

	item, exists := s.items[id]
	if !exists {
		return nil, errors.New("item not found")
	}

	return item, nil
}

// UpdateItem updates an existing grocery item
func (s *GroceryStore) UpdateItem(id int64, update GroceryItemCreate) (*GroceryItem, error) {
	s.Lock()
	defer s.Unlock()

	item, exists := s.items[id]
	if !exists {
		return nil, errors.New("item not found")
	}

	// Update the item with the provided data
	item.Name = update.Name
	item.Quantity = update.Quantity
	item.Unit = update.Unit
	item.Category = update.Category
	item.Notes = update.Notes
	item.UpdatedAt = time.Now()

	return item, nil
}

// DeleteItem removes a grocery item from the store
func (s *GroceryStore) DeleteItem(id int64) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.items[id]; !exists {
		return errors.New("item not found")
	}

	delete(s.items, id)
	return nil
}

// ListItemsParams defines parameters for filtering and sorting grocery items
type ListItemsParams struct {
	Name      *string
	Category  *string
	Purchased *bool
	SortBy    string
	SortDesc  bool
}

// ListItems returns all grocery items in the store with optional filtering and sorting
func (s *GroceryStore) ListItems(params ListItemsParams) []*GroceryItem {
	s.RLock()
	defer s.RUnlock()

	// First, filter items based on the provided parameters
	filteredItems := make([]*GroceryItem, 0, len(s.items))
	for _, item := range s.items {
		// Filter by name if provided
		if params.Name != nil && !strings.Contains(strings.ToLower(item.Name), strings.ToLower(*params.Name)) {
			continue
		}

		// Filter by category if provided
		if params.Category != nil && !strings.Contains(strings.ToLower(item.Category), strings.ToLower(*params.Category)) {
			continue
		}

		// Filter by purchased status if provided
		if params.Purchased != nil && item.Purchased != *params.Purchased {
			continue
		}

		// Add item to filtered list if it passed all filters
		filteredItems = append(filteredItems, item)
	}

	// Apply sorting if requested
	if params.SortBy != "" {
		sort.Slice(filteredItems, func(i, j int) bool {
			// Determine sort direction (ascending by default)
			less := comparator(filteredItems[i], filteredItems[j], params.SortBy)

			// Reverse the result if descending order is requested
			if params.SortDesc {
				return !less
			}
			return less
		})
	}

	return filteredItems
}

// comparator compares two grocery items based on the specified sort field
func comparator(a, b *GroceryItem, sortBy string) bool {
	switch sortBy {
	case "createdAt":
		return a.CreatedAt.Before(b.CreatedAt)
	case "updatedAt":
		return a.UpdatedAt.Before(b.UpdatedAt)
	default:
		// Default sort by ID if an invalid sort field is provided
		return a.Id < b.Id
	}
}
