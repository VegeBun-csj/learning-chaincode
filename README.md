# learning-chaincode
learn to write hyperledger fabric chaincode

- #### couchcar

  this is a chaincode writed by go language when using CouchDB  as the State Database in Fabric

  this contract is about  the car adapted by the marbles chaincode and Fabcar chaincode.I using the first-network sample to test the chaincode

  `./byfn.sh -c mychannel -s couchdb` 

  this chaincode using `four indexes`,so we can query by them using rich queries very easily:

  - index_Color
  - index_Owner
  - index_Price
  - index_Year

  > OK,the more using detail is writed in the chaincode annotation

- #### access-control

  using `Client Identity Chaincode Library` called `cid library` to reach access control. This chiancode modified a little by the fabcar chaincode.
  
- #### Using_FabricEvent

  This example is used to learn the event using in the Hyperledger Fabric Node SDK By changing the Fabcar example.

  > you can put the example `Event-Fabcar` in the directory of`fabric-samples`.And the next steps are same as the fabcar.Also you should replace the original fabcar chaincode with `chaincode ` in this directory .

  In this example , you can learn three types of events:

  - ContractEvent(chaincode)
  - BlockEvent(Block)
  - CommitEvent(Transaction)
  
- ### car_privateData

  This Example is used to display the usage of Private Data in Fabric 1.4.4. The details of it  are written in the chaincode.
  
- ### key-level-ep

  This chaincode demonstrated the usage of the `Key-level` endorsement policy,which is a particular endorsement policy unlike the `chaincode-level` endorsement policy.

  The using details are written in the chancode folder.