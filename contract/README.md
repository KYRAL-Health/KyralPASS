# Smart Contract

The hyperledger smart contract for Kyral.

# Build

You need these dependencies for this project.
```bash
sudo apt install make golang
```

# Smart Contract CRUD Functions

## <ins>CreateOrgVisit</ins>

Private function that adds private data (specified in the private input data) to a private collection that only a specific organization can see unless given permission to view. As well as adding data to a public collection specificed in input data.

|Input Data| Data Type|
|-------|-------|
| TxID | String |
| Description | String |
| OwnerOrg  | String |


|Transient Data| Data Type|
|-------|-------|
| TxID | String |
| Salt  | String |
| OperationCode | String |

## <ins>ChangeDescription</ins>

Updates the public description for an asset checking the users identity before confirmation

| Input Data | Data Type |
| -------| -------|
| txID | String |
| newDescription | String |


# Smart Contract Query Functions

## <ins>GetAllUserTx</ins>

Obtains all public user transactions associated with the UserData field

## <ins>GetAllOrgVisit</ins>

Obtain all public data associted with the OrgVisit Field

## <ins>ReadTransactionID</ins>

Obtain data of a specific transaction associted with primary key txid

|Input Data| Data Type|
|-------|-------|
| txid | String |


## <ins>CheckAsset</ins>

Check if an asset exists or not returns true or false

|Input Data| Data Type|
|-------|-------|
| txid | String |

## <ins>ReadAsset</ins>

Read public data from collection associated with OrgVisit wiht primary key txid

|Input Data| Data Type|
|-------|-------|
| txid | String |

## <ins>ReadAssetPrivateDetails</ins>

Read Private data from private collection of a specific org with primary key txid

|Input Data| Data Type|
|-------|-------|
| txid | String |

## <ins>VerifyAssetProperties</ins>

Verifies that the hash stored in the private collection is correct to validate any data stored in the private collection. 

|Input Data| Data Type|
|-------|-------|
| txid | String |

|Transient Data| Data Type|
|-------|-------|
| txid | String |
| patientID | String |
| OrgID | String |
| OperationCode | String |
| salt | String |