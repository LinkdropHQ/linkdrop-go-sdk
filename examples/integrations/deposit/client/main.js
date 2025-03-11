import { exec } from "child_process";
import { ethers, encodeBase58 } from "ethers";
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
    sender: wallet.address,
    amount: "100000",
    expiration: 1773159165,
}

const payload = {
    transferId: transferId,
    claimLink
}

function requestDataToSign(payload) {
    console.log(`Requesting data to sign`);
    exec(`go run ../server/main.go getDepositParams '${JSON.stringify(payload)}'`, (error, stdout, stderr) => {
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
    wallet.sendTransaction({
        to: to,
        value: ethers.parseEther(value),
        data: data
    }).then(txResponse => {
        console.log(`Transaction sent: ${txResponse.hash}`);
        txResponse.wait().then(receipt => {
            console.log(`Transaction confirmed: ${receipt.hash}`);
            registerDeposit(receipt.hash)
        }).catch(err => {
            console.error(`Error waiting for transaction confirmation: ${err.message}`);
        });
    }).catch(err => {
        console.error(`Error sending transaction: ${err.message}`);
    });
}

function registerDeposit(txHash) {
    console.log(`Registering deposit`);
    payload["txHash"] = txHash
    exec(`go run ../server/main.go registerDeposit '${JSON.stringify(payload)}'`, (error, stdout, stderr) => {
        if (error) {
            console.error(`Error: ${error.message}`);
            return;
        }
        if (stderr) {
            console.error(`Stderr: ${stderr}`);
            return;
        }
        console.log(`Status: ${stdout}`);
        console.log(`Claim Link (client-side generated): \nhttps://p2p.linkdrop.io/#/code?k=${encodeBase58(linkWallet.privateKey)}&c=8453&v=3&src=p2p`);
    });
}

requestDataToSign(payload);
