package main

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/rpc/client/http"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gammtypes "github.com/osmosis-labs/osmosis/v25/x/gamm/types"
)

type SwapEvent struct {
	TokensIn sdk.Coin
	// TODO: translate the tokens in amount into USDC value
	// Exampe API call to get the price:
	// curl -X 'GET' \ 'https://sqsprod.osmosis.zone/tokens/prices?base=uosmo&humanDenoms=false' \ -H 'accept: application/json'
	USDCValue     math.Int
	SenderAddress string
	PoolID        uint64
	Height        int64
}

const (
	archiveNodeAddress       = "https://rpc.archive.osmosis.zone:443"
	pricesAPIAddress         = "https://api.coingecko.com/api/v3/simple/price?ids=osmosis&vs_currencies=usd"
	startBlock         int64 = 17777000
	endBlock                 = 17777050
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in main", r)
		}
	}()

	fmt.Println("Hello, Osmosis!")

	rpcClient, err := http.New(archiveNodeAddress, "/websocket")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	for i := startBlock; i <= endBlock; i++ {
		curHeight := i
		results, err := rpcClient.BlockResults(context.Background(), &curHeight)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		swapEvent := SwapEvent{
			Height: curHeight,
		}

		for _, result := range results.TxsResults {
			events := result.GetEvents()
			for _, event := range events {

				// Find swap event
				if event.GetType() == gammtypes.TypeEvtTokenSwapped {

					attributes := event.GetAttributes()

					for _, attr := range attributes {
						// Parse tokens in
						if string(attr.GetKey()) == gammtypes.AttributeKeyTokensIn {

							inCoinStr := attr.GetValue()

							inCoin, err := sdk.ParseCoinNormalized(inCoinStr)
							if err != nil {
								fmt.Println(err)
								panic(err)
							}

							swapEvent.TokensIn = inCoin
						}

						// TODO:
						// - Parse pool id from swap event
						// - Parse sende address from swap event
					}
				}
			}
		}

		fmt.Println(swapEvent)

		time.Sleep(500 * time.Millisecond)
	}

	// TODO:
	// - Complete indexing logic
	//    * Ensure that indexing can start from a specified height, catch up to the tip and continue.
	//    according to the chain's progress.
	//    * Handle de-duplication as needed or overwrite existing data.
	//    * Choose a suitable storage solution for the data.
	//    * Convert tokens in amoun into USDC value for indexing volume.
	// - Expose web API
	//    * Implement an API to construct a merkle tree of swaps made from the indexed data
	// - Bonus
	//    * Infrastructure / cloud automation and a fully-functional service.
}
