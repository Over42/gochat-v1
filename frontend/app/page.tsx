"use client";

import { useState, useEffect, useContext } from "react";
import { useRouter } from "next/navigation";

import { AuthContext } from "@/context_providers/AuthContext";
import { WebSocketContext } from "@/context_providers/WebSocketContext";
import { API_URL, WEBSOCKET_URL } from "@/constants/constants";

export default function Home() {
  const [rooms, setRooms] = useState<{ id: string; name: string }[]>([]);
  const [roomName, setRoomName] = useState("");
  const { user } = useContext(AuthContext);
  const { setConn } = useContext(WebSocketContext);
  const router = useRouter();

  async function getRooms() {
    try {
      const res = await fetch(`${API_URL}/rooms`, {
        method: "GET",
      });

      const data = await res.json();
      if (res.ok) {
        setRooms(data);
      }
    } catch (err) {
      console.log(err);
    }
  }

  async function createRoom(e: React.SyntheticEvent) {
    e.preventDefault();
    try {
      const res = await fetch(`${API_URL}/rooms`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ name: roomName }),
      });

      if (res.ok) {
        getRooms();
      }
    } catch (err) {
      console.log(err);
    }
  }

  async function joinRoom(roomId: string) {
    const ws = new WebSocket(
      `${WEBSOCKET_URL}/rooms/${roomId}?userId=${user.id}&username=${user.username}`
    );
    if (ws.OPEN) {
      setConn(ws);
      router.push(`/rooms/${roomId}`);
      return;
    }
  }

  useEffect(() => {
    getRooms();
  }, []);

  return (
    <div className="my-8 px-4 md:mx-32 w-full h-full">
      <div className="flex justify-center mt-3 p-5">
        <input
          type="text"
          className="border border-grey p-2 rounded-md focus:outline-none focus:border-green"
          placeholder="Room name"
          value={roomName}
          onChange={(e) => setRoomName(e.target.value)}
        />
        <button
          className="bg-green border text-white rounded-md p-2 md:ml-4"
          onClick={createRoom}
        >
          Create
        </button>
      </div>
      <div className="mt-6">
        <div className="font-bold text-xl">Available Rooms:</div>
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mt-6">
          {rooms.map((room, index) => (
            <div
              key={index}
              className="border border-green p-4 flex items-center rounded-md w-full"
            >
              <div className="w-full">
                <div className="text-green font-bold text-lg">{room.name}</div>
              </div>
              <div className="">
                <button
                  className="px-4 text-white bg-green rounded-md"
                  onClick={() => joinRoom(room.id)}
                >
                  Join
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
