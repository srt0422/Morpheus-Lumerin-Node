package registries

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Lumerin-protocol/Morpheus-Lumerin-Node/proxy-router/internal/contracts/marketplace"
	"github.com/Lumerin-protocol/Morpheus-Lumerin-Node/proxy-router/internal/contracts/modelregistry"
	"github.com/Lumerin-protocol/Morpheus-Lumerin-Node/proxy-router/internal/internal/interfaces"
	"github.com/Lumerin-protocol/Morpheus-Lumerin-Node/proxy-router/internal/internal/lib"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Marketplace struct {
	// config
	marketplaceAddr common.Address

	// state
	nonce uint64
	mutex lib.Mutex
	mpABI *abi.ABI

	// deps
	marketplace *marketplace.Marketplace
	client      *ethclient.Client
	log         interfaces.ILogger
}

func NewMarketplace(marketplaceAddr common.Address, client *ethclient.Client, log interfaces.ILogger) *Marketplace {
	mp, err := marketplace.NewMarketplace(marketplaceAddr, client)
	if err != nil {
		panic("invalid marketplace ABI")
	}
	mpABI, err := modelregistry.ModelRegistryMetaData.GetAbi()
	if err != nil {
		panic("invalid marketplace ABI: " + err.Error())
	}
	return &Marketplace{
		marketplace:     mp,
		marketplaceAddr: marketplaceAddr,
		client:          client,
		mpABI:           mpABI,
		mutex:           lib.NewMutex(),
		log:             log,
	}
}

func (g *Marketplace) PostModelBid(ctx context.Context, provider string, model [32]byte, pricePerSecond uint64) error {
	tx, err := g.marketplace.PostModelBid(&bind.TransactOpts{Context: ctx}, common.HexToAddress(provider), model, big.NewInt(int64(pricePerSecond)))
	if err != nil {
		return err
	}

	// Wait for the transaction receipt
	receipt, err := bind.WaitMined(context.Background(), g.client, tx)

	if err != nil {
		return err
	}

	// Find the event log
	for _, log := range receipt.Logs {
		// Check if the log belongs to the OpenSession event
		_, err := g.marketplace.ParseBidPosted(*log)

		if err != nil {
			continue // not our event, skip it
		}

		return nil
	}

	return fmt.Errorf("OpenSession event not found in transaction logs")
}

func (g *Marketplace) GetBidsByProvider(ctx context.Context, provider common.Address, offset *big.Int, limit uint8) ([][32]byte, []marketplace.Bid, error) {
	adresses, bids, err := g.marketplace.GetBidsByProvider(&bind.CallOpts{Context: ctx}, provider, offset, limit)
	if err != nil {
		return nil, nil, err
	}

	return adresses, bids, nil
}

func (g *Marketplace) GetBidsByModelAgent(ctx context.Context, modelAgentId [32]byte, offset *big.Int, limit uint8) ([][32]byte, []marketplace.Bid, error) {
	addresses, bids, err := g.marketplace.GetBidsByModelAgent(&bind.CallOpts{Context: ctx}, modelAgentId, offset, limit)
	if err != nil {
		return nil, nil, err
	}

	return addresses, bids, nil
}
