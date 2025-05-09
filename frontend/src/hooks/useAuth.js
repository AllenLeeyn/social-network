import { useState } from "react";
import { login } from "../lib/apiAuth";

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

  const handleLogout = () => {
    setUser(null);
    // Optionally, clear cookies or session storage here
  };

  return { user, loading, error, handleLogin, handleLogout };
}
