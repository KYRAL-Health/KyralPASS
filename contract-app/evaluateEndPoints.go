package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func (contract *NetworkHandler) CheckAsset(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	inputID := params["id"]

	result, err := contract.method.EvaluateTransaction("CheckAsset", inputID)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(result))

}

func (contract *NetworkHandler) getAllOrgVisit(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: GetAllOrgVisit")
	results, err := contract.method.EvaluateTransaction("GetAllOrgVisit")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte(results))
}

func (contract *NetworkHandler) readPrivatePublicProperties(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: ReadAsset")

	params := mux.Vars(r)
	inputID := params["id"]

	results, err := contract.method.EvaluateTransaction("ReadAsset", inputID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(results))
	}
}

func (contract *NetworkHandler) readPrivateProperties(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: ReadPrivateProperties")

	params := mux.Vars(r)
	inputID := params["id"]

	results, err := contract.method.EvaluateTransaction("ReadPrivateProperties", inputID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(results))
	}
}

func (contract *NetworkHandler) queryField(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: queryField")

	w.Header().Set("Content-Type", "application/json")

	var query Field
	_ = json.NewDecoder(r.Body).Decode(&query)

	params := mux.Vars(r)
	inputID := params["id"]

	results, err := contract.method.EvaluateTransaction("QueryField", query.Field, inputID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(results))
	}
}

func (contract *NetworkHandler) verifyAssetProperties(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: VerifyAssetProperties")

	//params := mux.Vars(r)
	//inputID := params["id"]

	var tranM OrgVisitSec
	_ = json.NewDecoder(r.Body).Decode(&tranM)
	if tranM.KyralTxID == "" {
		w.Write([]byte("Input Something Please"))
		return
	}

	bytes, err := json.Marshal(tranM)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to Marshal"))
		return
	}

	//private data vvv
	transient := make(map[string][]byte)
	transient["asset_properties"] = []byte(string(bytes))

	results, err := contract.method.Submit("VerifyAssetProperties",
		client.WithArguments(tranM.KyralTxID),
		client.WithTransient(transient),
		// client.WithEndorsingOrganizations("Org1MSP", "Org2MSP"),
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(results))
	}
}
