package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

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
	chaincode, err := contractapi.NewChaincode(&contract)

	if err != nil {
		return errors.Wrap(err, "Failed to create chain code")
	}

	// TODO: we need proper naming
	chaincode.Info.Title = "Kyral-Contract"
	chaincode.Info.Version = "1.0.0"

	if err := chaincode.Start(); err != nil {
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
