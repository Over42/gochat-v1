"use client";

import { useState, useRef, useContext, useEffect } from "react";
import { useRouter } from "next/navigation";
import autosize from "autosize";

import { ChatBody, Message } from "./ChatBody";
import { WebSocketContext } from "@/context_providers/WebSocketContext";
import { AuthContext } from "@/context_providers/AuthContext";
import { API_URL } from "@/constants/constants";

export default function Room() {
  const [messages, setMessages] = useState<Array<Message>>([]);
  const textarea = useRef<HTMLTextAreaElement>(null);
  const { conn } = useContext(WebSocketContext);
  const [users, setUsers] = useState<Array<{ username: string }>>([]);
  const { user } = useContext(AuthContext);
  const router = useRouter();

  function sendMessage() {
    if (!textarea.current?.value) return;
    if (conn === null) {
      router.push("/");
      return;
    }

    conn.send(textarea.current.value);
    textarea.current.value = "";
  }

  // Get clients in the room
  useEffect(() => {
    if (conn === null) {
      router.push("/");
      return;
    }

    const roomId = conn.url.substring(
      conn.url.lastIndexOf("/") + 1,
      conn.url.indexOf("?")
    );

    async function getUsers() {
      try {
        const res = await fetch(`${API_URL}/rooms/${roomId}/clients`, {
          method: "GET",
          headers: { "Content-Type": "application/json" },
        });
        const data = await res.json();
        setUsers(data);
      } catch (err) {
        console.error(err);
      }
    }

    getUsers();
  }, []);

  // Handle websocket connection
  useEffect(() => {
    if (textarea.current) {
      autosize(textarea.current);
    }

    if (conn === null) {
      router.push("/");
      return;
    }

    conn.onmessage = (message) => {
      const msg: Message = JSON.parse(message.data);
      if (msg.content == "New user has joined") {
        setUsers([...users, { username: msg.username }]);
      }

      if (msg.content == "User left the chat") {
        const remainingUsers = users.filter(
          (user) => user.username != msg.username
        );
        setUsers([...remainingUsers]);
        setMessages([...messages, msg]);
        return;
      }

      if (user.username == msg.username) {
        msg.type = "self";
      } else {
        msg.type = "recv";
      }

      setMessages([...messages, msg]);
    };

    conn.onclose = () => {};
    conn.onerror = () => {};
    conn.onopen = () => {};
  }, [textarea, conn, messages, users]);

  return (
    <div className="flex flex-col w-full">
      <div className="p-4 mb-14">
        <ChatBody messages={messages} />
      </div>
      <div className="fixed bottom-0 mt-4 w-full">
        <div className="flex md:flex-row px-4 py-2 bg-grey rounded-md">
          <div className="flex w-full mr-4 rounded-md border border-green">
            <textarea
              ref={textarea}
              placeholder="Message"
              className="w-full h-10 p-2 rounded-md focus:outline-none"
              style={{ resize: "none" }}
            />
          </div>
          <div className="flex items-center">
            <button
              className="p-2 rounded-md bg-green text-white"
              onClick={sendMessage}
            >
              Send
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
