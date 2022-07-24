package main

import (
	"log"
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type serverConfig struct {
	CCID    string
	Address string
}

func run() error {
	contract := Contract{
		contractapi.Contract{

			Name: "Kyral-Contract",
			Info: metadata.InfoMetadata{
				Title: "Kyral-ChainCode",
			},
			TransactionContextHandler: &TransactionContext{},
		},
	}
	// chaincode, err := contractapi.NewChaincode(&contract)

	// if err != nil {
	// 	return errors.Wrap(err, "Failed to create chain code")
	// }

	// // TODO: we need proper naming
	// chaincode.Info.Title = "Kyral-Contract"
	// chaincode.Info.Version = "1.0.0"

	// if err := chaincode.Start(); err != nil {
	// 	return errors.Wrap(err, "Failed to start the chaincode")
	// }

	config := serverConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	}

	chaincode, err := contractapi.NewChaincode(&contract)

	if err != nil {
		log.Panicf("Failed to create chain code: %s", err)
	}

	// TODO: we need proper naming
	chaincode.Info.Title = "Kyral-Contract"
	chaincode.Info.Version = "1.0.0"

	server := &shim.ChaincodeServer{
		CCID:    config.CCID,
		Address: config.Address,
		CC:      chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	if err := server.Start(); err != nil {
		return errors.Wrap(err, "Failed to start the chaincode")
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		logrus.Fatalf("Failed to run contract chaincode: %v", err)
	}

	return
}
