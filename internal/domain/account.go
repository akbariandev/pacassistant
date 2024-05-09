package domain

import (
	"context"
)

type Account struct {
	Type                AccountType `json:"type"`
	Address             string      `json:"address"`
	PublicKey           string      `json:"publicKey"`
	Hash                string      `json:"hash"`
	Number              int32       `json:"number"`
	Balance             int64       `json:"balance"`
	Stake               int64       `json:"stake"`
	LastBondingHeight   uint32      `json:"lastBondingHeight"`
	LastSortitionHeight uint32      `json:"lastSortitionHeight"`
	UnbondingHeight     uint32      `json:"unbondingHeight"`
	AvailabilityScore   float64     `json:"availabilityScore"`
	//FirstTransaction    ShortTransaction `json:"firstTransaction"`
	//LastTransaction     ShortTransaction `json:"lastTransaction"`
}

type AccountUsecase interface {
	GetAccount(ctx context.Context, address string) (Account, error)
	ListCommitteeValidators(ctx context.Context) ([]Account, error)
}

type AccountType uint8

const (
	ACCOUNT AccountType = iota
	VALIDATOR
	CONTRACT
)

func (a AccountType) String() string {
	switch a {
	case ACCOUNT:
		return "account"
	case VALIDATOR:
		return "validator"
	case CONTRACT:
		return "contract"
	default:
		return "account"
	}
}

type AccountDetailsResponseData struct {
	Account
	TotalTxs   int64
	PacBalance string
	PacStake   string
	TimePassed string
	Version    string

	//Txs                    []ExtendedTransaction
	ExtendedLastTxAddress  ExtendedShortTransaction
	ExtendedFirstTxAddress ExtendedShortTransaction
}
