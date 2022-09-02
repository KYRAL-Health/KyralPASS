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
git clone "https://github.com/KYRAL-Health/KyralPASS.git"

##SSH
git clone "git@github.com:KYRAL-Health/KyralPASS.git"
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

## Plain Text User Struct
Use this struct for the user object. Encrypted User should be this object converted to json string and then encrypted using the function in crypt-utils.go
```go
type UserBlockchain struct {
	KyralID   uuid.UUID `json:"kyralID" gorm:"type:uuid;default:uuid_generate_v4();unique"`
	UpdatedAt time.Time
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Password  string         `json:"password"`
	Email     string         `json:"email" gorm:"unique"`
	Address   string         `json:"address"`
	Telecom   string         `json:"telecom"`
	Gender    string         `json:"gender"`
	BirthDate datatypes.Date `json:"birthDate"`
}

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
API_URL - URL to api to send event based updates
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
API_URL="http://localhost:8080"
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

## changes the description field for a specific KyralTxID
curl -X PUT -d '{"kyralTxID":"1329487","description":"newTest"}' "http://localhost:1234/changeDescription"

## Create user on blockchain - See crypto-utils.go for encryption and decryption format
curl -L -X POST 'localhost:1234/createUser' -H 'Content-Type: application/json' --data-raw '{
    "kyralUID": "f9092446-14fe-4c33-9028-7acd119e332e",
    "kyralEncryptedUser": "m1ysKKLrmWYpfIlFb6hBU6k64rXCRjxPS2mG33_eNZBDky-4vr-kLtTGeJdH3RNoNEZHMSt_2NHtVYBkAMbSo68vgiDv0xf8LAFOzLu6Ds-wWnvI8UVyVjeugootIjk4ls0soSZMz0Y9zbvXBpceGbGI5hlHgF8v78ux3IEbJysSZ1z6GOjK_tzsHJxoLH1NQJOFkyXsWSt8BzFz3k0wQLDbPwgksyBOomfApWQtVfGKaDqRI5thLwk23cJO0NNhiDrleEooWzbXrXkbIagH_21LZceaRfnu8bxUtBIf9L07BscFocyJQFCuPhDVAL4XMiz3VwDevrRTahC6tls_No0whelkKVlMNHO7uylJNejmezg1k62YK7LyhfhhGEMlh90ooO5Dktdwbfc87Q3TNd7JI8lmQqdaqfrWJSk8kEM=",
    "kyralEncryptedUserHash": "hWIZWbtFVEWyFbryfrcggP6_L-b9Wefj-YtLtnNsRu_EIvCZPB7mT87_28vLIaMiLtLCaefAuZZnQPp8tiDUsA==",
    "decryptKey": "hma7Uo4p2KYKhsplnQxoobdlRlzfegWM"
}'

## Read user from blockchain
curl -L -X POST 'localhost:1234/readUser' -H 'Content-Type: application/json' --data-raw '{
    "kyralUID":"f9092446-14fe-4c33-9028-7acd119e332e",
    "decryptKey": "hma7Uo4p2KYKhsplnQxoobdlRlzfegWM"
}'

## Update user on blockchain
curl -L -X PUT 'localhost:1234/updateUser' -H 'Content-Type: application/json' --data-raw '{
    "kyralUID": "f9092446-14fe-4c33-9028-7acd119e332e",
    "kyralEncryptedUser": "qoc7fHFWdslxix_fm6k6EpU3MbImJ_EHuaKwhG2wOcVEekKdgj6MRAuhznTLPriA_4L-BraJvoqqJFTz14fJ17c4iuEE8lfDjRbOnlPhc5dq1Z5heRqBtf_W5trX9g8LT8-gsLYO6CdNLMTmEAT6_5upINxCeoa3k9bRhYghM46RLX2-AGrmLBVrCRN92EV0Xn0BN1yAP5xvzqcCjO6NGEz3a4_YJXKMsVYSqkHq3-2xUlc97wwJRcY_wxnmyDWpMfUgKTRXdcneXAX8fuOb_xzH5YyueHOo2_jGCpDO6eP0_kDB8JcYsqgPGtOP8ULstsVLaHwkKlM478O67_0a",
    "kyralEncryptedUserHash": "n8IzST0d_FzSZwyQZe6Uh8Y7wbcDTt_LvZRQmEFACCf2Ws-BTbggXDIVc59Ta7UvV233yHDj2QPRYpxAFxby9Q==",
    "decryptKey": "hma7Uo4p2KYKhsplnQxoobdlRlzfegWM"
}'

## Transfer User to different organization
curl -L -X PUT 'localhost:1234/transferUser' -H 'Content-Type: application/json' --data-raw '{
    "kyralUID": "f9092446-14fe-4c33-9028-7acd119e332e",
    "decryptKey": "hma7Uo4p2KYKhsplnQxoobdlRlzfegWM",
    "orgID": "Org2MSP"
}'
```

## Cleanup 

1. Close the API

2. Remove the wallet folder in the application folder

3. In `fabric-samples/test-network` 

   ```bash
   ./network.sh down
   ```

   
