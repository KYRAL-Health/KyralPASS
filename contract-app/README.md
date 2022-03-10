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
git clone "https://gitlab.com/KYRAL-Health/KyralTEST.git"

##SSH
git clone "git@github.com:KYRAL-Health/KyralTEST.git"
```

## Deploying Contract

Back in the `fabric-samples/test-network` folder

```bash
./network.sh deployCC -ccn "KYRAL" -ccp ../../KyralTEST/contract/ -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"
```

MAKE SURE THAT THE `-ccn` TAG IS `KYRAL` , ALSO DOUBLE CHECK THAT THE PATH TO THE CONTRACT IS CORRECT

## Starting the APP

 Go into the `KyralTEST/contract-app` folder then type

```bash
go run ./
```

  This will make a keystore and wallet folder and start the API

NOTE - when starting the api again make sure to delete the wallet folder

```bash
rm -R ./wallet
```

If you are prompted with `les a go` message the api is running so open a new terminal to issue it some curl commands

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

   
