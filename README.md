# Decentralized Communication Protocol (DCP)

![Project Status](https://img.shields.io/badge/Status-Active--Development-green)
![Platform](https://img.shields.io/badge/Platform-Local--Network-blue)
![Security](https://img.shields.io/badge/Security-E2EE-red)
![Backend](https://img.shields.io/badge/Backend-Go-00ADD8?logo=go)
![Frontend](https://img.shields.io/badge/Frontend-React-61DAFB?logo=react)

## 📌 Overview
Decentralized Communication Protocol (DCP) is a serverless, Peer-to-Peer (P2P) framework designed for high-performance interaction within Local Area Networks (LAN). By decoupling the core networking logic from the user interface, DCP ensures high availability and cryptographic integrity.

The protocol follows a Modular Hybrid Architecture, where a headless Go node handles complex network operations, and a React-based dashboard provides a real-time visual interface via a local WebSocket bridge.

---

## 🏗 Technical Architecture

The system is engineered into distinct functional layers to ensure scalability from LAN to WAN:

*   Identity Layer (`/internal/peer`): Manages cryptographic identities using Ed25519. It generates persistent PeerIDs, ensuring users are recognized across sessions without relying on static IP addresses.
*   Discovery Layer (`/internal/discovery`): Implements mDNS (ZeroConf) for zero-configuration peer detection. It allows nodes to "announce" themselves and "discover" others on the same subnet automatically.
*   Transport Layer (`/internal/transport`): A multi-protocol engine that abstracts TCP, UDP, and QUIC connections, optimized for different data types (e.g., text vs. media).
*   Security Layer (`/internal/session`): Handles the X3DH key exchange and Double Ratchet algorithm to provide End-to-End Encryption (E2EE) with Perfect Forward Secrecy (PFS).
*   Gateway Layer (`/internal/gateway`): Acts as the IPC Bridge. It translates internal P2P events into JSON messages for the React frontend via WebSockets.
*   Sync Layer (`/internal/sync`): Utilizes Gossip Protocols to synchronize state and propagate messages across the decentralized network efficiently.

---
## 🚀 Key Capabilities

| Feature | Technical Implementation |
| :--- | :--- |
| Real-time Streaming | Multi-party WebRTC orchestration with dynamic coordinator election. |
| Decentralized Groups | Topic-based messaging using GossipSub for seamless data sync. |
| Secure Handshake | Elliptic-Curve Diffie-Hellman (ECDH) for safe key exchange over LAN. |
| High-Speed Transfer | Parallel TCP streams for lightning-fast file and media sharing. |
| Privacy First | Metadata-protected headers and encrypted local storage for every node. |

---

## 📂 Project Structure
'''text
decentralized-communication-protocol/

├── backend/                        # Golang P2P Core Engine
│   ├── cmd/
│   │   └── node/
│   │       └── main.go             # Entry point (Initializes Node & Gateway)
│   ├── configs/                    # [جديد] إدارة إعدادات النظام
│   │   ├── config.yaml             # Ports, Log levels, Retries
│   │   └── .env                    # Sensitive environment variables
│   ├── internal/                   # Core Logic (The "Brain")
│   │   ├── peer/                   # Identity & Routing
│   │   ├── transport/              # Network I/O (TCP/UDP/QUIC)
│   │   ├── discovery/              # mDNS Local Discovery
│   │   ├── session/                # Security (E2EE, Ratchet)
│   │   ├── gateway/                # WebSocket Bridge to React
│   │   ├── sync/                   # Gossip & State Sync
│   │   └── database/               # Local Storage (Encrypted)
│   ├── api/                        # API & Protocol Definitions
│   │   └── protocol.proto          # Protobuf / JSON Schemas
│   ├── pkg/                        # Global Utilities
│   │   ├── crypto/                 # Encryption wrappers
│   │   └── logger/                 # Unified system logging
│   ├── Makefile                    # [جديد] أتمتة المهام (Build, Run, Test)
│   └── go.mod                      # Go Dependencies
│
├── frontend/                       # React User Interface (Vite)
│   ├── src/
│   │   ├── components/             # UI Components (Chat, Peers)
│   │   ├── hooks/                  # useWebSocket, useAuth
│   │   ├── services/               # API & Bridge Handlers
│   │   ├── store/                  # [جديد] إدارة حالة التطبيق (Zustand/Redux)
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   └── .env.local                  # Frontend environment settings
│
├── scripts/                        # [جديد] سكربتات مساعدة للمهندسين
│   ├── deploy_test.sh              # تشغيل عدة عقد (Nodes) محلياً للاختبار
│   └── generate_certs.sh           # توليد شهادات الأمان الأولية
│
├── docs/                           # Documentation & Specifications
│   ├── architecture.md             # شرح تدفق البيانات (Data Flow)
│   └── api_spec.md                 # توثيق رسائل الـ WebSocket
│
├── Makefile                        # [جديد] لبناء وتشغيل المشروع بالكامل (Full-stack)
├── .gitignore                      # استثناء ملفات node_modules و Binaries
└── README.md                       # الوثيقة الرئيسية للمشروع

'''
---

## 🚀 Execution Workflow

### Prerequisites
*   Go 1.21+
*   Node.js (v18+) & NPM
*   Linux Environment (Optimized for Linux Mint/Debian)

### Launching the Node
1.  **Start the Backend Engine:**
    ```bash
    cd backend
    go run cmd/node/main.go
    ```
2.  **Start the Frontend Dashboard:**
    ```bash
    cd frontend
    npm run dev
    ```

---
📈 Development Roadmap

    [ ] Phase 1: Core Repository Architecture & Identity Generation.

    [ ] Phase 2: mDNS Peer Discovery & Secure Handshake implementation.

    [ ] Phase 3: P2P Encrypted Text Messaging & Gossip Sync (Pub/Sub).

    [ ] Phase 4: Real-time Audio/Video Streaming & Screen Sharing.

    [ ] Phase 5: Optimization for High-Density Local Networks.

    
## 🛡 Security & Reliability
*   **Forward Secrecy:** Message keys are rotated constantly; compromising one key does not compromise past messages.
*   **Zero-Trust Networking:** Every data packet is signed and verified at the protocol level.
*   **Headless Capability:** The backend can run as a background service on low-power devices (Raspberry Pi/Routers) without the UI.

## 🤝 License
This project is licensed under the MIT License.
