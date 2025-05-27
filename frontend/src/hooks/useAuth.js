import { useState } from "react"
import { login, logout, signup } from "../lib/apiAuth"

export function useAuth() {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const handleLogin = async (email, password) => {
    setLoading(true)
    setError(null)

    try {
      const userData = await login(email, password)
      console.log('Login response', userData)
      setUser(userData) // Save user data after successful login

      // localStorage.setItem('session-id', userData.sessionId);

      // Store UUID and sessionId in localStorage
      const uuid = userData?.data?.uuid;
      const sessionId = userData?.data?.sessionId;
      console.log('uuid', uuid)
      console.log(sessionId)
      if (uuid)      localStorage.setItem('user-uuid', uuid);
      if (sessionId) localStorage.setItem('session-id', sessionId);

      return userData.sessionId;
      
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleSignup = async (userData) => {
    setLoading(true)
    setError(null)

    try {
      const response = await signup(userData)
      alert("Signup successful! You can now log in.")
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = async () => {
    await logout()
    setUser(null)
    localStorage.removeItem('activeDM');
    localStorage.removeItem('session-id');
    localStorage.removeItem('user-uuid');
  }

  return { user, loading, error, handleLogin, handleSignup, handleLogout }
}
