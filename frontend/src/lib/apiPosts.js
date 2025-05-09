const API_URL = "/api"; // Proxy API base URL

export async function fetchPosts() {
  const response = await fetch(`${API_URL}/posts`, {
    method: "GET",
    credentials: "include", // Include cookies in the request
  });
  if (!response.ok) {
    console.log(response)
    throw new Error("Failed to fetch posts");
  }
  return response.json();
}
