"use client";

import { useEffect, useState } from "react";
import { submitFollowRequest, submitUnfollowRequest } from "../lib/apiFollow";
import { toast } from "react-toastify";
import { FaUserCircle } from 'react-icons/fa';
import Image from 'next/image';

export default function ProfileCard({ uuid, setPrivateProfile }) {
    const [user, setUser] = useState(null);
    const [isPrivateProfile, setIsPrivateProfile] = useState(false);
    const [followStatus, setFollowStatus] = useState(null);
    const [isOwnProfile, setIsOwnProfile] = useState(false);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
    async function fetchProfile() {
        try {
        const url = uuid
            ? `/frontend-api/profile/${uuid}`
            : "/frontend-api/profile";
        const res = await fetch(url, {
            credentials: "include",
        });

        if (res.status === 403) {
            const response = await res.json();

            const resData = response.message.split("|");
            // Profile is private - we have the UUID from the URL
            console.log("Private profile detected, UUID from URL:", uuid);

          // Create minimal user object with the UUID from URL
            setUser({
            uuid: uuid,
            nick_name: resData[0] || "Private User",
            first_name: "",
            last_name: "",
            visibility: "private",
            });

            setPrivateProfile && setPrivateProfile(true);
            setIsPrivateProfile(true);
            setFollowStatus(resData[1] || "none");
            return;
        }

        if (!res.ok) {
            throw new Error("Failed to fetch user data");
        }

        const data = await res.json();
        setUser(data.data);
        setPrivateProfile && setPrivateProfile(false);

        console.log(data.data)
        if (!uuid) {
            setIsOwnProfile(true);
        } else {
            setFollowStatus(data.data.status);
        }
        } catch (err) {
        console.error("Error fetching profile:", err);
        // Even on error, we have the UUID from the URL
        if (uuid) {
            setUser({ uuid: uuid, nick_name: "Unknown User" });
            }
        }
    }
    fetchProfile();
    }, [uuid, setPrivateProfile]);

    const handleFollowAction = async () => {
    if (!uuid || isOwnProfile) return; // Use uuid directly since it's always available

    setLoading(true);
    try {
        if (followStatus === "accepted") {
        await submitUnfollowRequest({ leader_uuid: uuid });
        setFollowStatus("none");
        toast.success("Unfollowed successfully");
        } else if (followStatus === "requested") {
        await submitUnfollowRequest({ leader_uuid: uuid });
        setFollowStatus("none");
        toast.success("Follow request cancelled");
        } else {
        await submitFollowRequest({ leader_uuid: uuid });
        // For private profiles, always set to "requested"
        setFollowStatus(
            isPrivateProfile
            ? "requested"
            : user?.visibility === "public"
            ? "accepted"
            : "requested"
        );
        toast.success(
            isPrivateProfile || user?.visibility !== "public"
            ? "Follow request sent"
            : "Now following"
        );
        }
    } catch (err) {
        toast.error(err.message || "Failed to perform action");
    } finally {
        setLoading(false);
    }
    };

    const getFollowButtonText = () => {
    if (loading) return "Loading...";
    switch (followStatus) {
        case "accepted":
        return "Unfollow";
        case "requested":
        return "Cancel Request";
        case "declined":
        case "inactive":
        case "none":
        default:
        return "Follow";
        }
    };

    const getFollowButtonClass = () => {
    switch (followStatus) {
        case "accepted":
        return "follow-btn following";
        case "requested":
        return "follow-btn pending";
        default:
        return "follow-btn";
    }
    };

    async function handleVisibilityChange(e, user, setUser) {
    const newVisibility = e.target.value;
    const prevVisibility = user.visibility;
    
    const updatedUser = { ...user, visibility: newVisibility };
    setUser(updatedUser);

    try {
        const res = await fetch(`/frontend-api/users/update`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(updatedUser),
        });

        if (!res.ok) {
        throw new Error("Failed to update visibility");
        }

        const data = await res.json();
        toast.success("Visibility successfully updated!");
    } catch (err) {
        setUser({ ...user, visibility: prevVisibility });
        toast.error("Error updating visibility");
    }
    }

    if (!user) return <div>Loading profile...</div>;

    if (isPrivateProfile) {
    return (
        <div className="profile-header">
        <FaUserCircle size={120} color="#aaa" />
        <div className="profile-info">
            <h2>{user.nick_name || "Private User"}</h2>
            <p>
            <strong>
                {user.first_name || ""} {user.last_name || ""}
            </strong>
            </p>
            <p>This profile is private.</p>
            {!isOwnProfile && uuid && (
            <button
                className={getFollowButtonClass()}
                onClick={handleFollowAction}
                disabled={loading}
            >
                {getFollowButtonText()}
            </button>
            )}
        </div>
        </div>
    );
    }

    return (
    <div className="profile-header">
        {user?.profile_image ? (
        <Image
            src={`/frontend-api/image/${user.profile_image}`}
            alt="User Avatar"
            width={120}
            height={120}
            className="radial-avatar"
        />
        ) : (
        <FaUserCircle size={120} color="#aaa" />
        )}
        <div className="profile-info">
        <h2>{user.nick_name || "Unknown User"}</h2>
        <p>
            <strong>
            {user.first_name || ""} {user.last_name || ""}
            </strong>
        </p>
        <p><span>Email: </span>{user.email || ""}</p>
        <p><span>Gender: </span>{user.gender || ""}</p>
        <p>
        <span>Date of Birth: </span>{" "}
        {user.birthday
            ? new Date(user.birthday).toLocaleDateString("en-US", {
                year: "numeric",
                month: "long",
                day: "numeric",
            })
            : "Not specified"}
        </p>
        <p><span>Bio: </span>{user.about_me || "Nothing is written."}</p>
        {isOwnProfile ? (
            <label>
            <span>Visibility: </span>
            <select
                value={user.visibility || "private"}
                onChange={(e) => handleVisibilityChange(e, user, setUser)}
            >
                <option value="public">Public</option>
                <option value="private">Private</option>
            </select>
            </label>
        ) : (
            <p>Visibility: {user.visibility || "Unknown"}</p>
        )}
        <p className="bio">{user.bio || ""}</p>

        {!isOwnProfile && uuid && (
            <button
            className={getFollowButtonClass()}
            onClick={handleFollowAction}
            disabled={loading}
            >
            {getFollowButtonText()}
            </button>
        )}
        </div>
    </div>
    );
}
