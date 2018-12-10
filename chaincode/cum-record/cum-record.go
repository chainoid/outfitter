// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
*/

package main

/* Imports
* 5 utility libraries for handling bytes, reading and writing JSON,
formatting, and string manipulation
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define Enum for classificate record type
type RecordType string

const (
	GroupType   RecordType = "G"
	StudentType RecordType = "S"
	TeacherType RecordType = "T" // For future use
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define UserRecord structure, with several properties.
Structure tags are used by encoding/json library
*/
type UserRecord struct {
	RecordType  string     `json:"recordType"`
	GroupName   string     `json:"groupName"`
	UserId      string     `json:"userId"`
	UserName    string     `json:"userName"`
	Description string     `json:"description"`
	RegisterTS  string     `json:"registerTS"`
	RecordList  []UserItem `json:"recordList"`
}

/* Define User structure, with several properties.
Structure tags are used by encoding/json library
*/
type User struct {
	StudentId   string `json:"studentId"`
	StudentName string `json:"studentName"`
	GroupName   string `json:"groupName"`
}

/* Define User Test structure, with several properties
Structure tags are used by encoding/json library
*/
type UserItem struct {
	StestId     string `json:"testId"`
	Group       string `json:"group"`
	Course      string `json:"course"`
	Teacher     string `json:"teacher"`
	AssignedTS  string `json:"assignedTS"`
	Rate        string `json:"rate"`
	ExecuteTS   string `json:"executeTS"`
	ExecuteDesc string `json:"executeDesc"`
}

/*
 * The Init method *
 called when the Smart Contract "elza-rec" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function
 -- see initLedger()
*/
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "posta-chaincode"
 The app also specifies the specific smart contract function to call with args
