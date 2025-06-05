"use client";

import React, { useState, useEffect } from "react";
import { fetchGroups, fetchFollowees } from "../../lib/apiAuth";
import { fetchMyPosts } from "../../lib/apiPosts";
import "./profile.css";
import ProfileCard from '../../components/ProfileCard';
import UsersList from "../../components/UsersList";
import GroupList from "../../components/GroupList";
import PostList from "../../components/PostList";
import FollowingsList from '../../components/FollowingsList';

export default function ProfilePage() {
  const userUUID = typeof window !== 'undefined' ? localStorage.getItem('user-uuid') : null;

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [groups, setGroups] = useState([]);
  const [groupsLoading, setGroupsLoading] = useState(true);
  const [groupsError, setGroupsError] = useState(null);
  useEffect(() => {
    async function loadGroups() {
      try {
        setGroupsLoading(true);
        const data = await fetchGroups();
        setGroups(data || []);
      } catch (err) {
        setGroupsError(err.message);
      } finally {
        setGroupsLoading(false);
      }
    }
    loadGroups();
  }, []);

  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const [myPosts, setMyPosts] = useState([]);
  useEffect(() => {
  async function loadMyData() {
    try {
      const [myPostsData, myFollowings] = await Promise.all(
        [fetchMyPosts(), fetchFollowees()]);
      setMyPosts(myPostsData.data);

      if (myFollowings && Array.isArray(myFollowings)) {
        const followersList = myFollowings.filter(item => item.leader_uuid === userUUID);
        const followingList = myFollowings.filter(item => item.follower_uuid === userUUID);

        setFollowers(followersList);
        setFollowing(followingList);
      }
    } catch (err) {
        setError(err.message);
    } finally {
        setLoading(false);
    }
  }
  loadMyData();
  }, []);

  return (
  <main>
      <div className="homepage-layout">

      {/* Profile Content */}
      <section className="main-post-section">
          <ProfileCard />

          <div className="follow-lists-row">
          <FollowingsList title="Followers" users={followers} 
          displayProperty={"follower_name"} linkProperty={"follower_uuid"}/>
          <FollowingsList title="Following" users={following}
          displayProperty={"leader_name"} linkProperty={"leader_uuid"}/>
          </div>

 
          <div className="user-posts-section">
          <h3>Groups</h3>
              <GroupList
              groups={groups}
              loading={groupsLoading}
              error={groupsError}
              />
          </div>
          
          {/* Main Post Content */}
          <div className="user-posts-section">
          <h3>Posts</h3>
              {loading && <div>Loading...</div>}
              {error && <div>Error: {error}</div>}
              {!loading && !error && <PostList posts={myPosts} />}
          </div>

      </section>

      {/* Right Sidebar */}
      <aside className="sidebar right-sidebar">
      <UsersList />
      </aside>
      </div>
  </main>
  );
}
