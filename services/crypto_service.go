package services

import (
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
	"github.com/spf13/viper"
)

func aptosClient() (cli *aptos.Client, err error) {
	client, err := aptos.NewClient(aptos.DevnetConfig)

	return client, nil
}

func getAccountAddress() (*aptos.Account, error) {
	privateKeyHex := viper.GetString("crypto.privateKey")
	privateKey := &crypto.Ed25519PrivateKey{}

	err := privateKey.FromHex(privateKeyHex)
	if err != nil {
		println("privateKey FromHex:", err.Error())
		return nil, err
	}

	account, err := aptos.NewAccountFromSigner(privateKey)

	return account, err
}

func Transfer(receive string, TransferAmount uint64) bool {
	client, err := aptosClient()
	if err != nil {
		println("Failed to create client:" + err.Error())
	}

	account, err := getAccountAddress()
	if err != nil {
		println("Failed to get account address:" + err.Error())
	}

	moduleAddress, err := aptos.ParseHex(viper.GetString("crypto.edenx"))

	receiveBytes, err := aptos.ParseHex(receive)
	if err != nil {
		println("Failed to parse receive address:" + err.Error())
	}

	// 将 amount 转换为字节数组
	amountBytes, err := bcs.SerializeU64(TransferAmount)
	if err != nil {
		println("Failed to serialize transfer amount:" + err.Error())
	}
	println("address:", receive)
	payload := &aptos.EntryFunction{
		Module: aptos.ModuleId{
			Address: aptos.AccountAddress(moduleAddress),
			Name:    "proof_of_achievement",
		},
		Function: "earn",
		ArgTypes: []aptos.TypeTag{},
		Args: [][]byte{
			receiveBytes,
			amountBytes,
		},
	}
	// 1. Build transaction
	rawTxn, err := client.BuildTransaction(account.AccountAddress(), aptos.TransactionPayload{Payload: payload})

	if err != nil {
		panic("Failed to build transaction:" + err.Error())
		return false
	}

	simulationResult, err := client.SimulateTransaction(rawTxn, account)
	if err != nil {
		panic("Failed to simulate transaction:" + err.Error())
		return false
	}
	println("simulationResult:", simulationResult)

	submitResponse, err := client.BuildSignAndSubmitTransaction(account, aptos.TransactionPayload{Payload: payload})
	if err != nil {
		println("Failed to submit transaction: " + err.Error())
		return false
	}

	txn, err := client.WaitForTransaction(submitResponse.Hash)
	if err != nil {
		println("Failed to wait for transaction: " + err.Error())
		return false
	}

	if !txn.Success {
		println("Transaction failed: " + submitResponse.Hash)
	}

	println("Transaction submitted:", txn.Success)
	return true
}
