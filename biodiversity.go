package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Specimen struct {
	Collection      string `json:"collection"`
	Updater         string `json:"updater"`
	CatalogNumber   string `json:"catalogNumber"`
	AccessionNumber string `json:"accessionNumber"`
	CatalogDate     string `json:"catalogDate"`
	Cataloger       string `json:"cataloger"`
	Taxon           string `json:"taxon"`
	Determiner      string `json:"determiner"`
	DetermineDate   string `json:"determineDate"`
	FieldNumber     string `json:"fieldNumber"`
	FieldDate       string `json:"fieldDate"`
	Collector       string `json:"collector"`
	Location        string `json:"location"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
	Habitat         string `json:"habitat"`
	Preparation     string `json:"preparation"`
	Condition       string `json:"condition"`
	Loans           string `json:"loans"`
	Grants          string `json:"grants"`
	Notes           string `json:"notes"`
	Image           string `json:"image"`
}

type Collection struct {
	Name            string `json:"name"`
	CreateSpecimen  string `json:"createSpecimen"`
	PrimaryUpdate   string `json:"primaryUpdate"`
	SecondaryUpdate string `json:"secondaryUpdate"`
	Georeference    string `json:"georeference"`
	LinkImages      string `json:"linkImages"`
	LinkAuxiliary   string `json:"linkAuxiliary"`
	TaxonName       string `json:"taxonName"`
	TaxonClass      string `json:"taxonClass"`
	SuggestTaxon    string `json:"suggestTaxon"`
	RegisterLoan    string `json:"registerLoan"`
	RegisterUse     string `json:"registerUse"`
	Query           string `json:"query"`
	FlagError       string `json:"flagError"`
}

type User struct {
	Username   string            `json:"username"`
	Membership map[string]string `json:"membership"`
}

type QueryResult struct {
	Guid   string    `json:"guid"`
	Record *Specimen `json:"specimen"`
}

type PendingTransaction struct {
	Transaction string   `json:"transaction"`
	Arguments   []string `json:"arguments"`
	Suggester   string   `json:"suggester"`
	Reason      string   `json:"reason"`
}

func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	sampleCollection := Collection{"KU Ornithology", "M", "MC", "MCA", "MCA", "MCAS", "MCA", "MC", "MC", "MCA", "MCAS", "MCAS", "MCASP", "MCASP"}
	collectionBytes, _ := json.Marshal(sampleCollection)
	err := ctx.GetStub().PutState("KU Ornithology", collectionBytes)

	if err != nil {
		return fmt.Errorf("Failed to put collection to world state. %s", err.Error())
	}

	managerMap := make(map[string]string)
	managerMap["KU Ornithology"] = "M"
	sampleManager := User{"manager", managerMap}
	managerBytes, _ := json.Marshal(sampleManager)
	err = ctx.GetStub().PutState("manager", managerBytes)

	if err != nil {
		return fmt.Errorf("Failed to put manager to world state. %s", err.Error())
	}

	curatorMap := make(map[string]string)
	curatorMap["KU Ornithology"] = "C"
	sampleCurator := User{"curator", curatorMap}
	curatorBytes, _ := json.Marshal(sampleCurator)
	err = ctx.GetStub().PutState("curator", curatorBytes)

	if err != nil {
		return fmt.Errorf("Failed to put curator to world state. %s", err.Error())
	}

	assistantMap := make(map[string]string)
	assistantMap["KU Ornithology"] = "A"
	sampleAssistant := User{"assistant", assistantMap}
	assistantBytes, _ := json.Marshal(sampleAssistant)
	err = ctx.GetStub().PutState("assistant", assistantBytes)

	if err != nil {
		return fmt.Errorf("Failed to put assistant to world state. %s", err.Error())
	}

	studentMap := make(map[string]string)
	studentMap["KU Ornithology"] = "S"
	sampleStudent := User{"student", studentMap}
	studentBytes, _ := json.Marshal(sampleStudent)
	err = ctx.GetStub().PutState("student", studentBytes)

	if err != nil {
		return fmt.Errorf("Failed to put student to world state. %s", err.Error())
	}

	publicMap := make(map[string]string)
	publicMap["KU Ornithology"] = "P"
	samplePublic := User{"public", publicMap}
	publicBytes, _ := json.Marshal(samplePublic)
	err = ctx.GetStub().PutState("public", publicBytes)

	if err != nil {
		return fmt.Errorf("Failed to put public to world state. %s", err.Error())
	}

	sampleSpecimen := Specimen{"KU Ornithology", "manager", "32581", "2002-IC-062", "06/19/2003", "Bentley, Andy C", "Pygoplites diacanthus", "Greenfield, David W", "", "G02-15", "01/27/2002", "", "Fiji, Viti Levu", "18.1483325958", "-178.3984985352", "Barrier reef off Suva Point north of wreck in main channel", "", "", "", "", "", ""}
	specimenBytes, _ := json.Marshal(sampleSpecimen)
	err = ctx.GetStub().PutState("0", specimenBytes)

	if err != nil {
		return fmt.Errorf("Failed to put specimen to world state. %s", err.Error())
	}

	return nil
}

func (s *SmartContract) RegisterCollection(ctx contractapi.TransactionContextInterface, name string, username string, createSpecimen string, primaryUpdate string, secondaryUpdate string, georeference string, linkImages string, linkAuxiliary string, taxonName string, taxonClass string, suggestTaxon string, registerLoan string, registerUse string, query string, flagError string) error {
	checkExistence, err := ctx.GetStub().GetState(name)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if checkExistence != nil {
		return fmt.Errorf("%s already exists", name)
	}

	checkUser, err := ctx.GetStub().GetState(username)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if checkUser == nil {
		return fmt.Errorf("%s does not exists", username)
	}

	attributionString := fmt.Sprintf("Registered Collection %s", name)
	attributionBytes := []byte(attributionString)
	err = ctx.GetStub().PutState(username+"|attribution", attributionBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	collection := Collection{name, createSpecimen, primaryUpdate, secondaryUpdate, georeference, linkImages, linkAuxiliary, taxonName, taxonClass, suggestTaxon, registerLoan, registerUse, query, flagError}
	collectionBytes, _ := json.Marshal(collection)
	err = ctx.GetStub().PutState(name, collectionBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	user := new(User)
	_ = json.Unmarshal(checkUser, user)

	user.Membership[name] = "M"
	userBytes, _ := json.Marshal(user)
	return ctx.GetStub().PutState(username, userBytes)

}

func (s *SmartContract) UpdateCollection(ctx contractapi.TransactionContextInterface, name string, username string, createSpecimen string, primaryUpdate string, secondaryUpdate string, georeference string, linkImages string, linkAuxiliary string, taxonName string, taxonClass string, suggestTaxon string, registerLoan string, registerUse string, query string, flagError string) error {
	checkExistence, err := ctx.GetStub().GetState(name)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if checkExistence == nil {
		return fmt.Errorf("%s does not exists", name)
	}

	checkUser, err := ctx.GetStub().GetState(username)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if checkUser == nil {
		return fmt.Errorf("%s does not exists", username)
	}

	oldCollection := new(Collection)
	_ = json.Unmarshal(checkExistence, oldCollection)

	user := new(User)
	_ = json.Unmarshal(checkUser, user)

	if role, ok := user.Membership[name]; ok {
		if role != "M" {
			return fmt.Errorf("%s is not the Manager for collection %s", username, name)
		}
	} else {
		return fmt.Errorf("%s is not registered with collection %s", username, name)
	}

	//Don't overwrite existing data with blank data
	if createSpecimen == "" {
		createSpecimen = oldCollection.CreateSpecimen
	}
	if primaryUpdate == "" {
		primaryUpdate = oldCollection.PrimaryUpdate
	}
	if secondaryUpdate == "" {
		secondaryUpdate = oldCollection.SecondaryUpdate
	}
	if georeference == "" {
		georeference = oldCollection.Georeference
	}
	if linkImages == "" {
		linkImages = oldCollection.LinkImages
	}
	if linkAuxiliary == "" {
		linkAuxiliary = oldCollection.LinkAuxiliary
	}
	if taxonName == "" {
		taxonName = oldCollection.TaxonName
	}
	if taxonClass == "" {
		taxonClass = oldCollection.TaxonClass
	}
	if suggestTaxon == "" {
		suggestTaxon = oldCollection.SuggestTaxon
	}
	if registerLoan == "" {
		registerLoan = oldCollection.RegisterLoan
	}
	if registerUse == "" {
		registerUse = oldCollection.RegisterUse
	}
	if query == "" {
		query = oldCollection.Query
	}
	if flagError == "" {
		flagError = oldCollection.FlagError
	}

	attributionString := fmt.Sprintf("Updated Collection %s access control policies", name)
	attributionBytes := []byte(attributionString)
	err = ctx.GetStub().PutState(username+"|attribution", attributionBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	collection := Collection{name, createSpecimen, primaryUpdate, secondaryUpdate, georeference, linkImages, linkAuxiliary, taxonName, taxonClass, suggestTaxon, registerLoan, registerUse, query, flagError}
	collectionBytes, _ := json.Marshal(collection)
	return ctx.GetStub().PutState(name, collectionBytes)
}

func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, username string) error {
	checkExistence, err := ctx.GetStub().GetState(username)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if checkExistence != nil {
		return fmt.Errorf("%s already exists", username)
	}

	emptyMap := make(map[string]string)
	user := User{username, emptyMap}
	userBytes, _ := json.Marshal(user)
	return ctx.GetStub().PutState(username, userBytes)

}

func (s *SmartContract) GrantPermission(ctx contractapi.TransactionContextInterface, granterName string, username string, collection string, permission string) error {
	if permission != "M" && permission != "C" && permission != "A" && permission != "S" && permission != "P" {
		return fmt.Errorf("%s is not a valid permission. Valid permissions are M, C, A, S, and P", permission)
	}

	checkUser, err := ctx.GetStub().GetState(username)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if checkUser == nil {
		return fmt.Errorf("%s does not exists", username)
	}

	checkGranter, err := ctx.GetStub().GetState(granterName)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if checkUser == nil {
		return fmt.Errorf("%s does not exists", granterName)
	}

	checkCollection, err := ctx.GetStub().GetState(collection)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if checkCollection == nil {
		return fmt.Errorf("%s does not exists", collection)
	}

	user := new(User)
	_ = json.Unmarshal(checkUser, user)

	granter := new(User)
	_ = json.Unmarshal(checkGranter, granter)

	role, ok := granter.Membership[collection]

	granteeRole, granteeOk := user.Membership[collection]

	if !granteeOk {
		granteeRole = "P"
	}

	if ok {
		if role != "M" && role != "C" {
			return fmt.Errorf("%s is not a Manager of Curator of collection %s", granterName, collection)
		}
		if role == "C" && permission == "M" {
			return fmt.Errorf("%s is a Curator for collection %s and cannot grant permission of Manager to %s", granterName, collection, username)
		}
		if role == "C" && granteeRole == "M" {
			return fmt.Errorf("%s is a Curator for collection %s and cannot change permission of Manager %s", granterName, collection, username)
		}
	} else {
		return fmt.Errorf("%s is not registered with collection %s", granterName, collection)
	}

	attributionString := fmt.Sprintf("Updated %s permission to %s in collection %s", username, permission, collection)
	attributionBytes := []byte(attributionString)
	err = ctx.GetStub().PutState(granterName+"|attribution", attributionBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	user.Membership[collection] = permission
	userBytes, _ := json.Marshal(user)
	return ctx.GetStub().PutState(username, userBytes)

}

func (s *SmartContract) Create(ctx contractapi.TransactionContextInterface, guid string, collection string, updater string, catalogNumber string, accessionNumber string, catalogDate string, cataloger string, taxon string, determiner string, determineDate string, fieldNumber string, fieldDate string, collector string, location string, latitude string, longitude string, habitat string, preparation string, condition string, notes string, image string) error {
	checkExistence, err := ctx.GetStub().GetState(guid)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkExistence != nil {
		return fmt.Errorf("%s already exists", guid)
	}

	checkUpdater, err := ctx.GetStub().GetState(updater)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkUpdater == nil {
		return fmt.Errorf("%s does not exist", updater)
	}

	checkCollection, err := ctx.GetStub().GetState(collection)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if checkCollection == nil {
		return fmt.Errorf("%s does not exists", collection)
	}

	collect := new(Collection)
	_ = json.Unmarshal(checkCollection, collect)

	user := new(User)
	_ = json.Unmarshal(checkUpdater, user)

	role, ok := user.Membership[collection]

	if !ok {
		role = "P"
	}

	if !strings.Contains(collect.CreateSpecimen, role) {
		return fmt.Errorf("%s has role %s but role %s is required to create specimen", updater, role, collect.CreateSpecimen)
	}

	attributionString := fmt.Sprintf("Created Specimen with GUID %s", guid)
	attributionBytes := []byte(attributionString)
	err = ctx.GetStub().PutState(updater+"|attribution", attributionBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	specimen := Specimen{collection, updater, catalogNumber, accessionNumber, catalogDate, cataloger, taxon, determiner, determineDate, fieldNumber, fieldDate, collector, location, latitude, longitude, habitat, preparation, condition, "", "", notes, image}

	specimenBytes, _ := json.Marshal(specimen)

	return ctx.GetStub().PutState(guid, specimenBytes)
}

func (s *SmartContract) Update(ctx contractapi.TransactionContextInterface, guid string, collection string, updater string, catalogNumber string, accessionNumber string, catalogDate string, cataloger string, taxon string, determiner string, determineDate string, fieldNumber string, fieldDate string, collector string, location string, latitude string, longitude string, habitat string, preparation string, condition string, conditionDate string, notes string, image string) error {
	checkExistence, err := ctx.GetStub().GetState(guid)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkExistence == nil {
		return fmt.Errorf("%s does not exists", guid)
	}

	checkUpdater, err := ctx.GetStub().GetState(updater)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkUpdater == nil {
		return fmt.Errorf("%s does not exist", updater)
	}

	checkCollection, err := ctx.GetStub().GetState(collection)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if checkCollection == nil {
		return fmt.Errorf("%s does not exists", collection)
	}

	collect := new(Collection)
	_ = json.Unmarshal(checkCollection, collect)

	user := new(User)
	_ = json.Unmarshal(checkUpdater, user)

	role, ok := user.Membership[collection]

	if !ok {
		role = "P"
	}

	oldSpecimen := new(Specimen)
	_ = json.Unmarshal(checkExistence, oldSpecimen)

	//Don't overwrite existing data with blank data
	if collection == "" {
		collection = oldSpecimen.Collection
	}
	if catalogNumber == "" {
		catalogNumber = oldSpecimen.CatalogNumber
	}
	if accessionNumber == "" {
		accessionNumber = oldSpecimen.AccessionNumber
	}
	if catalogDate == "" {
		catalogDate = oldSpecimen.CatalogDate
	}
	if cataloger == "" {
		cataloger = oldSpecimen.Cataloger
	}
	if taxon == "" {
		taxon = oldSpecimen.Taxon
	}
	if determiner == "" {
		determiner = oldSpecimen.Determiner
	}
	if determineDate == "" {
		determineDate = oldSpecimen.DetermineDate
	}
	if fieldNumber == "" {
		fieldNumber = oldSpecimen.FieldNumber
	}
	if fieldDate == "" {
		fieldDate = oldSpecimen.FieldDate
	}
	if collector == "" {
		collector = oldSpecimen.Collector
	}
	if location == "" {
		location = oldSpecimen.Location
	}
	if latitude == "" {
		latitude = oldSpecimen.Latitude
	}
	if longitude == "" {
		longitude = oldSpecimen.Longitude
	}
	if habitat == "" {
		habitat = oldSpecimen.Habitat
	}
	if preparation == "" {
		preparation = oldSpecimen.Preparation
	}
	if condition != "" {
		condition = oldSpecimen.Condition + condition + " " + conditionDate + "\n"
	} else {
		condition = oldSpecimen.Condition
	}
	if notes != "" {
		notes = oldSpecimen.Notes + notes + "\n"
	} else {
		notes = oldSpecimen.Notes
	}
	if image == "" {
		image = oldSpecimen.Image
	}

	if collection != oldSpecimen.Collection {
		return fmt.Errorf("collection %s does not match existing specimen collection %s", collection, oldSpecimen.Collection)
	}

	if catalogNumber != oldSpecimen.CatalogNumber || accessionNumber != oldSpecimen.AccessionNumber || catalogDate != oldSpecimen.CatalogDate || cataloger != oldSpecimen.Cataloger || fieldNumber != oldSpecimen.FieldNumber || fieldDate != oldSpecimen.FieldDate || collector != oldSpecimen.Collector {
		if !strings.Contains(collect.PrimaryUpdate, role) {
			return fmt.Errorf("%s has role %s but role %s is required to update primary info", updater, role, collect.PrimaryUpdate)
		}
	}

	if location != oldSpecimen.Location || latitude != oldSpecimen.Latitude || longitude != oldSpecimen.Longitude || habitat != oldSpecimen.Habitat {
		if !strings.Contains(collect.Georeference, role) {
			return fmt.Errorf("%s has role %s but role %s is required to update geolocation info", updater, role, collect.Georeference)
		}
	}

	if preparation != oldSpecimen.Preparation || condition != oldSpecimen.Condition || notes != oldSpecimen.Notes {
		if !strings.Contains(collect.SecondaryUpdate, role) {
			return fmt.Errorf("%s has role %s but role %s is required to update secondary info", updater, role, collect.SecondaryUpdate)
		}
	}

	if taxon != oldSpecimen.Taxon || determiner != oldSpecimen.Determiner || determineDate != oldSpecimen.DetermineDate {
		if !strings.Contains(collect.TaxonName, role) {
			return fmt.Errorf("%s has role %s but role %s is required to update taxon name", updater, role, collect.TaxonName)
		}
	}

	if image != oldSpecimen.Image {
		if !strings.Contains(collect.LinkImages, role) {
			return fmt.Errorf("%s has role %s but role %s is required to link images", updater, role, collect.LinkImages)
		}
	}

	specimen := Specimen{collection, updater, catalogNumber, accessionNumber, catalogDate, cataloger, taxon, determiner, determineDate, fieldNumber, fieldDate, collector, location, latitude, longitude, habitat, preparation, condition, oldSpecimen.Loans, oldSpecimen.Grants, notes, image}

	//Check if an actual change was made
	specimen.Updater = oldSpecimen.Updater
	if cmp.Equal(specimen, *oldSpecimen) {
		return fmt.Errorf("Updated specimen is equivalent to old specimen. Operation aborted to conserve blockchain resources")
	}
	specimen.Updater = updater

	attributionString := fmt.Sprintf("Updated Specimen with GUID %s", guid)
	attributionBytes := []byte(attributionString)
	err = ctx.GetStub().PutState(updater+"|attribution", attributionBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	specimenBytes, _ := json.Marshal(specimen)

	return ctx.GetStub().PutState(guid, specimenBytes)
}

func (s *SmartContract) SuggestUpdate(ctx contractapi.TransactionContextInterface, guid string, collection string, updater string, catalogNumber string, accessionNumber string, catalogDate string, cataloger string, taxon string, determiner string, determineDate string, fieldNumber string, fieldDate string, collector string, location string, latitude string, longitude string, habitat string, preparation string, condition string, conditionDate string, notes string, image string, reason string) error {
	checkExistence, err := ctx.GetStub().GetState(guid)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkExistence == nil {
		return fmt.Errorf("%s does not exists", guid)
	}

	checkUpdater, err := ctx.GetStub().GetState(updater)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkUpdater == nil {
		return fmt.Errorf("%s does not exist", updater)
	}

	if collection == "" {
		specimen := new(Specimen)
		_ = json.Unmarshal(checkExistence, specimen)
		collection = specimen.Collection
	}

	checkCollection, err := ctx.GetStub().GetState(collection)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if checkCollection == nil {
		return fmt.Errorf("%s does not exists", collection)
	}

	collect := new(Collection)
	_ = json.Unmarshal(checkCollection, collect)

	user := new(User)
	_ = json.Unmarshal(checkUpdater, user)

	role, ok := user.Membership[collection]

	if !ok {
		role = "P"
	}

	if !strings.Contains(collect.FlagError, role) {
		return fmt.Errorf("%s has role %s but role %s is required to suggest updates", updater, role, collect.FlagError)
	}

	checkPendingTransactions, err := ctx.GetStub().GetState("pending" + guid)
	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	pendingTransactions := []PendingTransaction{}

	if checkPendingTransactions != nil {
		_ = json.Unmarshal(checkPendingTransactions, &pendingTransactions)
	}

	pendingTransaction := PendingTransaction{"Update", []string{guid, collection, updater, catalogNumber, accessionNumber, catalogDate, cataloger, taxon, determiner, determineDate, fieldNumber, fieldDate, collector, location, latitude, longitude, habitat, preparation, condition, conditionDate, notes, image}, updater, reason}

	pendingTransactions = append(pendingTransactions, pendingTransaction)

	pendingTransactionsBytes, _ := json.Marshal(pendingTransactions)

	attributionString := fmt.Sprintf("Suggested update to specimen with GUID %s", guid)
	attributionBytes := []byte(attributionString)
	err = ctx.GetStub().PutState(updater+"|attribution", attributionBytes)

	return ctx.GetStub().PutState("pending"+guid, pendingTransactionsBytes)

}

func (s *SmartContract) ApproveTransaction(ctx contractapi.TransactionContextInterface, guid string, username string, transactionIndex string) error {
	checkTransactions, err := ctx.GetStub().GetState("pending" + guid)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkTransactions == nil {
		return fmt.Errorf("Pending transaction for %s does not exists", guid)
	}

	transactions := []PendingTransaction{}
	_ = json.Unmarshal(checkTransactions, &transactions)

	index, err := strconv.Atoi(transactionIndex)

	if err != nil {
		return fmt.Errorf("Error. Provided pending transaction index is not an integer. %s", err.Error())
	}

	if index < 0 || index >= len(transactions) {
		return fmt.Errorf("Error. Provided pending transaction index does not correspond to an existing pending transaction.")
	}

	transaction := transactions[index]

	if transaction.Transaction == "Update" {
		args := transaction.Arguments

		notes := args[20]
		if notes == "" {
			notes = "Approved update suggested by user " + args[2]
		} else {
			notes = notes + "\n" + "Approved update suggested by user " + args[2]
		}

		attributionString := fmt.Sprintf("Approved suggested update to specimen with GUID %s", guid)
		attributionBytes := []byte(attributionString)
		err = ctx.GetStub().PutState(username+"|attribution", attributionBytes)

		//remove the approved transaction from the list of pending transactions
		transactions = append(transactions[:index], transactions[index+1:]...)
		transactionsBytes, _ := json.Marshal(transactions)
		ctx.GetStub().PutState("pending"+guid, transactionsBytes)

		return s.Update(ctx, args[0], args[1], username, args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14], args[15], args[16], args[17], args[18], args[19], notes, args[21])

	}

	return fmt.Errorf("Error, pending transaction name not valid.")
}

func (s *SmartContract) DenyTransaction(ctx contractapi.TransactionContextInterface, guid string, username string, transactionIndex string) error {
	checkTransactions, err := ctx.GetStub().GetState("pending" + guid)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkTransactions == nil {
		return fmt.Errorf("Pending transaction for %s does not exists", guid)
	}

	transactions := []PendingTransaction{}
	_ = json.Unmarshal(checkTransactions, &transactions)

	index, err := strconv.Atoi(transactionIndex)

	if err != nil {
		return fmt.Errorf("Error. Provided pending transaction indexd is not an integer. %s", err.Error())
	}

	if index < 0 || index >= len(transactions) {
		return fmt.Errorf("Error. Provided pending transaction index does not correspond to an existing pending transaction.")
	}

	checkUser, err := ctx.GetStub().GetState(username)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkUser == nil {
		return fmt.Errorf("%s does not exist", username)
	}

	user := new(User)
	_ = json.Unmarshal(checkUser, user)

	checkSpecimen, err := ctx.GetStub().GetState(guid)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkSpecimen == nil {
		return fmt.Errorf("specimen with guid %s does not exist", guid)
	}

	specimen := new(Specimen)
	_ = json.Unmarshal(checkSpecimen, specimen)

	collection := specimen.Collection

	checkCollection, err := ctx.GetStub().GetState(collection)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkCollection == nil {
		return fmt.Errorf("Collection %s does not exist", collection)
	}

	collect := new(Collection)
	_ = json.Unmarshal(checkCollection, collect)

	role, ok := user.Membership[collection]

	if !ok {
		role = "P"
	}

	if !strings.Contains(collect.PrimaryUpdate, role) {
		return fmt.Errorf("%s has role %s but role %s is required to update primary info", username, role, collect.PrimaryUpdate)
	}

	transactions = append(transactions[:index], transactions[index+1:]...)
	transactionsBytes, _ := json.Marshal(transactions)
	return ctx.GetStub().PutState("pending"+guid, transactionsBytes)
}

func (s *SmartContract) Override(ctx contractapi.TransactionContextInterface, guid string, username string, condition string, loans string, grants string, notes string) error {
	checkExistence, err := ctx.GetStub().GetState(guid)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkExistence == nil {
		return fmt.Errorf("%s does not exists", guid)
	}

	checkUser, err := ctx.GetStub().GetState(username)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkUser == nil {
		return fmt.Errorf("%s does not exist", username)
	}

	specimen := new(Specimen)
	_ = json.Unmarshal(checkExistence, specimen)

	collectionBytes, _ := ctx.GetStub().GetState(specimen.Collection)
	collect := new(Collection)
	_ = json.Unmarshal(collectionBytes, collect)

	user := new(User)
	_ = json.Unmarshal(checkUser, user)

	role, ok := user.Membership[specimen.Collection]

	if !ok {
		role = "P"
	}

	if !strings.Contains(collect.PrimaryUpdate, role) {
		return fmt.Errorf("%s has role %s but role %s is required to update and override primary info", username, role, collect.PrimaryUpdate)
	}

	if condition != "" {
		specimen.Condition = condition + "\n"
	}
	if loans != "" {
		specimen.Loans = loans + "\n"
	}
	if grants != "" {
		specimen.Grants = grants + "\n"
	}
	if notes != "" {
		specimen.Notes = notes
	}

	if condition == "None" {
		specimen.Condition = ""
	}
	if loans == "None" {
		specimen.Loans = ""
	}
	if grants == "None" {
		specimen.Grants = ""
	}
	if notes == "None" {
		specimen.Notes = ""
	}

	attributionString := fmt.Sprintf("Overrode condition, loan, grant, and/or notes history for specimen with guid %s", guid)
	attributionBytes := []byte(attributionString)
	err = ctx.GetStub().PutState(username+"|attribution", attributionBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	specimenBytes, _ := json.Marshal(specimen)

	return ctx.GetStub().PutState(guid, specimenBytes)
}

func (s *SmartContract) RegisterLoan(ctx contractapi.TransactionContextInterface, guid string, username string, description string, loanee string, date string) error {
	checkExistence, err := ctx.GetStub().GetState(guid)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkExistence == nil {
		return fmt.Errorf("%s does not exists", guid)
	}

	checkUser, err := ctx.GetStub().GetState(username)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkUser == nil {
		return fmt.Errorf("%s does not exist", username)
	}

	specimen := new(Specimen)
	_ = json.Unmarshal(checkExistence, specimen)

	collectionBytes, _ := ctx.GetStub().GetState(specimen.Collection)
	collect := new(Collection)
	_ = json.Unmarshal(collectionBytes, collect)

	user := new(User)
	_ = json.Unmarshal(checkUser, user)

	role, ok := user.Membership[specimen.Collection]

	if !ok {
		role = "P"
	}

	if !strings.Contains(collect.RegisterLoan, role) {
		return fmt.Errorf("%s has role %s but role %s is required to register loans", username, role, collect.RegisterLoan)
	}

	specimen.Loans = specimen.Loans + "Loaned: " + description + " to " + loanee + " on " + date + "\n"

	attributionString := fmt.Sprintf("Registered loan for specimen with GUID %s", guid)
	attributionBytes := []byte(attributionString)
	err = ctx.GetStub().PutState(username+"|attribution", attributionBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	specimenBytes, _ := json.Marshal(specimen)

	return ctx.GetStub().PutState(guid, specimenBytes)
}

func (s *SmartContract) ReturnLoan(ctx contractapi.TransactionContextInterface, guid string, username string, description string, loanee string, date string) error {
	checkExistence, err := ctx.GetStub().GetState(guid)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkExistence == nil {
		return fmt.Errorf("%s does not exists", guid)
	}

	checkUser, err := ctx.GetStub().GetState(username)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkUser == nil {
		return fmt.Errorf("%s does not exist", username)
	}

	specimen := new(Specimen)
	_ = json.Unmarshal(checkExistence, specimen)

	collectionBytes, _ := ctx.GetStub().GetState(specimen.Collection)
	collect := new(Collection)
	_ = json.Unmarshal(collectionBytes, collect)

	user := new(User)
	_ = json.Unmarshal(checkUser, user)

	role, ok := user.Membership[specimen.Collection]

	if !ok {
		role = "P"
	}

	if !strings.Contains(collect.RegisterLoan, role) {
		return fmt.Errorf("%s has role %s but role %s is required to register loans", username, role, collect.RegisterLoan)
	}

	specimen.Loans = specimen.Loans + "Returned: " + description + " to " + loanee + " on " + date + "\n"

	specimenBytes, _ := json.Marshal(specimen)

	attributionString := fmt.Sprintf("Returned loan for specimen with GUID %s", guid)
	attributionBytes := []byte(attributionString)
	err = ctx.GetStub().PutState(username+"|attribution", attributionBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	return ctx.GetStub().PutState(guid, specimenBytes)
}

func (s *SmartContract) RegisterGrant(ctx contractapi.TransactionContextInterface, guid string, username string, description string, grantee string, date string) error {
	checkExistence, err := ctx.GetStub().GetState(guid)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkExistence == nil {
		return fmt.Errorf("%s does not exists", guid)
	}

	checkUser, err := ctx.GetStub().GetState(username)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkUser == nil {
		return fmt.Errorf("%s does not exist", username)
	}

	specimen := new(Specimen)
	_ = json.Unmarshal(checkExistence, specimen)

	collectionBytes, _ := ctx.GetStub().GetState(specimen.Collection)
	collect := new(Collection)
	_ = json.Unmarshal(collectionBytes, collect)

	user := new(User)
	_ = json.Unmarshal(checkUser, user)

	role, ok := user.Membership[specimen.Collection]

	if !ok {
		role = "P"
	}

	if !strings.Contains(collect.RegisterUse, role) {
		return fmt.Errorf("%s has role %s but role %s is required to register usage grants", username, role, collect.RegisterUse)
	}

	specimen.Grants = specimen.Grants + "Granted: " + description + " to " + grantee + " on " + date + "\n"

	attributionString := fmt.Sprintf("Registered grant for specimen with GUID %s", guid)
	attributionBytes := []byte(attributionString)
	err = ctx.GetStub().PutState(username+"|attribution", attributionBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	specimenBytes, _ := json.Marshal(specimen)

	return ctx.GetStub().PutState(guid, specimenBytes)
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, guid string, username string) (*Specimen, error) {
	specimenBytes, err := ctx.GetStub().GetState(guid)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if specimenBytes == nil {
		return nil, fmt.Errorf("%s does not exist", guid)
	}

	checkUser, err := ctx.GetStub().GetState(username)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkUser == nil {
		return nil, fmt.Errorf("%s does not exist", username)
	}

	specimen := new(Specimen)
	_ = json.Unmarshal(specimenBytes, specimen)

	collectionBytes, err := ctx.GetStub().GetState(specimen.Collection)

	collection := new(Collection)
	_ = json.Unmarshal(collectionBytes, collection)

	user := new(User)
	_ = json.Unmarshal(checkUser, user)

	role, ok := user.Membership[specimen.Collection]

	if !ok {
		role = "P"
	}

	if !strings.Contains(collection.Query, role) {
		return nil, fmt.Errorf("%s has role %s but role %s is required to query specimens", username, role, collection.Query)
	}

	return specimen, nil
}

func (s *SmartContract) GetHistory(ctx contractapi.TransactionContextInterface, guid string) (string, error) {
	recordIterator, err := ctx.GetStub().GetHistoryForKey(guid)

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

func (s *SmartContract) QueryAllSpecimens(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	recordIterator, err := ctx.GetStub().GetStateByRange("0", "999999999999")

	if err != nil {
		return nil, fmt.Errorf("Failed to get record iterator. %s", err.Error())
	}

	defer recordIterator.Close()

	results := []QueryResult{}

	for recordIterator.HasNext() {
		response, err := recordIterator.Next()

		if err != nil {
			return nil, fmt.Errorf("Error. %s", err.Error())
		}

		specimen := new(Specimen)

		err = json.Unmarshal(response.Value, specimen)
		if err == nil {
			result := QueryResult{response.Key, specimen}
			results = append(results, result)
		}

	}

	return results, nil
}

func (s *SmartContract) UpdateTaxonClass(ctx contractapi.TransactionContextInterface, collection string, username string, oldTaxon string, newTaxon string) (int, error) {
	checkUser, err := ctx.GetStub().GetState(username)

	if err != nil {
		return 0, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if checkUser == nil {
		return 0, fmt.Errorf("%s does not exist", username)
	}

	checkCollection, err := ctx.GetStub().GetState(collection)

	if err != nil {
		return 0, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if checkCollection == nil {
		return 0, fmt.Errorf("%s does not exists", collection)
	}

	collect := new(Collection)
	_ = json.Unmarshal(checkCollection, collect)

	user := new(User)
	_ = json.Unmarshal(checkUser, user)

	role, ok := user.Membership[collection]

	if !ok {
		role = "P"
	}

	if !strings.Contains(collect.TaxonClass, role) {
		return 0, fmt.Errorf("%s has role %s but role %s is required to update taxon class", username, role, collect.TaxonClass)
	}

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return 0, fmt.Errorf("Failed to get results Iterator for all specimens. %s", err.Error())
	}
	defer resultsIterator.Close()

	specimensChanged := 0

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return 0, fmt.Errorf("Failed to get specimen. %s", err.Error())
		}

		specimen := new(Specimen)
		_ = json.Unmarshal(queryResponse.Value, specimen)

		if specimen.Collection == collection && specimen.Taxon == oldTaxon {
			specimen.Taxon = newTaxon

			specimenBytes, _ := json.Marshal(specimen)

			err = ctx.GetStub().PutState(queryResponse.Key, specimenBytes)

			if err != nil {
				return 0, fmt.Errorf("Failed to put to world state. %s", err.Error())
			}

			specimensChanged += 1

		}
	}

	attributionString := fmt.Sprintf("Updated all %s taxons to %s in collection %s", oldTaxon, newTaxon, collection)
	attributionBytes := []byte(attributionString)
	err = ctx.GetStub().PutState(username+"|attribution", attributionBytes)

	if err != nil {
		return 0, fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	return specimensChanged, nil

}

func (s *SmartContract) CouchQuery(ctx contractapi.TransactionContextInterface, queryString string) ([]Specimen, error) {
	recordIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("Failed to get record iterator from query string. %s", err.Error())
	}
	defer recordIterator.Close()

	results := []Specimen{}

	for recordIterator.HasNext() {
		record, err := recordIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Failed to get record from record iterator. %s", err.Error())
		}

		specimen := new(Specimen)

		err = json.Unmarshal(record.Value, specimen)

		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal specimen. %s", err.Error())
		}

		results = append(results, *specimen)

	}

	return results, nil
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
