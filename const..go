package main

const (
	// SatInBTC represents number of satoshis in 1 bitcoin
	SatInBTC = uint64(100000000)

	// AHRKDec represents number of decimals in AHRK
	AHRKDec = uint64(1000000)

	// WavesNodeURL is an URL for Waves Node
	WavesNodeURL = "https://nodes.wavesnodes.com"

	// WavesMonitorTick interval in seconds
	WavesMonitorTick = 10

	// PricesURL is URL for crypo prices
	PricesURL = "https://min-api.cryptocompare.com/data/price?fsym=WAVES&tsyms=BTC,ETH,HRK,USD,EUR,NGN,JPY"

	// PricesHNBURL is URL for fiat prices
	PricesHNBURL = "https://api.hnb.hr/tecajn/v1?valuta=USD"

	// AHRK Address
	AHRKAddress = "3PPc3AP75DzoL8neS4e53tZ7ybUAVxk2jAb"

	// AHRKId is AHRK asset id
	AHRKId = "Gvs59WEEXVAQiRZwisUosG7fVNr8vnzS8mjkgqotrERT"

	// WavesFee represents fee amount in Waves
	WavesFee = 100000

	// AHRKFee represents fee amount in AHRK
	AHRKFee = 50000

	// TelPollerTimeout is Telegram poller timeout in seconds
	TelPollerTimeout = 30

	// TelAnonOps group for error logging
	TelAnonOps = -1001213539865

	// TelAnonTeam group for team messages
	TelAnonTeam = -1001280228955

	// TelKriptokuna group for Kriptokuna messages
	TelKriptokuna = -1001456424919

	// TelAnonEuro group for AnonEuro messages
	TelAnonEuro = -1001735622646
)
