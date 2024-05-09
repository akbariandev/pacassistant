package config

import "go.mongodb.org/mongo-driver/bson/primitive"

// database collection.
const (
	CollectionBlock            = "block"
	CollectionTransaction      = "transaction"
	CollectionSettings         = "settings.html"
	CollectionCommitteeHistory = "committee_history"
	CollectionPeer             = "peer"
)

// cache keys
const (
	TotalCountCacheKey       = "TotalCountCacheKey"
	TotalTransactionCacheKey = "TotalTransactionCacheKey"
	TransactionsCacheKey     = "TransactionsCacheKey"
	BlocksCacheKey           = "BlockCacheKey"
	CirculationCacheKey      = "CirculationCacheKey"
	PriceCacheKey            = "PriceCacheKey"
)

const (
	// MaxSupply is base on https://pactus.org/about/faq/#genesis_allocation
	MaxSupply = 42_000_000
	// ReservedCoin Treasury
	ReservedCoin = 21_000_000
)

type ExtraData struct {
	MainSettingsID primitive.ObjectID `yaml:"main_settings_id" json:"main_settings_id"`
	RootNode       []string           `yaml:"root_node" json:"root_node"`
}
