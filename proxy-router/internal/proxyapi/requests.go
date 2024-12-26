package proxyapi

import (
	"encoding/json"

	"github.com/MorpheusAIs/Morpheus-Lumerin-Node/proxy-router/internal/lib"
	"github.com/ethereum/go-ethereum/common"
)

type PingReq struct {
	ProviderAddr common.Address `json:"providerAddr" validate:"required,eth_addr"`
	ProviderURL  string         `json:"providerUrl"  validate:"required,hostname_port"`
}

type PingRes struct {
	PingMs int64 `json:"ping,omitempty"`
}

type InitiateSessionReq struct {
	User        common.Address `json:"user"        validate:"required,eth_addr"`
	Provider    common.Address `json:"provider"    validate:"required,eth_addr"`
	Spend       lib.BigInt     `json:"spend"       validate:"required,number"`
	ProviderUrl string         `json:"providerUrl" validate:"required,hostname_port"`
	BidID       common.Hash    `json:"bidId"       validate:"required,hex32"`
}

type PromptReq struct {
	Signature string          `json:"signature" validate:"required,hexadecimal"`
	Message   json.RawMessage `json:"message"   validate:"required"`
	Timestamp string          `json:"timestamp" validate:"required,timestamp"`
}

type PromptHead struct {
	SessionID lib.Hash `header:"session_id" validate:"hex32"`
	ModelID   lib.Hash `header:"model_id"   validate:"hex32"`
	ChatID    lib.Hash `header:"chat_id"    validate:"hex32"`
}

type InferenceRes struct {
	Signature lib.HexString   `json:"signature,omitempty" validate:"required,hexadecimal"`
	Message   json.RawMessage `json:"message" validate:"required"`
	Timestamp uint64          `json:"timestamp" validate:"required,timestamp"`
}

type UpdateChatTitleReq struct {
	Title string `json:"title" validate:"required"`
}

type ResultResponse struct {
	Result bool `json:"result"`
}

type ChatCompletionRequestSwaggerExample struct {
	Stream   bool `json:"stream"`
	Messages []struct {
		Role    string `json:"role" example:"user"`
		Content string `json:"content" example:"tell me a joke"`
	} `json:"messages"`
}
