# decentralized-communication-protocol
General Objective: An integrated decentralized communication protocol that operates within local networks, enabling high-quality media sharing and calls without requiring an internet connection or a central server.

---

1. Networking Stack

· Communication Model: Pure Peer-to-Peer (P2P); every user is a "node" that acts as both a server and a client simultaneously.
· User Discovery: Uses mDNS (ZeroConf) protocol, allowing devices to find each other immediately upon connecting to the same Wi-Fi network without any manual setup.
· Transport Protocol:
  · TCP: For text messages and files to ensure complete, in-order delivery.
  · UDP (WebRTC): For voice and video calls to ensure minimal latency.

2. Core Features

· Chats:
  · 1-on-1 private messaging.
  · Decentralized group chats based on data synchronization between active members.
· Media Sharing: Send and receive images, videos, and large files at lightning speed (limited only by your router's speed).
· Calls: HD-quality voice and video calls over the local network.
· Heartbeating: Checks the status of connected users every second to maintain an accurate, real-time "currently online" list.
· Voice Messages: Send voice recordings with a waveform visualizer feature.

3. Security Layer

· End-to-End Encryption (E2EE): Uses X3DH and Double Ratchet algorithms (same as Signal) for absolute privacy.
· Key Exchange: Uses elliptic-curve cryptography for key exchange between devices upon discovery.
· Encrypted Local Storage: Local database on each device is encrypted, preventing message access even if the device is physically accessed without the app.
· Digital Signatures: Every message is digitally signed to prevent impersonation within the network.

4. User Experience (UX) & Professionalism

· Smart Sync: When a new device joins a group, the "oldest" available member automatically provides it with the last 100 messages.
· Typing & Read Indicators: Real-time typing indicators and read receipts.
· User Interface: Modern design with Dark Mode support and a drag-and-drop interface for file sending.
· Background Operation: The app continues to run and receive notifications even when the main interface is closed.
