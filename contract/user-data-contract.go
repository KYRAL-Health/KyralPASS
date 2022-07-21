package main

import (
	"encoding/json"
	"fmt"
)

func (c *Contract) CreateNewUser(ctx TransactionContextInterface, kyralUID, kyralEncryptedUser, hash_ string) error {
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: %v", err)
	}
	decryptPasswordRaw, ok := transientMap["decryptPassword"]
	if !ok {
		return fmt.Errorf("decryptPassword not found in the transient map")
	}

	existingUser, err := ctx.GetStub().GetState(kyralUID)
	if err != nil {
		return fmt.Errorf("Error getting state: %v", err)
	}

	if existingUser != nil {
		return fmt.Errorf("User already exists")
	}

	decryptPassword := string(decryptPasswordRaw[:])

	_, err = GetClientOrgID(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to get verified OrgID: %v", err)
	}

	_, matches, err := decrypt(decryptPassword, kyralEncryptedUser, hash_)
	if err != nil {
		return fmt.Errorf("failed to decrypt user: %v", err)
	}

	if !matches {
		return fmt.Errorf("decrypted user does not match the hash")
	}

	user := User{
		KyralUID:               kyralUID,
		KyralEncryptedUser:     kyralEncryptedUser,
		KyralEncryptedUserHash: hash_,
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v\nUser: %v", err, user)
	}

	if err := ctx.GetStub().PutState(kyralUID, userBytes); err != nil {
		return err
	}

	if err := ctx.GetStub().SetEvent("create_user", []byte(kyralUID)); err != nil {
		return err
	}

	return nil

}

func (c *Contract) ReadUser(ctx TransactionContextInterface, kyralUID, decryptPassword string) (*User, error) {

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
	if _, correctKey, err := decrypt(decryptPassword, user.KyralEncryptedUser, user.KyralEncryptedUserHash); err != nil {
		return nil, fmt.Errorf("failed to decrypt user: %v", err)
	} else if !correctKey {
		return nil, fmt.Errorf("decrypted user does not match the hash")
	}
	return user, nil
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
