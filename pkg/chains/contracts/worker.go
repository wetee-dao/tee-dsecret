package contracts

import (
	"errors"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/subnet"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func (c *Contract) GetWorkerId(user types.AccountID) (uint64, error) {
	id, _, err := c.subnet.QueryMintWorker(user, chain.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return 0, err
	}

	if id.IsNone() {
		return 0, errors.New("worker not found")
	}

	return id.UnWrap()
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

func (c *Contract) GetPodsVersionByWorker(workerId uint64) ([]model.PodVersion, error) {
	data, _, err := c.cloud.QueryWorkerPodsVersion(workerId, chain.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, err
	}

	list := make([]model.PodVersion, 0, len(*data))
	for _, v := range *data {
		list = append(list, model.PodVersion{
			PodId:   v.F0,
			Version: v.F1,
			Status:  v.F2,
		})
	}

	return list, nil
}

func (c *Contract) ResigerCluster(name []byte, p2p_id [32]byte, ip model.Ip, port uint32, level byte, region_id uint32) error {
	return c.subnet.CallWorkerRegister(name, p2p_id, subnet.Ip(ip), port, level, region_id, chain.CallParams{
		Signer:    c.signer,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}
