package usecase

import (
	"context"
	"encoding/hex"
	"encoding/json"

	"github.com/akbariandev/pacassistant/transport/client"
	pactusPB "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/akbariandev/pacassistant/pkg/logger"

	"github.com/akbariandev/pacassistant/internal/domain"
)

const (
	defaultAccountTestnetPrefix   = "tpc1z"
	defaultValidatorTestnetPrefix = "tpc1p"
	defaultValidatorPrefix        = "pc1p"
	defaultAccountPrefix          = "pc1z"
	defaultTreasuryAccount        = "000000000000000000000000000000000000000000"
)

type AccountUsecase struct {
	pactusClient    *client.Pactus
	transactionRepo domain.TransactionRepo
	logger          logger.Logger
}

func NewAccount(pactusClient *client.Pactus, logger logger.Logger,
	transactionRepo domain.TransactionRepo,
) domain.AccountUsecase {
	return &AccountUsecase{
		pactusClient:    pactusClient,
		logger:          logger,
		transactionRepo: transactionRepo,
	}
}

func (a *AccountUsecase) GetAccount(ctx context.Context, address string) (domain.Account, error) {
	var acc domain.Account

	if len(address) < 4 {
		return domain.Account{}, status.Error(codes.InvalidArgument, "address is invalid")
	}

	/*
		first, last, err := a.transactionRepo.GetFirstAndLastTransaction(ctx, address)
		if err != nil {
			a.logger.Error(true, err.Error())
			return domain.Account{}, err
		}

		acc.FirstTransaction = first
		acc.LastTransaction = last
	*/

	isAccount := false

	switch {
	case address[:5] == defaultAccountTestnetPrefix:
		isAccount = true
	case address[:4] == defaultAccountPrefix:
		isAccount = true
	case address == defaultTreasuryAccount:
		isAccount = true
	}

	if isAccount {
		account, err := a.pactusClient.Blockchain.GetAccount(ctx, &pactusPB.GetAccountRequest{
			Address: address,
		})
		if err != nil {
			return domain.Account{}, err
		}

		acc.Type = domain.ACCOUNT
		acc.Hash = hex.EncodeToString(account.Account.Hash)
		acc.Number = account.Account.Number
		acc.Balance = account.Account.Balance
		acc.Address = address

		return acc, nil
	}

	validator, err := a.pactusClient.Blockchain.GetValidator(ctx, &pactusPB.GetValidatorRequest{
		Address: address,
	})
	if err != nil {
		return domain.Account{}, err
	}

	acc.Type = domain.VALIDATOR
	acc.Hash = hex.EncodeToString(validator.Validator.Hash)
	acc.PublicKey = validator.Validator.PublicKey
	acc.Number = validator.Validator.Number
	acc.Stake = validator.Validator.Stake
	acc.LastSortitionHeight = validator.Validator.LastSortitionHeight
	acc.UnbondingHeight = validator.Validator.UnbondingHeight
	acc.Address = address
	acc.AvailabilityScore = validator.Validator.AvailabilityScore

	return acc, nil
}

func (a *AccountUsecase) ListCommitteeValidators(ctx context.Context) ([]domain.Account, error) {
	block, err := a.pactusClient.Blockchain.GetBlockchainInfo(ctx, &pactusPB.GetBlockchainInfoRequest{})
	if err != nil {
		return nil, err
	}

	if len(block.CommitteeValidators) == 0 {
		return nil, nil
	}

	accounts := make([]domain.Account, 0, len(block.CommitteeValidators))

	b, err := json.Marshal(block.CommitteeValidators)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}
