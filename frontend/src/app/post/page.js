"use client";
import { useSearchParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import SidebarSection from "../../components/SidebarSection";
import CommentsSection from "../../components/CommentSection";
import "./post.css";
import "../../styles/PostList.css";
import { usePosts } from "../../hooks/usePosts";
import { fetchPostById } from "../../lib/apiPosts";
import CategoriesList from "../../components/CategoriesList";
import ConnectionList from "../../components/ConnectionList";
import { fetchFollowees, fetchGroups } from "../../lib/apiAuth";
import { toast } from 'react-toastify';
import PostCard from '../../components/PostCard';
import UsersList from "../../components/UsersList";
import GroupList from "../../components/GroupList";

export default function PostPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const id = searchParams.get("id");
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
        setComments(postData.data.Comments);
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
      setComments(postData.data.Comments);
    } catch (err) {
      setError(err.message);
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <main>
      <div className="homepage-layout">
        {/* Left Sidebar */}
        <aside className="sidebar left-sidebar">
          <SidebarSection title="Categories">
            <CategoriesList
              categories={categories}
              loading={categoriesLoading}
              error={categoriesError}
              onCategoryClick={(cat) =>
                router.push(`/?category=${encodeURIComponent(cat)}`)
              }
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

        {/* Main Post Content */}
        <section className="main-post-section">
          {post ? (
            <PostCard post={post} />
          ) : (
            <div className="post-item">
              <h3>Post not found</h3>
            </div>
          )}

          <CommentsSection
            title="Comments"
            comments={comments || []}
            postId={post.id}
            onCommentSubmitted={refreshComments}
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

