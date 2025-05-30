"use client";
import React from 'react';
import Link from 'next/link';
import Image from 'next/image';
import '../styles/Navbar.css';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHouse } from '@fortawesome/free-solid-svg-icons';

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
                <Link href="/" className="nav-link"><FontAwesomeIcon icon="fa-solid fa-house" /></Link>
                <Link href="/followers" className="nav-link">Groups</Link>
                <Link href="/messages" className="nav-link">Messages</Link>
                <Link href="/notification" className="nav-link">Notification</Link>
            </div>
            <div className="right-links">
                <Link href="/profile" className="nav-link">Profile</Link>
                <Link href="/" className="nav-link">Login/Logout</Link>
            </div>
        </div>
    );
}
