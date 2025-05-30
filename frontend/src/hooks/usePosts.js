import { useState, useEffect } from "react";
import { fetchPosts, fetchCategories } from "../lib/apiPosts";

export function usePosts() {
  const [posts, setPosts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [users, setUsers] = useState([]);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadData() {
      try {
        const [postsData, categoriesData] = await Promise.all([
          fetchPosts(),
          fetchCategories(),
        ]);
        console.log(postsData)
        setPosts(postsData.data);
        setCategories(categoriesData.data);
        setUsers(postsData.users);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    loadData();
  }, []);

  return { posts, categories, users, loading, error };
}
