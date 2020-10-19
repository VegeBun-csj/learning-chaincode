/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
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

        //监听区块事件
        await network.addBlockListener("myblockListener",(err,block)=>{
            if (err) {
                console.log(err);
                return;
              }
        
              console.log('*************** start block header **********************')
              console.log(util.inspect(block.header, {showHidden: false, depth: 5}))    //输出区块头数据，设置递归深度为5
              console.log('*************** end block header **********************')
              console.log('*************** start block data **********************')    
              let data = block.data.data[0];                                            //输出区块体数据
              console.log(util.inspect(data, {showHidden: false, depth: 5}))
              console.log('*************** end block data **********************')
              console.log('*************** start block metadata ****************')      //输出区块元数据
              console.log(util.inspect(block.metadata, {showHidden: false, depth: 5}))
              console.log('*************** end block metadata ****************')
        })


        // Get the contract from the network.
        const contract = network.getContract('fabcar');

        await contract.submitTransaction('createCar', 'CAR16', 'Honda', 'Accord', 'Black', 'Tom');
        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
