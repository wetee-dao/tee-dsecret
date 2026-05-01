package model

type ChainConfigType struct {
	ChainType     string
	Urls          []string
	SubnetAddress string
	CloudAddress  string
}

var LocalChainConfig = ChainConfigType{
	Urls:          []string{"wss://xiaobai.asyou.me:30001/ws"},
	SubnetAddress: "0x496806883725e8544340dd35fe743b3b8af67b19",
	CloudAddress:  "0x9faed02b7624207dc0fedfded4842be33cad4eb3",
}

var TestChainConfig = ChainConfigType{
	Urls:          []string{"wss://asset-hub-paseo.ibp.network", "wss://asset-hub-paseo-rpc.n.dwellir.com"},
	SubnetAddress: "0x4ec58c127786d767fdd968ced12f30be5f9a4bac",
	CloudAddress:  "0x239a33348e8a698cee4eb06af3780dee2ea227e3",
}

var MainChainConfig = ChainConfigType{
	Urls:          []string{"wss://polkadot-asset-hub-rpc.polkadot.io"},
	SubnetAddress: "",
	CloudAddress:  "",
}

func GetChainConfig(chainType string) *ChainConfigType {
	switch chainType {
	case "local":
		return &LocalChainConfig
	case "test":
		return &TestChainConfig
	case "main":
		return &MainChainConfig
	}
	return nil
}
