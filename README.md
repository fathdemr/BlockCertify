# 🎓 BlockCertify

BlockCertify is a backend service that enables **secure diploma storage and verification** using **Arweave** for permanent file storage and **Polygon blockchain** for on-chain verification.

The goal of this project is to ensure that academic diplomas cannot be tampered with and can be verified transparently by anyone.

---

## 🚀 How It Works

1. A diploma (PDF) is uploaded to the API
2. The file is hashed (SHA-256)
3. The diploma file is uploaded to **Arweave**
4. The diploma hash + Arweave Transaction ID are stored on **Polygon**
5. Anyone can later verify a diploma by uploading the same file

---

## 🧱 Architecture

Client (Postman / Frontend)
|
v
Go API (BlockCertify)
|
+– Hash Diploma (SHA-256)
|
+– Store File → Arweave
|
+– Store Hash → Polygon Smart Contract

---

## 🛠️ Tech Stack

- **Go (Golang)** – Backend API
- **Arweave** – Permanent decentralized file storage
- **Polygon (Amoy / Mainnet)** – Blockchain verification layer
- **Solidity** – Smart contract
- **Postman** – API testing

---

## 📦 API Endpoints

### Upload Diploma

POST /api/upload-diploma

**Request (multipart/form-data):**
- `diploma` → PDF file

**Response:**
```json
{
  "arweaveTxId": "arweave_Tx_ID",
  "polygonTxHash": "polygon_Tx_Hash"
}
```

### Verify Diploma

POST /api/verify-diploma

**Request (multipart/form-data):**
- `diploma` → PDF file

**Response:**
```json
{
  "exists": "true",
  "arweaTxID": "arweave_Tx_ID"
}
```

### 📄 `.env.example`

```env
PORT=8080

# Arweave
ARWEAVE_KEY=./arweave_keyfile.json

# Polygon
POLYGON_RPC_URL=https://rpc-amoy.polygon.technology
PRIVATE_KEY=your_private_key_here
CONTRACT_ADDRESS=your_contract_address_here
CHAIN_ID=80002
```
