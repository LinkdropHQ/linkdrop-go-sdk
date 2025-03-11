import { Wallet } from "ethers";

const linkKey = Wallet.createRandom().privateKey;
console.log(linkKey)