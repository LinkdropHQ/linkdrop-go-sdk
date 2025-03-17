import { exec } from "child_process";
import { ethers } from "ethers";
import dotenv from "dotenv";
dotenv.config();

const provider = new ethers.JsonRpcProvider(process.env.RPC_URL);
const wallet = new ethers.Wallet(process.env.PRIVATE_KEY, provider);

const chainId = 8453;

// Let's assume we need to generate a recovered link (link key) for a known transferId
const transferId = ethers.Wallet.createRandom().address

// We need to generate a new linkKey and linkKeyId
const linkRecovery = ethers.Wallet.createRandom()
const linkKey = linkRecovery.privateKey // Stays on the client
const linkKeyId = linkRecovery.address

const payload = {
    transferId,
    linkKeyId,
    claimLink: {
        token: {
            type: "ERC20",
            chainId,
            address: "0x833589fcd6edb6e08f4c7c32d4f71b54bda02913",
        },
    }
}

function requestDataToSign(payload) {
    console.log(`Requesting data to sign`);
    exec(`go run ../server/main.go getRecoveredLinkTypedData '${JSON.stringify(payload)}'`, (error, stdout, stderr) => {
        if (error) {
            console.error(`Error: ${error.message}`);
            return;
        }
        if (stderr) {
            console.error(`Stderr: ${stderr}`);
            return;
        }
        const typedData = JSON.parse(stdout)
        console.log("Received data to sign:", typedData)
        signTypedData(typedData)
    });
}

async function signTypedData(typedData) {
    console.log("Signing typed data")
    delete typedData.types["EIP712Domain"] // It's important to remove "EIP712Domain" from types when using ethers.js
    if (typedData.domain["salt"] === '') {
        delete typedData.domain["salt"]
    }
    const signature = await wallet.signTypedData(
        typedData.domain, typedData.types, typedData.message
    );
    console.log("Signature:", signature)
    buildLink(signature)
}

function buildLink(signature) {
    const signatureLength = ((signature).length - 2) / 2 // without 0x (-2) prefix in bytes (/2)
    console.log(`Recovered Link:`)
    console.log(`https://p2p.linkdrop.io/#/code?k=${linkKey}&sg=${ethers.encodeBase58(signature)}&i=${transferId}&c=${chainId}&v=3&sgl=${signatureLength}&src=p2p`)
}

requestDataToSign(payload)