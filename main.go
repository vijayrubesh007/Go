package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing mutual fund investments
type SmartContract struct {
	contractapi.Contract
}

// MutualFund represents a mutual fund
type MutualFund struct {
	FundID     string  `json:"fundId"`
	Name       string  `json:"name"`
	Investor   string  `json:"investor"`
	Amount     float64 `json:"amount"`
	ReturnRate float64 `json:"returnRate"`
	Duration   int     `json:"duration"` // in months
}

// InitLedger initializes the ledger with some sample mutual funds
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	mutualFunds := []MutualFund{
		{FundID: "FUND1", Name: "Equity Growth Fund", Investor: "Alice", Amount: 10000.00, ReturnRate: 12.5, Duration: 12},
		{FundID: "FUND2", Name: "Debt Income Fund", Investor: "Bob", Amount: 5000.00, ReturnRate: 8.0, Duration: 24},
		{FundID: "FUND3", Name: "Balanced Fund", Investor: "Charlie", Amount: 20000.00, ReturnRate: 10.0, Duration: 18},
	}

	for i, fund := range mutualFunds {
		fundAsBytes, _ := json.Marshal(fund)
		err := ctx.GetStub().PutState("FUND"+strconv.Itoa(i+1), fundAsBytes)
		if err != nil {
			return fmt.Errorf("failed to add mutual fund %d: %v", i+1, err)
		}
	}

	return nil
}

// AddMutualFund adds a new mutual fund to the ledger
func (s *SmartContract) AddMutualFund(ctx contractapi.TransactionContextInterface, fundID, name, investor string, amount float64, returnRate float64, duration int) error {
	mutualFund := MutualFund{
		FundID:     fundID,
		Name:       name,
		Investor:   investor,
		Amount:     amount,
		ReturnRate: returnRate,
		Duration:   duration,
	}

	fundAsBytes, _ := json.Marshal(mutualFund)
	return ctx.GetStub().PutState(fundID, fundAsBytes)
}

// QueryMutualFund retrieves a mutual fund's details from the ledger
func (s *SmartContract) QueryMutualFund(ctx contractapi.TransactionContextInterface, fundID string) (*MutualFund, error) {
	fundAsBytes, err := ctx.GetStub().GetState(fundID)
	if err != nil {
		return nil, fmt.Errorf("failed to read mutual fund %s: %v", fundID, err)
	}
	if fundAsBytes == nil {
		return nil, fmt.Errorf("mutual fund %s does not exist", fundID)
	}

	var mutualFund MutualFund
	_ = json.Unmarshal(fundAsBytes, &mutualFund)
	return &mutualFund, nil
}

// UpdateMutualFund updates an existing mutual fund's details
func (s *SmartContract) UpdateMutualFund(ctx contractapi.TransactionContextInterface, fundID, name, investor string, amount float64, returnRate float64, duration int) error {
	fundAsBytes, err := ctx.GetStub().GetState(fundID)
	if err != nil {
		return fmt.Errorf("failed to read mutual fund %s: %v", fundID, err)
	}
	if fundAsBytes == nil {
		return fmt.Errorf("mutual fund %s does not exist", fundID)
	}

	mutualFund := MutualFund{
		FundID:     fundID,
		Name:       name,
		Investor:   investor,
		Amount:     amount,
		ReturnRate: returnRate,
		Duration:   duration,
	}

	updatedFundAsBytes, _ := json.Marshal(mutualFund)
	return ctx.GetStub().PutState(fundID, updatedFundAsBytes)
}

// DeleteMutualFund deletes a mutual fund from the ledger
func (s *SmartContract) DeleteMutualFund(ctx contractapi.TransactionContextInterface, fundID string) error {
	return ctx.GetStub().DelState(fundID)
}

// QueryAllMutualFunds retrieves all mutual funds from the ledger
func (s *SmartContract) QueryAllMutualFunds(ctx contractapi.TransactionContextInterface) ([]*MutualFund, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var mutualFunds []*MutualFund
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var mutualFund MutualFund
		_ = json.Unmarshal(queryResponse.Value, &mutualFund)
		mutualFunds = append(mutualFunds, &mutualFund)
	}

	return mutualFunds, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating mutual fund chaincode: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting mutual fund chaincode: %v", err)
	}
}
