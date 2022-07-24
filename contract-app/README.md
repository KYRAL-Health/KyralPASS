# How to use the API

## MAKE SURE GO, DOCKER-COMPOSE & DOCKER ARE INSTALLED

## Start up local network

Start in the GOPATH folder (please google this) if u dont have a GOPATH folder just make any folder

Install fabric-samples and binaries etc. in that folder

```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s
```

Then in the `fabric-samples/test-network` folder start the network by typing in 

```bash
./network.sh up createChannel -ca -s couchdb
```

let this load then you have the network set up 

## Setting up the APP

clone the repo in the GOPATH folder or watever folder that was made

```bash
##HTTPS
git clone "https://gitlab.com/z3fc/shf_continuity.git"

##SSH
git clone "git@gitlab.com:z3fc/shf_continuity.git"
```

go into the `contract` and `contract-app` folders and type `go get ./...`

## Changes to contract for running locally

The contract has been modified to work on kubernetes, which doesn't work nicely when trying to run locally for development.

You will need to comment out lines 44-70 (current line numbers. May change in the future) in main.go

The run function should go from:

```go
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
```

to:

```go
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

	// config := serverConfig{
	// 	CCID:    os.Getenv("CHAINCODE_ID"),
	// 	Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	// }

	// chaincode, err := contractapi.NewChaincode(&contract)

	// if err != nil {
	// 	log.Panicf("Failed to create chain code: %s", err)
	// }

	// // TODO: we need proper naming
	// chaincode.Info.Title = "Kyral-Contract"
	// chaincode.Info.Version = "1.0.0"

	// server := &shim.ChaincodeServer{
	// 	CCID:    config.CCID,
	// 	Address: config.Address,
	// 	CC:      chaincode,
	// 	TLSProps: shim.TLSProperties{
	// 		Disabled: true,
	// 	},
	// }

	// if err := server.Start(); err != nil {
	// 	return errors.Wrap(err, "Failed to start the chaincode")
	// }

	return nil
}
```

## Deploying Contract

Back in the `fabric-samples/test-network` folder

```bash
./network.sh deployCC -ccn "KYRAL" -ccp ../../shf_continuity/contract/ -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"
```

MAKE SURE THAT THE `-ccn` TAG IS `KYRAL` , ALSO DOUBLE CHECK THAT THE PATH TO THE CONTRACT IS CORRECT

## Starting the APP

First we need to set some environment variables to set parameters for the blockchain. This can either be done by exporting environment variables or using a .env file.

```
MSP_ID - MSP_ID of blockchain organization
CRYPTO_PATH - full path to user msp folder
TLS_CERT_PATH - full path to peer tls ca.crt
PEER_ENDPOINT - ip/dns name and port of peer
GATEWAY_PEER - fqdn of peer
CHANNEL_NAME - name of fabric channel
CHAINCODE_NAME - name of chaincode submitted to fabric channel
EVENT_UPDATE_API_KEY - api key to communicate events to kyral pass api
```

example environment variables
```
CHAINCODE_NAME=KYRAL
CHANNEL_NAME=mychannel
MSP_ID=Org1MSP
CRYPTO_PATH=/home/lindsay/go/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp
TLS_CERT_PATH=/home/lindsay/go/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
PEER_ENDPOINT=localhost:7051
GATEWAY_PEER=peer0.org1.example.com
EVENT_UPDATE_API_KEY="testAPIKey"
```

If you want to use a .env file instead, make these changes in application.go.  
at the top of the main function adds

```go
	godotenv.Load()

	mspID = os.Getenv("MSP_ID")
	certPath = filepath.Join(os.Getenv("CRYPTO_PATH"), "signcerts/cert.pem")
	keyPath = filepath.Join(os.Getenv("CRYPTO_PATH"), "keystore")
	tlsCertPath = os.Getenv("TLS_CERT_PATH")
	peerEndpoint = os.Getenv("PEER_ENDPOINT")
	gatewayPeer = os.Getenv("GATEWAY_PEER")
	channelName = os.Getenv("CHANNEL_NAME")
	chaincodeName = os.Getenv("CHAINCODE_NAME")
	eventUpdateApiKey = os.Getenv("EVENT_UPDATE_API_KEY")
```

 Go into the `shf_continuity/contract-app` folder then type

```bash
go run ./
```

  This will make a keystore and wallet folder and start the API

NOTE - when starting the api again make sure to delete the wallet folder

```bash
rm -R ./wallet
```

If you are prompted with `application.go: Started` message the api is running so open a new terminal to issue it some curl commands

make sure both terminals can be viewed

```bash

## adds data to the private and public blockchain for a specific user
curl -d '{"kyralTxID":"1329487","kyralPatientID":"1329487","operationCode":"76453","description":"test"}' "http://localhost:1234/createPrivate"

## reads the private data of a specified txID
curl "http://localhost:1234/readPrivateProperties/1329487"

## reads the public data of a specified txID
curl "http://localhost:1234/readPrivatePublicAsset/1329487"

## verifies that data sent from a host is the correct asset
curl -d '(insert data returned by `curl "http://localhost:1234/readPrivateProperties/1329487"`)' "http://localhost:1234/verifyAssetProperties"

##changes the description field for a specific KyralTxID
curl -X PUT -d '{"kyralTxID":"1329487","description":"newTest"}' "http://localhost:1234/changeDescription"



```

## Cleanup 

1. Close the API

2. Remove the wallet folder in the application folder

3. In `fabric-samples/test-network` 

   ```bash
   ./network.sh down
   ```

   
