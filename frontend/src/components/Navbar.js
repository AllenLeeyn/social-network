"use client"
import Link from "next/link"
import Image from "next/image"
import "../styles/Navbar.css"
import { useAuth } from "../hooks/useAuth";
import { toast } from 'react-toastify';
import NotificationBell from './notifications/NotificationBell';
import { useNotifications } from "../contexts/NotificationsContext";

export default function Navbar() {
  // Logout handler
  const { handleLogout } = useAuth();
    const { notifications } = useNotifications();

  const onLogoutClick = async (e) => {
    e.preventDefault();
    toast.success("logging out..");
    await handleLogout();
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
        {/*<Link href="/notifications" className="nav-link">
          Notifications
        </Link>*/}
        <NotificationBell notifications={notifications} />
      </div>
      <div className="right-links">
        <Link href="/profile" className="nav-link">
          Profile
        </Link>
        <a href="/" className="nav-link" onClick={onLogoutClick}>
          Logout
        </a>
      </div>
    </div>
  )
}
