package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"strings"
)


//./byfn.sh -c mychannel -s couchdb
//进入cli容器进行安装
//docker exec -it cli bash
//peer chaincode install -n couchcar -v 1.0 -p github.com/chaincode/couchcar/go
//export CHANNEL_NAME=mychannel
//peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
//检测是否安装成功(如果输出索引信息即为成功)
//docker logs peer0.org1.example.com  2>&1 | grep "CouchDB index"
//在peer0org2上安装
//export CORE_PEER_LOCALMSPID="Org2MSP" 
//export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
//export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
//export CORE_PEER_ADDRESS=peer0.org2.example.com:9051


//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -c '{"Args":["createCar","car0","奇瑞","suv","dark","200000","2015","songjian"]}'
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -c '{"Args":["createCar","car1","奔驰","高端车身","blue","400000","2016","ccc"]}'
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -c '{"Args":["createCar","car2","凯美瑞","中端车身","green","150000","2017","songjian"]}'
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -c '{"Args":["createCar","car3","红旗","超高端车身","white","1000000","2016","ccc"]}'
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -c '{"Args":["createCar","car4","奇瑞","高端车身","dark","300000","2016","songjian"]}'
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -c '{"Args":["createCar","car5","奇瑞","豪华SUV","white","400000","2017","ccc"]}'
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -c '{"Args":["createCar","car6","桑塔纳","高端车身","red","400000","2018","jerry"]}'
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -c '{"Args":["createCar","car7","东风标致","高端车身","gray","400000","2019","jerry"]}'

//基于couchDB使用不同的索引进行的富查询
//查询所有者为songjian的汽车													{"selector":{"docType":"couch_car","owner":"songjian"},"use_index":["_design/indexOwnerDoc","indexOwner"]}
//peer chaincode query -C $CHANNEL_NAME -n couchcar -c '{"Args":["queryCars", "{\"selector\":{\"docType\":\"couch_car\",\"owner\":\"songjian\"}, \"use_index\":[\"_design/indexOwnerDoc\", \"indexOwner\"]}"]}'
//查询颜色为蓝色的汽车 
//peer chaincode query -C $CHANNEL_NAME -n couchcar -c '{"Args":["queryCars", "{\"selector\":{\"docType\":\"couch_car\",\"color\":\"blue\"}, \"use_index\":[\"_design/indexColoeDoc\", \"indexColor\"]}"]}'
// 查询所有的汽车  
//peer chaincode query -C $CHANNEL_NAME -n couchcar -c '{"Args":["queryCars", "{\"selector\":{\"docType\":\"couch_car\"}, \"use_index\":[\"_design/indexColoeDoc\", \"indexColor\"]}"]}'

//改变汽车的拥有者
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -c '{"Args":["changeCarOwner","car2","jerry"]}'

//单个键(汽车编号)查询
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -c '{"Args":["readCar","car2"]}'

//基于键的范围查询
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -c '{"Args":["getCarsByRange","car1","car4"]}'




//链码的升级
//peer chaincode install -n couchcar -v 2.0 -p github.com/chaincode/couchcar/go
//peer chaincode upgrade -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n couchcar -v 2.0 -c '{"Args":["init"]}' -P "OR ('Org0MSP.peer','Org1MSP.peer')"




type CarContract struct {

}

type Car struct{
	ObjectType string `json:"docType"`	//namespace
	CarNum string `json:"car_num"`		//汽车编号(唯一的)
	Make string `json:"make"`			//制造厂商
	Model string `json:"model"`			//模型
	Color string `json:"color"`			//颜色
	Price int `json:"price"`			//价格
	Year int `json:"year"`				//生产年份
	Owner string `json:"owner"`			//拥有者
}


func (c *CarContract)Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


func (c *CarContract)Invoke(stub shim.ChaincodeStubInterface) pb.Response  {
	function,args := stub.GetFunctionAndParameters()
	if function == "queryCar"{
		return c.queryCar(stub,args)
	}else if function == "createCar"{
		return c.createCar(stub,args)
	}else if function == "changeCarOwner"{
		return c.changeCarOwner(stub,args)
	}else if function == "queryCars"{
		return c.queryCars(stub,args)
	}else if function == "getCarsByRange"{
		return c.getCarsByRange(stub,args)
	}else if  function == "readCar"{
		return c.readCar(stub,args)
	}

	fmt.Println("没有找到这个的函数方法"+ function)
	return shim.Error("未知的方法名")
}

