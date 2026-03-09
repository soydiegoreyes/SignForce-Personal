# System Overview — SignForce

## What is SignForce?

SignForce is a SaaS platform for **advanced electronic signatures** with legal validity in Mexico under NOM-151-SCFI-2016. It allows businesses to upload documents, invite signers, apply cryptographic signatures using SAT credentials (e.firma/FIEL), and package everything into legally compliant containers.

## Core Capabilities

### 1. Document Management
Users upload PDF documents, classify them (contract, memorandum, invitation), and organize them into folders. Documents are encrypted at rest with AES-256-GCM.

### 2. Signature Workflow
- Upload document → Add signers/observers → Place signature boxes → Send invitations
- Each signer receives a notification and uses their RSA private key + X.509 certificate to sign
- Signatures are generated in XAdES-BES format

### 3. Cryptographic Operations
- **Encryption:** AES-256-GCM with unique nonce per 64KB chunk
- **Signing:** RSA + SHA-256 using SAT-issued certificates
- **Format:** XAdES-BES XML signatures (ETSI EN 319 132)
- **Packaging:** ASiC-E containers (ETSI EN 319 162)
- **Verification:** QR codes for instant signature verification

### 4. Identity Verification
- Biometric facial verification with liveness detection (TensorFlow)
- X.509 certificate chain validation against SAT
- CURP/RFC validation (currently commented out)

### 5. Team Management
- Root users create companies (institutions)
- Admins create teams within companies
- Team members are invited via email
- Role-based access: root → admin → member

## User Roles

| Role | Access |
|------|--------|
| Root | Create teams, view all company data, approve clients |
| Admin | Manage team members, configure settings |
| Member | Upload/sign documents, manage personal keys |
| External Guest | Sign specific documents via invitation link |

## Data Flow

```
User uploads PDF
  → sfmiddle receives file + metadata
  → sfmiddle stores file on disk
  → sfback encrypts file (AES-256-GCM)
  → sfback stores hash in MySQL
  → User adds signers via UI
  → sfmiddle creates folder + invitations
  → emailServ sends notification emails
  → Signer opens invitation link
  → sfia verifies identity (biometric)
  → Signer uploads RSA key + certificate
  → sfback generates XAdES-BES signature
  → sfmiddle packages ASiC-E container
  → QR code generated for verification
```
