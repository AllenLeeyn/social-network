"use client";
import React, { useState } from "react";
import { useAuth } from "../../hooks/useAuth"; // Import the useAuth hook
import "./login.css";

export default function AuthPage() {
    const { handleLogin, error, loading } = useAuth(); // Use the hook for login logic
    const [mode, setMode] = useState("login"); // 'login' or 'register'

    // Login form state
    const [loginEmail, setLoginEmail] = useState("");
    const [loginPassword, setLoginPassword] = useState("");

    // Registration form state
    const [registerFirstName, setFirstName] = useState("");
    const [registerLastName, setLastName] = useState("");
    const [registerNickname, setNickname] = useState("");
    const [registerEmail, setRegisterEmail] = useState("");
    const [registerPassword, setRegisterPassword] = useState("");
    const [registerConfirmPassword, setConfirmPassword] = useState("");
    const [registerDateOfBirth, setDateOfBirth] = useState("");
    const [registerAvatar, setAvatar] = useState(null);
    const [registerAboutMe, setAboutMe] = useState("");
    const [formError, setFormError] = useState("");

    // Handle login submission
    const handleLoginSubmit = async (e) => {
        e.preventDefault();
        try {
            await handleLogin("nickname", loginEmail, loginPassword); // Pass the credentials
            alert("Login successful! Redirecting...");
            window.location.href = "/posts"; // Redirect to posts page
        } catch (err) {
            console.error("Login failed:", err.message);
        }
    };

    // Handle registration submission (placeholder logic)
    const handleRegister = (e) => {
        e.preventDefault();

        if (registerPassword !== registerConfirmPassword) {
            setFormError("Passwords do not match.");
            return;
        }
        setFormError("");

        alert(`Registering:
            First Name: ${registerFirstName}
            Last Name: ${registerLastName}
            Nickname: ${registerNickname}
            Avatar: ${registerAvatar ? registerAvatar.name : "No File Selected"}
            About Me: ${registerAboutMe}`);
    };

    return (
        <div className={`auth-container ${mode}`}>
            {/* Left Side */}
            <div
                className="auth-left"
                style={{
                    background: mode === "login" ? "#1e293b" : "#059669",
                    color: "#fff",
                    transition: "background 0.5s",
                }}
            >
                <div className="auth-left-content">
                    {mode === "login" ? (
                        <>
                            <h2>Welcome to grit:Hub!</h2>
                            <p>Connect instantly. Log in to continue.</p>
                            <button
                                onClick={() => setMode("register")}
                                className="auth-toggle-btn"
                            >
                                New here? Register
                            </button>
                        </>
                    ) : (
                        <>
                            <h2>Join grit:Hub!</h2>
                            <p>Sign up and start connecting now.</p>
                            <button
                                onClick={() => setMode("login")}
                                className="auth-toggle-btn"
                            >
                                Already have an account? Log In
                            </button>
                        </>
                    )}
                </div>
            </div>

            {/* Right Side */}
            <div className="auth-right">
                <div className="auth-form-container">
                    {mode === "login" ? (
                        <form onSubmit={handleLoginSubmit} className="auth-form">
                            <h3>Login</h3>
                            <input
                                type="email"
                                placeholder="Email"
                                value={loginEmail}
                                onChange={(e) => setLoginEmail(e.target.value)}
                                required
                            />
                            <input
                                type="password"
                                placeholder="Password"
                                value={loginPassword}
                                onChange={(e) => setLoginPassword(e.target.value)}
                                required
                            />
                            {error && <p style={{ color: "red" }}>{error}</p>}
                            <button type="submit" disabled={loading}>
                                {loading ? "Logging in..." : "Log In"}
                            </button>
                        </form>
                    ) : (
                        <form onSubmit={handleRegister} className="auth-form">
                            {formError && <div className="error-message">{formError}</div>}
                            <h3>Register</h3>
                            <fieldset className="auth-form">
                                <legend>Required Information</legend>
                                <input
                                    type="text"
                                    placeholder="First Name"
                                    value={registerFirstName}
                                    onChange={(e) => setFirstName(e.target.value)}
                                    required
                                />
                                <input
                                    type="text"
                                    placeholder="Last Name"
                                    value={registerLastName}
                                    onChange={(e) => setLastName(e.target.value)}
                                    required
                                />
                                <input
                                    type="email"
                                    placeholder="Email"
                                    value={registerEmail}
                                    onChange={(e) => setRegisterEmail(e.target.value)}
                                    required
                                />
                                <input
                                    type="password"
                                    placeholder="Password"
                                    value={registerPassword}
                                    onChange={(e) => setRegisterPassword(e.target.value)}
                                    required
                                />
                                <input
                                    type="password"
                                    placeholder="Confirm Password"
                                    value={registerConfirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    required
                                />
                                <input
                                    type="date"
                                    value={registerDateOfBirth}
                                    onChange={(e) => setDateOfBirth(e.target.value)}
                                    required
                                />
                            </fieldset>

                            <fieldset className="auth-form">
                                <legend>Optional Information</legend>
                                <input
                                    type="text"
                                    placeholder="Nickname - Optional"
                                    value={registerNickname}
                                    onChange={(e) => setNickname(e.target.value)}
                                />
                                <input
                                    type="file"
                                    accept="image/*"
                                    onChange={(e) => setAvatar(e.target.files[0])}
                                />
                                <textarea
                                    placeholder="About Me - Optional"
                                    value={registerAboutMe}
                                    onChange={(e) => setAboutMe(e.target.value)}
                                />
                            </fieldset>

                            <button type="submit">Sign Up</button>
                        </form>
                    )}
                </div>
            </div>
        </div>
    );
}