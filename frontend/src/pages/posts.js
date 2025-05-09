import { usePosts } from "../hooks/usePosts";

export default function PostsPage() {
  const { posts, error, loading } = usePosts();

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Posts</h1>
      <ul>
        {posts.map((post) => (
          <li key={post.id}>
            <h2>{post.title}</h2>
            <p>{post.content}</p>
            <small>By {post.userName}</small>
          </li>
        ))}
      </ul>
    </div>
  );
}
