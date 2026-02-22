# BlockCertify API Documentation

**Base URL:** `http://localhost:8080/api/v1`

All protected routes require a valid JWT sent as a `jwt` cookie (set automatically after login).

---

## 🔓 Public Routes

---

### Health Check

#### `GET /ping`

Simple health-check to confirm the server is running.

**Response `200`**
```json
{ "message": "pong" }
```

---

### Auth — `/auth/user`

> ⚠️ Note: The auth routes are nested under `/auth/user`, so the full paths are `/api/v1/auth/user/...`

---

#### `POST /auth/user/register/admin`

Registers a new admin user for a university.

**Request Body** `application/json`

| Field          | Type   | Required | Description                          |
|----------------|--------|----------|--------------------------------------|
| `firstName`    | string | ✅        | Min 2 characters                     |
| `lastName`     | string | ✅        | Min 2 characters                     |
| `email`        | string | ✅        | Must be a valid email address        |
| `universityId` | string | ✅        | UUID of the university               |
| `password`     | string | ✅        | Min 8 characters                     |

**Example Request**
```json
{
  "firstName": "Fatih",
  "lastName": "Demir",
  "email": "fatih@university.edu",
  "universityId": "550e8400-e29b-41d4-a716-446655440000",
  "password": "securepass123"
}
```

**Response `201`**
```json
{ "message": "Success" }
```

**Response `400`**
```json
{ "error": "Invalid request body" }
```

---

#### `POST /auth/user/login`

Authenticates a user and returns a JWT token. The token is also set as an `httpOnly` cookie named `jwt`.

**Request Body** `application/json`

| Field      | Type   | Required | Description           |
|------------|--------|----------|-----------------------|
| `email`    | string | ✅        | Registered email      |
| `password` | string | ✅        | Min 8 characters      |

**Example Request**
```json
{
  "email": "fatih@university.edu",
  "password": "securepass123"
}
```

**Response `200`**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "tokenType": "Bearer",
  "expiresIn": 3600,
  "role": "admin"
}
```

**Response `400`**
```json
{ "error": "Invalid credentials", "details": "..." }
```

---

#### `POST /auth/user/logout`

Clears the `jwt` cookie.

**Response `200`**
```json
{ "message": "Logged out successfully" }
```

---

### Universities

#### `GET /universities`

Returns a list of all registered universities. Used to populate dropdowns during registration.

**Response `200`**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Istanbul Technical University"
  }
]
```

---

## 🔒 Protected Routes (JWT Required)

All routes below require the `jwt` cookie to be present. Requests without a valid token will receive `401 Unauthorized`.

---

### Diploma — `/diploma`

---

#### `POST /diploma/prepare`

**Phase 1 of diploma issuance.** Accepts the diploma PDF and student metadata, uploads the file to Arweave, and returns data for the frontend to sign on the Polygon blockchain via MetaMask.

**Request** `multipart/form-data`

| Field            | Type    | Required | Description                         |
|------------------|---------|----------|-------------------------------------|
| `diploma`        | file    | ✅        | PDF file only                       |
| `firstName`      | string  | ✅        | Student's first name                |
| `lastName`       | string  | ✅        | Student's last name                 |
| `email`          | string  | ✅        | Student's email address             |
| `university`     | string  | ✅        | University name                     |
| `faculty`        | string  | ✅        | Faculty name                        |
| `department`     | string  | ✅        | Department name                     |
| `graduationYear` | integer | ✅        | Year between 1950 and current year+1|
| `studentNumber`  | string  | ✅        | Student ID number                   |
| `nationality`    | string  | ✅        | Student's nationality               |

**Response `200`**
```json
{
  "diplomaHash": "sha256-hash-of-the-pdf",
  "arweaveTxID": "arweave-transaction-id",
  "arweaveUrl": "https://arweave.net/arweave-transaction-id"
}
```

**Response `400`**
```json
{ "error": "Invalid metadata", "details": "firstName is required" }
```

---

#### `POST /diploma/confirm`

**Phase 2 of diploma issuance.** Called after MetaMask successfully signs and submits the Polygon transaction. Saves the diploma record to the database.

**Request Body** `application/json`

