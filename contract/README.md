# Smart Contract

The hyperledger smart contract for Kyral. Built to run as a container on kubernetes

# Build

You need these dependencies for this project.

```bash
sudo apt install make golang
```

# Running Locally

This contract is built to be used running in a container on a kubernetes based network. To use it with the test network a few changes need to be made.

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

# Smart Contract CRUD Functions

## <ins>CreateOrgVisit</ins>

Private function that adds private data (specified in the private input data) to a private collection that only a specific organization can see unless given permission to view. As well as adding data to a public collection specificed in input data.

| Input Data  | Data Type |
| ----------- | --------- |
| TxID        | String    |
| Description | String    |
| OwnerOrg    | String    |

| Transient Data | Data Type |
| -------------- | --------- |
| TxID           | String    |
| Salt           | String    |
| OperationCode  | String    |

## <ins>ChangeDescription</ins>

Updates the public description for an asset checking the users identity before confirmation

| Input Data     | Data Type |
| -------------- | --------- |
| txID           | String    |
| newDescription | String    |

## <ins>CreateNewUser</ins>

Creates a new user on chain. This user object needs to be encrypted with hashes to verify correct encryption/decryption. It is the clients job to ensure that data is encrypted properly and that decrypt key is able to decrypt ciphertext and hash matches. Encryption needs to match the functions found in crypto-utils.go

| Input Data         | Data Type |
| ------------------ | --------- |
| kyralUID           | String    |
| kyralEncryptedUser | String    |
| hash\_             | String    |

## <ins>UpdateUserData</ins>

Update user object on chain provided that the object belongs to the same organization. It is the clients job to ensure that data is encrypted properly and that decrypt key is able to decrypt ciphertext and hash matches.  
Encryption needs to match the functions found in crypto-utils.go

| Input Data         | Data Type |
| ------------------ | --------- |
| kyralUID           | String    |
| kyralEncryptedUser | String    |
| hash\_             | String    |

## <ins>TransferUser</ins>

Transfer on chain user object to a different organization. orgID field is an organizations MSPID. This changes the endorsement policy on a key level. If the MSPID does not exist, unfortunately the record is locked as no organization will match the organization policy.  
This can be modified to include another trusted organization to be a part of the endorsement policy so this situation does not occur.

| Input Data | Data Type |
| ---------- | --------- |
| kyralUID   | String    |
| orgID      | String    |

# Smart Contract Query Functions

## <ins>GetAllUserTx</ins>

Obtains all public user transactions associated with the UserData field

## <ins>GetAllOrgVisit</ins>

Obtain all public data associted with the OrgVisit Field

## <ins>ReadTransactionID</ins>

Obtain data of a specific transaction associted with primary key txid

| Input Data | Data Type |
| ---------- | --------- |
| txid       | String    |

## <ins>CheckAsset</ins>

Check if an asset exists or not returns true or false

| Input Data | Data Type |
| ---------- | --------- |
| txid       | String    |

## <ins>ReadAsset</ins>

Read public data from collection associated with OrgVisit wiht primary key txid

| Input Data | Data Type |
| ---------- | --------- |
| txid       | String    |

## <ins>ReadAssetPrivateDetails</ins>

Read Private data from private collection of a specific org with primary key txid

| Input Data | Data Type |
| ---------- | --------- |
| txid       | String    |

## <ins>VerifyAssetProperties</ins>

Verifies that the hash stored in the private collection is correct to validate any data stored in the private collection.

| Input Data | Data Type |
| ---------- | --------- |
| txid       | String    |

| Transient Data | Data Type |
| -------------- | --------- |
| txid           | String    |
| patientID      | String    |
| OrgID          | String    |
| OperationCode  | String    |
| salt           | String    |

## <ins>ReadUser</ins>

Retrieve encrypted user object from the blockchain. Client responsibility to check if decryption key is able to decypt ciphertext and plaintext hash matches on chain hash field.

| Input Data | Data Type |
| ---------- | --------- |
| kyralUID   | String    |

## <ins>CheckUser</ins>

Used to check if a KyralID record exists

| Input Data | Data Type |
| ---------- | --------- |
| kyralUID   | String    |
