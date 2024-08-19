package main

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"os"

	"wetee.app/dsecret/chain"
	"wetee.app/dsecret/dkg"
	"wetee.app/dsecret/graph"
	p2p "wetee.app/dsecret/peer"
	"wetee.app/dsecret/store"
	types "wetee.app/dsecret/type"
	"wetee.app/dsecret/util"
)

var DefaultChainUrl string = "ws://wetee-node.worker-addon.svc.cluster.local:9944"

func main() {
	// 获取环境变量
	peerSecret := util.GetEnv("PEER_PK", "")
	tcpPort := util.GetEnvInt("TCP_PORT", 61000)
	udpPort := util.GetEnvInt("UDP_PORT", 61000)
	chainAddr := util.GetEnv("CHAIN_ADDR", DefaultChainUrl)
	password := util.GetEnv("PASSWORD", "")

	// 初始化数据库
	err := store.InitDB(password)
	if err != nil {
		fmt.Println("Init db error:", err)
		os.Exit(1)
	}

	// 初始化加密套件。
	nodeSecret, err := types.PrivateKeyFromLibp2pHex(peerSecret)
	if err != nil {
		fmt.Println("Marshal PKG_PK error:", err)
		os.Exit(1)
	}

	// 链接区块链
	err = chain.InitChain(chainAddr, nodeSecret)
	if err != nil {
		fmt.Println("Connect to chain error:", err)
		os.Exit(1)
	}

	// Get boot peers from chain
	bootPeers, err := chain.ChainIns.GetBootPeers()
	if err != nil {
		fmt.Println("Get node list error:", err)
		os.Exit(1)
	}
	if len(bootPeers) == 0 {
		fmt.Println("No boot peers found")
		os.Exit(1)
	}

	secretNodes, _, nodes, err := chain.ChainIns.GetNodes()
	if err != nil {
		fmt.Println("Get node list error:", err)
		os.Exit(1)
	}

	// 获取阈值参数
	threshold := len(secretNodes) * 2 / 3

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	boots := make([]string, 0, len(bootPeers))
	for _, b := range bootPeers {
		var gopub ed25519.PublicKey = b.Id[:]
		pub, _ := types.PubKeyFromStdPubKey(gopub)
		n := &types.Node{
			ID: pub.String(),
		}
		d := util.GetUrlFromIp1(b.Ip)
		url := d + "/tcp/" + fmt.Sprint(b.Port) + "/p2p/" + n.PeerID().String()
		boots = append(boots, url)
	}

	// 启动 P2P 网络
	peer, err := p2p.NewP2PNetwork(ctx, nodeSecret, boots, nodes, uint32(tcpPort), uint32(udpPort))
	if err != nil {
		fmt.Println("Start P2P peer error:", err)
		os.Exit(1)
	}

	// 创建 DKG 实例。
	dkg, err := dkg.NewRabinDKG(nodeSecret, nodes, threshold, peer)
	if err != nil {
		fmt.Println("Create DKG error:", err)
		os.Exit(1)
	}

	// 启动节点
	go peer.Start(ctx)

	// 运行 DKG 协议。
	if err := dkg.Start(ctx); err != nil {
		fmt.Println("Start DKG error:", err)
		os.Exit(1)
	}

	graph.StartServer(dkg)
}

// import (
// 	"fmt"

// 	"github.com/consensys/gnark/frontend"

// 	"go.dedis.ch/kyber/v3"
// 	"go.dedis.ch/kyber/v3/share"
// 	dkg "go.dedis.ch/kyber/v3/share/dkg/pedersen"
// 	"go.dedis.ch/kyber/v3/suites"
// )

// // 门限方案参数
// type ThresholdScheme struct {
// 	Threshold int          // 恢复密钥所需的最小碎片数量
// 	Total     int          // 碎片总数
// 	Suite     suites.Suite // 密码学套件
// }

