export async function submitFollowResponse(feedbackData) {
  const response = await fetch("/frontend-api/follower/response", {
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

export async function submitFollowRequest(userData) {
  const response = await fetch("/frontend-api/follower/request", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(userData),
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.message || "Failed to submit follow request");
  }
  return data;
}

export async function submitUnfollowRequest(userData) {
  const response = await fetch("/frontend-api/follower/unfollow", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(userData),
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.message || "Failed to unfollow user");
  }
  return data;
}
