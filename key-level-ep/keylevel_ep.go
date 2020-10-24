/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */


//docker exec -it cli bash
//peer chaincode install -n keylevelep -v 1.0 -p github.com/chaincode/key-level-ep
//export CHANNEL_NAME=mychannel
//我们设置链码的背书策略为org1和org2的成员同时背书，以方便下面进行测试
//peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n keylevelep -v 1.0 -c '{"Args":["init"]}' -P "AND ('Org1MSP.member','Org2MSP.member')"

//export CORE_PEER_LOCALMSPID="Org2MSP" 
//export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
//export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
//export CORE_PEER_ADDRESS=peer0.org2.example.com:9051
//peer chaincode install -n keylevelep -v 1.0 -p github.com/chaincode/key-level-ep


//注意这里交易调用的时候需要同时写peer的grpc通信地址以及它们的TLS根证书，因为我们的背书策略是AND
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n keylevelep --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt  --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["initLedger"]}'
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n keylevelep -c '{"Args":["queryAllCars"]}'
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n keylevelep -c '{"Args":["queryCar","CAR2"]}'


//将CAR2的owner改为songjian,由peer0org1和peer0org2背书,也就是默认的chaincode-level的背书
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n keylevelep --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt  --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["changeCarOwner","CAR2","songjian"]}'
//查看CAR2的owner，可以发现成功改变
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n keylevelep -c '{"Args":["queryCar","CAR2"]}'
//Chaincode invoke successful. result: status:200 payload:"{\"colour\":\"green\",\"make\":\"Hyundai\",\"model\":\"Tucson\",\"owner\":\"songjian\"}"


//再来试试只用org1节点进行背书,并尝试将CAR2的拥有者改为ccc
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n keylevelep --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt  -c '{"Args":["changeCarOwner","CAR2","ccc"]}'
//再来查询CAR2的owner
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n keylevelep -c '{"Args":["queryCar","CAR2"]}'
//可以看到owner仍然为songjian,并没有被改变，所以是背书策略没有满足，交易不合法没有通过校验
//result: status:200 payload:"{\"colour\":\"green\",\"make\":\"Hyundai\",\"model\":\"Tucson\",\"owner\":\"songjian\"}" 


