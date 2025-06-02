"use client";

import { useEffect, useState } from "react";

export default function ProfileCard() {
  const [user, setUser] = useState(null);
  const [followers, setFollowers] = useState([]);
  const [followings, setFollowings] = useState([]);

  useEffect(() => {
    async function fetchProfile() {
      try {
        const res = await fetch("/frontend-api/profile", {
          credentials: "include", // if cookies/session are used
        });
        if (!res.ok) {
          throw new Error("Failed to fetch user data");
        }
        const data = await res.json();
        setUser(data.data);

        const followRes = await fetch("/frontend-api/followers", {
          credentials: "include",
        });
        if (!followRes.ok) throw new Error("Failed to fetch followers");
        const followJson = await followRes.json();
        const followData = Array.isArray(followJson.data)
          ? followJson.data
          : [];

        const followersList = followData.filter(
          (entry) => entry.leader_uuid === data.data.uuid
        );
        const followingList = followData.filter(
          (entry) => entry.follower_uuid === data.data.uuid
        );

        setFollowers(followersList);
        setFollowings(followingList);
      } catch (err) {
        console.error(err);
      }
    }

    fetchProfile();
  }, []);

  if (!user) return <div>Loading profile...</div>;

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
