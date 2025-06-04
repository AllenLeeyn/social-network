"use client";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import SidebarSection from "../../../components/SidebarSection";
import CommentsSection from "../../../components/CommentSection";
// import "../../../styles/globals.css";
import "../../../styles/PostList.css";
import { usePosts } from "../../../hooks/usePosts";
import { fetchPostById } from "../../../lib/apiPosts";
import CategoriesList from "../../../components/CategoriesList";
import { fetchFollowees, fetchGroups } from "../../../lib/apiAuth";
import { toast } from "react-toastify";
import PostCard from "../../../components/PostCard";
import UsersList from "../../../components/UsersList";

export default function PostPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id;
  const [post, setPost] = useState(null);
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [connections, setConnections] = useState([]);
  const [connectionsLoading, setConnectionsLoading] = useState(true);
  const [connectionsError, setConnectionsError] = useState(null);

  const [groups, setGroups] = useState([]);
  const [groupsLoading, setGroupsLoading] = useState(true);
  const [groupsError, setGroupsError] = useState(null);

  // Fetch categories (and optionally users, etc.)
  const {
    categories,
    loading: categoriesLoading,
    error: categoriesError,
  } = usePosts();

  useEffect(() => {
    async function fetchData() {
      try {
        const postData = await fetchPostById(id);
        setPost(postData.data.Post);
        // This is the key fix - ensure comments is always an array
        setComments(postData.data.Comments || []);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }
    if (id) fetchData();
  }, [id]);

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

  const refreshComments = async () => {
    try {
      const postData = await fetchPostById(id);
      // Also fix this function to handle null comments
      setComments(postData.data.Comments || []);
    } catch (err) {
      setError(err.message);
    }
  };

  if (loading) return <div className="loading">Loading...</div>;
  if (error) return <div className="error">Error: {error}</div>;
  if (!post) return <div className="error">Post not found</div>;

  return (
    <main>
      <div className="post-page-layout">

        {/* Main Content */}
        <section className="main-feed">
          <div className="single-page">
          <PostCard post={post} />
          </div>
          <CommentsSection
            comments={comments}
            postId={id}
            onCommentAdded={refreshComments}
          />
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
