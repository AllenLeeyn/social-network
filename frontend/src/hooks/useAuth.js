import { useState } from "react";
import { login, logout, signup } from "../lib/apiAuth";

export function useAuth() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleLogin = async (nickName, email, password) => {
    setLoading(true);
    setError(null);

    try {
      const userData = await login(nickName, email, password);
      setUser(userData); // Save user data after successful login
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleSignup = async (userData) => {
    setLoading(true);
    setError(null);

    try {
      const response = await signup(userData);
      alert("Signup successful! You can now log in.");
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    await logout()
    setUser(null);
    // Optionally, clear cookies or session storage here
  };

  return { user, loading, error, handleLogin, handleSignup, handleLogout };
}