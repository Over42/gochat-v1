"use client";

import { useState, useContext, useEffect } from "react";
import { useRouter } from "next/navigation";

import { API_URL } from "@/constants/constants";
import { UserInfo, AuthContext } from "@/context_providers/AuthContext";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const { authenticated } = useContext(AuthContext);
  const router = useRouter();

  useEffect(() => {
    if (authenticated) {
      router.push("/");
      return;
    }
  }, [authenticated]);

  async function loginHandler(e: React.SyntheticEvent) {
    e.preventDefault();
    try {
      const res = await fetch(`${API_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      const data = await res.json();
      if (res.ok) {
        const user: UserInfo = {
          username: data.username,
          id: data.id,
        };

        localStorage.setItem("user_info", JSON.stringify(user));
        return router.push("/");
      }
    } catch (err) {
      console.log(err);
    }
  }

  return (
    <div className="flex items-center justify-center min-w-full min-h-screen">
      <form className="flex flex-col md:w-1/5">
        <input
          type="email"
          placeholder="email"
          className="p-3 mt-8 rounded-md border-2 border-grey focus:outline-none focus:border-green"
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          type="password"
          placeholder="password"
          className="p-3 mt-8 rounded-md border-2 border-grey focus:outline-none focus:border-green"
          onChange={(e) => setPassword(e.target.value)}
        />
        <button
          className="p-3 mt-6 rounded-md bg-green font-bold text-white"
          type="submit"
          onClick={loginHandler}
        >
          Login
        </button>
      </form>
    </div>
  );
}
