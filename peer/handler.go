package peer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p/core/network"
	types "wetee.app/dsecret/type"
	"wetee.app/dsecret/util"
)

func genStream(handler func(*types.Message) error) func(network.Stream) {
	return func(stream network.Stream) {
		buf, err := io.ReadAll(stream)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("read stream: %s", err)
			}

			err = stream.Reset()
			if err != nil {
				fmt.Printf("reset stream: %s", err)
			}

			return
		}

		protocolID := stream.Protocol()

		err = stream.Close()
		if err != nil {
			fmt.Printf("close stream: %s", err)
			return
		}

		data := &types.Message{}
		err = json.Unmarshal(buf, data)
		if err != nil {
			fmt.Printf("unmarshal data: %s", err)
			return
		}

		util.LogRevmsg("<<<<<< P2P  Rev()", "from ", stream.Conn().RemotePeer(), "| type:", data.Type+", ProtocolID =", protocolID)
		err = handler(data)
		if err != nil {
			fmt.Printf("handle data: %s \n", err)
			return
		}
	}
}
