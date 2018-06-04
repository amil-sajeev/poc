package main

import (
  "bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {

}

type Patient struct {
	Id   string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_lame"`
}

type Physician struct {
	Id   string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}

type PatientRecord struct {
	Id   string `json:"id"`
	Owner string `json:"owner"`
	CurrentPhysician   string `json:"current_physician"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "getAllPatient" {
		return s.getAllPatient(APIstub)
	} else if function == "getPatient" {
		return s.getPatient(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} 
	/* else if function == "createPatient" {
		return s.createPatient(APIstub, args)
	} else if function == "createPhysician" {
		return s.createPhysician(APIstub, args)
	}  else if function == "getAllRecords" {
		return s.getAllRecords(APIstub)
	} else if function == "getRecord" {
		return s.getRecord(APIstub, args)
	} else if function == "changeCurrentPhysician" {
		return s.changeCurrentPhysician(APIstub, args)
	}*/

	return shim.Error("Invalid Smart Contract function name.")
}


func (s *SmartContract) getAllPatient(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "P1"
	endKey := "P99"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


func (s *SmartContract) getPatient(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	patientAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(patientAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	patients := []Patient{
		Patient{FirstName: "Rincy", LastName: "Yohannan"},
		Patient{FirstName: "Amil", LastName: "Sajeev"},
	}

	i := 1
	for i < len(patients) {
		fmt.Println("i is ", i)
		patientAsBytes, _ := json.Marshal(patients[i])
		APIstub.PutState("P"+strconv.Itoa(i), patientAsBytes)
		fmt.Println("Added", patients[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
