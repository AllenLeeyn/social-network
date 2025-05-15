"use client"
import React from "react"
import Link from "next/link"
import Image from "next/image"
import "../styles/Navbar.css"
import { logout } from "../lib/apiAuth" // Import logout

export default function Navbar() {
  // Logout handler
  const handleLogout = async (e) => {
    e.preventDefault()
    await logout()
    window.location.href = "/"
  }

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
          Connections
        </Link>
        <Link href="/messages" className="nav-link">
          Messages
        </Link>
        <Link href="/notifications" className="nav-link">
          Notifications
        </Link>
      </div>
      <div className="right-links">
        <Link href="/profile" className="nav-link">
          Profile
        </Link>
        <a href="/" className="nav-link" onClick={handleLogout}>
          Logout
        </a>
      </div>
    </div>
  )
}
