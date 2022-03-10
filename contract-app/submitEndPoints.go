package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
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
		//client.WithEndorsingOrganizations("Org1MSP", "Org2MSP"),
	)

	if err != nil {
		fmt.Printf("error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Println(string(result))
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
	} else {
		w.Write([]byte("Success"))
	}

}
