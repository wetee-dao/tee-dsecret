package subnet

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
)

type Subnet struct {
	ChainClient *chain.ChainClient
	Address     types.H160
}

func (c *Subnet) Client() *chain.ChainClient {
	return c.ChainClient
}
func (c *Subnet) ContractAddress() types.H160 {
	return c.Address
}

func (c *Subnet) DryRunSetCode(
	code_hash types.H256, params chain.DryRunCallParams,
) (util.Result[util.NullTuple, Error], error) {
	v, err := chain.DryRun[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "set_code",
			Args:     []any{code_hash},
		},
	)
	return *v, err
}

func (c *Subnet) CallSetCode(
	code_hash types.H256, params chain.CallParams,
) error {
	err := chain.Call(
		c,
		params.Signer,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "set_code",
			Args:     []any{code_hash},
		},
	)
	return err
}

func (c *Subnet) QueryBootNodes(
	params chain.DryRunCallParams,
) (util.Result[[]SecretNode, Error], error) {
	v, err := chain.DryRun[util.Result[[]SecretNode, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "boot_nodes",
			Args:     []any{},
		},
	)
	return *v, err
}

func (c *Subnet) DryRunSetBootNodes(
	nodes []types.U128, params chain.DryRunCallParams,
) (util.Result[util.NullTuple, Error], error) {
	v, err := chain.DryRun[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "set_boot_nodes",
			Args:     []any{nodes},
		},
	)
	return *v, err
}

func (c *Subnet) CallSetBootNodes(
	nodes []types.U128, params chain.CallParams,
) error {
	err := chain.Call(
		c,
		params.Signer,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "set_boot_nodes",
			Args:     []any{nodes},
		},
	)
	return err
}

func (c *Subnet) QueryWorkers(
	params chain.DryRunCallParams,
) (util.Result[[]K8sCluster, Error], error) {
	v, err := chain.DryRun[util.Result[[]K8sCluster, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "workers",
			Args:     []any{},
		},
	)
	return *v, err
}

func (c *Subnet) DryRunWorkerRegister(
	name []byte, validator_id AccountId, p2p_id AccountId, ip Ip, port uint32, level byte, params chain.DryRunCallParams,
) (util.Result[types.U128, Error], error) {
	v, err := chain.DryRun[util.Result[types.U128, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "worker_register",
			Args:     []any{name, validator_id, p2p_id, ip, port, level},
		},
	)
	return *v, err
}

func (c *Subnet) CallWorkerRegister(
	name []byte, validator_id AccountId, p2p_id AccountId, ip Ip, port uint32, level byte, params chain.CallParams,
) error {
	err := chain.Call(
		c,
		params.Signer,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "worker_register",
			Args:     []any{name, validator_id, p2p_id, ip, port, level},
		},
	)
	return err
}

func (c *Subnet) DryRunWorkerMortgage(
	id types.U128, cpu uint32, mem uint32, cvm_cpu uint32, cvm_mem uint32, disk uint32, gpu uint32, deposit types.U256, params chain.DryRunCallParams,
) (util.Result[types.U128, Error], error) {
	v, err := chain.DryRun[util.Result[types.U128, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "worker_mortgage",
			Args:     []any{id, cpu, mem, cvm_cpu, cvm_mem, disk, gpu, deposit},
		},
	)
	return *v, err
}

func (c *Subnet) CallWorkerMortgage(
	id types.U128, cpu uint32, mem uint32, cvm_cpu uint32, cvm_mem uint32, disk uint32, gpu uint32, deposit types.U256, params chain.CallParams,
) error {
	err := chain.Call(
		c,
		params.Signer,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "worker_mortgage",
			Args:     []any{id, cpu, mem, cvm_cpu, cvm_mem, disk, gpu, deposit},
		},
	)
	return err
}

func (c *Subnet) DryRunWorkerUnmortgage(
	id types.U128, mortgage_id types.U128, params chain.DryRunCallParams,
) (util.Result[types.U128, Error], error) {
	v, err := chain.DryRun[util.Result[types.U128, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "worker_unmortgage",
			Args:     []any{id, mortgage_id},
		},
	)
	return *v, err
}

func (c *Subnet) CallWorkerUnmortgage(
	id types.U128, mortgage_id types.U128, params chain.CallParams,
) error {
	err := chain.Call(
		c,
		params.Signer,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "worker_unmortgage",
			Args:     []any{id, mortgage_id},
		},
	)
	return err
}

