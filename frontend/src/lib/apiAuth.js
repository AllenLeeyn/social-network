const API_URL = "/frontend-api";

export async function login(email, password) {
  const response = await fetch(`${API_URL}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.message || "Login failed");
  }
  const data = await response.json();

  return data;
}

export async function signup(userData) {
  const response = await fetch(`${API_URL}/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(userData),
  });
  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.message || "Signup failed");
  }
  const data = await response.json();

  return data;
}

export async function logout() {
  const response = await fetch(`${API_URL}/logout`, {
    method: "GET",
    credentials: "include", // Ensure cookies are included
  });

  if (response.status === 200) {
    window.location.href = "/login"; // Redirect to login
    return;
  }

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.message || "Logout failed");
  }
}

export async function fetchFollowees() {
  const response = await fetch("/frontend-api/followers/", {
    method: "GET",
    credentials: "include",
  });
  if (!response.ok) throw new Error("Failed to fetch followees");
  const result = await response.json();
  return result.data;
}

export async function fetchGroups() {
  const response = await fetch("/frontend-api/groups?joined=true", {
    method: "GET",
    credentials: "include",
  });
  if (!response.ok) throw new Error("Failed to fetch groups");
  const result = await response.json();
  return result.data;
}
