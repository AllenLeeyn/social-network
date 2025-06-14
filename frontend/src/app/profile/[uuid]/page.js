"use client";

import React, { useState, useEffect, use } from "react";
import { useRouter } from 'next/navigation';
import { fetchUserPosts } from "../../../lib/apiPosts";
import { fetchGroups, fetchFollowees } from "../../../lib/apiAuth";
import "./profile.css";
import ProfileCard from "../../../components/ProfileCard";
import GroupList from "../../../components/GroupList";
import UsersList from "../../../components/UsersList";
import PostList from "../../../components/PostList";
import FollowingsList from '../../../components/FollowingsList';
import SidebarSection from "../../../components/SidebarSection";

export default function ProfilePage({ params }) {
    const { uuid } = use(params);
    const router = useRouter();

    const userUUID = typeof window !== 'undefined' ? localStorage.getItem('user-uuid') : null;

    useEffect(() => {
        if (uuid === userUUID) {
        router.replace('/profile');
        }
    }, [uuid, userUUID, router]);

    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const [groups, setGroups] = useState([]);
    const [groupsLoading, setGroupsLoading] = useState(true);
    const [groupsError, setGroupsError] = useState(null);
    useEffect(() => {
      async function loadGroups() {
        try {
          setGroupsLoading(true);
          const data = await fetchGroups(uuid);
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
    const [isPrivateProfile, setIsPrivateProfile] = useState(false);
    useEffect(() => {
        if (isPrivateProfile) {
            setMyPosts([]);
            setLoading(false);
            return;
        }

        async function loadMyData() {
            try {
            const [userPostsData, myFollowings] = await Promise.all(
                [fetchUserPosts(uuid), fetchFollowees(uuid)]); // need to add fetch by UUID
            setMyPosts(userPostsData.data);

            if (myFollowings && Array.isArray(myFollowings)) {
                const followersList = myFollowings.filter(item => item.leader_uuid === uuid);
                const followingList = myFollowings.filter(item => item.follower_uuid === uuid);

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
    }, [isPrivateProfile]);

    return (
    <main>
        <div className="homepage-layout">

        {/* Profile Content */}
        <section className="main-post-section">
            <ProfileCard uuid={uuid} setPrivateProfile={setIsPrivateProfile} />

          {/* Main Post Content */}
            {!isPrivateProfile && (
                <>

                    <div className="follow-lists-row">
                        {/* Followers Modal */}
                        <FollowingsList title="Followers" users={followers} 
                        displayProperty={"follower_name"} linkProperty={"follower_uuid"}/>

                        {/* Following Modal */}
                        <FollowingsList title="following" users={following}
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

                </>
            )}
        </section>

        {/* Right Sidebar */}
        <aside className="sidebar right-sidebar">
            <SidebarSection title="Chat list">
                <UsersList />
            </SidebarSection>
        </aside>
        </div>
    </main>
    );
}
