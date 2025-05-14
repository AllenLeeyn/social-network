"use client";
import React from 'react';
import Link from 'next/link';

import Image from 'next/image';
import '../styles/Navbar.css';


export default function Navbar() {

return (
        <div className="navbar">
            <div className="logo-title">
                <Image
                src="/logo.png"
                alt="Site Logo"
                width={35}
                height={35}
                className="logo-img"
                />
                <span className="site-title">grit:Hub</span>
            </div>

            <div className="center-links">
                <Link href="/" className="nav-link">Home</Link>
                <Link href="/groups" className="nav-link">Connections</Link>
                <Link href="/messages" className="nav-link">Messages</Link>
                <Link href="/notifications" className="nav-link">Notifications</Link>
            </div>
            <div className="right-links">
                <Link href="/profile" className="nav-link">Profile</Link>
                <Link href="/" className="nav-link">Logout</Link>
            </div>
        </div>
    );
}
