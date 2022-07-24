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
	_ = json.NewDecoder(r.Body).Decode(&user)

	if user.KyralUID == "" || user.KyralEncryptedUser == "" || user.KyralEncryptedUserHash == "" || user.DecryptKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing Fields."))
		log.Println("Failed to create user. Empty Fields")
		return
	}

	//private data vvv
	transient := make(map[string][]byte)
	transient["decryptPassword"] = []byte(user.DecryptKey)

	//Private Data to be Submitted in function "CreateOrgVisit"
	result, err := contract.method.Submit("CreateNewUser",
		client.WithArguments(user.KyralUID, user.KyralEncryptedUser, user.KyralEncryptedUserHash),
		client.WithTransient(transient),
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
