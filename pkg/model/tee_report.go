package model

// 请求开始TEE容器
type TeeParam struct {
	// sign address
	Address string
	// report time
	Time int64
	// 0: sgx, 1: sev 2: tdx 3: sev-snp
	TeeType uint8
	// report data
	Data []byte
	// report
	Report []byte
}

type TeeReport struct {
	// 0: sgx, 1: sev 2: tdx 3: sev-snp
	TeeType uint8
	// report code signer
	CodeSigner []byte
	// report code signature
	CodeSignature []byte
	// report ProductID
	CodeProductID []byte
}
