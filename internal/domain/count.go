package domain

import (
	"context"
)

type Counts struct {
	TotalBlocks       uint32 `json:"totalBlocks"`
	TotalTransactions uint64 `json:"totalTransactions"`
	TotalAccounts     uint64 `json:"totalAccounts"`
	TotalValidators   uint64 `json:"totalValidators"`
	TotalPeers        uint64 `json:"totalPeers"`
}

type TotalTransactionInfo struct {
	TotalCount        uint64 `bson:"totalCount"`
	Total24HoursCount uint64 `bson:"total24HoursCount"`
	TotalValue        uint64 `bson:"totalValue"`
	TotalFee          uint64 `bson:"totalFee"`
}

type CountUsecase interface {
	GetTotalCount(ctx context.Context) (Counts, error)
	GetTotalTransactionInfo(ctx context.Context) (TotalTransactionInfo, error)
}
