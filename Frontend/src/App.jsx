import React, { useEffect, useState } from 'react';
import { useWebSocket } from './hooks/useWebSocket';

function App() {
    // استخراج المنفذ من الرابط لتشغيل أكثر من نسخة (مثلاً ?port=8081)
    const queryParams = new URLSearchParams(window.location.search);
    const backendPort = queryParams.get('port') || '8080';
    
    const { data, connected, sendMsg } = useWebSocket(`ws://localhost:${backendPort}/ws`);
    
    const [peers, setPeers] = useState([]);
    const [myInfo, setMyInfo] = useState({ id: 'Initializing...', port: backendPort });
    const [logs, setLogs] = useState([]);

    // جلب هوية العقدة المحلية عند التشغيل
    useEffect(() => {
        fetch(`http://localhost:${backendPort}/status`)
            .then(res => res.json())
            .then(data => setMyInfo({ id: data.peer_id, port: data.port }))
            .catch(err => console.error("Status Error:", err));
    }, [backendPort]);

    // معالجة تدفق البيانات من الباكند
    useEffect(() => {
        if (data) {
            switch (data.type) {
                case 'peer_discovered':
                    setPeers(prev => {
                        if (prev.some(p => p.peer_id === data.data.peer_id)) return prev;
                        addLog(`🌐 Node Discovery: ${data.data.peer_id.substring(0, 8)} Joined`);
                        return [...prev, data.data];
                    });
                    break;

                case 'peer_removed':
                    setPeers(prev => prev.filter(p => p.peer_id !== data.data.peer_id));
                    addLog(`🗑️ Node Removal: ${data.data.peer_id.substring(0, 8)} Timed out`);
                    break;

                case 'ping_sent_status':
                    const icon = data.data.status === 'Delivered' ? '✅' : '❌';
                    addLog(`${icon} Ping to ${data.data.target.substring(0, 6)}: ${data.data.status}`);
                    break;

                case 'ping_received':
                    addLog(`📥 INCOMING PING: Connection established from peer!`);
                    // تنبيه بصري بسيط
                    document.title = "⚠️ New Ping!";
                    setTimeout(() => document.title = "DCP Dashboard", 2000);
                    break;

                default:
                    break;
            }
        }
    }, [data]);

    const addLog = (msg) => {
        const time = new Date().toLocaleTimeString();
        setLogs(prev => [`[${time}] ${msg}`, ...prev].slice(0, 15));
    };

    const handlePing = (peer) => {
        sendMsg("send_ping", { 
            peer_id: peer.peer_id,
            ip: peer.ip,
            port: peer.port
        });
    };

    return (
        <div style={styles.container}>
            {/* Header */}
            <header style={styles.header}>
                <h1 style={styles.title}>DCP PROTOCOL <span style={styles.version}>v1.5</span></h1>
                <div style={{ ...styles.status, color: connected ? '#00ff41' : '#ff3b3b' }}>
                    {connected ? 'CONNECTED TO BACKEND' : 'BACKEND OFFLINE'}
                </div>
            </header>

            {/* Local Node Info */}
            <div style={styles.card}>
                <h2 style={styles.cardTitle}>LOCAL NODE INFO</h2>
                <div style={styles.infoGrid}>
                    <div><strong>ID:</strong> <span style={styles.idText}>{myInfo.id}</span></div>
                    <div><strong>PORT:</strong> <span style={styles.highlight}>{myInfo.port}</span></div>
                </div>
            </div>

            <div style={styles.mainGrid}>
                {/* Peer List */}
                <div style={styles.card}>
                    <h2 style={styles.cardTitle}>NETWORK PEERS ({peers.length})</h2>
                    <div style={styles.peerList}>
                        {peers.length === 0 && <p style={styles.muted}>No neighbors detected...</p>}
                        {peers.map(peer => (
                            <div key={peer.peer_id} style={styles.peerItem}>
                                <div>
                                    <div style={styles.peerIdText}>{peer.peer_id.substring(0, 12)}...</div>
                                    <div style={styles.peerSubText}>{peer.ip}:{peer.port}</div>
                                </div>
                                <button onClick={() => handlePing(peer)} style={styles.pingBtn}>
                                    PING
                                </button>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Console / Logs */}
                <div style={styles.card}>
                    <h2 style={styles.cardTitle}>SYSTEM LOGS</h2>
                    <div style={styles.console}>
                        {logs.map((log, i) => (
                            <div key={i} style={styles.logLine}>{log}</div>
                        ))}
                        {logs.length === 0 && <div style={styles.muted}>Awaiting network traffic...</div>}
                    </div>
                </div>
            </div>
        </div>
    );
}

// CSS-in-JS Styles
const styles = {
    container: { padding: '40px', backgroundColor: '#050505', color: '#e0e0e0', minHeight: '100vh', fontFamily: "'Fira Code', monospace" },
    header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '30px', borderBottom: '2px solid #1a1a1a', paddingBottom: '15px' },
    title: { margin: 0, fontSize: '22px', fontWeight: 'bold', color: '#00d4ff' },
    version: { fontSize: '10px', color: '#555' },
    status: { fontSize: '12px', fontWeight: 'bold' },
    card: { backgroundColor: '#0f0f12', padding: '20px', borderRadius: '4px', border: '1px solid #222', marginBottom: '20px' },
    cardTitle: { marginTop: 0, fontSize: '14px', color: '#888', marginBottom: '15px', borderLeft: '3px solid #00d4ff', paddingLeft: '10px' },
    infoGrid: { display: 'flex', gap: '40px', fontSize: '14px' },
    idText: { color: '#ffcc00' },
    highlight: { color: '#00d4ff' },
    mainGrid: { display: 'grid', gridTemplateColumns: '1.2fr 0.8fr', gap: '20px' },
    peerList: { maxHeight: '400px', overflowY: 'auto' },
    peerItem: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '15px', backgroundColor: '#16161a', marginBottom: '8px', border: '1px solid #1a1a1a' },
    peerIdText: { fontSize: '15px', fontWeight: 'bold' },
    peerSubText: { fontSize: '11px', color: '#666' },
    pingBtn: { backgroundColor: '#00d4ff', color: '#000', border: 'none', padding: '6px 12px', cursor: 'pointer', fontWeight: 'bold', fontSize: '12px' },
    console: { backgroundColor: '#000', padding: '15px', borderRadius: '4px', height: '350px', overflowY: 'auto', border: '1px solid #111' },
    logLine: { fontSize: '12px', color: '#00ff41', marginBottom: '4px', borderBottom: '1px solid #080808' },
    muted: { color: '#444', fontSize: '12px', fontStyle: 'italic' }
};

export default App;
