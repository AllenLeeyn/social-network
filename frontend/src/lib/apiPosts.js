const API_URL = "/frontend-api"; // Proxy API base URL

export async function fetchPosts() {
  const response = await fetch(`${API_URL}/allPosts`, {
    method: "GET",
    credentials: "include", // Include cookies in the request
  });
  if (!response.ok) {
    throw new Error("Failed to fetch posts");
  }
  return response.json();
}

export async function fetchMyPosts() {
  const response = await fetch(`${API_URL}/myCreatedPosts`, {
    method: "GET",
    credentials: "include",
  });
  if (!response.ok) {
    throw new Error("Failed to fetch posts");
  }
  return response.json();
}

export async function fetchUserPosts(uuid) {
  const response = await fetch(`${API_URL}/userPosts/${uuid}`, {
    method: "GET",
    credentials: "include",
  });
  if (!response.ok) {
    throw new Error("Failed to fetch posts");
  }
  return response.json();
}

export async function fetchGroupPosts(uuid) {
  const response = await fetch(`${API_URL}/groupPosts/${uuid}`, {
    method: "GET",
    credentials: "include",
  });
  if (!response.ok) {
    throw new Error("Failed to fetch posts");
  }
  return response.json();
}

export async function fetchPostById(id) {
  const response = await fetch(`${API_URL}/post/${id}`, {
    method: "GET",
    credentials: "include",
  });
  if (!response.ok) throw new Error("Failed to fetch post");
  return response.json();
}

export async function createPost(postData) {
  const response = await fetch(`${API_URL}/submitPost`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(postData),
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.message || "Failed to create post");
  }
  return data;
}

export async function fetchCategories() {
  const response = await fetch("/frontend-api/categories", {
    credentials: "include",
  });
  if (!response.ok) throw new Error("Failed to fetch categories");
  return response.json();
}

export async function submitComment(commentData) {
  const response = await fetch("/frontend-api/submitComment", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(commentData),
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.message || "Failed to submit comment");
  }
  return data;
}

export async function submitCommentFeedback(feedbackData) {
  const response = await fetch("/frontend-api/commentFeedback", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(feedbackData),
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.message || "Failed to submit comment feedback");
  }
  return data;
}

export async function submitPostFeedback(feedbackData) {
  const response = await fetch("/frontend-api/postFeedback", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(feedbackData),
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.message || "Failed to submit post feedback");
  }
  return data;
}

export async function fetchPostsByCategory(categoryName) {
  const response = await fetch(
    `/frontend-api/posts/${encodeURIComponent(categoryName)}`,
    {
      method: "GET",
      credentials: "include",
    }
  );
  if (!response.ok) throw new Error("Failed to fetch posts by category");
  return response.json();
}
