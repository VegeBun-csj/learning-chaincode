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

- #### access-control

  using `Client Identity Chaincode Library` called `cid library` to reach access control. This chiancode modified a little by the fabcar chaincode.