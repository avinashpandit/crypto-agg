package bybit

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

const (
	DEFAULT_ID           = 52
	DEFAULT_TAKER_FEE    = 0.0075
	DEFAULT_MAKER_FEE    = 0.0025
	DEFAULT_TXFEE        = 0.005
	DEFAULT_WITHDRAW     = true
	DEFAULT_DEPOSIT      = true
	DEFAULT_CONFIRMATION = 2
	DEFAULT_LISTED       = true

	Name    = "bybit.api.go"
	Version = "1.0.4"
	// Https
	MAINNET       = "https://api.bybit.com"
	MAINNET_BACKT = "https://api.bytick.com"
	TESTNET       = "https://api-testnet.bybit.com"
	DEMO_ENV      = "https://api-demo.bybit.com"

	// WebSocket public channel - Mainnet
	SPOT_MAINNET    = "wss://stream.bybit.com/v5/public/spot"
	LINEAR_MAINNET  = "wss://stream.bybit.com/v5/public/linear"
	INVERSE_MAINNET = "wss://stream.bybit.com/v5/public/inverse"
	OPTION_MAINNET  = "wss://stream.bybit.com/v5/public/option"

	// WebSocket public channel - Testnet
	SPOT_TESTNET    = "wss://stream-testnet.bybit.com/v5/public/spot"
	LINEAR_TESTNET  = "wss://stream-testnet.bybit.com/v5/public/linear"
	INVERSE_TESTNET = "wss://stream-testnet.bybit.com/v5/public/inverse"
	OPTION_TESTNET  = "wss://stream-testnet.bybit.com/v5/public/option"

	// WebSocket private channel
	WEBSOCKET_PRIVATE_MAINNET = "wss://stream.bybit.com/v5/private"
	WEBSOCKET_TRADE_MAINNET   = "wss://stream.bybit.com/v5/trade"
	WEBSOCKET_PRIVATE_TESTNET = "wss://stream-testnet.bybit.com/v5/private"
	WEBSOCKET_TRADE_TESTNET   = "wss://stream-testnet.bybit.com/v5/trade"
	WEBSOCKET_PRIVATE_DEMO    = "wss://wss://stream-demo.bybit.com/v5/private"
	WEBSOCKET_TRADE_DEMO      = "wss://wss://stream-demo.bybit.com/v5/trade"

	// Deprecated: V3 is deprecated and replaced by v5
	V3_CONTRACT_PRIVATE = "wss://stream.bybit.com/contract/private/v3"
	V3_UNIFIED_PRIVATE  = "wss://stream.bybit.com/unified/private/v3"
	V3_SPOT_PRIVATE     = "wss://stream.bybit.com/spot/private/v3"

	// Globals
	timestampKey  = "X-BAPI-TIMESTAMP"
	signatureKey  = "X-BAPI-SIGN"
	apiRequestKey = "X-BAPI-API-KEY"
	recvWindowKey = "X-BAPI-RECV-WINDOW"
	signTypeKey   = "X-BAPI-SIGN-TYPE"
)
