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
        setPosts(data.data);

        // Extract unique categories from all posts
        const allCategories = data.data
          .flatMap(post => post.categories || [])
          .filter(cat => cat && cat.id && cat.name);

        // Remove duplicates by id
        const uniqueCategories = [];
        const seen = new Set();
        for (const cat of allCategories) {
          if (!seen.has(cat.id)) {
            seen.add(cat.id);
            uniqueCategories.push(cat);
          }
        }
        setCategories(uniqueCategories);

        setUsers(data.users); // This may also need adjustment if users are not top-level
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
