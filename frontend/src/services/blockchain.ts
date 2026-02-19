// services/blockchain.ts
import { ethers } from 'ethers';
import CONTRACT_ABI from '../contractABI.json'; // extract your ABI to a JSON file

const CONTRACT_ADDRESS = import.meta.env.VITE_CONTRACT_ADDRESS;

export async function storeDiplomaWithMetaMask(
    diplomaHash: string,
    arweaveTxId: string
): Promise<{ txHash: string; blockNumber: number }> {
    if (!window.ethereum) throw new Error('MetaMask not installed');

    if (!CONTRACT_ADDRESS) {
        throw new Error(
            'Contract address is not configured. ' +
            'Add VITE_CONTRACT_ADDRESS to your frontend/.env file and restart the dev server.'
        );
    }

    // Use MetaMask as the provider AND signer
    const provider = new ethers.BrowserProvider(window.ethereum);
    const signer = await provider.getSigner(); // This is the MetaMask account

    const contract = new ethers.Contract(CONTRACT_ADDRESS, CONTRACT_ABI, signer);

    // MetaMask will pop up asking the user to sign & pay gas
    const tx = await contract.storeDiploma(diplomaHash, arweaveTxId);
    const receipt = await tx.wait();

    return {
        txHash: receipt.hash,
        blockNumber: receipt.blockNumber,
    };
}

export async function verifyDiplomaOnChain(
    diplomaHash: string
): Promise<{ exists: boolean; arweaveTxId: string }> {
    // Read-only — no signer needed
    const provider = new ethers.BrowserProvider(window.ethereum);
    const contract = new ethers.Contract(CONTRACT_ADDRESS, CONTRACT_ABI, provider);

    const [exists, arweaveTxId] = await contract.verifyDiploma(diplomaHash);
    return { exists, arweaveTxId };
}