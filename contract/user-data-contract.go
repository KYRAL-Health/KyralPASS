package main

import (
	"encoding/json"
	"fmt"
)

func (c *Contract) CreateNewUser(ctx TransactionContextInterface, kyralUID, kyralEncryptedUser, hash_ string) error {
	existingUser, err := ctx.GetStub().GetState(kyralUID)
	if err != nil {
		return fmt.Errorf("Error getting state: %v", err)
	}

	if existingUser != nil {
		return fmt.Errorf("User already exists")
	}

	clientOrgID, err := GetClientOrgID(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to get verified OrgID: %v", err)
	}

	user := User{
		KyralUID:               kyralUID,
		KyralEncryptedUser:     kyralEncryptedUser,
		KyralEncryptedUserHash: hash_,
		OwnerOrg:               clientOrgID,
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v\nUser: %v", err, user)
	}

	if err := ctx.GetStub().PutState(kyralUID, userBytes); err != nil {
		return err
	}

	err = SetAssetStateBasedEndorsement(ctx, kyralUID, clientOrgID)
	if err != nil {
		return fmt.Errorf("failed setting state based endorsement for owner: %v", err)
	}

	if err := ctx.GetStub().SetEvent("create_user", []byte(kyralUID)); err != nil {
		return err
	}

	return nil

}

func (c *Contract) ReadUser(ctx TransactionContextInterface, kyralUID string) (*User, error) {

	userBytes, err := ctx.GetStub().GetState(kyralUID)
	if err != nil {
		return nil, fmt.Errorf("Error getting state: %v", err)
	}

	if userBytes == nil {
		return nil, fmt.Errorf("User does not exist")
	}
	var user *User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}

	return user, nil
}

func (c *Contract) UpdateUserData(ctx TransactionContextInterface, kyralUID, kyralEncryptedUser, hash_ string) error {

	existingUser, err := ctx.GetStub().GetState(kyralUID)
	if err != nil {
		return fmt.Errorf("Error getting state: %v", err)
	}

	if existingUser == nil {
		return fmt.Errorf("User does not exist")
	}

	clientOrgID, err := GetClientOrgID(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to get verified OrgID: %v", err)
	}

	var userCheck *User

	err = json.Unmarshal(existingUser, &userCheck)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user: %v", err)
	}

	if userCheck.OwnerOrg != clientOrgID {
		return fmt.Errorf("This user is managed by another organization. Contact %s to make changes", clientOrgID)
	}

	user := User{
		KyralUID:               kyralUID,
		KyralEncryptedUser:     kyralEncryptedUser,
		KyralEncryptedUserHash: hash_,
		OwnerOrg:               clientOrgID,
	}

	userBytes, err := json.Marshal(user)

	if err != nil {
		return fmt.Errorf("failed to marshal user: %v\nUser: %v", err, user)
	}

	return ctx.GetStub().PutState(kyralUID, userBytes)
}

func (c *Contract) CheckUser(ctx TransactionContextInterface, kyralUID string) (bool, error) {

	userBytes, err := ctx.GetStub().GetState(kyralUID)
	if err != nil {
		return false, fmt.Errorf("Error getting state: %v", err)
	}

	if userBytes == nil {
		return false, nil
	}

	return true, nil
}

func (c *Contract) TransferUser(ctx TransactionContextInterface, kyralUID, orgID string) error {

	if kyralUID == "" || orgID == "" {
		return fmt.Errorf("kyralUID and orgID are required")
	}

	userBytes, err := ctx.GetStub().GetState(kyralUID)
	if err != nil {
		return fmt.Errorf("Error getting state: %v", err)
	}

	if userBytes == nil {
		return fmt.Errorf("User does not exist")
	}

	var user *User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user: %v", err)
	}

	clientOrgID, err := GetClientOrgID(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to get verified OrgID: %v", err)
	}

	if user.OwnerOrg != clientOrgID {
		return fmt.Errorf("User does not belong to your organization")
	}

	user.OwnerOrg = orgID

	userBytes2, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v\nUser: %v", err, user)
	}

	if err := ctx.GetStub().PutState(kyralUID, userBytes2); err != nil {
		return err
	}

	err = SetAssetStateBasedEndorsement(ctx, kyralUID, orgID)
	if err != nil {
		if err := ctx.GetStub().PutState(kyralUID, userBytes); err != nil {
			return err
		}
		return fmt.Errorf("failed setting state based endorsement for owner: %v", err)
	}

	if err := ctx.GetStub().SetEvent("create_user", []byte(kyralUID)); err != nil {
		return err
	}

	return nil
}
