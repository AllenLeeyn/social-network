"use client";

import { useEffect, useState } from "react";

export default function ProfileCard({ uuid, setPrivateProfile }) {
  const [user, setUser] = useState(null);
  const [followers, setFollowers] = useState([]);
  const [followings, setFollowings] = useState([]);
  const [isPrivateProfile, setIsPrivateProfile] = useState(false);

  useEffect(() => {
    async function fetchProfile() {
      try {
        // Use the uuid prop if provided, otherwise fetch current user
        const url = uuid
          ? `/frontend-api/profile?uuid=${uuid}`
          : "/frontend-api/profile";
        const res = await fetch(url, {
          credentials: "include",
        });
        if (res.status === 403) {
          // Profile is private, but try to get minimal info if possible
          const data = await res.json();
          setUser(data.data || {}); // backend should return at least nick_name, profile_image, etc.
          setPrivateProfile && setPrivateProfile(true);
          return;
        }
        if (!res.ok) {
          throw new Error("Failed to fetch user data");
        }
        const data = await res.json();
        setUser(data.data);
        setPrivateProfile && setPrivateProfile(false);
      } catch (err) {
        console.error(err);
      }
    }
    fetchProfile();
  }, [uuid]);

  if (!user) return <div>Loading profile...</div>;

  if (isPrivateProfile) {
    return (
      <div className="profile-header">
        <img
          src={`/frontend-api/image/${user.profile_image}`}
          alt={user.nick_name}
          className="profile-avatar"
        />
        <div className="profile-info">
          <h2>{user.nick_name}</h2>
          <p>
            <strong>
              {user.first_name} {user.last_name}
            </strong>
          </p>
          <p>This profile is private.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="profile-header">
      <img
        src={`/frontend-api/image/${user.profile_image}`}
        alt={user.nick_name}
        className="profile-avatar"
      />
      <div className="profile-info">
        <h2>{user.nick_name}</h2>
        <p>
          <strong>
            {user.first_name} {user.last_name}
          </strong>
        </p>
        <p>{user.email}</p>
        <p>{user.gender}</p>
        <p>
          <span>Date of Birth:</span> {user.birthday}
        </p>
        <p>{user.about_me || "Nothing is written."}</p>
        <p>Visibility: {user.visibility}</p>
        <p className="bio">{user.bio}</p>
        {
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
        }
      </div>
    </div>
  );
}