//我们将CAR2的背书策略设置为只需要org1的成员背书（需要注意的是，键级别的背书策略的设置其实也是以交易的形式发起的，所以它本身也是一笔交易，所以调用的时候必须满足AND策略，即两个组织的成员同时背书才行）
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n keylevelep --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt  --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["setKeyLevelEndorsement","CAR2","Org1MSP"]}'
//我们再来改变CAR2的owner为ccc,看看能否成功(仅向org1的成员发起背书)
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n keylevelep --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt  -c '{"Args":["changeCarOwner","CAR2","ccc"]}'
//查看CAR2的owner
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n keylevelep -c '{"Args":["queryCar","CAR2"]}'
//可以看到修改成功
//Chaincode invoke successful. result: status:200 payload:"{\"colour\":\"green\",\"make\":\"Hyundai\",\"model\":\"Tucson\",\"owner\":\"ccc\"}" 
//到此，键级别的背书策略的设置就成功了


 package main

 import (
	 "bytes"
	 "encoding/json"
	 "fmt"
	 "strconv"
	 "github.com/hyperledger/fabric/core/chaincode/shim/ext/statebased"
	 "github.com/hyperledger/fabric/core/chaincode/shim"
	 sc "github.com/hyperledger/fabric/protos/peer"
 )
 
 // Define the Smart Contract structure
 type SmartContract struct {
 }
 
 // Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
 type Car struct {
	 Make   string `json:"make"`
	 Model  string `json:"model"`
	 Colour string `json:"colour"`
	 Owner  string `json:"owner"`
 }
 
 func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	 // Retrieve the requested Smart Contract function and arguments
	 function, args := APIstub.GetFunctionAndParameters()
	 // Route to the appropriate handler function to interact with the ledger appropriately
	 if function == "queryCar" {
		 return s.queryCar(APIstub, args)
	 } else if function == "initLedger" {
		 return s.initLedger(APIstub)
	 } else if function == "createCar" {
		 return s.createCar(APIstub, args)
	 } else if function == "queryAllCars" {
		 return s.queryAllCars(APIstub)
	 } else if function == "changeCarOwner" {
		 return s.changeCarOwner(APIstub, args)
	 } else if function == "setKeyLevelEndorsement"{
		 return s.setKeyLevelEndorsement(APIstub,args)
	 }
 
	 return shim.Error("Invalid Smart Contract function name.")
 }
 
 func (s *SmartContract) queryCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 1 {
		 return shim.Error("Incorrect number of arguments. Expecting 1")
	 }
 
	 carAsBytes, _ := APIstub.GetState(args[0])
	 return shim.Success(carAsBytes)
 }
 
 func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	 cars := []Car{
		 Car{Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko"},
		 Car{Make: "Ford", Model: "Mustang", Colour: "red", Owner: "Brad"},
		 Car{Make: "Hyundai", Model: "Tucson", Colour: "green", Owner: "Jin Soo"},
		 Car{Make: "Volkswagen", Model: "Passat", Colour: "yellow", Owner: "Max"},
		 Car{Make: "Tesla", Model: "S", Colour: "black", Owner: "Adriana"},
		 Car{Make: "Peugeot", Model: "205", Colour: "purple", Owner: "Michel"},
		 Car{Make: "Chery", Model: "S22L", Colour: "white", Owner: "Aarav"},
		 Car{Make: "Fiat", Model: "Punto", Colour: "violet", Owner: "Pari"},
		 Car{Make: "Tata", Model: "Nano", Colour: "indigo", Owner: "Valeria"},
		 Car{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
	 }
 
	 i := 0
	 for i < len(cars) {
		 fmt.Println("i is ", i)
		 carAsBytes, _ := json.Marshal(cars[i])
		 APIstub.PutState("CAR"+strconv.Itoa(i), carAsBytes)
		 fmt.Println("Added", cars[i])
		 i = i + 1
	 }
 
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) createCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 5 {
		 return shim.Error("Incorrect number of arguments. Expecting 5")
	 }
 
	 var car = Car{Make: args[1], Model: args[2], Colour: args[3], Owner: args[4]}
 
	 carAsBytes, _ := json.Marshal(car)
	 APIstub.PutState(args[0], carAsBytes)
 
	 res := "the car has been created,carnum:"+args[0]+",make:"+args[1]+",model:"+args[2]+",color:"+args[3]+",owner:"+args[4]
	 err := APIstub.SetEvent("createCarEvent", []byte(res))
	 if err != nil{
		 return shim.Error(fmt.Sprintf("Failed to emit event"))
	 }
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) queryAllCars(APIstub shim.ChaincodeStubInterface) sc.Response {
	 startKey := "CAR0"
	 endKey := "CAR999"
 
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
 
 func (s *SmartContract) changeCarOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 2 {
		 return shim.Error("Incorrect number of arguments. Expecting 2")
	 }
 
	 carAsBytes, _ := APIstub.GetState(args[0])
	 car := Car{}
 
	 json.Unmarshal(carAsBytes, &car)
	 car.Owner = args[1]
 
	 carAsBytes, _ = json.Marshal(car)
	 APIstub.PutState(args[0], carAsBytes)
 
	 res := "the owner of the car"+ args[0]+"has been changed to "+ args[1]
	 err := APIstub.SetEvent("changeOwnerEvent", []byte(res))
	 if err != nil{
		 return shim.Error(fmt.Sprintf("Failed to emit event"))
	 }
	 return shim.Success(nil)
 }
 
//set key-level endorsement policy
func (s *SmartContract)setKeyLevelEndorsement(API shim.ChaincodeStubInterface,args [] string)sc.Response{
	if len(args) != 2{
		return shim.Error("Incorrect number of arguments.Expecting key and EP(MSP)")
	}
	key := args[0]
	EP := args[1]
	newEP, err := statebased.NewStateEP(nil)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = newEP.AddOrgs(statebased.RoleTypeMember, EP)
	if err != nil {
		return shim.Error(err.Error())
	}

	policyByte, err := newEP.Policy()
	if err != nil {
		return shim.Error(err.Error())
	}
	err = API.SetStateValidationParameter(key, policyByte)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}


 // The main function is only relevant in unit test mode. Only included here for completeness.
 func main() {
 
	 // Create a new Smart Contract
	 err := shim.Start(new(SmartContract))
	 if err != nil {
		 fmt.Printf("Error creating new Smart Contract: %s", err)
	 }
 }
 