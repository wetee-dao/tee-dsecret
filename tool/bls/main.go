package main

// import (
// 	"fmt"

// 	"go.dedis.ch/kyber/v4/pairing/bn256"
// 	"go.dedis.ch/kyber/v4/sign/bls"
// 	"go.dedis.ch/kyber/v4/util/random"
// )

// func main() {
// 	msg := []byte("Hello Boneh-Lynn-Shacham")
// 	suite := bn256.NewSuite()
// 	private1, public1 := bls.NewKeyPair(suite, random.New())
// 	private2, public2 := bls.NewKeyPair(suite, random.New())
// 	sig1, err := bls.Sign(suite, private1, msg)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("sig1:", sig1)
// 	sig2, err := bls.Sign(suite, private2, msg)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("sig2:", sig2)
// 	aggregatedSig, err := bls.AggregateSignatures(suite, sig1, sig2)
// 	fmt.Println("aggregatedSig:", aggregatedSig)
// 	if err != nil {
// 		panic(err)
// 	}

// 	aggregatedKey := bls.AggregatePublicKeys(suite, public1, public2)

// 	err = bls.Verify(suite, aggregatedKey, msg, aggregatedSig)
// 	if err != nil {
// 		panic(err)
// 	}
// }
