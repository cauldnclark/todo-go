import { useState, useEffect, useCallback, useRef } from "react";

export interface WebSocketMessage {
  type: string;
  content: string;
  user_id: number;
}

const useWebsocket = (
  url: string,
  token: string,
  onMessage: (event: MessageEvent<string>) => void,
  onOpen: (event: Event) => void,
  onClose: (event: CloseEvent) => void
) => {
  const [readyState, setReadyState] = useState(WebSocket?.CONNECTING || 0);
  const wsRef = useRef<WebSocket>(null);

  const connect = useCallback(() => {
    if (!token) {
      console.warn("No token provided");
      return;
    }

    const wsUrl = `${url}?token=${token}`;
    const ws = new WebSocket(wsUrl);

    wsRef.current = ws;

    ws.onopen = (event) => {
      console.log("WebSocket connected");
      onOpen(event);
      setReadyState(WebSocket.OPEN);
      
      // Send initial heartbeat message
      ws.send(JSON.stringify({ event: "heartbeat", data: {} }));
    };

    ws.onmessage = (event) => {
      console.log("Received message:", event);
      onMessage(event);
    };

    ws.onclose = (event) => {
      console.log("WebSocket closed:", event.code, event.reason);
      onClose(event);
      setReadyState(WebSocket.CLOSED);
    };

    ws.onerror = (event) => {
      console.error("WebSocket error:", event);
    };

    return () => {
      ws.close();
    };
  }, [url, token]);

  useEffect(() => {
    connect();

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [connect]);

  // Update event handlers when callbacks change
  useEffect(() => {
    if (wsRef.current) {
      wsRef.current.onmessage = (event) => {
        console.log("Received message:", event);
        onMessage(event);
      };
      wsRef.current.onopen = (event) => {
        console.log("WebSocket connected");
        onOpen(event);
        setReadyState(WebSocket.OPEN);
      };
      wsRef.current.onclose = (event) => {
        console.log("WebSocket closed:", event.code, event.reason);
        onClose(event);
        setReadyState(WebSocket.CLOSED);
      };
    }
  }, [onMessage, onOpen, onClose]);

  const sendMessage = useCallback((message: WebSocketMessage) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    } else {
      console.warn("WebSocket not open. Cannot send message.");
    }
  }, []);

  return { readyState, sendMessage };
};

export default useWebsocket;
