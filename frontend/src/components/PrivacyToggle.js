"use client";

import React from "react";

export default function PrivacyToggle({ isPrivate, setIsPrivate }) {
    return (
        <div className="privacy-toggle">
            <span className={`account-status ${isPrivate ? "private" : "public"}`}>
            </span>
            <button onClick={() => setIsPrivate(!isPrivate)} className="toggle-privacy-btn">
                Switch to {isPrivate ? "Public" : "Private"}
            </button>
        </div>
    );
}

