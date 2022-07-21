package main

import (
	"encoding/json"
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
		log.Printf("Failed to evaluate transaction: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))

	log.Println("check asset success", result)

}

func (contract *NetworkHandler) getAllOrgVisit(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: GetAllOrgVisit")
	results, err := contract.method.EvaluateTransaction("GetAllOrgVisit")
	if err != nil {
		log.Printf("Failed to getAllOrgVisit: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(results))
}

func (contract *NetworkHandler) readPrivatePublicProperties(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: ReadAsset")

	params := mux.Vars(r)
	inputID := params["id"]

	results, err := contract.method.EvaluateTransaction("ReadAsset", inputID)
	if err != nil {
		log.Printf("failed to read public properties: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(results))

}

func (contract *NetworkHandler) readPrivateProperties(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: ReadPrivateProperties")

	params := mux.Vars(r)
	inputID := params["id"]

	results, err := contract.method.EvaluateTransaction("ReadPrivateProperties", inputID)
	if err != nil {
		log.Printf("failed to read private properties: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(results))

}

func (contract *NetworkHandler) queryField(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: queryField")

	w.Header().Set("Content-Type", "application/json")

	var query Field
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		log.Printf("failed to decode query: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	params := mux.Vars(r)
	inputID := params["id"]

	results, err := contract.method.EvaluateTransaction("QueryField", query.Field, inputID)
	if err != nil {
		log.Printf("failed to query field: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(results))

}

func (contract *NetworkHandler) verifyAssetProperties(w http.ResponseWriter, r *http.Request) {

	log.Println("--> Submit: VerifyAssetProperties")

	//params := mux.Vars(r)
	//inputID := params["id"]

	var tranM OrgVisitSec
	_ = json.NewDecoder(r.Body).Decode(&tranM)
	if tranM.KyralTxID == "" {
		log.Printf("failed to decode query: %v\n", "KyralTxID is empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Input Something Please"))
		return
	}

	bytes, err := json.Marshal(tranM)

	if err != nil {
		log.Printf("failed to marshal: %v\n", err)
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
		log.Printf("failed to verify asset properties: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(results))

}

func (contract *NetworkHandler) readUser(w http.ResponseWriter, r *http.Request) {
	log.Println("--> Submit: ReadUser")
	var userPost GetUser
	_ = json.NewDecoder(r.Body).Decode(&userPost)
	if userPost.KyralUID == "" {
		log.Printf("failed to decode query: %v\n", "KyralUID is empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("KyralUID is required"))
		return
	}
	if userPost.DecryptKey == "" {
		log.Printf("failed to decode query: %v\n", "DecryptKey is empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("DecryptKey is required"))
		return
	}

	transient := make(map[string][]byte)
	transient["decryptPassword"] = []byte(userPost.DecryptKey)

	results, err := contract.method.EvaluateTransaction("ReadUser",
		userPost.KyralUID, userPost.DecryptKey,
	)
	if err != nil {
		log.Printf("failed to verify asset properties: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// json header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func (contract *NetworkHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	log.Println("--> Submit: CheckUser")

	params := mux.Vars(r)
	inputID := params["id"]

	results, err := contract.method.EvaluateTransaction("CheckUser",
		inputID,
	)
	if err != nil {
		log.Printf("failed to check user: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(results)
}
