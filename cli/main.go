package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/Lumerin-protocol/Morpheus-Lumerin-Node/api-gateway/client"
	chat "github.com/Lumerin-protocol/Morpheus-Lumerin-Node/cli/chat"
	"github.com/ethereum/go-ethereum/common"

	dotenv "github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

const httpErrorMessage string = "internal error: %v; http status: %v"

var sessionId string

func main() {
	api_host := "http://localhost:8082"
	dotenv.Load(".env")

	if v := os.Getenv("API_HOST"); v != "" {
		api_host = v
	}

	actions := NewActions(client.NewApiGatewayClient(api_host, http.DefaultClient))
	app := &cli.App{
		Usage: "A client to call the Morpheus Lumerin API",
		Commands: []*cli.Command{
			{
				Name:    "healthcheck",
				Aliases: []string{"he"},
				Usage:   "check application health",
				Action:  actions.healthcheck,
			},
			{
				Name:    "wallet",
				Aliases: []string{"w"},
				Usage:   "morpheus wallet --privateKey <private-key>",
				Action:  actions.setupWallet,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "privateKey",
						Required: true,
					},
				},
			},
			{
				Name:    "chat",
				Aliases: []string{},
				Usage:   "open interactive chat client",
				Action:  actions.startChat,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "sessionId",
						Required: false,
					},
				},
			},
			{
				Name:    "blockchainModels",
				Aliases: []string{"bm"},
				Usage:   "list models",
				Action:  actions.blockchainModels,
			},
			{
				Name:    "openBlockchainSession",
				Usage:   "open a blockchain session",
				Aliases: []string{"obs"},
				Action:  actions.openBlockchainSession,
			},
			{
				Name:    "closeBlockchainSession",
				Aliases: []string{"cbs"},
				Usage:   "close a blockchain session",
				Action:  actions.closeBlockchainSession,
			},
			{
				Name:    "getAllowance",
				Aliases: []string{"ga"},
				Usage:   "get allowance",
				Action:  actions.getAllowance,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "spender",
						Required: true,
					},
				},
			},
			{
				Name:    "approveAllowance",
				Aliases: []string{"aa"},
				Usage:   "approve allowance",
				Action:  actions.approveAllowance,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "spender",
						Required: true,
					},
					&cli.Uint64Flag{
						Name:     "amount",
						Required: true,
					},
				},
			},
			{
				Name:    "proxyRouterConfig",
				Aliases: []string{"prc"},
				Usage:   "view proxy router config",
				Action:  actions.proxyRouterConfig,
			},
			{
				Name:    "proxyRouterFiles",
				Aliases: []string{"prf"},
				Usage:   "get the files associated with the proxy router pid",
				Action:  actions.proxyRouterFiles,
			},
			{
				Name:    "createChatCompletions",
				Aliases: []string{"ccc"},
				Usage:   "create chat completions by sending a prompt to the AI engine",
				Action:  actions.createChatCompletions,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "prompt",
						Required: true,
					},
					&cli.StringFlag{
						Name: "messages",
					},
				},
			},
			{
				Name:    "initiateProxySession",
				Aliases: []string{"ips"},
				Usage:   "initiate a proxy session",
				Action:  actions.initiateProxySession,
			},
			{
				Name:    "blockchainProviders",
				Aliases: []string{"bp"},
				Usage:   "list blockchain providers",
				Action:  actions.blockchainProviders,
			},
			{
				Name:    "createBlockchainProvider",
				Aliases: []string{"bpc"},
				Usage:   "create a blockchain provider",
				Action:  actions.createBlockchainProvider,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "address",
						Required: true,
					},
					&cli.Uint64Flag{
						Name:     "stake",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "endpoint",
						Required: true,
					},
				},
			},
			{
				Name:    "blockchainProviderBid",
				Aliases: []string{"bpb"},
				Usage:   "list provider bids",
				Action:  actions.blockchainProvidersBids,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "address",
						Required: true,
					},
					&cli.Int64Flag{
						Name:  "offset",
						Value: 0,
					},
					&cli.UintFlag{
						Name:  "limit",
						Value: 10,
					},
				},
			},
			{
				Name:    "createBlockchainProviderBid",
				Aliases: []string{"cbpb"},
				Usage:   "createBlockchainProviderBid {{model}}",
				Action:  actions.createBlockchainProviderBid,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "model",
						Required: true,
					},
					&cli.Uint64Flag{
						Name:     "pricePerSecond",
						Required: true,
					},
				},
			},
			{
				Name:    "createAndStreamChatCompletions",
				Aliases: []string{"csc"},
				Usage:   "create and stream chat completions",
				Action:  actions.createAndStreamChatCompletions,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "prompt",
						Required: true,
					},
					&cli.StringFlag{
						Name: "messages",
					},
				},
			},
			{
				Name:    "createAndStreamSessionChatCompletions",
				Aliases: []string{"cssc"},
				Usage:   "create and stream session chat completions",
				Action:  actions.createChatCompletions,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "prompt",
						Required: true,
					},
					&cli.StringFlag{
						Name: "messages",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type actions struct {
	client *client.ApiGatewayClient
}

func NewActions(c *client.ApiGatewayClient) *actions {
	return &actions{client: c}
}

var localProxyRouterUrl string
var contractAddress string
var userWalletAddress string
var bidId string
var provider string
var providerEndpoint string
var stake int

func (a *actions) setupWallet(cCtx *cli.Context) error {
	return a.client.CreateWallet(cCtx.Context, cCtx.String("privateKey"))
}

func (a *actions) startChat(cCtx *cli.Context) error {
	bidId := cCtx.String("bid")

	if bidId == "" {
		bidId = "0x921d3b1f984a7df85f9a640f834adb49f1c058a6564f914699fa63c1518a919e"
	}

	providerEndpoint := cCtx.String("providerEndpoint")

	if providerEndpoint == "" {
		providerEndpoint = "ironman1.dev.lumerin.io:3333"
	}

	sessionDuration := cCtx.Uint("sessionDuration")

	if sessionDuration == 0 {
		sessionDuration = 300
	}

	sessionRequest := &client.SessionRequest{
		BidId:            bidId,
		ProviderEndpoint: providerEndpoint,
		SessionDuration:  sessionDuration,
	}

	session, err := a.client.OpenSession(cCtx.Context, sessionRequest)

	if err != nil {
		chat.Run("")
	}

	// chat.Run("")
	chat.Run(session.SessionId)

	return nil
}

func (a *actions) getAllowance(cCtx *cli.Context) error {
	spender := cCtx.String("spender")
	res, err := a.client.GetAllowance(cCtx.Context, spender)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func (a *actions) approveAllowance(cCtx *cli.Context) error {
	spender := cCtx.String("spender")
	amount := cCtx.Uint64("amount")
	res, err := a.client.ApproveAllowance(cCtx.Context, spender, amount)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func (a *actions) healthcheck(cCtx *cli.Context) error {
	res, err := a.client.HealthCheck(cCtx.Context)
	if err != nil {
		return err
	}
	// Output the result of the healthcheck
	jsonData, err := json.Marshal(res)
	fmt.Println(string(jsonData))
	return nil
}

func (a *actions) proxyRouterConfig(cCtx *cli.Context) error {
	config, err := a.client.GetProxyRouterConfig(cCtx.Context)
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(config)
	fmt.Println(string(jsonData))
	return nil
}

func (a *actions) proxyRouterFiles(cCtx *cli.Context) error {
	files, err := a.client.GetProxyRouterFiles(cCtx.Context)

	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(files)
	fmt.Println(string(jsonData))
	return nil
}

func (a *actions) createChatCompletions(cCtx *cli.Context) error {
	prompt := cCtx.String("prompt")
	var messages []client.ChatCompletionMessage
	json.Unmarshal([]byte(cCtx.String("messages")), &messages)

	completion, err := a.client.Prompt(cCtx.Context, prompt, messages)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(completion)
	fmt.Println(string(jsonData))
	return nil
}

// todo: retrieve session id
func (a *actions) createAndStreamChatCompletions(cCtx *cli.Context) error {
	prompt := cCtx.String("prompt")
	var messages []*client.ChatCompletionMessage
	json.Unmarshal([]byte(cCtx.String("messages")), &messages)

	completion, err := a.client.PromptStream(cCtx.Context, prompt, "", func(msg *client.ChatCompletionStreamResponse) error {
		fmt.Println(msg)
		return nil
	})

	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(completion)
	fmt.Println(string(jsonData))
	return nil
}

func (a *actions) initiateProxySession(cCtx *cli.Context) error {
	session, err := a.client.InitiateSession(cCtx.Context)
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(session)
	fmt.Println(string(jsonData))
	return nil
}

func (a *actions) blockchainProviders(cCtx *cli.Context) error {
	providers, err := a.client.GetAllProviders(cCtx.Context)

	if err != nil {
		return err
	}

	for _, item := range providers["providers"].([]interface{}) {
		provider := item.(map[string]interface{})
		// fmt.Println(provider)
		// fmt.Println(reflect.TypeOf(provider))
		fmt.Println(provider["Address"], " - ", provider["Endpoint"])
	}

	// jsonData, err := json.Marshal(providers)

	// if err != nil {
	// 	return err
	// }

	return nil
}

func (a *actions) createBlockchainProvider(cCtx *cli.Context) error {
	stake := cCtx.Uint64("stake")
	endpoint := cCtx.String("endpoint")

	providers, err := a.client.CreateNewProvider(cCtx.Context, stake, endpoint)

	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(providers)
	fmt.Println(string(jsonData))
	return nil
}

func (a *actions) createBlockchainProviderBid(cCtx *cli.Context) error {
	model := cCtx.String("model")
	pricePerSecond := cCtx.Uint64("pricePerSecond")

	_, err := a.client.CreateNewProviderBid(cCtx.Context, model, pricePerSecond)

	if err != nil {
		return err
	}

	fmt.Println("bid created for model ", model)
	return nil
}

type Bid struct {
	Id             string
	Provider       common.Address
	ModelAgentId   string
	PricePerSecond *big.Int
	Nonce          *big.Int
	CreatedAt      *big.Int
	DeletedAt      *big.Int
}

func (a *actions) blockchainProvidersBids(cCtx *cli.Context) error {
	address := cCtx.String("address")
	offset := cCtx.Int64("offset")
	limit := cCtx.Uint("limit")

	bidsResult, err := a.client.GetBidsByProvider(cCtx.Context, address, big.NewInt(offset), uint8(limit))

	if err != nil {
		return err
	}

	bidsMap := bidsResult.(map[string]interface{})

	bids := bidsMap["bids"].([]interface{})

	for _, item := range bids {
		bid := item.(map[string]interface{})
		fmt.Println("Bid: ", bid["Id"])
		fmt.Println("\t- price per second: ", bid["PricePerSecond"])
		fmt.Println("\t- model: ", bid["ModelAgentId"])
	}

	return nil
}

func (a *actions) blockchainModels(cCtx *cli.Context) error {
	modelsResult, err := a.client.GetAllModels(cCtx.Context)
	if err != nil {
		return err
	}
	modelsMap := modelsResult.(map[string]interface{})

	models := modelsMap["models"].([]interface{})

	// fmt.Println("models: ", models)
	for _, item := range models {
		model := item.(map[string]interface{})
		fmt.Println("Model: ", model["Name"], " - ", model["Id"])
		fmt.Println("\t- provider: ", model["Owner"], "\n")
	}

	return nil
}

func (a *actions) openBlockchainSession(cCtx *cli.Context) error {

	session, err := a.client.OpenStakeSession(cCtx.Context, &client.SessionStakeRequest{
		Approval:    cCtx.String("approval"),
		ApprovalSig: cCtx.String("approvalSig"),
		Stake:       cCtx.Uint64("stake"),
	})

	if err != nil {
		return err
	}

	fmt.Println("session opened: ", session.SessionId)
	return nil
}

func (a *actions) closeBlockchainSession(cCtx *cli.Context) error {
	err := a.client.CloseSession(cCtx.Context)
	if err != nil {
		return err
	}

	// TODO: Output a message indicating the blockchain session was closed and showing relevant data for the session
	// jsonData, err := json.Marshal(bids)
	// fmt.Println(string(jsonData))
	return nil
}
