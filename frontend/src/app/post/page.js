"use client";
import { Suspense } from "react";
import PostContent from "../../components/PostContent";

export default function PostPage() {
  return (
    <main>
      <Suspense fallback={<div>Loading post...</div>}>
        <PostContent />
      </Suspense>
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
