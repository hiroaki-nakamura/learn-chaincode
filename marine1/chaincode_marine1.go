/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("state", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "send" {
    	var oldState, event, newState string
        var err error

    	oldStateAsBytes, err := stub.GetState("state")
      	if err != nil {
  		  return nil, err
        }
        n := bytes.IndexByte(oldStateAsBytes, 0)
        oldState = string(oldStateAsBytes[:n])

        event = args[0]
        newState = transition(oldState, event)
        	fmt.Println("send " + oldState + "->(" + event + ")->" + newState)
        err = stub.PutState("state", []byte(newState))
      	if err != nil {
  		  return nil, err
        }
 	    return nil, nil
    }
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
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

// read - query function to read key/value pair
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

func transition(oldState string, event string) string {
var newState string
if oldState=="initial" {
  if event=="Imp_PO_Template_Open" {newState="Imp_PO_Draft"}
}
if oldState=="Imp_PO_Draft" {
  if event=="Imp_PO_Draft_Submit" {newState="Exp_PO_Draft"}
  if event=="Imp_PO_Agree" {newState="Exp_PO_Agreed"}
}
if oldState=="Exp_PO_Draft" {
  if event=="Exp_PO_EXW_Agree" {newState="Imp_PO_EXWFinalized"}
  if event=="Exp_PO_CIF_Agree" {newState="Imp_PO_CIFFinalized"}
  if event=="Exp_PO_CFR_Agree" {newState="Imp_PO_CFRFinalized"}
  if event=="Exp_PO_FOB_Agree" {newState="Imp_PO_FOBFinalized"}
  if event=="Exp_PO_DDP_Agree" {newState="Imp_PO_DDPFinalized"}
  if event=="Exp_PO_Draft_Amend" {newState="Imp_PO_Draft"}
}
if oldState=="Exp_InsApp_Draft" {
  if event=="Exp_InsApp_Exp_Submit" {newState="InsCo_InsApp_Draft"}
}
if oldState=="InsCo_InsApp_Draft" {
  if event=="InsCo_InsApp_Exp_Approve" {newState="Exp_InsApp_OpenCover"}
  if event=="InsCo_InsApp_Imp_Approve" {newState="Imp_InsApp_OpenCover"}
  if event=="InsCo_InsApp_Exp_Deny" {newState="Exp_InsApp_Draft"}
  if event=="InsCo_InsApp_Imp_Deny" {newState="Imp_InsApp_Draft"}
}
if oldState=="Exp_InsApp_OpenCover" {
  if event=="Exp_ShipInfo_Exp_Received" {newState="Exp_ShipInfo_Draft"}
}
if oldState=="Exp_ShipInfo_Draft" {
  if event=="Exp_ShipInfo_Draft_Submit" {newState="InsCo_ShipInfo_Received"}
}
if oldState=="InsCo_ShipInfo_Received" {
  if event=="InsCo_ShipInfo_Exp_Approve" {newState="Exp_InsApp_Finalized"}
  if event=="InsCo_ShipInfo_Imp_Approve" {newState="Imp_InsApp_Finalized"}
  if event=="InsCo_ShipInfo_Imp_Deny" {newState="Imp_ShipInfo_Draft"}
  if event=="InsCo_ShipInfo_Exp_Deny" {newState="Exp_ShipInfo_Draft"}
}
if oldState=="Exp_InsApp_Finalized" {
  if event=="Exp_InsPremFee_Exp_Submit" {newState="InsCo_InsPremFee_Received"}
}
if oldState=="InsCo_InsPremFee_Received" {
  if event=="InsCo_InsPremFee_Exp_Approve" {newState="Exp_InsPremFee_Finalized"}
  if event=="InsCo_InsPremFee_Imp_Approve" {newState="Imp_InsPremFee_Finalized"}
  if event=="InsCo_InsPremFee_Imp_Deny" {newState="Imp_InsApp_Finalized"}
  if event=="InsCo_InsPremFee_Exp_Deny" {newState="Exp_InsApp_Finalized"}
}
if oldState=="Exp_IncidentInfo_Received" {
  if event=="Exp_ClaimDoc_Template_Open" {newState="Exp_ClaimDoc_Draft"}
}
if oldState=="Exp_ClaimDoc_Draft" {
  if event=="Exp_ClaimDoc_Draft_Submit" {newState="InsCo_ClaimDoc_Received"}
}
if oldState=="Exp_InsPremFee_Finalized" {
  if event=="Exp_Product_Ship_Request" {newState="Shipper_Product_PlantReady"}
}
if oldState=="InsCo_ClaimDoc_Received" {
  if event=="InsCo_ClaimDoc_Surveyor_Send" {newState="Surveyor_Survey_Received"}
}
if oldState=="Shipper_Product_ExpPortReceived" {
  if event=="Shipper_Product_ExpPortExp_Incident" {newState="Exp_IncidentInfo_Received"}
  if event=="Shipper_Product_Sea_Ship" {newState="Shipper_Product_SeaTransport"}
  if event=="Shipper_Product_ExpPortImp_Incident" {newState="Imp_IncidentInfo_Received"}
}
if oldState=="Shipper_Product_PlantReady" {
  if event=="Shipper_Product_Plant_Ship" {newState="Shipper_Product_ExpLandTransport"}
}
if oldState=="Shipper_Product_ImpPortReceived" {
  if event=="Shipper_Product_ImpPortExp_Incident" {newState="Exp_IncidentInfo_Received"}
  if event=="Shipper_Product_ImpLand_Ship" {newState="Shipper_Product_ImpLandTransport"}
  if event=="Shipper_Product_ImpPortImp_Incident" {newState="Imp_IncidentInfo_Received"}
}
if oldState=="Importer_Product_Destination" {
  if event=="Shipper_Product_DestinationExp_Incident" {newState="Exp_IncidentInfo_Received"}
  if event=="Imp_Product_Desitination_OK" {newState="Imp_Product_DestinationApproved"}
  if event=="Shipper_Product_DestinationImp_Incident" {newState="Imp_IncidentInfo_Received"}
}
if oldState=="Surveyor_Survey_Received" {
  if event=="Surveyor_ClaimDoc_Survery_Deny" {newState="InsCo_Survey_Denied"}
  if event=="Sureveyor_ClaimDoc_Survey_Approve" {newState="InsCo_Survery_Approved"}
}
if oldState=="InsCo_Survey_Denied" {
  if event=="InsCo_ClaimDoc_Exp_Deny" {newState="Exp_ClaimDoc_Denied"}
  if event=="InsCo_ClaimDoc_Exp_Deny" {newState="Imp_ClaimDoc_Denied"}
}
if oldState=="InsCo_Survery_Approved" {
  if event=="InsCo_ClaimDoc_Exp_Approve" {newState="Exp_ClaimDoc_Approved"}
  if event=="InsCo_ClaimDoc_Imp_Approve" {newState="Imp_ClaimDoc_Approved"}
}
if oldState=="Exp_ClaimDoc_Approved" {
  if event=="Exp_Claim_Imp_Inform" {newState="Final_ClaimDoc_Approved"}
}
if oldState=="Exp_ClaimDoc_Denied" {
  if event=="Exp_ClaimDoc_Deny_Send" {newState="Final_ClaimDoc_Denied"}
}
if oldState=="Imp_Product_DestinationApproved" {
  if event=="Imp_Product_Exp_Approve" {newState="Final_Product_Approved"}
}
return newState
}