// // 生成密钥碎片
// func GenerateSecretShares(scheme *ThresholdScheme, secret kyber.Scalar) []*share.PriShare {
// 	pri := share.NewPriPoly(scheme.Suite, scheme.Threshold, secret, scheme.Suite.RandomStream())

// 	return pri.Shares(scheme.Total)
// }

// // 恢复密钥
// func RecoverSecret(scheme *ThresholdScheme, shares []*share.PriShare) (kyber.Scalar, error) {
// 	return share.RecoverSecret(scheme.Suite, shares, scheme.Threshold, scheme.Total)
// }

// // Gnark 算术电路
// type ThresholdCircuit struct {
// 	SecretShare share.PriShare
// 	Secret      kyber.Scalar
// }

// func (circuit *ThresholdCircuit) Define(curve frontend.API) error {
// 	// 将密钥碎片和密钥恢复过程定义为算术电路
// 	// ...

// 	return nil
// }

// type DkgNode struct {
// 	dkg         *dkg.DistKeyGenerator
// 	pubKey      kyber.Point
// 	privKey     kyber.Scalar
// 	deals       []*dkg.Deal
// 	resps       []*dkg.Response
// 	secretShare *share.PriShare
// }

// func main() {
// 	// 初始化密码学套件
// 	suite := suites.MustFind("Ed25519")

// 	// 设置门限方案参数
// 	// scheme := &ThresholdScheme{
// 	// 	Threshold: 3, // 至少需要 3 个碎片才能恢复密钥
// 	// 	Total:     5, // 总共 5 个碎片
// 	// 	Suite:     suite,
// 	// }

// 	// 生成密钥
// 	privKey := suite.Scalar().Pick(suite.RandomStream())
// 	// 生成公钥
// 	pubKey := suite.Point().Mul(privKey, nil)
// 	fmt.Printf("原始密钥: %x\n", privKey.String())

// 	// 创建DKG结点
// 	node := &DkgNode{
// 		pubKey:  pubKey,
// 		privKey: privKey,
// 		deals:   make([]*dkg.Deal, 0),
// 		resps:   make([]*dkg.Response, 0),
// 	}
// 	d, err := dkg.NewDistKeyGenerator(suite, privKey, []kyber.Point{pubKey}, 3)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	node.dkg = d

// 	// // 生成密钥碎片
// 	// shares := GenerateSecretShares(scheme, privKey)

// 	// // 打印密钥碎片
// 	// fmt.Println("密钥碎片:")
// 	// for _, share := range shares {
// 	// 	fmt.Printf("碎片 %d: %x\n", share.I, share.V.String())
// 	// }

// 	// // 要加密的数据
// 	// data := []byte("这是一个秘密信息")

// 	// // 加密数据
// 	// ciphertext, err := ecies.Encrypt(suite, pubKey, data, suite.Hash)
// 	// if err != nil {
// 	// 	fmt.Println("加密失败:", err)
// 	// 	return
// 	// }

// 	// // 解密数据
// 	// decryptedData, err := ecies.Decrypt(suite, privKey, ciphertext, suite.Hash)
// 	// if err != nil {
// 	// 	fmt.Println("解密失败:", err)
// 	// 	return
// 	// }

// 	// if string(decryptedData) != string(data) {
// 	// 	fmt.Println("解密数据与原始数据不一致")
// 	// 	return
// 	// }

// 	// // 生成 zk-SNARK 证明
// 	// // ...

// 	// // 恢复密钥
// 	// recoveredSecret, err := RecoverSecret(scheme, shares[0:3])
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	return
// 	// }
// 	// fmt.Printf("恢复的密钥: %x\n", recoveredSecret.String())

// 	// // 验证密钥是否恢复成功
// 	// if recoveredSecret.Equal(privKey) {
// 	// 	fmt.Println("密钥恢复成功")
// 	// } else {
// 	// 	fmt.Println("密钥恢复失败")
// 	// }

// 	// 在 Substrate 中验证证明
// }
