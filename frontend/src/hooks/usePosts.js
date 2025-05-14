import { useState, useEffect } from "react";
import { fetchPosts } from "../lib/apiPosts";

export function usePosts() {
  const [posts, setPosts] = useState([]);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadPosts() {
      try {
        const data = await fetchPosts();
        setPosts(data.posts); // Assuming backend returns { posts: [...] }
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    loadPosts();
  }, []);

  return { posts, error, loading };
}
