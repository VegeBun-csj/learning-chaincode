/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const { stat } = require('fs');
const path = require('path');

const util = require('util');
const ccpPath = path.resolve(__dirname, '..', '..', 'first-network', 'connection-org1.json');

async function main() {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');
        // Get the contract from the network.
        const contract = network.getContract('fabcar');

        const transaction = contract.createTransaction("createCar");
    
        await transaction.addCommitListener((err,transactionID,status,blockNum)=>{
            if(err){
                console.log(err);
                return;
            }
            if(status === 'VALID'){
                console.log('transaction committed');
                console.log(`transactionID is : ${transactionID}`)
                console.log(`status is: ${status}`);
                console.log(`blockNum is:${blockNum}`);
                console.log('transaction committed end');
            }else{
                console.log('transation commit failed');
                console.log(`status is ${status}`);
            }

        })

        await transaction.submit('CAR19', 'Honda', 'Accord', 'Black', 'Tom');
        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
