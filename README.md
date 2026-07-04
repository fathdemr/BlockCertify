# 🎓 BlockCertify

BlockCertify is a backend service that enables **secure diploma storage and verification** using **Arweave** for permanent file storage and **Polygon blockchain** for on-chain verification.

The goal of this project is to ensure that academic diplomas cannot be tampered with and can be verified transparently by anyone.

---

## 🚀 How It Works

1. An admin uploads a diploma (PDF) to the API
2. The file is hashed (SHA-256)
3. The diploma file is uploaded to **Arweave**
4. The diploma hash + Arweave Transaction ID are stored on **Polygon**
5. Anyone can later verify a diploma: the API checks that the hash recorded in the database also exists on the Polygon smart contract and points to the same Arweave transaction

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
|
+– Store Metadata → PostgreSQL

---

## 🛠️ Tech Stack

- **Go (Golang)** – Backend API (Gin)
- **Arweave** – Permanent decentralized file storage
- **Polygon (Amoy / Mainnet)** – Blockchain verification layer
- **Solidity** – Smart contract
- **PostgreSQL** – Diploma metadata & users
- **Redis** – Caching
- **JWT (RS256)** – RSA-signed authentication tokens

---

## 🔐 Authentication

- Login: `POST /exapi/user/login` → returns an RS256-signed JWT (also set as `jwt` cookie)
- Send the token as `Authorization: Bearer <token>` (or rely on the cookie)
- Diploma issuance endpoints require an **admin** role
- Admin registration (`POST /exapi/user/register/admin`) itself requires an existing admin; bootstrap the first admin via `data/seed.sql`

---

## 📦 API Endpoints

### Diplomas (`/api`)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/upload` | admin | Full issuance: hash → Arweave → Polygon → DB |
| POST | `/api/prepare-upload` | admin | Phase 1 (MetaMask flow): hash + Arweave upload |
| POST | `/api/confirm-upload` | admin | Phase 2 (MetaMask flow): verify on-chain, save to DB |
| POST | `/api/verify` | public | Verify a diploma against DB **and** the Polygon contract |
| GET | `/api/records` | JWT | List diplomas belonging to the authenticated user |
| GET | `/api/records/:diplomaId` | public | Stream the diploma PDF from Arweave |
| GET | `/api/wallet/status` | public | Platform wallet status |

Upload requests are `multipart/form-data` with a `diploma` PDF (max 25 MB) plus
`firstName`, `lastName`, `email`, `university`, `faculty`, `department`,
`graduationYear`, `studentNumber`, `nationality` fields.

**Verify request:**
```json
{ "diplomaId": "DIP-..." }
```

**Verify response:**
```json
{
  "verified": true,
  "diplomaID": "DIP-...",
  "diplomaHash": "…",
  "arweaveTxID": "…",
  "arweaveUrl": "https://arweave.net/…",
  "polygonTxHash": "0x…",
  "studentName": "…",
  "university": "…",
  "degree": "…",
  "issueDate": "2026-07-04"
}
```

`verified: false` (HTTP 200) means the diploma could not be validated.
HTTP 502 means verification infrastructure (Polygon RPC) is unavailable — a
different situation than a fake diploma.

### Users (`/exapi/user`)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/login` | public | Login, returns JWT |
| POST | `/logout` | public | Clears the JWT cookie |
| POST | `/register/admin` | admin | Register a new admin |
| GET | `/me` | JWT | Current user profile |
| PUT | `/me` | JWT | Update profile / password |

### Health

- `GET /healthz` – liveness probe
- `GET /readyz` – readiness probe

---

## ⚙️ Configuration

Configuration lives in `internal/config.yaml` (gitignored — never commit it).
Keys used for JWT signing are RSA (PEM, PKCS#8):

```yaml
crypto:
  rsa_private_key: |
    -----BEGIN PRIVATE KEY-----
    ...
    -----END PRIVATE KEY-----
  rsa_public_key: |
    -----BEGIN PUBLIC KEY-----
    ...
    -----END PUBLIC KEY-----
  token_expire_duration_hour: 86400
arweave:
  walletKey: '{...}'   # Arweave JWK
polygon:
  rpcUrl: 'https://rpc-amoy.polygon.technology'
  privateKey: '...'
  contractAddress: '0x...'
  chainID: 80002
```

Environment variables override config values (`DB_LIVE_HOST` → `db.live.host`).

---

## 🧪 Development

```bash
go build ./...   # build
go vet ./...     # static checks
go test ./...    # unit tests
```

CI (GitHub Actions) runs build + vet + test on every push and pull request.
