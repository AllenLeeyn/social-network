"use client";

import { useState, useEffect } from "react";

import SidebarSection from "../components/SidebarSection";
import CategoriesList from "../components/CategoriesList";
import PostList from "../components/PostList";
import GroupList from "../components/GroupList";
import CreatePost from "../components/CreatePost";
import Modal from "../components/Modal";
import UsersList from "../components/UsersList";

import { fetchPostsByCategory } from "../lib/apiPosts";
import { usePosts } from "../hooks/usePosts";
import ConnectionList from "../components/ConnectionList";
import { fetchFollowees, fetchGroups } from "../lib/apiAuth";

import { useWebsocketContext } from "../contexts/WebSocketContext";

export default function HomePage() {
  const [showModal, setShowModal] = useState(false);
  const { posts, categories, loading, error } = usePosts([]);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [filteredPosts, setFilteredPosts] = useState([]);
  const [categoryLoading, setCategoryLoading] = useState(false);
  const [categoryError, setCategoryError] = useState(null);

  const [connections, setConnections] = useState([]);
  const [connectionsLoading, setConnectionsLoading] = useState(true);
  const [connectionsError, setConnectionsError] = useState(null);

  const [groups, setGroups] = useState([]);
  const [groupsLoading, setGroupsLoading] = useState(true);
  const [groupsError, setGroupsError] = useState(null);

  const { isConnected, connect } = useWebsocketContext();

  useEffect(() => {
    if (!selectedCategory) setFilteredPosts(posts);
  }, [posts, selectedCategory]);

  useEffect(() => {
    async function loadConnections() {
      try {
        setConnectionsLoading(true);
        const data = await fetchFollowees();
        setConnections(data || []);
      } catch (err) {
        setConnectionsError(err.message);
      } finally {
        setConnectionsLoading(false);
      }
    }
    loadConnections();
  }, []);

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

  const handleCategoryClick = async (cat) => {
    if (selectedCategory === cat) {
      setSelectedCategory(null);
      setFilteredPosts(posts);
      setCategoryError(null);
      return;
    }
    setSelectedCategory(cat);
    setCategoryLoading(true);
    setCategoryError(null);
    try {
      const data = await fetchPostsByCategory(cat);
      setFilteredPosts(data.data.Posts || []);
    } catch (err) {
      setCategoryError(err.message);
      setFilteredPosts([]);
    } finally {
      setCategoryLoading(false);
    }
  };

  // Filter logic
  const displayedPosts = selectedCategory ? filteredPosts : posts;

  return (
    <main>
      <div className="homepage-layout">
        {/* Left Sidebar */}
        <aside className="sidebar left-sidebar">
          <SidebarSection title="Categories">
            <CategoriesList
              categories={categories}
              loading={loading}
              error={error}
              onCategoryClick={handleCategoryClick}
            />
          </SidebarSection>
          <SidebarSection title="Groups">
            <GroupList
              groups={groups}
              loading={groupsLoading}
              error={groupsError}
            />
          </SidebarSection>
          <SidebarSection title="Connections">
            <ConnectionList
              connections={connections}
              loading={connectionsLoading}
              error={connectionsError}
            />
          </SidebarSection>
        </aside>

        {/* Center / Main View */}
        <section className="main-feed post-list-section">
          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
            }}
          >
            <h2>Latest Posts</h2>
            <button
              className="create-post-btn"
              onClick={() => setShowModal(true)}
              aria-label="Create a new post"
            >
              + Create Post
            </button>
          </div>
          {loading && <div>Loading...</div>}
          {error && <div>Error: {error}</div>}
          {!loading && !error && <PostList posts={displayedPosts} />}

          {showModal && (
            <Modal onClose={() => setShowModal(false)} title="Create Post">
              <CreatePost
                categories={categories}
                onClose={() => setShowModal(false)}
              />
            </Modal>
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
