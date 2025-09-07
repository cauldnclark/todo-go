import { useEffect } from "react";

const API_URL = import.meta.env.VITE_API_URL;

export function GoogleCallback() {
  useEffect(() => {
    async function handleGoogleCallback() {
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get("code");
      const state = urlParams.get("state");
      const error = urlParams.get("error");

      if (error) {
        setTimeout(() => {
          window.location.href = "/login";
        }, 2000);
        alert("Google login failed");
        return;
      }

      if (!code) {
        setTimeout(() => {
          window.location.href = "/login";
        }, 2000);
        alert("Google login failed");
        return;
      }

      const response = await fetch(`${API_URL}/auth/google`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          code,
          state,
        }),
      });

      const data = await response.json();

      console.log(data);

      if (data.token) {
        localStorage.setItem("todo-auth-token", data.token);
      }

      if (data.user) {
        localStorage.setItem("todo-auth-user", JSON.stringify(data.user));
      }

      setTimeout(() => {
        window.location.href = "/";
      }, 800);
    }

    handleGoogleCallback();
  }, []);

  return null;
}
