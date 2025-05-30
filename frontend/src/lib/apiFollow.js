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