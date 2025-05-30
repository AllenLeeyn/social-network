<<<<<<< HEAD
"use client"
import Link from "next/link"
import Image from "next/image"
import "../styles/Navbar.css"
import { useAuth } from "../hooks/useAuth";
import { toast } from 'react-toastify';
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { FaUserCircle } from 'react-icons/fa';
=======
"use client";
import React from 'react';
import Link from 'next/link';
import Image from 'next/image';
import '../styles/Navbar.css';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHouse } from '@fortawesome/free-solid-svg-icons';
>>>>>>> gigi

export default function Navbar() {
  // Logout handler
  const { handleLogout } = useAuth();
  const router = useRouter(); 

  const [userName, setUserName] = useState("");
  const [profileImage, setProfileImage] = useState("");

<<<<<<< HEAD
  useEffect(() => {
    const storedUserName = localStorage.getItem("user-nick_name");
    if (storedUserName) {
      setUserName(storedUserName);
    }

    const storedProfileImage = localStorage.getItem("user-profile_image") || null;
    const imageUrl = storedProfileImage ? `/frontend-api/image/${storedProfileImage}` : null;
    setProfileImage(imageUrl);
  }, []);

  const onLogoutClick = async (e) => {
    e.preventDefault();
    toast.success("logging out..");
    await handleLogout();
    router.push("/login");
  };

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
        <Link href="/" className="nav-link">
          Home
        </Link>
        <Link href="/groups" className="nav-link">
          Groups
        </Link>
        <Link href="/messages" className="nav-link">
          Messages
        </Link>
        <Link href="/notifications" className="nav-link">
          Notifications
        </Link>
      </div>
      <div className="right-links">
        {profileImage ? (
          <Image
            src={profileImage}
            alt="Profile Image"
            width={30}
            height={30}
            style={{ borderRadius: '50%' }}
          />
        ) : (
          <FaUserCircle size={30} color="#aaa" style={{ verticalAlign: 'middle' }} />
        )}
        <Link href="/profile" className="nav-link">
          {userName || "Profile"}
        </Link>
        <a href="/" className="nav-link" onClick={onLogoutClick}>
          Logout
        </a>
      </div>
    </div>
  )
=======
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
>>>>>>> gigi
}
