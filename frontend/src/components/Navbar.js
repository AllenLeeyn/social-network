"use client";
import Link from "next/link";
import Image from "next/image";
import "../styles/Navbar.css";
import { useAuth } from "../hooks/useAuth";
import { toast } from "react-toastify";
import { useRouter } from "next/navigation";
import { usePathname } from "next/navigation";
import { useState, useEffect } from "react";
import { FaUserCircle } from "react-icons/fa";
import { FaHouse } from "react-icons/fa6";
import { FaUserGroup } from "react-icons/fa6";
import { FaArrowRightFromBracket } from "react-icons/fa6";
import { FaMessage } from "react-icons/fa6";
import { FaBell } from "react-icons/fa6";
import { FaUsers } from "react-icons/fa6";
import NotificationBell from "./notifications/NotificationBell";
import { useNotifications } from "../contexts/NotificationsContext";
import { fetchUsers } from "../lib/apiAuth";

import "../styles/notifications/Bell.css";

export default function Navbar() {
  const { handleLogout } = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const [userName, setUserName] = useState("");
  const [profileImage, setProfileImage] = useState("");
  const [isAuthChecked, setIsAuthChecked] = useState(false);

  // Check authentication status on component mount
  useEffect(() => {
    async function checkAuthentication() {
      // Skip auth check for login page
      if (pathname === "/login") {
        setIsAuthChecked(true);
        return;
      }

      try {
        await fetchUsers(); // Use your existing API function
        setIsAuthChecked(true);
      } catch (error) {
        console.error(
          "Authentication check failed, redirecting to login:",
          error
        );
        router.push("/login");
      }
    }

    checkAuthentication();
  }, [pathname, router]);

  useEffect(() => {
    const storedUserName = localStorage.getItem("user-nick_name");
    if (storedUserName) {
      setUserName(storedUserName);
    }

    const storedProfileImage =
      localStorage.getItem("user-profile_image") || null;
    const imageUrl = storedProfileImage
      ? `/frontend-api/image/${storedProfileImage}`
      : null;
    setProfileImage(imageUrl);
  }, []);

  const { notifications } = useNotifications();

  const onLogoutClick = async (e) => {
    e.preventDefault();
    toast.success("logging out..");
    await handleLogout();
    router.push("/login");
  };

  // Don't render navbar on login page or while checking auth
  if (pathname === "/login" || !isAuthChecked) {
    return null;
  }

  return (
    <div className="navbar">
      <Link href="/" className="logo-title">
        <div className="logo-title">
          <Image
            src="/logo.png"
            alt="Site Logo"
            width={35}
            height={35}
            className="logo-img"
          />
          <span className="site-title">grit:hub</span>
        </div>
      </Link>

      <div className="center-links">
        <Link
          href="/"
          className={`nav-link ${pathname === "/" ? "active" : ""}`}
        >
          <FaHouse /> Home
        </Link>
        <Link
          href="/groups"
          className={`nav-link ${pathname === "/groups" ? "active" : ""}`}
        >
          <FaUserGroup /> Groups
        </Link>
        <Link
          href="/users"
          className={`nav-link ${pathname === "/users" ? "active" : ""}`}
        >
          <FaUsers /> Users
        </Link>
        <Link
          href="/messages"
          className={`nav-link ${pathname === "/messages" ? "active" : ""}`}
        >
          <FaMessage /> Messages
        </Link>
        <NotificationBell notifications={notifications} />
      </div>
      <div className="right-links">
        {profileImage ? (
          <Image
            src={profileImage}
            alt="Profile Image"
            width={30}
            height={30}
            style={{ borderRadius: "50%" }}
          />
        ) : (
          <FaUserCircle
            size={30}
            color="#aaa"
            style={{ verticalAlign: "middle" }}
          />
        )}
        <Link
          href="/profile"
          className={`nav-link ${pathname === "/profile" ? "active" : ""}`}
        >
          {userName || "Profile"}
        </Link>
        <a href="/" className="nav-link" onClick={onLogoutClick}>
          <FaArrowRightFromBracket />
        </a>
      </div>
    </div>
  );
}
