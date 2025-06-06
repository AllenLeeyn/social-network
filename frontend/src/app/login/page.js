"use client";
import React, { useState, useEffect } from "react";
import { useAuth } from "../../hooks/useAuth";
import { useRouter } from "next/navigation";
import "./login.css";
import { toast } from "react-toastify";
import { handleImage } from "../../lib/handleImage";
import { fetchUsers } from "../../lib/apiAuth";

export default function AuthPage() {
  const { handleLogin, handleSignup, setError, error, loading } = useAuth();
  const [mode, setMode] = useState("login");
  const router = useRouter();

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
  const [registerGender, setGender] = useState("");
  const [registerVisibility, setVisibility] = useState("");
  const [formError, setFormError] = useState("");
  const [previewUrl, setPreviewUrl] = useState(null);
  const [isLoggingIn, setIsLoggingIn] = useState(false);

  // Check if user is already authenticated
  useEffect(() => {
    async function checkAuthentication() {
      try {
        await fetchUsers(); // Use the same function as Navbar
        // If successful (200), user is already authenticated
        router.push("/");
      } catch (error) {
        // If failed, user is not authenticated, stay on login page
        console.log("User not authenticated, staying on login page");
      }
    }

    checkAuthentication();
  }, [router]);

  // Handle login submission
  const handleLoginSubmit = async (e) => {
    e.preventDefault();

    if (isLoggingIn) return;
    setIsLoggingIn(true);

    try {
      const response = await handleLogin(loginEmail, loginPassword);
      toast.success("Login successful! Redirecting...");
      setTimeout(() => {
        window.location.href = "/";
      }, 1500);
    } catch (err) {
      console.error("Login failed:", err.message);
      setIsLoggingIn(false);
    }
  };

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (!file) return;
    setAvatar(e.target.files[0]);

    const url = URL.createObjectURL(file);
    setPreviewUrl(url);
  };

  // Handle registration submission
  const handleRegister = async (e) => {
    e.preventDefault();

    if (isLoggingIn) return;
    setIsLoggingIn(true);

    if (registerPassword !== registerConfirmPassword) {
      setIsLoggingIn(false);
      setFormError("Passwords do not match.");
      return;
    }
    setFormError("");

    let imageUUID = null;
    if (registerAvatar) {
      try {
        imageUUID = await handleImage([registerAvatar]);
      } catch (err) {
        setFormError("Image upload failed: " + err.message);
        setIsLoggingIn(false);
        return;
      }
    }

    const userData = {
      first_name: registerFirstName,
      last_name: registerLastName,
      nick_name: registerNickname,
      email: registerEmail,
      password: registerPassword,
      confirmPassword: registerConfirmPassword,
      birthday: new Date(registerDateOfBirth).toISOString(),
      about_me: registerAboutMe,
      gender: registerGender,
      visibility: registerVisibility,
      profile_image: imageUUID ? Object.values(imageUUID)[0] : null,
    };

    try {
      await handleSignup(userData);
      toast.success("Signup successful! Redirecting...");
      setTimeout(() => {
        window.location.href = "/";
      }, 1500);
    } catch (err) {
      console.error("Signup failed:", err.message);
      setIsLoggingIn(false);
    }
  };

  return (
    <div className={`auth-container ${mode}`}>
      {/* Left Side */}
      <div
        className="auth-left"
        style={{
          background: mode === "login" ? "#1e293bCC" : "transparent",
          color: "#fff",
          transition: "background 0.5s",
        }}
      >
        <div className="auth-left-content">
          {mode === "login" ? (
            <>
              <h2>Welcome to grit:hub!</h2>
              <p>Connect instantly. Log in to continue.</p>
              <button
                onClick={() => {
                  setMode("register");
                  setFormError("");
                  setError("");
                }}
                className="auth-toggle-btn"
              >
                New here? Register
              </button>
            </>
          ) : (
            <>
              <h2>Join grit:hub!</h2>
              <p>Sign up and start connecting now.</p>
              <button
                onClick={() => {
                  setMode("login");
                  setFormError("");
                  setError("");
                }}
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
                <label className="field-label">Birthday</label>
                <input
                  type="date"
                  value={registerDateOfBirth}
                  onChange={(e) => setDateOfBirth(e.target.value)}
                  max={new Date().toISOString().split("T")[0]}
                  required
                />
                {/* Gender Radio Group */}
                <div className="form-field">
                  <label className="field-label">Gender</label>
                  <div className="radio-group">
                    <label className="radio-option">
                      <input
                        type="radio"
                        name="gender"
                        value="Male"
                        checked={registerGender === "Male"}
                        onChange={(e) => setGender(e.target.value)}
                        required
                      />
                      <span className="radio-custom"></span>
                      Male
                    </label>
                    <label className="radio-option">
                      <input
                        type="radio"
                        name="gender"
                        value="Female"
                        checked={registerGender === "Female"}
                        onChange={(e) => setGender(e.target.value)}
                        required
                      />
                      <span className="radio-custom"></span>
                      Female
                    </label>
                  </div>
                </div>

                {/* Visibility Radio Group */}
                <div className="form-field">
                  <label className="field-label">Profile Visibility</label>
                  <div className="radio-group">
                    <label className="radio-option">
                      <input
                        type="radio"
                        name="visibility"
                        value="public"
                        checked={registerVisibility === "public"}
                        onChange={(e) => setVisibility(e.target.value)}
                        required
                      />
                      <span className="radio-custom"></span>
                      Public
                    </label>
                    <label className="radio-option">
                      <input
                        type="radio"
                        name="visibility"
                        value="private"
                        checked={registerVisibility === "private"}
                        onChange={(e) => setVisibility(e.target.value)}
                        required
                      />
                      <span className="radio-custom"></span>
                      Private
                    </label>
                  </div>
                </div>
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
                  onChange={handleFileChange}
                />
                {previewUrl && (
                  <img
                    src={previewUrl}
                    alt="Avatar Preview"
                    style={{ width: 100, height: 100, objectFit: "cover", marginTop: 10 }}
                  />
                )}
                <textarea
                  placeholder="About Me - Optional"
                  value={registerAboutMe}
                  onChange={(e) => setAboutMe(e.target.value)}
                />
              </fieldset>

              {error && <p style={{ color: "red" }}>{error}</p>}
              <button type="submit">Sign Up</button>
            </form>
          )}
        </div>
      </div>
    </div>
  );
}

// let isToggling = false;

// function smoothToggleAuthForm() {
//     if (isToggling) return;

//     isToggling = true;
//     const container = document.querySelector('.auth-form-container');

//     // Fade out
//     container.style.transition = 'all 1.5s cubic-bezier(0.23, 1, 0.32, 1)';
//     container.style.opacity = '0';
//     container.style.transform = 'translateY(-35px) scale(0.94)';
//     container.style.filter = 'blur(1px)';

//     setTimeout(() => {
//         // Your form switching logic here
//         // updateFormContent(); // Call your existing function

//         // Fade in
//         container.style.opacity = '1';
//         container.style.transform = 'translateY(0) scale(1)';
//         container.style.filter = 'blur(0px)';

//         setTimeout(() => {
//             isToggling = false;
//         }, 1800);

//     }, 1500);
// }

// // Attach to button
// document.querySelector('.auth-toggle-btn').addEventListener('click', smoothToggleAuthForm);
