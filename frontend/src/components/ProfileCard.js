"use client";

import { useEffect, useState } from "react";
import { submitFollowRequest, submitUnfollowRequest } from "../lib/apiFollow";
import { toast } from "react-toastify";

export default function ProfileCard({ uuid, setPrivateProfile }) {
  const [user, setUser] = useState(null);
  const [followers, setFollowers] = useState([]);
  const [followings, setFollowings] = useState([]);
  const [isPrivateProfile, setIsPrivateProfile] = useState(false);
  const [followStatus, setFollowStatus] = useState(null);
  const [isOwnProfile, setIsOwnProfile] = useState(false);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    async function fetchProfile() {
      try {
        const url = uuid
          ? `/frontend-api/profile?uuid=${uuid}`
          : "/frontend-api/profile";
        const res = await fetch(url, {
          credentials: "include",
        });

        if (res.status === 403) {
          // Profile is private - we have the UUID from the URL
          console.log("Private profile detected, UUID from URL:", uuid);

          // Create minimal user object with the UUID from URL
          setUser({
            uuid: uuid,
            nick_name: "Private User",
            first_name: "",
            last_name: "",
            visibility: "private",
          });

          setPrivateProfile && setPrivateProfile(true);
          setIsPrivateProfile(true);
          setFollowStatus("none");
          return;
        }

        if (!res.ok) {
          throw new Error("Failed to fetch user data");
        }

        const data = await res.json();
        setUser(data.data);
        setPrivateProfile && setPrivateProfile(false);

        if (!uuid) {
          setIsOwnProfile(true);
        } else {
          setFollowStatus("none");
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

  if (!user) return <div>Loading profile...</div>;

  if (isPrivateProfile) {
    return (
      <div className="profile-header">
        <img
          src={
            user.profile_image
              ? `/frontend-api/image/${user.profile_image}`
              : ""
          }
          alt={user.nick_name || "Private User"}
          className="profile-avatar"
        />
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
      <img
        src={
          user.profile_image ? `/frontend-api/image/${user.profile_image}` : ""
        }
        alt={user.nick_name || "User"}
        className="profile-avatar"
      />
      <div className="profile-info">
        <h2>{user.nick_name || "Unknown User"}</h2>
        <p>
          <strong>
            {user.first_name || ""} {user.last_name || ""}
          </strong>
        </p>
        <p>{user.email || ""}</p>
        <p>{user.gender || ""}</p>
        <p>
          <span>Date of Birth:</span> {user.birthday || "Not specified"}
        </p>
        <p>{user.about_me || "Nothing is written."}</p>
        <p>Visibility: {user.visibility || "Unknown"}</p>
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

        {isOwnProfile && (
          <div className="connection-buttons">
            <div className="followers-list">
              <h3>Followers</h3>
              {followers.length === 0 ? (
                <p>No followers yet.</p>
              ) : (
                <ul>
                  {followers.map((entry) => (
                    <li key={entry.follower_uuid}>{entry.follower_name}</li>
                  ))}
                </ul>
              )}
            </div>
            <div className="following-list">
              <h3>Following</h3>
              {followings.length === 0 ? (
                <p>Not following anyone.</p>
              ) : (
                <ul>
                  {followings.map((entry) => (
                    <li key={entry.leader_uuid}>{entry.leader_name}</li>
                  ))}
                </ul>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