//查询汽车（通过car_num汽车编号查找）
func (c *CarContract)queryCar(stub shim.ChaincodeStubInterface,args []string) pb.Response  {
	if len(args) != 1{
		return shim.Error("参数错误")
	}
	carInfo, _ := stub.GetState(args[0])
	//如果通过编号查询到没有汽车，就报错
	if carInfo == nil{
		errresp := "{\"Error\":\"汽车编号：" + args[0] + "的汽车信息不存在，请输入正确的汽车编号\"}"
		return shim.Error(errresp)
	}
	return shim.Success(carInfo)
}


//创建汽车

func (c *CarContract)createCar(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 7{
		return  shim.Error("参数错误")
	}
	if len(args[0]) <= 0{
		return  shim.Error("参数不能为空")
	}
	if len(args[1]) <= 0{
		return shim.Error("参数不能为空")
	}
	if len(args[2]) <= 0 {
		return shim.Error("参数不能为空")
	}
	if len(args[3]) <= 0{
		return shim.Error("参数不能为空")
	}

	car_num := args[0]
	make :=  strings.ToLower(args[1])
	model := strings.ToLower(args[2])
	color := strings.ToLower(args[3])
	price, err := strconv.Atoi(args[4])
	if err != nil{
		return shim.Error("第5个参数必须是字符串")
	}
	year, _ := strconv.Atoi(args[5])
	if err != nil{
		return shim.Error("第6个参数必须是字符串")
	}
	owner := strings.ToLower(args[6])

	//在创建新的汽车之前，需要先从状态数据库中查看是否已经存在了相应汽车编号的汽车
	oldcarInfo, _ := stub.GetState(car_num)
	if oldcarInfo != nil{
		return shim.Error("该汽车编号的启辰已经存在，请加入不同编号的汽车")
	}

	fmt.Println("开始增加汽车...........")

	//docType是固定的，为couch_car，标识不同的状态数据
	objectName := "couch_car"
	car := &Car{objectName,car_num,make,model,color,price,year,owner}
	carInfo,err := json.Marshal(car)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(car_num, carInfo)
	if err != nil{
		return shim.Error(err.Error())
	}

	fmt.Println("...........增加汽车成功")

	return shim.Success(nil)

}



//改变汽车的拥有者changeOwner

func (c *CarContract)changeCarOwner(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 2{
		return shim.Error("参数错误，请输入两个参数：汽车编号，新的所有者")
	}

	newowner := strings.ToLower(args[1])

	oldcarinfo, _ := stub.GetState(args[0])
	if oldcarinfo == nil {
		errresp := "{\"Error\":\"汽车编号：" + args[0] + "的汽车信息不存在，请输入正确的汽车编号\"}"
		return shim.Error(errresp)
	}

	newcarInfo := Car{}
	err := json.Unmarshal(oldcarinfo, &newcarInfo)
	if err != nil{
		return shim.Error(err.Error())
	}

	fmt.Println("开始更新指定汽车编号的拥有者")
	//将原来的汽车的owner字段该为新的owner
	newcarInfo.Owner = newowner

	carinfo, _ := json.Marshal(newcarInfo)
	err = stub.PutState(args[0], carinfo)
	if err != nil{
		return shim.Error(err.Error())
	}

	fmt.Println("更新拥有者成功.........")
	return  shim.Success(nil)
}



//处理迭代器，返回结果数据
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
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

	return &buffer, nil
}



//处理querystring查询语句，生成迭代器对象，调用constructQueryResponseFromIterator处理迭代器对象，返回查询结果数据
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}


//查询所有的汽车(富查询)queryAllCars
//自定义所有的查询
func (c *CarContract)queryCars(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	//只需要传入一个参数：就是couchDB查询的querystring
	if len(args) != 1{
		return shim.Error("参数错误，请输入selector查询语句")
	}

	queryStrings := args[0]

	queryResults,err := getQueryResultForQueryString(stub,queryStrings)

	if err != nil{
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}


//基于键的范围查询
func (c *CarContract)getCarsByRange(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 2{
		return  shim.Error("参数错误，请输入两个参数：startkey,endkey")
	}

	resultIterator, err := stub.GetStateByRange(args[0], args[1])
	if err != nil{
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	result, err := constructQueryResponseFromIterator(resultIterator)
	if err != nil {
		return shim.Error(err.Error())
	}

	return  shim.Success(result.Bytes())
}

//基于单个键的查询
func (c *CarContract)readCar(stub shim.ChaincodeStubInterface,args []string) pb.Response  {
	if len(args) != 1{
		return shim.Error("参数错误，请输入一个参数：汽车编号")
	}

	result, _ := stub.GetState(args[0])
	if result == nil {
		return shim.Error("该编号的汽车信息不存在，请重新输入汽车编号")
	}

	return shim.Success(result)
}



// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(CarContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}















