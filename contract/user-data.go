package main

type User struct {
	KyralUID               string `json:"kyralUID"`
	KyralEncryptedUser     string `json:"kyralEncryptedUser"`
	KyralEncryptedUserHash string `json:"kyralEncryptedUserHash"`
}
