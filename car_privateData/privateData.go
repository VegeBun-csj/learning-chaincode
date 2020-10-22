package main
import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)


type SmartContract struct {
}

//比如我们有三个组织，分别是 org1\org2\org3
//org1代表车管所
//org2代表停车场公司
//org3代表公安局

//假设有这样的场景，车管所登记的车型信息需要公安局来审核，停车场的信息需要给公安机关来进行审核，但是车管所和停车场公司之间不想共享信息
//所以org1和org3之间有一个隐私数据集合CarInfo_Collection，org2和org3之间有一个隐私数据集合ParkInfo_Collection
//这是在同一个通道中更细粒度的隐私保护



//docker exec -it cli bash
//peer chaincode install -n car_private -v 1.0 -p github.com/chaincode/car_private
//export CHANNEL_NAME=mychannel
//peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n car_private -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.member','Org2MSP.member','Org3MSP.member')" --collections-config=/opt/gopath/src/github.com/chaincode/car_private/collections_config.json
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n car_private -c '{"Args":["createCar","苏F666v8","奔驰","black","songjian","3421312311231213x"]}'
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n car_private -c '{"Args":["createParkInfo","苏F666v8","2020-10-20","安徽芜湖","2020-10-21","2h"]}'
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n car_private -c '{"Args":["queryCar","苏F666v8"]}'


//在org2peer0上安装
//export CORE_PEER_LOCALMSPID="Org2MSP"
//export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
//export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
//export CORE_PEER_ADDRESS=peer0.org2.example.com:9051
//peer chaincode install -n car_private -v 1.0 -p github.com/chaincode/car_private

//在org3peer0上安装
//dokcer exec -it org3cli bash
//peer chaincode install -n car_private -v 1.0 -p github.com/chaincode/car_private


//下面来检测一下，三个组织用户的访问权限
//首先是对Car信息的访问
//-----------------------------------------------------------------------------------------test collection_CarInfo------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//首先在cli中以org1的身份来查询Car的信息：
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n car_private -c '{"Args":["queryCar","苏F666v8"]}'
//输出结果为：Chaincode invoke successful. result: status:200 payload:"{\"make\":\"\345\245\224\351\251\260\",\"colour\":\"black\",\"owner\":\"songjian\",\"owner_id_car_num\":\"3421312311231213x\"}"

//然后以org2的身份来查询Car的信息
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n car_private -c '{"Args":["queryCar","苏F666v8"]}'
//可以看到出现了背书错误，org2身份不满足隐私数据集合中定义的访问策略,不能访问隐私集合collection_CarInfo
//Error: endorsement failure during invoke. response: status:500 message:"GET_STATE failed: transaction ID: d109eb004ca6d77c32afda9d5f46605817ed10338e0d1bc8fa6930370f9bc751: tx creator does not have read access permission on privatedata in chaincodeName:car_private collectionName: collection_carInfo"


//以org3的身份查询访问：
//docker exec -it Org3cli bash
//export CHANNEL_NAME=mychannel
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n car_private -c '{"Args":["queryCar","苏F666v8"]}'
//Chaincode invoke successful. result: status:200 payload:"{\"make\":\"\345\245\224\351\251\260\",\"colour\":\"black\",\"owner\":\"songjian\",\"owner_id_car_num\":\"3421312311231213x\"}"
//可以看到org1和org3都是可以访问隐私集合collection_carInfo的，Org2是不能访问的，符合我们合约的预期效果


//-----------------------------------------------------------------------------------------test collection_ParkInfo------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//下面来测试另一个集合collection_parkInfo的隐私效果,它应该是只让org2和org3访问，但是org1不能访问

//先来测试org1的身份进行查询
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n car_private -c '{"Args":["queryParkInfo","苏F666v8"]}'
//可以发现出现了访问权限错误：不能访问隐私数据集合collection_parkInfo
//Error: endorsement failure during invoke. response: status:500 message:"GET_STATE failed: transaction ID: 33c5f8f2999f47f1ea2931de98707d8370c04d8a94ac96025f60ff11a6130bdb: tx creator does not have read access permission on privatedata in chaincodeName:car_private collectionName: collection_parkInfo"

//测试org2的身份
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n car_private -c '{"Args":["queryParkInfo","苏F666v8"]}'
//Chaincode invoke successful. result: status:200 payload:"{\"enter_time\":\"2020-10-20\",\"park_location\":\"\345\256\211\345\276\275\350\212\234\346\271\226\",\"leave_time\":\"2020-10-21\",\"duration_time\":\"2h\"}"


//测试org3的身份
//peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n car_private -c '{"Args":["queryParkInfo","苏F666v8"]}'
//Chaincode invoke successful. result: status:200 payload:"{\"enter_time\":\"2020-10-20\",\"park_location\":\"\345\256\211\345\276\275\350\212\234\346\271\226\",\"leave_time\":\"2020-10-21\",\"duration_time\":\"2h\"}"



//车管所的汽车信息（这里只是模拟，实际中的信息更加复杂）
//创建汽车信息的时候是	CarNumber 车牌号为key
type Car struct {
	Make   string `json:"make"`						//汽车制造厂商
	Colour string `json:"colour"`					//汽车颜色
	Owner  string `json:"owner"`					//车主姓名
	OwnerIdCarNum string `json:"owner_id_car_num"`	//车主身份证号码	有的身份证是有字母的
}


//停车场的停车信息
//创建停车信息的时候也是以车牌号码为key
type ParkCarInfo struct {
	EnterTime string `json:"enter_time""`		//驶入停车场的时间
	ParkLocation string `json:"park_location"`	//停车地点
	LeaveTime string `json:"leave_time"`		//离开停车场的时间
	DurationTime string `json:"duration_time"`	//停留时间
}




func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	if function == "queryCar" {
		return s.queryCar(APIstub, args)
	} else if function == "createCar" {
		return s.createCar(APIstub, args)
	} else if function == "createParkInfo"{
		return s.createParkInfo(APIstub,args)
	} else if function == "queryParkInfo"{
		return s.queryParkInfo(APIstub,args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}



func (s *SmartContract) createCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var car = Car{Make: args[1], Colour: args[2], Owner: args[3], OwnerIdCarNum: args[4]}

	carAsBytes, _ := json.Marshal(car)
	err := APIstub.PutPrivateData("collection_carInfo", args[0], carAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}


func (s *SmartContract) queryCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	carAsBytes, err := APIstub.GetPrivateData("collection_carInfo", args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(carAsBytes)
}

func (s *SmartContract)createParkInfo(APIstub shim.ChaincodeStubInterface,args []string)sc.Response  {
	if len(args) != 5{
		return shim.Error("请输入正确的参数个数（5）")
	}

	parkInfo := ParkCarInfo{EnterTime:args[1],ParkLocation:args[2],LeaveTime:args[3],DurationTime:args[4]}

	parkBytes, _ := json.Marshal(parkInfo)
	//以车牌号码为key
	err := APIstub.PutPrivateData("collection_parkInfo", args[0], parkBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}

func (s *SmartContract)queryParkInfo(APIstub shim.ChaincodeStubInterface,args []string) sc.Response  {
	if len(args) != 1{
		return shim.Error("请输入正确的参数（1）")
	}
	parkBytes, err := APIstub.GetPrivateData("collection_parkInfo", args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(parkBytes)
}




// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
