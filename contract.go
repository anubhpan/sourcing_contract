package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
// ============================================================================================================================
// Asset Definitions - The ledger will store contract and its attributes
// ============================================================================================================================

type Contract struct {
	ObjectType string        `json:"docType"` //field for couchdb
	ContractId string         `json:"id"`      //the contract ID
	LabEqp	LabEquipment   `json:"labeqp"`  // field capturing terms and agreements for supply of Lab Equipments
	LabChem	LabChemicals   `json:"labchem"`    //field capturing terms and agreements for supply of Lab Chemicals
	InfSer	Infraservice   `json:"infraser"`  //field capturing terms and agreements for supply of Infrastructure services
	Contr	Contractor	  `json:"contractor"` // contractor or sub-contrator who is executing the contract at any stage
}

type Contractor struct {
	ContctrID	string 	`json:"ContrId"`
	ContctrName	string 	`json:"ContrName"`   
}


type LabEquipment struct {
	EqId	string `json:"eqid"`  // equipment ID
	EqType	string `json:"eqtyp"`  // equipment type
	EqQuant	string `json:"eqquant"`   // quantity
	SourcedBy	string `json:"eqsrcby"`     // who will source/ main supplier of contract or sub-contractor
}

type LabChemicals struct {
	ChemId	string `json:"cmid"`  // equipment ID
	ChemName	string `json:"cmname"`  // equipment type
	ChemQuant	string `json:"cmquant"`   // quantity
	SourcedBy	string `json:"cmsrcby"`     // who will source/ main supplier of contract or sub-contractor
}

type Infraservice struct {
	SLAattr1	string `json:"sla1"` // sla agreement attributes
	SLAattr2	string `json:"sla2"`   // sla agreement attributes
	ManagedBy	string `json:"imngdby"`   // Managed by (self/procured)   
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// initialize the chaincode
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Contract Is Starting Up")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// this is a very simple dumb test.  let's write to the ledger and error on any errors

	err := stub.PutState("hello_insidetrack", []byte(args[0]))
	if err != nil {
		return nil, err
	}


	fmt.Println(" - ready for action")                          //self-test pass
	return nil, nil
}

// ============================================================================================================================
// Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error)  {

	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}
// ============================================================================================================================
// Query function
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error)  {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}
//==========================
// Write function, invoke function to write key/value pair
//==========================
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}
//==========================
// read, query function to read key/value pair
//=============================
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
