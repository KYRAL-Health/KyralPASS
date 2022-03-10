package main

import (
	"encoding/json"
	"fmt"

	//"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const OrgVisitCollection = "orgVisitCollection"

const transferAgreementObjectType = "transferAgreement"

type Contract struct {
	contractapi.Contract
}

// TransactionContext implementation of TransactionContextInterface.
type TransactionContext struct {
	contractapi.TransactionContext
}

// TransactionContextInterface an interface to describe the minimum required
// functions for a transaction context.
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
}

/************************************

	PRIVATE DATA FUNCTIONS

*************************************/

//CreateOrgVisit2 creates a new orgVisit can be read by all organisations
//Private data is stored in the owners organisation speicfic collection
func (c *Contract) CreateOrgVisit(ctx TransactionContextInterface, txId, description string) error {

	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient: %v", err)
	}

	// Asset properties must be retrieved from the transient field as they are private
	immutablePropertiesJSON, ok := transientMap["transient"]
	if !ok {
		return fmt.Errorf("transient key not found in the transient map")
	}

	// Get client org id and verify it matches peer org id.
	// In this scenario, client is only authorized to read/write private data from its own peer.
	clientOrgID, err := GetClientOrgID(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to get verified OrgID: %v", err)
	}

	asset := OrgVisit{
		KyralTxID:   txId,
		Description: description,
		OwnerOrg:    clientOrgID,
	}

	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to create asset JSON: %v", err)
	}

	err = ctx.GetStub().PutState(asset.KyralTxID, assetBytes)
	if err != nil {
		return fmt.Errorf("failed to put asset in public data: %v", err)
	}

	// Set the endorsement policy such that an owner org peer is required to endorse future updates
	err = SetAssetStateBasedEndorsement(ctx, asset.KyralTxID, clientOrgID)
	if err != nil {
		return fmt.Errorf("failed setting state based endorsement for owner: %v", err)
	}

	// Persist private immutable asset properties to owner's private data collection
	collection := BuildCollectionName(clientOrgID)
	err = ctx.GetStub().PutPrivateData(collection, asset.KyralTxID, immutablePropertiesJSON)
	if err != nil {
		return fmt.Errorf("failed to put Asset private details: %v", err)
	}

	return nil
}

//Public Description
func (c *Contract) ChangeDescription(ctx TransactionContextInterface, txID string, newDescription string) error {

	clientOrgID, err := GetClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("failed to get verified OrgID: %v", err)
	}

	asset, err := c.ReadAsset(ctx, txID)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}

	// Auth check to ensure that client's org actually owns the asset
	if clientOrgID != asset.OwnerOrg {
		return fmt.Errorf("a client from %s cannot update the description of a asset owned by %s", clientOrgID, asset.OwnerOrg)
	}

	asset.Description = newDescription
	updatedAssetJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal asset: %v", err)
	}

	return ctx.GetStub().PutState(txID, updatedAssetJSON)
}

/*
func (c *Contract) DeletePrivateVisit(ctx TransactionContextInterface) error {

	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: %v", err)
	}

	// Asset properties are private, therefore they get passed in transient field
	transientDeleteJSON, ok := transientMap["transient"]
	if !ok {
		return fmt.Errorf("asset to delete not found in the transient map")
	}

	type assetDelete struct {
		ID string `json:"tID"`
	}

	var assetDeleteInput assetDelete
	err = json.Unmarshal(transientDeleteJSON, &assetDeleteInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if len(assetDeleteInput.ID) == 0 {
		return fmt.Errorf("tID field must be a non-empty string")
	}

	// Verify that the client is submitting request to peer in their organization
	err = VerifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("DeleteAsset cannot be performed: Error %v", err)
	}

	log.Printf("Deleting Asset: %v", assetDeleteInput.ID)
	valAsbytes, err := ctx.GetStub().GetPrivateData(OrgVisitCollection, assetDeleteInput.ID) //get the asset from chaincode state
	if err != nil {
		return fmt.Errorf("failed to read asset: %v", err)
	}
	if valAsbytes == nil {
		return fmt.Errorf("asset not found: %v", assetDeleteInput.ID)
	}

	ownerCollection, err := GetCollectionName(ctx) // Get owners collection
	if err != nil {
		return fmt.Errorf("failed to infer private collection name for the org: %v", err)
	}

	//check the asset is in the caller org's private collection
	valAsbytes, err = ctx.GetStub().GetPrivateData(ownerCollection, assetDeleteInput.ID)
	if err != nil {
		return fmt.Errorf("failed to read asset from owner's Collection: %v", err)
	}
	if valAsbytes == nil {
		return fmt.Errorf("asset not found in owner's private Collection %v: %v", ownerCollection, assetDeleteInput.ID)
	}

	// delete the asset from state
	err = ctx.GetStub().DelPrivateData(OrgVisitCollection, assetDeleteInput.ID)
	if err != nil {
		return fmt.Errorf("failed to delete state: %v", err)
	}

	// Finally, delete private details of asset
	err = ctx.GetStub().DelPrivateData(ownerCollection, assetDeleteInput.ID)
	if err != nil {
		return err
	}

	return nil

}
*/
