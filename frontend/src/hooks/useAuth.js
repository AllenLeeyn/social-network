import { useState } from "react"
import { login, logout, signup } from "../lib/apiAuth"
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

export function useAuth() {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const handleLogin = async (email, password) => {
    setLoading(true)
    setError(null)

    try {
      const userData = await login(email, password)
      setUser(userData)

      localStorage.setItem('user-uuid', userData.data.uuid);
      localStorage.setItem('user-nick_name', userData.data.nick_name);

    } catch (err) {
      setError(err.message)
      throw err;

    } finally {
      setLoading(false)
    }
  }

  const handleSignup = async (userData) => {
    setLoading(true)
    setError(null)

    try {
      const response = await signup(userData)
      setUser(response)

      localStorage.setItem('user-uuid', response.data.uuid);
      localStorage.setItem('user-nick_name', response.data.nick_name);

    } catch (err) {
      setError(err.message)
      throw err;

    } finally {
      setLoading(false)
    }
  }

  const handleLogout = async () => {
    await logout()
    setUser(null)
    localStorage.removeItem('activeDM');
    localStorage.removeItem('user-uuid');
    localStorage.removeItem('user-nick_name');
  }

  return { user, loading, error, handleLogin, handleSignup, handleLogout }
}
