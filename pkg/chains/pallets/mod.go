package pallets

import (
	"errors"
	"fmt"

	chain "github.com/wetee-dao/ink.go"

	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/app"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/dsecret"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/gpu"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/task"
	gtypes "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/worker"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// Chain
type Chain struct {
	*chain.ChainClient
	signer *chain.Signer
}

func NewContract(url string, pk *model.PrivKey) (*Chain, error) {
	client, err := chain.ClientInit(url, true)
	if err != nil {
		return nil, err
	}

	p, err := pk.ToSigner()
	if err != nil {
		return nil, err
	}

	return &Chain{
		ChainClient: client,
		signer:      p,
	}, nil
}

func (c *Chain) GetClient() *chain.ChainClient {
	return c.ChainClient
}

func (c *Chain) GetSignerAddress() string {
	return c.signer.SS58Address(42)
}

func (c *Chain) GetAccount(workID gtypes.WorkId) ([]byte, error) {
	if workID.Wtype.IsAPP {
		return c.GetAppAccount(workID.Id)
	} else if workID.Wtype.IsTASK {
		return c.GetTaskAccount(workID.Id)
	} else if workID.Wtype.IsGPU {
		return c.GetGpuAccount(workID.Id)
	}
	return nil, errors.New("unknow work type")
}

func (c *Chain) GetVersion(client *chain.ChainClient, workID gtypes.WorkId) (ret uint64, err error) {
	if workID.Wtype.IsAPP {
		return c.GetAppVersionLatest(workID.Id)
	} else if workID.Wtype.IsTASK {
		return c.GetTaskVersionLatest(workID.Id)
	} else if workID.Wtype.IsGPU {
		return c.GetGpuVersionLatest(workID.Id)
	}

	return 0, errors.New("unknow work type")
}

func (c *Chain) GetSecretEnv(client *chain.ChainClient, workID gtypes.WorkId) (ret []byte, isSome bool, err error) {
	if workID.Wtype.IsAPP {
		return c.GetAppSecretEnv(workID.Id)
	} else if workID.Wtype.IsTASK {
		return c.GetTaskSecretEnv(workID.Id)
	} else if workID.Wtype.IsGPU {
		return c.GetGpuSecretEnv(workID.Id)
	}

	return nil, false, errors.New("unknow work type")
}

// GetWorkCodeSignature 函数根据工作 ID 获取相应的代码签名
func (c *Chain) GetWorkCodeSignature(client *chain.ChainClient, workID gtypes.WorkId) (ret []byte, err error) {
	// 判断工作类型是否为 APP
	if workID.Wtype.IsAPP {
		// 调用 weteeapp 获取代码签名的最新数据
		return app.GetCodeSignatureLatest(client.Api.RPC.State, workID.Id)
	} else if workID.Wtype.IsTASK {
		// 调用 weteetask 获取代码签名的最新数据
		return task.GetCodeSignatureLatest(client.Api.RPC.State, workID.Id)
	} else if workID.Wtype.IsGPU {
		// 调用 weteegpu 获取代码签名的最新数据
		return gpu.GetCodeSignatureLatest(client.Api.RPC.State, workID.Id)
	}

	// 如果工作类型未知，返回错误信息
	return nil, errors.New("unknow work type")
}

// 获取全网当前程序的代码版本
// Get CodeSignature/SIgner
func (c *Chain) GetDsecretCode(client *chain.ChainClient) ([]byte, []byte, error) {
	// 检查节点代码是否和 wetee 上要求的版本一致
	codeSignature, err := dsecret.GetCodeSignatureLatest(client.Api.RPC.State)
	if err != nil {
		fmt.Println("Get code signature error:", err)
		return nil, nil, err
	}
	codeSigner, err := dsecret.GetCodeSignerLatest(client.Api.RPC.State)
	if err != nil {
		fmt.Println("Get code signer error:", err)
		return nil, nil, err
	}

	return codeSignature, codeSigner, nil
}

// GetWorkerCode 函数用于获取 weteeworker 的代码签名和签名者
func (c *Chain) GetWorkerCode(client *chain.ChainClient) ([]byte, []byte, error) {
	// 获取 weteeworker 的最新代码签名
	codeSignature, err := worker.GetCodeSignatureLatest(client.Api.RPC.State)
	// 处理获取代码签名过程中的错误
	if err != nil {
		fmt.Println("Get code signature error:", err)
		return nil, nil, err
	}

	// 获取 weteeworker 的最新代码签名者
	codeSigner, err := worker.GetCodeSignerLatest(client.Api.RPC.State)
	// 处理获取代码签名者过程中的错误
	if err != nil {
		fmt.Println("Get code signer error:", err)
		return nil, nil, err
	}

	// 返回获取到的代码签名和签名者
	return codeSignature, codeSigner, nil
}

// GetWorkCode 函数用于获取工作代码签名和代码签名者
func (c *Chain) GetWorkCode(workID gtypes.WorkId) ([]byte, []byte, error) {
	// 获取工作签名
	sig, err := c.GetWorkSignature(workID)
	if err != nil {
		return nil, nil, err
	}

	// 获取工作代码签名者
	signer, err := c.GetWorkCodeSigner(workID)
	if err != nil {
		return nil, nil, err
	}

	// 返回获取到的签名和签名者，以及一个 nil 错误
	return sig, signer, nil
}

// GetWorkSignature 函数根据工作 ID 获取相应的代码签名
func (c *Chain) GetWorkSignature(workID gtypes.WorkId) ([]byte, error) {
	// 判断工作类型是否为 APP
	if workID.Wtype.IsAPP {
		// 调用 weteeapp 获取代码签名的最新数据
		return app.GetCodeSignatureLatest(c.Api.RPC.State, workID.Id)
	} else if workID.Wtype.IsTASK {
		// 调用 weteetask 获取代码签名的最新数据
		return task.GetCodeSignatureLatest(c.Api.RPC.State, workID.Id)
	} else if workID.Wtype.IsGPU {
		// 调用 weteegpu 获取代码签名的最新数据
		return gpu.GetCodeSignatureLatest(c.Api.RPC.State, workID.Id)
	}

	// 如果工作类型未知，返回错误信息
	return nil, errors.New("unknow work type")
}

// GetWorkCodeSigner 函数根据工作 ID 获取相应的代码签名者
func (c *Chain) GetWorkCodeSigner(workID gtypes.WorkId) ([]byte, error) {
	// 判断工作类型是否为 APP
	if workID.Wtype.IsAPP {
		// 调用 weteeapp 获取代码签名者的最新数据
		return app.GetCodeSignerLatest(c.Api.RPC.State, workID.Id)
	} else if workID.Wtype.IsTASK {
		// 调用 weteetask 获取代码签名者的最新数据
		return task.GetCodeSignerLatest(c.Api.RPC.State, workID.Id)
	} else if workID.Wtype.IsGPU {
		// 调用 weteegpu 获取代码签名者的最新数据
		return gpu.GetCodeSignerLatest(c.Api.RPC.State, workID.Id)
	}

	// 如果工作类型未知，返回错误信息
	return nil, errors.New("unknow work type")
}
