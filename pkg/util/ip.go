package util

import (
	"fmt"
	"math/big"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func GetUrlFromIp(ip model.Ip) string {
	url := ""
	if !ip.Domain.IsNone {
		url = "/dns4/" + string(ip.Domain.V)
	} else if !ip.Ipv4.IsNone {
		ipv4 := ip.Ipv4.V
		url = "/ip4/" + fmt.Sprintf("%d.%d.%d.%d",
			(ipv4>>24)&0xFF,
			(ipv4>>16)&0xFF,
			(ipv4>>8)&0xFF,
			ipv4&0xFF)
	} else if !ip.Ipv6.IsNone {
		ipv6 := ip.Ipv6.V
		ipv6Int128 := big.NewInt(0)
		ipv6Int128.SetBytes(ipv6.Bytes())
		url = "/ip6/" + fmt.Sprintf("%04x:%04x:%04x:%04x:%04x:%04x:%04x:%04x",
			ipv6Int128.Rsh(ipv6Int128, 112).Uint64(),
			ipv6Int128.Rsh(ipv6Int128, 96).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 80).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 64).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 48).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 32).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 16).Uint64()&0xFFFF,
			ipv6Int128.Uint64()&0xFFFF)
	}
	return url
}
