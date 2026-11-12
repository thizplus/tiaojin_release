package main

import (
	"sync"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IDMapping stores the mapping between MongoDB ObjectID and PostgreSQL UUID
type IDMapping struct {
	mu       sync.RWMutex
	mappings map[string]map[string]uuid.UUID // collection -> objectId -> uuid
}

var idMapping *IDMapping

// InitMapping initializes the ID mapping storage
func InitMapping() {
	idMapping = &IDMapping{
		mappings: make(map[string]map[string]uuid.UUID),
	}
}

// SetMapping stores a mapping between MongoDB ObjectID and PostgreSQL UUID
func SetMapping(collection string, objectID primitive.ObjectID, newUUID uuid.UUID) {
	idMapping.mu.Lock()
	defer idMapping.mu.Unlock()

	if idMapping.mappings[collection] == nil {
		idMapping.mappings[collection] = make(map[string]uuid.UUID)
	}
	idMapping.mappings[collection][objectID.Hex()] = newUUID
}

// GetMapping retrieves the UUID for a given MongoDB ObjectID
func GetMapping(collection string, objectID primitive.ObjectID) (uuid.UUID, bool) {
	idMapping.mu.RLock()
	defer idMapping.mu.RUnlock()

	if collMap, ok := idMapping.mappings[collection]; ok {
		if id, ok := collMap[objectID.Hex()]; ok {
			return id, true
		}
	}
	return uuid.Nil, false
}

// GetMappingByHex retrieves the UUID for a given MongoDB ObjectID hex string
func GetMappingByHex(collection string, objectIDHex string) (uuid.UUID, bool) {
	idMapping.mu.RLock()
	defer idMapping.mu.RUnlock()

	if collMap, ok := idMapping.mappings[collection]; ok {
		if id, ok := collMap[objectIDHex]; ok {
			return id, true
		}
	}
	return uuid.Nil, false
}

// GetOrCreateMapping retrieves existing mapping or creates a new one
func GetOrCreateMapping(collection string, objectID primitive.ObjectID) uuid.UUID {
	if existing, ok := GetMapping(collection, objectID); ok {
		return existing
	}

	newUUID := uuid.New()
	SetMapping(collection, objectID, newUUID)
	return newUUID
}

// GetMappingCount returns the count of mappings for a collection
func GetMappingCount(collection string) int {
	idMapping.mu.RLock()
	defer idMapping.mu.RUnlock()

	if collMap, ok := idMapping.mappings[collection]; ok {
		return len(collMap)
	}
	return 0
}

// Collection names for mapping
const (
	CollectionUsers    = "users"
	CollectionCities   = "cities"
	CollectionTours    = "tours"
	CollectionPeriods  = "periods"
	CollectionBookings = "bookings"
	CollectionSales    = "sales"
	CollectionBanners  = "banners"
)