*/
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryAllGroups" {
		return s.queryAllGroups(APIstub)
	} else if function == "addGroup" {
		return s.addGroup(APIstub, args)
	} else if function == "queryAllUsers" {
		return s.queryAllUsers(APIstub)
	} else if function == "generateSetForGroup" {
		return s.generateSetForGroup(APIstub, args)
	} else if function == "prepareForDelivery" {
		return s.prepareForDelivery(APIstub, args)
	} else if function == "deliveryItem" {
		return s.deliveryItem(APIstub, args)
	} else if function == "addUser" {
		return s.addUser(APIstub, args)
	} else if function == "getUserRecord" {
		return s.getUserRecord(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The initLedger method
 * Will add group and student records to our network
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	stests := []UserRecord{
		UserRecord{RecordType: "G", GroupName: "AB17", Description: "Desc AB17"},
		UserRecord{RecordType: "G", GroupName: "AB18", Description: "Desc AB18"},
		UserRecord{RecordType: "G", GroupName: "AB19", Description: "Desc AB19"},
		UserRecord{RecordType: "G", GroupName: "AB20", Description: "Desc AB20"},

		UserRecord{RecordType: "S", GroupName: "AB17", UserId: "AB1701", UserName: "Fighter 1701", RegisterTS: time.Now().Format(time.RFC3339), Description: "Desc 1701"},
		UserRecord{RecordType: "S", GroupName: "AB17", UserId: "AB1702", UserName: "Fighter 1702", RegisterTS: time.Now().Format(time.RFC3339), Description: "Desc 1702"},
		UserRecord{RecordType: "S", GroupName: "AB17", UserId: "AB1703", UserName: "Fighter 1703", RegisterTS: time.Now().Format(time.RFC3339), Description: "Desc 1703"},
		UserRecord{RecordType: "S", GroupName: "AB17", UserId: "AB1704", UserName: "Fighter 1704", RegisterTS: time.Now().Format(time.RFC3339), Description: "Desc 1704"},
	}

	i := 0
	for i < len(stests) {
		fmt.Println("i is ", i)
		stestAsBytes, _ := json.Marshal(stests[i])
		APIstub.PutState(fmt.Sprintf("%016d", rand.Int63n(1e16)), stestAsBytes)
		fmt.Println("Added", stests[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The queryAllGroups method *
allows for assessing all group records added to the ledger(all groups in the cumulative system)
This method does not take any arguments. Returns JSON string containing results.
*/
func (s *SmartContract) queryAllGroups(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "9999"

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

		// Create an object
		userRecord := UserRecord{}
		// Unmarshal record to stest object
		json.Unmarshal(queryResponse.Value, &userRecord)

		// Add only filtered by RecordType as Group records
		if userRecord.RecordType == "G" {

			// Add comma before array members,suppress it for the first array member
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
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllGroups:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The addGroup method
 * This method takes in four arguments (attributes to be saved in the ledger).
 */
func (s *SmartContract) addGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var groupRecord = UserRecord{RecordType: "G", GroupName: args[0], Description: args[1]}

	groupRecordAsBytes, _ := json.Marshal(groupRecord)
	err := APIstub.PutState(fmt.Sprintf("%016d", rand.Int63n(1e16)), groupRecordAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record new group: %s", args[0]))
	}

	return shim.Success(nil)
}

/* The addUser method *
   Generate initial user record
   This method takes in four arguments (attributes to be saved in the ledger).
*/
func (s *SmartContract) addUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Fill the user date
	var userRecord = UserRecord{RecordType: "S", UserId: args[0], UserName: args[1], GroupName: args[2], Description: args[3], RegisterTS: time.Now().Format(time.RFC3339)}

	userRecordAsBytes, _ := json.Marshal(userRecord)

	err := APIstub.PutState(fmt.Sprintf("%016d", rand.Int63n(1e16)), userRecordAsBytes)

	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to add user to system with id: %s", args[0]))
	}

	fmt.Println("Added user with id: ", args[0])

	return shim.Success(nil)
}

/*
 * The queryAllUsers method *
   allows for assessing all user records added to the ledger
   This method does not take any arguments. Returns JSON string containing results.
*/
func (s *SmartContract) queryAllUsers(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "9999"

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

		// Create an object
		userRecord := UserRecord{}
		// Unmarshal record to stest object
		json.Unmarshal(queryResponse.Value, &userRecord)

		// Add only filtered by RecordType as Group records
		if userRecord.RecordType == "S" {

			// Add comma before array members,suppress it for the first array member
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
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllUsers:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The queryTestById method *
 Used to view the records of one particular parsel
 It takes one argument -- the key for the parsel in question
*/

/*
func (s *SmartContract) queryTestById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	stestAsBytes, _ := APIstub.GetState(args[0])
	if stestAsBytes == nil {
		return shim.Error("Could not locate test record")
	}
	return shim.Success(stestAsBytes)
}
|*/

/* The generateSetForGroup method *
   Generate list of set for one group
   This method takes in four arguments (attributes to be saved in the ledger).
*/
func (s *SmartContract) generateSetForGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	startKey := "0"
	endKey := "9999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Create an object
		userRecord := UserRecord{}
		// Unmarshal record to stest object
		json.Unmarshal(queryResponse.Value, &userRecord)

		if userRecord.GroupName == args[0] {
			var studentTest = UserItem{StestId: fmt.Sprintf("%X", rand.Int()), Group: args[0], Course: args[1], Teacher: args[2],
				AssignedTS: time.Now().Format(time.RFC1123Z), ExecuteDesc: ""}

			userRecord.RecordList = append(userRecord.RecordList, studentTest)

			userRecordAsBytes, _ := json.Marshal(userRecord)

			APIstub.PutState(queryResponse.Key, userRecordAsBytes)

			fmt.Println("Added", studentTest)
		}

	}

	return shim.Success(nil)
}

/*
 * The getUserRecord method *
   allows for assessing all the records from selected student

    Returns JSON string containing results.
*/

func (s *SmartContract) getUserRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	studentRecordAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Could not locate user data")
	}

	return shim.Success(studentRecordAsBytes)
}

/*
 * The prepareForDelivery method *
   allows for assessing all the records from selected group/item

	Returns JSON string containing results.
*/

func (s *SmartContract) prepareForDelivery(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	startKey := "0"
	endKey := "9999"

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	fmt.Printf("- prepareForDelivery param args[0] :%s  args[1] :%s\n", args[0], args[1])

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing Results
	var buffer bytes.Buffer

	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	// Iteration my the Student Record List
	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Create an object
		userRecord := UserRecord{}
		// Unmarshal record to stest object
		json.Unmarshal(queryResponse.Value, &userRecord)

		// Selection by the Record type and Group name
		if userRecord.RecordType == "S" && userRecord.GroupName == args[0] {

			// Iteration by the Course

			for i := 0; i < len(userRecord.RecordList); i++ {

				if userRecord.RecordList[i].Course == args[1] {

					// Add comma before array members,suppress it for the first array member
					if bArrayMemberAlreadyWritten == true {
						buffer.WriteString(",")
					}

					buffer.WriteString("{\"Key\":")
					buffer.WriteString("\"")
					buffer.WriteString(queryResponse.Key)
					buffer.WriteString("\"")
					buffer.WriteString(", \"Record\":")

					// Put only selected fields
					buffer.WriteString("{\"userId\":\"")
					buffer.WriteString(queryResponse.Key)
					buffer.WriteString("\",")

					buffer.WriteString("\"testId\":\"")
					buffer.WriteString(userRecord.RecordList[i].StestId)
					buffer.WriteString("\",")

					buffer.WriteString("\"userName\":\"")
					buffer.WriteString(userRecord.UserName)
					buffer.WriteString("\",")

					buffer.WriteString("\"group\":\"")
					buffer.WriteString(userRecord.GroupName)
					buffer.WriteString("\",")

					buffer.WriteString("\"course\":\"")
					buffer.WriteString(userRecord.RecordList[i].Course)
					buffer.WriteString("\",")

					buffer.WriteString("\"assignedTS\":\"")
					buffer.WriteString(userRecord.RecordList[i].AssignedTS)
					buffer.WriteString("\",")

					buffer.WriteString("\"teacher\":\"")
					buffer.WriteString(userRecord.RecordList[i].Teacher)
					buffer.WriteString("\",")

					buffer.WriteString("\"executeTS\":\"")
					buffer.WriteString(userRecord.RecordList[i].ExecuteTS)
					buffer.WriteString("\",")

					buffer.WriteString("\"rate\":\"")
					buffer.WriteString(userRecord.RecordList[i].Rate)
					buffer.WriteString("\"")

					buffer.WriteString("}")

					buffer.WriteString("}")
					bArrayMemberAlreadyWritten = true
				}
			}
		}
	}

	buffer.WriteString("]")

	if bArrayMemberAlreadyWritten == false {
		return shim.Error("No group/item found")
	}

	fmt.Printf("- prepareForDelivery:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())

}

/*
 * The deliveryItem method *
 * The data in the state can be updated .
 */
func (s *SmartContract) deliveryItem(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	fmt.Printf("- deliveryItem param args[0] :%s  args[1] :%s args[3] :%s\n", args[0], args[1], args[2])

	userRecordAsBytes, _ := APIstub.GetState(args[0])

	if userRecordAsBytes == nil {
		return shim.Error("Could not locate selected student record")
	}

	userRecord := UserRecord{}

	json.Unmarshal(userRecordAsBytes, &userRecord)

	// Looking for selected test by course name
RecordListLoop:

	for i := 0; i < len(userRecord.RecordList); i++ {

		if userRecord.RecordList[i].Course == args[1] {

			if userRecord.RecordList[i].Rate != "" {
				return shim.Error("Selected item already delivered.")
			}

			userRecord.RecordList[i].Rate = args[2]
			userRecord.RecordList[i].ExecuteTS = time.Now().Format(time.RFC1123Z)

			userRecordAsBytes, _ := json.Marshal(userRecord)

			err := APIstub.PutState(args[0], userRecordAsBytes)

			if err != nil {
				return shim.Error(fmt.Sprintf("Failed to change status of item record %s", args[0]))
			}

			break RecordListLoop
		}
	}

	return shim.Success(userRecordAsBytes)
}

/*
 * main function  - calls the Start function
   The main function starts the chaincode in the container during instantiation.
*/
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