| Field            | Type    | Required | Description                                  |
|------------------|---------|----------|----------------------------------------------|
| `diplomaHash`    | string  | ✅        | Hash returned from `/prepare`                |
| `arweaveTxID`    | string  | ✅        | Arweave TX ID returned from `/prepare`       |
| `polygonTxHash`  | string  | ✅        | Transaction hash from MetaMask/ethers.js     |
| `blockNumber`    | integer | ❌        | Block number from the Polygon transaction    |
| `firstName`      | string  | ✅        | Student's first name                         |
| `lastName`       | string  | ✅        | Student's last name                          |
| `email`          | string  | ✅        | Student's email address                      |
| `university`     | string  | ✅        | University name                              |
| `faculty`        | string  | ❌        | Faculty name                                 |
| `department`     | string  | ✅        | Department name                              |
| `graduationYear` | integer | ✅        | Graduation year                              |
| `studentNumber`  | string  | ❌        | Student ID number                            |
| `nationality`    | string  | ❌        | Student's nationality                        |

**Example Request**
```json
{
  "diplomaHash": "abc123...",
  "arweaveTxID": "txid123...",
  "polygonTxHash": "0xabc...",
  "blockNumber": 12345678,
  "firstName": "Fatih",
  "lastName": "Demir",
  "email": "fatih@university.edu",
  "university": "Istanbul Technical University",
  "department": "Computer Engineering",
  "graduationYear": 2024
}
```

**Response `200`**
```json
{
  "diplomaHash": "abc123...",
  "arweaveTxID": "txid123..."
}
```

---

#### `POST /diploma/verify`

Verifies a diploma's authenticity by its Diploma ID. Checks the on-chain Polygon record and Arweave storage.

**Request Body** `application/json`

| Field       | Type   | Required | Description                  |
|-------------|--------|----------|------------------------------|
| `DiplomaID` | string | ✅        | Public ID of the diploma     |

**Example Request**
```json
{
  "DiplomaID": "diploma-public-id"
}
```

**Response `200`**
```json
{
  "verified": true,
  "diplomaHash": "abc123...",
  "arweaveTxID": "txid123...",
  "arweaveUrl": "https://arweave.net/txid123",
  "studentName": "Fatih Demir",
  "university": "Istanbul Technical University",
  "degree": "Computer Engineering",
  "issueDate": "2024-06-15T00:00:00Z",
  "polygonTxHash": "0xabc...",
  "diplomaID": "diploma-public-id"
}
```

**Response `400`**
```json
{ "error": "Diploma not found or verification failed" }
```

---

#### `GET /diploma/records`

Returns all diploma records stored in the database.

**Response `200`**
```json
[
  {
    "diplomaId": "diploma-public-id",
    "userName": "Fatih Demir",
    "department": "Computer Engineering",
    "createDate": "2024-06-15T10:30:00Z",
    "diplomaPdf": "https://arweave.net/txid123"
  }
]
```

---

#### `GET /diploma/records/:diplomaId`

Streams the diploma PDF directly from Arweave for a given diploma ID.

**Path Parameter**

| Parameter   | Description              |
|-------------|--------------------------|
| `diplomaId` | Public ID of the diploma |

**Response `200`**  
Binary PDF stream with headers:
```
Content-Type: application/pdf
Content-Disposition: inline; filename=diploma.pdf
```

**Response `400`**
```json
{ "error": "Failed to fetch diploma" }
```

---

### Wallet — `/wallet`

---

#### `POST /wallet/upload-key-file`

Connects an Arweave wallet by uploading its JSON key file. Returns the wallet address and AR balance.

**Request** `multipart/form-data`

| Field    | Type | Required | Description              |
|----------|------|----------|--------------------------|
| `wallet` | file | ✅        | Arweave JSON key file    |

**Response `200`**
```json
{
  "address": "arweave-wallet-address",
  "balance": "1.234",
  "status": "wallet connected"
}
```

**Response `400`**
```json
{ "error": "...", "details": "wallet file is missing" }
```

---

## Error Format

All error responses follow this structure:

```json
{
  "error": "Human-readable error message",
  "details": "Optional technical details"
}
```

| HTTP Status | Meaning                             |
|-------------|-------------------------------------|
| `200`       | Success                             |
| `201`       | Resource created                    |
| `400`       | Bad request / validation error      |
| `401`       | Unauthorized (missing/invalid JWT)  |
| `415`       | Unsupported media type              |
| `500`       | Internal server error               |
