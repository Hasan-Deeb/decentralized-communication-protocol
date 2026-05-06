import { useEffect, useState, useRef } from 'react';

export const useWebSocket = (url) => {
    const [data, setData] = useState(null);
    const [connected, setConnected] = useState(false);
    const socketRef = useRef(null);

    const connect = () => {
        const socket = new WebSocket(url);

        socket.onopen = () => {
            console.log("✅ Connected to DCP Backend");
            setConnected(true);
        };

        socket.onmessage = (event) => {
            try {
                const parsedData = JSON.parse(event.data);
                setData(parsedData);
            } catch (err) {
                console.error("❌ WS Parse Error:", err);
            }
        };

        socket.onclose = () => {
            console.warn("⚠️ WS Connection Lost. Retrying in 3s...");
            setConnected(false);
            // إعادة الاتصال التلقائي
            setTimeout(connect, 3000);
        };

        socket.onerror = (err) => {
            console.error("❌ WS Socket Error:", err);
            socket.close();
        };

        socketRef.current = socket;
    };

    // دالة إرسال الرسائل من الواجهة إلى الباكند
    const sendMsg = (type, payload) => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            socketRef.current.send(JSON.stringify({ type, data: payload }));
        } else {
            console.error("❌ Cannot send: WebSocket is not connected");
        }
    };

    useEffect(() => {
        connect();
        return () => socketRef.current?.close();
    }, [url]);

    return { data, connected, sendMsg };
};
