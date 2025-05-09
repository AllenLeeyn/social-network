"use client";
import React, { useState } from 'react';
import './Login.css'; 

export default function AuthPage() {
    const [mode, setMode] = useState('login'); // 'login' or 'register'

    // Simple form state (replace with validation/handlers later)
    const [loginEmail, setLoginEmail] = useState('');
    const [loginPassword, setLoginPassword] = useState('');

    const [registerFirstName, setFirstName] = useState('');
    const [registerLastName, setLastName] = useState('');
    const [registerNickname, setNickname] = useState('');
    const [registerEmail, setRegisterEmail] = useState('');
    const [registerPassword, setRegisterPassword] = useState('');
    const [registerConfirmPassword, setConfirmPassword] = useState('')
    const [registerDateOfBirth, setDateOfBirth] = useState('');
    const [registerAvatar, setAvatar] = useState(null)
    const [registerAboutMe, setAboutMe] = useState('')
    const [error, setError] = useState('')

    // Handlers (stubbed for now)
    const handleLogin = (e) => {
        e.preventDefault();
        // Add validation/auth logic here
        alert(`Logging in as ${loginEmail}`);
    };

    const handleRegister = (e) => {
        e.preventDefault();

        if (registerPassword !== registerConfirmPassword) {
            setError("Password do not match.")
            return;
        }
        setError('')

        // Add validation/registration logic here

        // Change the below when implementation is done to something like
        // alert(`Registration successful! Welcome, ${registerFirstName}.`);
        alert(`Registering 
            (${registerFirstName})
            (${registerLastName})
            (${registerNickname})
            (${registerAvatar ? registerAvatar.name : 'No File Selected'})
            (${registerAboutMe})
            `);

    };

    return (
        <div className={`auth-container ${mode}`}>
            {/* Left Side */}
            <div className="auth-left" style={{
                background: mode === 'login' ? '#1e293b' : '#059669',
                color: '#fff',
                transition: 'background 0.5s'
            }}>
                <div className="auth-left-content">
                {mode === 'login' ? (
                    <>
                    <h2>Welcome to grit:Hub!</h2>
                    <p>Connect instantly. Log in to continue.</p>
                    <button onClick={() => setMode('register')} className="auth-toggle-btn">
                        New here? Register
                    </button>
                    </>
                ) : (
                    <>
                    <h2>Join grit:Hub!</h2>
                    <p>Sign up and start connecting now.</p>
                    <button onClick={() => setMode('login')} className="auth-toggle-btn">
                        Already have an account? Log In
                    </button>
                    </>
                )}
                </div>
            </div>

            {/* Right Side */}
            <div className="auth-right">
                <div className="auth-form-container">
                {mode === 'login' ? (
                    <form onSubmit={handleLogin} className="auth-form">
                    <h3>Login</h3>
                    <input
                        type="email"
                        placeholder="Email"
                        value={loginEmail}
                        onChange={e => setLoginEmail(e.target.value)}
                        required
                    />
                    <input
                        type="password"
                        placeholder="Password"
                        value={loginPassword}
                        onChange={e => setLoginPassword(e.target.value)}
                        required
                    />
                    <button type="submit">Log In</button>
                    </form>
                ) : (
                    <form onSubmit={handleRegister} className="auth-form">
                        {error && <div className="error-message">{error}</div>}
                        <h3>Register</h3>
                        <fieldset className="auth-form">
                            <legend>Required Information</legend>
                            <input
                                type="text"
                                placeholder="First Name"
                                value={registerFirstName}
                                onChange={e => setFirstName(e.target.value)}
                                required
                            />
                            <input
                                type="text"
                                placeholder="Last Name"
                                value={registerLastName}
                                onChange={e => setLastName(e.target.value)}
                                required
                            />
                            <input
                                type="email"
                                placeholder="Email"
                                value={registerEmail}
                                onChange={e => setRegisterEmail(e.target.value)}
                                required
                            />
                            <input
                                type="password"
                                placeholder="Password"
                                value={registerPassword}
                                onChange={e => setRegisterPassword(e.target.value)}
                                required
                            />
                            <input
                                type="password"
                                placeholder="Confirmed Password"
                                value={registerConfirmPassword}
                                onChange={e => setConfirmPassword(e.target.value)}
                                required
                            />
                            <input
                                type="date"
                                placeholder=""
                                value={registerDateOfBirth}
                                onChange={e => setDateOfBirth(e.target.value)}
                                required
                            />
                        </fieldset>

                        <fieldset className="auth-form">
                            <legend>Optional Information</legend>
                            <input
                                type="text"
                                placeholder="Nickname - Optional"
                                value={registerNickname}
                                onChange={e => setNickname(e.target.value)}
                            />
                            <input
                                type="file"
                                accept="image/*"
                                onChange={e => setAvatar(e.target.files[0])}
                            />
                            <textarea
                                placeholder="About Me - Optional"
                                value={registerAboutMe}
                                onChange={e => setAboutMe(e.target.value)}
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
