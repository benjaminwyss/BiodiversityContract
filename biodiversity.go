package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type BiodiversityRecord struct {
	Kind     string `json:"kind"`
	Genus    string `json:"genus"`
	Species  string `json:"species"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	sampleRecord := BiodiversityRecord{"Bird", "Topaza", "T. Pella", "Crimson Topaz", "Northern South America"}
	sampleRecordBytes, _ := json.Marshal(sampleRecord)
	err := ctx.GetStub().PutState("RECORD"+strconv.Itoa(1), sampleRecordBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}
	return nil
}

func (s *SmartContract) Create(ctx contractapi.TransactionContextInterface, recordKey string, kind string, genus string, species string, name string, location string) error {
	record := BiodiversityRecord{kind, genus, species, name, location}

	recordBytes, _ := json.Marshal(record)

	return ctx.GetStub().PutState(recordKey, recordBytes)
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, recordKey string) (*BiodiversityRecord, error) {
	recordBytes, err := ctx.GetStub().GetState(recordKey)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if recordBytes == nil {
		return nil, fmt.Errorf("%s does not exist", recordKey)
	}

	record := new(BiodiversityRecord)
	_ = json.Unmarshal(recordBytes, record)

	return record, nil
}

func (s *SmartContract) GetHistory(ctx contractapi.TransactionContextInterface, recordKey string) (string, error) {
	recordIterator, err := ctx.GetStub().GetHistoryForKey(recordKey)

	if err != nil {
		return "", fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	defer recordIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	isArrayMemberWritten := false
	for recordIterator.HasNext() {
		response, err := recordIterator.Next()

		if err != nil {
			return "", fmt.Errorf("Error. %s", err.Error())
		}

		if isArrayMemberWritten {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the corresponding value null.
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		isArrayMemberWritten = true
	}

	buffer.WriteString("]")

	return buffer.String(), nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating biodiversity chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting biodiversity chaincode: %s", err.Error())
	}
}
