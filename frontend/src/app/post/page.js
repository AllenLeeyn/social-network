"use client";
import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import SidebarSection from "../../components/SidebarSection";
import CommentsSection from "../../components/CommentSection";
import "./post.css";
import {
  sampleCategories,
  sampleUsers,
  sampleGroups,
  sampleConnections,
} from "../../data/mockData";
import { usePosts } from "../../hooks/usePosts";
import { fetchPostById, submitPostFeedback } from "../../lib/apiPosts";

export default function PostPage() {
  const searchParams = useSearchParams();
  const id = searchParams.get("id");
  const [post, setPost] = useState(null);
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

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

  const refreshComments = async () => {
    try {
      const postData = await fetchPostById(id);
      setComments(postData.data.Comments);
    } catch (err) {
      setError(err.message);
    }
  };

  const handlePostFeedback = async (rating) => {
    try {
      await submitPostFeedback({ parent_id: post.id, rating });
      // Refresh post data to update like/dislike counts
      const postData = await fetchPostById(id);
      setPost(postData.data.Post);
    } catch (err) {
      alert(err.message || "Failed to submit feedback");
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
            {categoriesLoading ? (
              <div>Loading...</div>
            ) : categoriesError ? (
              <div>Error: {categoriesError}</div>
            ) : (
              <ul className="categories">
                {categories.map((cat) => (
                  <li key={cat.id} className="category-item">
                    <strong>{cat.name}</strong>
                  </li>
                ))}
              </ul>
            )}
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

        {/* Main Post Content */}
        <section className="main-post-section">
          {post ? (
            <div key={post.ID} className="post-item">
              <h3>{post.title}</h3>
              <p>
                <em>by {post.user.nick_name}</em>
              </p>
              <p>{post.content}</p>
              <div className="post-actions" style={{ marginTop: "1em" }}>
                <button
                  onClick={() => handlePostFeedback(1)}
                  disabled={post.liked}
                  aria-label="Like"
                >
                  👍 Like {post.like_count}
                </button>
                <button
                  onClick={() => handlePostFeedback(-1)}
                  disabled={post.disliked}
                  aria-label="Dislike"
                  style={{ marginLeft: "1em" }}
                >
                  👎 Dislike {post.dilike_count}
                </button>
              </div>
            </div>
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
            <ul className="users">
              {sampleUsers.map((user) => (
                <li
                  key={user.id}
                  className={`user-item${user.online ? " online" : ""}${
                    user.unread ? " unread" : ""
                  }`}
                >
                  <img src={user.avatar} alt={user.username} />
                  <span>
                    {user.fullName} ({user.username})
                  </span>
                </li>
              ))}
            </ul>
          </SidebarSection>
        </aside>
      </div>
    </main>
  );
}



// "use client";
// import { useSearchParams } from "next/navigation";
// import { useEffect, useState } from "react";
// import SidebarSection from "../../components/SidebarSection";
// import CommentsSection from "../../components/CommentSection";
// import "./post.css";
// import {
//   sampleCategories,
//   sampleUsers,
//   sampleGroups,
//   sampleConnections,
// } from "../../data/mockData";

// export default function PostPage() {
//   const searchParams = useSearchParams();
//   const id = searchParams.get("id");
//   const [post, setPost] = useState(null);
//   const [comments, setComments] = useState([]);
//   const [loading, setLoading] = useState(true);
//   const [error, setError] = useState(null);

//   useEffect(() => {
//     async function fetchData() {
//       try {
//         const postRes = await fetch(`/api/post?id=${id}`);
//         if (!postRes.ok) throw new Error("Post not found");
//         const postData = await postRes.json();
//         setPost(postData.post);
//         setComments(postData.comments);
//       } catch (err) {
//         setError(err.message);
//       } finally {
//         setLoading(false);
//       }
//     }
//     if (id) fetchData();
//   }, [id]);

//   if (loading) return <div>Loading...</div>;
//   if (error) return <div>Error: {error}</div>;

//   return (
//     <main>
//       <div className="homepage-layout">
//         {/* Left Sidebar */}
//         <aside className="sidebar left-sidebar">
//           <SidebarSection title="Categories">
//             <ul className="categories">
//               {sampleCategories.map((cat) => (
//                 <li key={cat.id} className="category-item">
//                   <strong>{cat.name}</strong>
//                 </li>
//               ))}
//             </ul>
//           </SidebarSection>
//           <SidebarSection title="Groups">
//             <ul className="groups">
//               {sampleGroups.map((group) => (
//                 <li key={group.id} className="group-item">
//                   <strong>{group.name}</strong>
//                 </li>
//               ))}
//             </ul>
//           </SidebarSection>
//           <SidebarSection title="Connections">
//             <ul className="connections">
//               {sampleConnections.map((conn) => (
//                 <li key={conn.id} className="connection-item">
//                   <span>
//                     <strong>
//                       {conn.fullName} ({conn.username})
//                     </strong>
//                   </span>
//                 </li>
//               ))}
//             </ul>
//           </SidebarSection>
//         </aside>

//         {/* Main Post Content */}
//         <section className="main-post-section">
//           {post ? (
//             <div key={post.ID} className="post-item">
//               <h3>{post.title}</h3>
//               <p>
//                 <em>by {post.userName}</em>
//               </p>
//               <p>{post.content}</p>
//             </div>
//           ) : (
//             <div className="post-item">
//               <h3>Post not found</h3>
//             </div>
//           )}

//           <CommentsSection title="Comments" comments={comments || []} />
//         </section>

//         {/* Right Sidebar */}
//         <aside className="sidebar right-sidebar">
//           <SidebarSection title="Active Users">
//             <ul className="users">
//               {sampleUsers.map((user) => (
//                 <li
//                   key={user.id}
//                   className={`user-item${user.online ? " online" : ""}${
//                     user.unread ? " unread" : ""
//                   }`}
//                 >
//                   <img src={user.avatar} alt={user.username} />
//                   <span>
//                     {user.fullName} ({user.username})
//                   </span>
//                 </li>
//               ))}
//             </ul>
//           </SidebarSection>
//         </aside>
//       </div>
//     </main>
//   );
// }
