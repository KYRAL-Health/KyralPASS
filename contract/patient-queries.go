package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//gets all OrgVisit Data in the world state
func (c *Contract) GetAllOrgVisit(ctx contractapi.TransactionContextInterface) ([]*OrgVisit, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var patients []*OrgVisit
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var patient OrgVisit
		err = json.Unmarshal(queryResponse.Value, &patient)
		if err != nil {
			return nil, err
		}
		patients = append(patients, &patient)
	}

	return patients, nil
}

func (c *Contract) CheckAsset(ctx contractapi.TransactionContextInterface, assetID string) (bool, error) {

	assetJSON, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return false, err
	}
	if assetJSON == nil {
		return false, fmt.Errorf("%s does not exist", assetID)
	}

	return true, nil
}

// ReadAsset returns the public asset data
func (c *Contract) ReadAsset(ctx contractapi.TransactionContextInterface, assetID string) (*OrgVisit, error) {

	assetJSON, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("%s does not exist", assetID)
	}

	var asset *OrgVisit
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

// GetPrivateProperties returns the immutable asset properties from owner's private data collection
func (c *Contract) ReadPrivateProperties(ctx contractapi.TransactionContextInterface, txID string) (string, error) {
	// In this scenario, client is only authorized to read/write private data from its own peer.
	collection, err := GetClientImplicitCollectionName(ctx)
	if err != nil {
		return "", err
	}

	immutableProperties, err := ctx.GetStub().GetPrivateData(collection, txID)
	if err != nil {
		return "", fmt.Errorf("failed to read private properties from client org's collection: %v", err)
	}
	if immutableProperties == nil {
		return "", fmt.Errorf("private details does not exist in client org's collection: %s", txID)
	}

	return string(immutableProperties), nil
}

func (c *Contract) VerifyAssetProperties(ctx contractapi.TransactionContextInterface, txid string) (bool, error) {
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return false, fmt.Errorf("error getting transient: %v", err)
	}

	/// Asset properties must be retrieved from the transient field as they are private
	immutablePropertiesJSON, ok := transMap["asset_properties"]
	if !ok {
		return false, fmt.Errorf("asset_properties key not found in the transient map")
	}

	asset, err := c.ReadAsset(ctx, txid)
	if err != nil {
		return false, fmt.Errorf("failed to get asset: %v", err)
	}

	collectionOwner := BuildCollectionName(asset.OwnerOrg)
	immutablePropertiesOnChainHash, err := ctx.GetStub().GetPrivateDataHash(collectionOwner, txid)
	if err != nil {
		return false, fmt.Errorf("failed to read asset private properties hash from seller's collection: %v", err)
	}
	if immutablePropertiesOnChainHash == nil {
		return false, fmt.Errorf("asset private properties hash does not exist: %s", txid)
	}

	hash := sha256.New()
	hash.Write(immutablePropertiesJSON)
	calculatedPropertiesHash := hash.Sum(nil)

	// verify that the hash of the passed immutable properties matches the on-chain hash
	if !bytes.Equal(immutablePropertiesOnChainHash, calculatedPropertiesHash) {
		return false, fmt.Errorf("hash %x for passed immutable properties %s does not match on-chain hash %x",
			calculatedPropertiesHash,
			immutablePropertiesJSON,
			immutablePropertiesOnChainHash,
		)
	}

	return true, nil
}

//TODO
//Complex query where you can specifiy the field and data
func (c *Contract) QueryField(ctx TransactionContextInterface, field string, input string) ([]*OrgVisit, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":\"%s\"}}", field, input)

	return getQueryResultForQueryString(ctx, queryString)
}
