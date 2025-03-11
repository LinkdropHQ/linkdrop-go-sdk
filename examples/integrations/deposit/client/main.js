import { exec } from "child_process";
import { ethers } from "ethers";
import dotenv from "dotenv";
dotenv.config();

const provider = new ethers.JsonRpcProvider(process.env.RPC_URL);
const wallet = new ethers.Wallet(process.env.PRIVATE_KEY, provider);

const linkWallet = ethers.Wallet.createRandom();
const claimLinkKey = linkWallet.privateKey;
const transferId = linkWallet.address;
const claimLink = {
    token: {
        type: "ERC20",
        chainId: 8453,
        address: "0x833589fcd6edb6e08f4c7c32d4f71b54bda02913",
    },
    sender: "0x5659A8557FdBA11AA04cfCfcc59EeF9FA412A7dD",
    amount: "100000",
    expiration: 1773159165,
}

const payload = {
    command: "getDepositParams",
    transferId: transferId,
    claimLink
}

function requestDataToSign(payload) {
    console.log(`Requesting data to sign`);
    exec(`go run ../server/main.go '${JSON.stringify(payload)}'`, (error, stdout, stderr) => {
        if (error) {
            console.error(`Error: ${error.message}`);
            return;
        }
        if (stderr) {
            console.error(`Stderr: ${stderr}`);
            return;
        }
        console.log(`Received data to sign: ${stdout}`);
        const depositParams = JSON.parse(stdout)
        sendTransaction(depositParams.to, depositParams.value, depositParams.data)
    });
}


function sendTransaction(to, value, data) {
    console.log(`Sending transaction to ${to} with value ${value} and data ${data}`);
    // wallet.sendTransaction({
    //     to: to,
    //     value: ethers.utils.parseEther(value),
    //     data: data
    // }).then(txResponse => {
    //     console.log(`Transaction sent: ${txResponse.hash}`);
    //     txResponse.wait().then(receipt => {
    //         console.log(`Transaction confirmed: ${receipt.transactionHash}`);
    //     }).catch(err => {
    //         console.error(`Error waiting for transaction confirmation: ${err.message}`);
    //     });
    // }).catch(err => {
    //     console.error(`Error sending transaction: ${err.message}`);
    // });
}


requestDataToSign(payload);

