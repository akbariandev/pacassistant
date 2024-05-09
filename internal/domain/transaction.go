package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Hash        string             `bson:"hash" json:"hash"`
	BlockHeight uint32             `bson:"block_height" json:"blockHeight"`
	Version     int32              `bson:"version,omitempty" json:"version"`
	Type        TransactionType    `bson:"type" json:"type"`
	From        string             `bson:"from,omitempty" json:"from"`
	To          string             `bson:"to,omitempty" json:"to"`
	Value       int64              `bson:"value,omitempty" json:"value"`
	Fee         int64              `bson:"fee,omitempty" json:"fee"`
	Memo        string             `bson:"memo,omitempty" json:"memo"`
	Signature   string             `bson:"signature,omitempty" json:"signature"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
}

type ShortTransaction struct {
	Hash      string    `bson:"hash" json:"hash"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
}

type ExtendedTotalTransactionInfo struct {
	TotalTransactionInfo
	HTotalValue        string
	HTotalFee          string
	HTotalCount        string
	HTotal24HoursCount string
}

type TransactionsDataResponse struct {
	TxsInfo ExtendedTotalTransactionInfo
	//Txs     []ExtendedTransaction
	Version string
}

type TransactionDetailDataResponse struct {
	Transaction
	PacValue   string
	PacFee     string
	TimePassed string
	Version    string
}

type ExtendedShortTransaction struct {
	ShortHash  string
	TimePassed string
}

type TransactionRepo interface {
	Repository[Transaction]
	GetByHash(ctx context.Context, hash string) (Transaction, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetTotalInfo(ctx context.Context) (TotalTransactionInfo, error)
	GetLatest(ctx context.Context, size int32) ([]Transaction, error)
	GetAccountTotalTransactions(ctx context.Context, address string) (int64, error)
	GetFirstAndLastTransaction(ctx context.Context, address string) (first, last ShortTransaction, err error)
}

type TransactionUsecase interface {
	GetTransaction(ctx context.Context, id primitive.ObjectID) (Transaction, error)
	GetTransactionByHash(ctx context.Context, hash string) (Transaction, error)
	TotalTransactions(ctx context.Context) (int64, error)
	TotalAccountTransactions(ctx context.Context, address string) (int64, error)
	GetCachedTransactions(ctx context.Context) ([]Transaction, error)
}

type TransactionType uint8

const (
	TRANSFER TransactionType = iota + 1
	BOND
	SORTITION
	UNBOND
	WITHDRAW
)

func (t TransactionType) String() string {
	switch t {
	case TRANSFER:
		return "transfer"
	case BOND:
		return "bond"
	case SORTITION:
		return "sortition"
	case UNBOND:
		return "unbond"
	case WITHDRAW:
		return "withdraw"
	default:
		return "unspecified"
	}
}
