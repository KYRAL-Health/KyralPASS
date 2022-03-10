package main

type OrgVisitSec struct {
	KyralTxID      string `json:"kyralTxID"`
	KyralPatientID string `json:"kyralPatientID"`
	OperationCode  string `json:"operationCode"`
	Salt           string `json:"salt"`
}

type OrgVisit struct {
	KyralTxID      string `json:"kyralTxID"`
	Description    string `json:"description"`
	KyralPatientID string `json:"kyralPatientID"`
	OrgID          string `json:"orgID"`
	OperationCode  string `json:"operationCode"`
	Salt           string `json:"salt"`
}

type UserData struct {
	KyralUID  string `json:"kyralUID"`
	KyralTxID string `json:"kyralTxID"`
}

type Field struct {
	Field string `json:"field"`
}
