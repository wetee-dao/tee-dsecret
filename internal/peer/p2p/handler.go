package p2p

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	types "wetee.app/dsecret/type"
	"wetee.app/dsecret/util"
)

// genStream 函数创建了一个回调函数，用于处理通过网络接收到的消息
func genStream(handler func(*types.Message) error) func(network.Stream) {
	// 返回的函数将会作为接收网络流的处理逻辑
	return func(stream network.Stream) {
		// 读取流中的所有数据
		buf, err := io.ReadAll(stream)
		if err != nil {
			// 如果读取过程中发生错误，且错误不是因为流结束
			if err != io.EOF {
				fmt.Printf("read stream: %s", err)
			}
			// 尝试重置流，以便可以重新读取
			err = stream.Reset()
			if err != nil {
				fmt.Printf("reset stream: %s", err)
			}
			// 如果读取过程中出错，直接返回，不进行后续处理
			return
		}

		// 读取流的协议识别码
		protocolID := stream.Protocol()
		// 关闭流，释放资源
		err = stream.Close()
		if err != nil {
			fmt.Printf("close stream: %s", err)
			// 如果关闭流时出错，直接返回，不进行后续处理
			return
		}

		// 将读取到的字节数据反序列化为对象
		data := &types.Message{}
		err = json.Unmarshal(buf, data)
		if err != nil {
			// 打印解析 JSON 数据时的错误信息
			fmt.Printf("unmarshal data: %s", err)
			// 如果反序列化过程中出错，直接返回，不进行后续处理
			return
		}

		pids := protocol.ConvertToStrings([]protocol.ID{protocolID})
		// 记录接收消息的日志，包括来源和类型
		util.LogRevmsg("<<<<<< P2P  Rev()", "from ", stream.Conn().RemotePeer(), "| type:", pids[0]+"."+data.Type)
		// 将解析后的消息传递给外部提供的处理函数
		err = handler(data)
		if err != nil {
			// 打印处理消息时的错误信息
			fmt.Printf("handle data: %s \n", err)
			// 如果消息处理过程中出错，直接返回，不进行后续处理
			return
		}
	}
}
