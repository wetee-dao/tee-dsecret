package types

type ReencryptSecretRequest struct {
	OrgId    string  `json:"org_id,omitempty"`
	SecretId string  `json:"secret_id,omitempty"`
	RdrPk    *PubKey `json:"rdr_pk,omitempty"`
	AcpProof []byte  `json:"acp_proof,omitempty"`
}

type ReencryptedSecretShare struct {
	OrgId    string  `json:"org_id,omitempty"`
	SecretId string  `json:"secret_id,omitempty"`
	RdrPk    *PubKey `json:"rdr_pk,omitempty"`
	Index    int32   `json:"index,omitempty"`
	XncSki   []byte  `json:"xnc_ski,omitempty"` // reencrypted share
	Chlgi    []byte  `json:"chlgi,omitempty"`   // challenge
	Proofi   []byte  `json:"proofi,omitempty"`  // proof
}
