# Architecture — SignForce

## Service Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    INTERNET                                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│              Nginx (Host - SSL Termination)                  │
│              Let's Encrypt certificates                      │
│              testsignforce.luxspace.org                       │
│              Ports: 80/443 → 127.0.0.1:8090                │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│              Docker Network: signforce_default               │
│                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ signforce-  │  │   sfmiddle  │  │   sfback    │         │
│  │ web (Nginx) │  │  Gateway    │  │  Crypto API │         │
│  │  :8090→:80  │  │   :5002     │  │   :5001     │         │
│  │             │  │             │  │             │         │
│  │  Static     │  │  Auth/JWT   │  │  AES-256    │         │
│  │  HTML/CSS   │  │  Routing    │  │  XAdES-BES  │         │
│  │  JS files   │  │  Sessions   │  │  RSA/SHA256 │         │
│  └─────────────┘  └──────┬──────┘  └──────┬──────┘         │
│                          │                │                  │
│                    ┌─────▼────────────────▼─────┐           │
│                    │       MySQL 8.0             │           │
│                    │       :3306                 │           │
│                    │  users, documents, keys,    │           │
│                    │  signatures, teams, etc.    │           │
│                    └────────────────────────────┘           │
│                                                              │
│  ┌─────────────┐  ┌─────────────┐                           │
│  │  emailServ  │  │    sfia     │                           │
│  │  SMTP       │  │  Python AI  │                           │
│  │  Notif.     │  │  Biometric  │                           │
│  │             │  │  OCR        │                           │
│  └─────────────┘  └─────────────┘                           │
└─────────────────────────────────────────────────────────────┘
```

## API Endpoints

### sfback (:5001) — Cryptographic Core
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/login` | POST | Authenticate user |
| `/getuser` | GET | Get user data |
| `/signdocument` | POST | Generate XAdES signature |
| `/hashdatab64` | POST | Hash data in Base64 |
| `/uploadKeys` | POST | Upload RSA keys |
| `/logout` | POST | End session |

### sfmiddle (:5002) — Gateway (30 endpoints)
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/register` | POST | Register new institution |
| `/loginUser` | POST | User login |
| `/logoutUser` | POST | User logout |
| `/uploadDocs` | POST | Upload documents |
| `/downloadDoc` | POST | Download document |
| `/uploadk` | POST | Upload RSA keys |
| `/statusk` | GET | Get key status |
| `/newSignFolder` | POST | Start signing process |
| `/closeInvite` | POST | Close folder + send invitations |
| `/getinvite` | GET | Get invitation data |
| `/signDocument` | POST | Sign a document |
| `/getAsice` | GET | Build ASiC-E container |
| `/getfolders` | GET | List user folders |
| `/approvalsDash` | GET | Admin approvals dashboard |
| `/viewSign` | GET | View signature details |
| `/inviteuser` | POST | Invite user to institution |
| `/createuser` | POST | Create user in institution |

## Database Schema (Key Tables)
- `users` — User accounts with hashed passwords
- `institutions` — Companies/organizations  
- `teams` — Teams within institutions
- `documents` — Document metadata + encrypted file paths
- `folders` — Document groupings for signing
- `folderdocuments` — Many-to-many folder↔document
- `invitations` — Signing invitations
- `signatures` — XAdES signature records with IP, user agent, timestamp
- `userkeys` — RSA key pairs (file paths, cert metadata)
- `faceembeddings` — Biometric face vectors
- `payment` — Payment records ⚠️ stores card data in plain text

## Technology Decisions
- **Go** for backend: performance, strong typing, native concurrency
- **Python** for AI: TensorFlow ecosystem, OCR libraries
- **Static frontend**: no framework dependency, fast loading, easy deployment
- **Docker**: consistent deployment, isolation between services
- **MySQL**: relational integrity for document/signature audit trail
