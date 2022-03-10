package main

type OrgVisit struct {
	KyralTxID   string `json:"kyralTxID"`
	Description string `json:"description"`
	OwnerOrg    string `json:"OwnerOrg"`
}

type PublicData struct {
	KyralTxID string `json:"kyralTxID"`
	KyralUID  string `json:"kyralUID"`
}

type UserData struct {
	KyralUID  string `json:"kyralUID"`
	KyralTxID string `json:"kyralTxID"`
}
