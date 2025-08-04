package contracts

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/subnet"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func (c *Contract) GetMintWorker(user types.AccountID) (*model.K8sCluster, error) {
	workerWrap, _, err := c.subnet.QueryMintWorker(user, chain.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, err
	}

	if workerWrap.IsNone() {
		return nil, errors.New("worker not found")
	}

	tuple, _ := workerWrap.UnWrap()
	worker := tuple.F1
	return &model.K8sCluster{
		Id:            tuple.F0,
		Name:          worker.Name,
		Owner:         worker.Owner,
		Level:         worker.Level,
		RegionId:      worker.RegionId,
		StartBlock:    worker.StartBlock,
		StopBlock:     worker.StopBlock,
		TerminalBlock: worker.TerminalBlock,
		P2pId:         worker.P2pId,
		Ip:            model.Ip(worker.Ip),
		Port:          worker.Port,
		Status:        worker.Status,
	}, nil
}

func (c *Contract) GetWorker(workerId uint64) (*model.K8sCluster, error) {
	workerWrap, _, err := c.subnet.QueryWorker(workerId, chain.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, err
	}

	if workerWrap.IsNone() {
		return nil, errors.New("worker not found")
	}

	worker, _ := workerWrap.UnWrap()

	return &model.K8sCluster{
		Id:            workerId,
		Name:          worker.Name,
		Owner:         worker.Owner,
		Level:         worker.Level,
		RegionId:      worker.RegionId,
		StartBlock:    worker.StartBlock,
		StopBlock:     worker.StopBlock,
		TerminalBlock: worker.TerminalBlock,
		P2pId:         worker.P2pId,
		Ip:            model.Ip(worker.Ip),
		Port:          worker.Port,
		Status:        worker.Status,
	}, nil
}

func (c *Contract) ResigerCluster(name []byte, p2p_id [32]byte, ip model.Ip, port uint32, level byte, region_id uint32) error {
	return c.subnet.ExecWorkerRegister(name, p2p_id, subnet.Ip(ip), port, level, region_id, chain.ExecParams{
		Signer:    c.signer,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}

func (c *Contract) GetPodsVersionByWorker(workerId uint64) ([]model.PodVersion, error) {
	data, _, err := c.cloud.QueryWorkerPodsVersion(workerId, chain.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, err
	}

	list := make([]model.PodVersion, 0, len(*data))
	for _, v := range *data {
		list = append(list, model.PodVersion{
			PodId:    v.F0,
			Version:  v.F1,
			LastMint: v.F2,
			Status:   v.F3,
		})
	}

	return list, nil
}

func (c *Contract) GetPodsByIds(podIds []uint64) ([]model.Pod, error) {
	data, _, err := c.cloud.QueryPodsByIds(podIds, chain.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, err
	}

	pods := make([]model.Pod, 0, len(*data))
	for _, v := range *data {
		contariners := make([]model.Container, 0, len(v.F2))
		for i := range v.F2 {
			bt, _ := json.Marshal(v.F2[i].F1)
			contariner := model.Container{}
			json.Unmarshal(bt, &contariner)
			contariners = append(contariners, contariner)
		}
		pods = append(pods, model.Pod{
			PodId:               v.F0,
			Owner:               v.F1.Owner,
			Ptype:               model.PodType(v.F1.Ptype),
			TeeType:             ConvertTEEType(v.F1.TeeType),
			Containers:          contariners,
			Version:             v.F3,
			LastMintBlockNumber: v.F4,
			Status:              v.F5,
		})
	}

	return pods, nil
}

func (c *Contract) TxCallOfStartPod(nodeId uint64, pod_key types.AccountID, signer types.AccountID) (*types.Call, error) {
	return c.cloud.CallOfStartPod(nodeId, pod_key, chain.DryRunParams{
		Origin:    signer,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}
func (c *Contract) DryStartPod(nodeId uint64, pod_key types.AccountID, signer types.AccountID) error {
	_, _, err := c.cloud.DryRunStartPod(nodeId, pod_key, chain.DefaultParamWithOrigin(signer))
	return err
}

func (c *Contract) TxCallOfMintPod(nodeId uint64, hash types.H256, signer types.AccountID) (*types.Call, error) {
	return c.cloud.CallOfMintPod(nodeId, hash, chain.DryRunParams{
		Origin:    signer,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}
func (c *Contract) DryMintPod(nodeId uint64, hash types.H256, signer types.AccountID) error {
	_, _, err := c.cloud.DryRunMintPod(nodeId, hash, chain.DefaultParamWithOrigin(signer))
	return err
}

func (c *Contract) TxCallOfUploadSecret(user types.H160, index uint64, data types.H256, signer types.AccountID) (*types.Call, error) {
	return c.cloud.CallOfUpdateSecret(user, index, data, chain.DryRunParams{
		Origin:    signer,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}
func (c *Contract) DryUploadSecret(user types.H160, index uint64, data types.H256, signer types.AccountID) error {
	_, _, err := c.cloud.DryRunUpdateSecret(user, index, data, chain.DefaultParamWithOrigin(signer))
	return err
}
