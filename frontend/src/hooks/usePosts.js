import { useState, useEffect } from "react";
import { fetchPosts } from "../lib/apiPosts";

export function usePosts() {
  const [posts, setPosts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [users, setUsers] = useState([]);

  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadPosts() {
      try {
        const data = await fetchPosts();
        setPosts(data.data); // <-- changed from data.posts to data.data
        setCategories(data.categories);
        setUsers(data.users);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    loadPosts();
  }, []);

  return { posts, categories, users, loading, error };
}
