"use client";

import { useState, createContext } from "react";

type Conn = WebSocket | null;

export const WebSocketContext = createContext<{
  conn: Conn;
  setConn: (c: Conn) => void;
}>({
  conn: null,
  setConn: () => {},
});

export default function WebSocketProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [conn, setConn] = useState<Conn>(null);

  return (
    <WebSocketContext.Provider
      value={{
        conn: conn,
        setConn: setConn,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
}
