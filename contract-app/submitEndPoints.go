package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	gwproto "github.com/hyperledger/fabric-protos-go/gateway"
	"google.golang.org/grpc/status"
)

//Contract Apis

func (contract *NetworkHandler) createPrivate(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: CreatePrivate")
	w.Header().Set("Content-Type", "application/json")

	var asset OrgVisit
	_ = json.NewDecoder(r.Body).Decode(&asset)

	var tranM OrgVisitSec

	tranM.OperationCode = asset.OperationCode
	tranM.KyralTxID = asset.KyralTxID
	tranM.KyralPatientID = asset.KyralPatientID
	tranM.Salt = SaltGen()

	bytes, err := json.Marshal(tranM)

	if err != nil {
		log.Println("Failed to marshal private", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	//private data vvv
	transient := make(map[string][]byte)
	transient["transient"] = []byte(string(bytes))

	//Private Data to be Submitted in function "CreateOrgVisit"
	result, err := contract.method.Submit("CreateOrgVisit",
		client.WithArguments(asset.KyralTxID, asset.Description),
		client.WithTransient(transient),

		//Do not enable this line of code vvv
		client.WithEndorsingOrganizations(mspID),
	)

	if err != nil {
		log.Println("Failed to submit private", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
	log.Println("Successfully submitted private", result)
}

func (contract *NetworkHandler) changeDescription(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: ChangeDescription")
	w.Header().Set("Content-Type", "application/json")

	var asset OrgVisit
	_ = json.NewDecoder(r.Body).Decode(&asset)

	_, err := contract.method.SubmitTransaction("ChangeDescription", asset.KyralTxID, asset.Description)

	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed"))
		log.Println("Failed to change description", err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}

}

func (contract *NetworkHandler) createUser(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: CreateUser")
	w.Header().Set("Content-Type", "application/json")

	var user EncryptedUserSubmit
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if user.KyralUID == "" || user.KyralEncryptedUser == "" || user.KyralEncryptedUserHash == "" || user.DecryptKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing Fields."))
		log.Println("Failed to create user. Empty Fields")
		return
	}

	_, matches, err := decrypt(user.DecryptKey, user.KyralEncryptedUser, user.KyralEncryptedUserHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to decrypt user: %v", err)))
		return
	}

	if !matches {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("decrypted user does not match the hash")))
		return
	}

	//Private Data to be Submitted in function "CreateOrgVisit"
	result, err := contract.method.Submit("CreateNewUser",
		client.WithArguments(user.KyralUID, user.KyralEncryptedUser, user.KyralEncryptedUserHash),
	)

	if err != nil {
		switch err := err.(type) {
		case *client.EndorseError:
			fmt.Printf("Endorse error with gRPC status %v: %s\n", status.Code(err), err)
		case *client.SubmitError:
			fmt.Printf("Submit error with gRPC status %v: %s\n", status.Code(err), err)
		case *client.CommitStatusError:
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Printf("Timeout waiting for transaction %s commit status: %s", err.TransactionID, err)
			} else {
				fmt.Printf("Error obtaining commit status with gRPC status %v: %s\n", status.Code(err), err)
			}
		case *client.CommitError:
			fmt.Printf("Transaction %s failed to commit with status %d: %s\n", err.TransactionID, int32(err.Code), err)
		}

		// Any error that originates from a peer or orderer node external to the gateway will have its details
		// embedded within the gRPC status error. The following code shows how to extract that.
		statusErr := status.Convert(err)
		for _, detail := range statusErr.Details() {
			errDetail := detail.(*gwproto.ErrorDetail)
			fmt.Printf("Error from endpoint: %s, mspId: %s, message: %s\n", errDetail.Address, errDetail.MspId, errDetail.Message)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// if err != nil {
	// 	log.Println("Failed to create new user", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
	log.Println("Successfully created user", result)

}

func (contract *NetworkHandler) transferUser(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: TransferUser")
	w.Header().Set("Content-Type", "application/json")

	var userTransfer UserTransferPost
	if err := json.NewDecoder(r.Body).Decode(&userTransfer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if userTransfer.KyralUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("kyralID is required"))
		return
	}

	if userTransfer.OrgID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("orgID is required"))
		return
	}

	if userTransfer.DecryptKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("decryptKey is required"))
		return
	}

	results, err := contract.method.EvaluateTransaction("ReadUser",
		userTransfer.KyralUID,
	)
	if err != nil {
		log.Printf("failed to verify asset properties: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var returnedUser EncryptedUserSubmit
	if err := json.Unmarshal(results, &returnedUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if _, match, err := decrypt(userTransfer.DecryptKey, returnedUser.KyralEncryptedUser, returnedUser.KyralEncryptedUserHash); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else if !match {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect decrypt key"))
		return
	}

	result, err := contract.method.Submit("TransferUser",
		client.WithArguments(userTransfer.KyralUID, userTransfer.OrgID),
	)

	if err != nil {
		switch err := err.(type) {
		case *client.EndorseError:
			fmt.Printf("Endorse error with gRPC status %v: %s\n", status.Code(err), err)
		case *client.SubmitError:
			fmt.Printf("Submit error with gRPC status %v: %s\n", status.Code(err), err)
		case *client.CommitStatusError:
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Printf("Timeout waiting for transaction %s commit status: %s", err.TransactionID, err)
			} else {
				fmt.Printf("Error obtaining commit status with gRPC status %v: %s\n", status.Code(err), err)
			}
		case *client.CommitError:
			fmt.Printf("Transaction %s failed to commit with status %d: %s\n", err.TransactionID, int32(err.Code), err)
		}

		// Any error that originates from a peer or orderer node external to the gateway will have its details
		// embedded within the gRPC status error. The following code shows how to extract that.
		statusErr := status.Convert(err)
		for _, detail := range statusErr.Details() {
			errDetail := detail.(*gwproto.ErrorDetail)
			fmt.Printf("Error from endpoint: %s, mspId: %s, message: %s\n", errDetail.Address, errDetail.MspId, errDetail.Message)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// if err != nil {
	// 	log.Println("Failed to create new user", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
	log.Println("Successfully transfered user to "+userTransfer.OrgID, result)

}

func (contract *NetworkHandler) updateUser(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: UpdateUser")
	w.Header().Set("Content-Type", "application/json")

	var updatedUser EncryptedUserSubmit
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if updatedUser.KyralUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("kyralID is required"))
		return
	}

	if updatedUser.KyralEncryptedUser == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("kyralEncryptedUser is required"))
		return
	}

	if updatedUser.KyralEncryptedUserHash == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("kyralEncryptedUserHash is required"))
		return
	}

	if updatedUser.DecryptKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("decryptKey is required"))
		return
	}

	results, err := contract.method.EvaluateTransaction("ReadUser",
		updatedUser.KyralUID,
	)
	if err != nil {
		log.Printf("failed to verify asset properties: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var returnedUser EncryptedUserSubmit
	if err := json.Unmarshal(results, &returnedUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if _, match, err := decrypt(updatedUser.DecryptKey, returnedUser.KyralEncryptedUser, returnedUser.KyralEncryptedUserHash); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else if !match {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect decrypt key"))
		return
	}

	if _, match, err := decrypt(updatedUser.DecryptKey, updatedUser.KyralEncryptedUser, updatedUser.KyralEncryptedUserHash); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else if !match {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect hash"))
		return
	}

	result, err := contract.method.Submit("UpdateUserData",
		client.WithArguments(updatedUser.KyralUID, updatedUser.KyralEncryptedUser, updatedUser.KyralEncryptedUserHash),
	)

	if err != nil {
		switch err := err.(type) {
		case *client.EndorseError:
			fmt.Printf("Endorse error with gRPC status %v: %s\n", status.Code(err), err)
		case *client.SubmitError:
			fmt.Printf("Submit error with gRPC status %v: %s\n", status.Code(err), err)
		case *client.CommitStatusError:
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Printf("Timeout waiting for transaction %s commit status: %s", err.TransactionID, err)
			} else {
				fmt.Printf("Error obtaining commit status with gRPC status %v: %s\n", status.Code(err), err)
			}
		case *client.CommitError:
			fmt.Printf("Transaction %s failed to commit with status %d: %s\n", err.TransactionID, int32(err.Code), err)
		}

		// Any error that originates from a peer or orderer node external to the gateway will have its details
		// embedded within the gRPC status error. The following code shows how to extract that.
		statusErr := status.Convert(err)
		for _, detail := range statusErr.Details() {
			errDetail := detail.(*gwproto.ErrorDetail)
			fmt.Printf("Error from endpoint: %s, mspId: %s, message: %s\n", errDetail.Address, errDetail.MspId, errDetail.Message)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
	log.Println("Successfully updated user "+updatedUser.KyralUID, result)

}