func (c *Subnet) DryRunWorkerStop(
	id types.U128, params chain.DryRunCallParams,
) (util.Result[types.U128, Error], error) {
	v, err := chain.DryRun[util.Result[types.U128, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "worker_stop",
			Args:     []any{id},
		},
	)
	return *v, err
}

func (c *Subnet) CallWorkerStop(
	id types.U128, params chain.CallParams,
) error {
	err := chain.Call(
		c,
		params.Signer,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "worker_stop",
			Args:     []any{id},
		},
	)
	return err
}

func (c *Subnet) QuerySecrets(
	params chain.DryRunCallParams,
) (util.Result[[]SecretNode, Error], error) {
	v, err := chain.DryRun[util.Result[[]SecretNode, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "secrets",
			Args:     []any{},
		},
	)
	return *v, err
}

func (c *Subnet) DryRunSecretRegister(
	name []byte, validator_id AccountId, p2p_id AccountId, ip Ip, port uint32, params chain.DryRunCallParams,
) (util.Result[types.U128, Error], error) {
	v, err := chain.DryRun[util.Result[types.U128, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "secret_register",
			Args:     []any{name, validator_id, p2p_id, ip, port},
		},
	)
	return *v, err
}

func (c *Subnet) CallSecretRegister(
	name []byte, validator_id AccountId, p2p_id AccountId, ip Ip, port uint32, params chain.CallParams,
) error {
	err := chain.Call(
		c,
		params.Signer,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "secret_register",
			Args:     []any{name, validator_id, p2p_id, ip, port},
		},
	)
	return err
}

func (c *Subnet) DryRunSecretDeposit(
	id types.U128, deposit types.U256, params chain.DryRunCallParams,
) (util.Result[util.NullTuple, Error], error) {
	v, err := chain.DryRun[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "secret_deposit",
			Args:     []any{id, deposit},
		},
	)
	return *v, err
}

func (c *Subnet) CallSecretDeposit(
	id types.U128, deposit types.U256, params chain.CallParams,
) error {
	err := chain.Call(
		c,
		params.Signer,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "secret_deposit",
			Args:     []any{id, deposit},
		},
	)
	return err
}

func (c *Subnet) DryRunSecretJoin(
	id types.U128, params chain.DryRunCallParams,
) (util.Result[util.NullTuple, Error], error) {
	v, err := chain.DryRun[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "secret_join",
			Args:     []any{id},
		},
	)
	return *v, err
}

func (c *Subnet) CallSecretJoin(
	id types.U128, params chain.CallParams,
) error {
	err := chain.Call(
		c,
		params.Signer,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "secret_join",
			Args:     []any{id},
		},
	)
	return err
}

func (c *Subnet) DryRunSecretDelete(
	id types.U128, params chain.DryRunCallParams,
) (util.Result[util.NullTuple, Error], error) {
	v, err := chain.DryRun[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "secret_delete",
			Args:     []any{id},
		},
	)
	return *v, err
}

func (c *Subnet) CallSecretDelete(
	id types.U128, params chain.CallParams,
) error {
	err := chain.Call(
		c,
		params.Signer,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "secret_delete",
			Args:     []any{id},
		},
	)
	return err
}

func (c *Subnet) QueryEpoch(
	params chain.DryRunCallParams,
) (Tuple_69, error) {
	v, err := chain.DryRun[Tuple_69](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "epoch",
			Args:     []any{},
		},
	)
	return *v, err
}

func (c *Subnet) DryRunNextEpoch(
	params chain.DryRunCallParams,
) (util.Result[util.NullTuple, Error], error) {
	v, err := chain.DryRun[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "next_epoch",
			Args:     []any{},
		},
	)
	return *v, err
}

func (c *Subnet) CallNextEpoch(
	params chain.CallParams,
) error {
	err := chain.Call(
		c,
		params.Signer,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "next_epoch",
			Args:     []any{},
		},
	)
	return err
}

func (c *Subnet) DryRunNextEpochWithGov(
	params chain.DryRunCallParams,
) (util.Result[util.NullTuple, Error], error) {
	v, err := chain.DryRun[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "next_epoch_with_gov",
			Args:     []any{},
		},
	)
	return *v, err
}

func (c *Subnet) CallNextEpochWithGov(
	params chain.CallParams,
) error {
	err := chain.Call(
		c,
		params.Signer,
		params.Amount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "next_epoch_with_gov",
			Args:     []any{},
		},
	)
	return err
}
