// This is our homepage

"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation"; 

import SidebarSection from "../components/SidebarSection";
import CategoriesList from "../components/CategoriesList";
import PostList from "../components/PostList";
import CreatePost from "../components/CreatePost";
import Modal from "../components/Modal";
import UsersList from "../components/UsersList";

import { fetchPosts } from "../lib/apiPosts";
import { usePosts } from "../hooks/usePosts";

import { useWebsocketContext } from '../contexts/WebSocketContext';

import {
  sampleGroups,
  sampleConnections,
} from "../data/mockData";

export default function HomePage() {
  const [showModal, setShowModal] = useState(false);
  const { posts, categories, loading, error } = usePosts([]);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [isAuthorized, setIsAuthorized] = useState(false);
  const router = useRouter(); 

  const { isConnected, connect } = useWebsocketContext();
  const [sessionId, setSessionId] = useState(null);
  useEffect(() => {
    if (typeof window !== 'undefined') {
      setSessionId(localStorage.getItem('session-id'));
    }
  }, []);
  //const sessionId = localStorage.getItem('session-id'); // Or from cookies/auth context

  useEffect(() => {
    async function checkAccess() {
      try {
        await fetchPosts();
        setIsAuthorized(true);
      } catch (error) {
        console.error("Access denied, redirecting to login:", error);
        router.push("/login");
      }
    }
    console.log('Session ID:', localStorage.getItem('session-id'));
    checkAccess();
  }, [router]);


  // Filter logic
  const filteredPosts = selectedCategory
    ? posts.filter(
        (post) =>
          // If your post has a catNames string (e.g. "General, sqlite3")
          post.catNames &&
          post.catNames
            .split(",")
            .map((s) => s.trim())
            .includes(selectedCategory)
        // If your post has categories as array of names: post.categories.includes(selectedCategory)
      )
    : posts;

  const handleCategoryClick = (cat) => {
    setSelectedCategory((prev) => (prev === cat ? null : cat));
  };

    useEffect(() => {
    if (sessionId && !isConnected) {
      connect(sessionId); // Trigger connection
    }
  }, [sessionId, connect, isConnected]);

  if (!isAuthorized) {
    // Prevent rendering the homepage until access is verified
    return null;
  }


  return (
    <main>
      <h1>WebSocket Status: {isConnected ? '✅ Connected' : '❌ Disconnected'}</h1>
      <div className="homepage-layout">
        {/* Left Sidebar */}
        <aside className="sidebar left-sidebar">
          <SidebarSection title="Categories">
            <CategoriesList
              categories={categories}
              loading={loading}
              error={error}
              selectedCategory={selectedCategory}
              onCategoryClick={handleCategoryClick}
            />
          </SidebarSection>
          <SidebarSection title="Groups">
            <ul className="groups">
              {sampleGroups.map((group) => (
                <li key={group.id} className="group-item">
                  <strong>{group.name}</strong>
                </li>
              ))}
            </ul>
          </SidebarSection>
          <SidebarSection title="Connections">
            <ul className="connections">
              {sampleConnections.map((conn) => (
                <li key={conn.id} className="connection-item">
                  <span>
                    <strong>
                      {conn.fullName} ({conn.username})
                    </strong>
                  </span>
                </li>
              ))}
            </ul>
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
          {!loading && !error && <PostList posts={filteredPosts} />}

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
          <SidebarSection title="Active Users">
            <UsersList />
          </SidebarSection>
        </aside>
      </div>
    </main>
  );
}
