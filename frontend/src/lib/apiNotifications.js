const API_URL = "/frontend-api"; // Proxy API base URL

export async function fetchNotifications() {
  const response = await fetch(`${API_URL}/notifications`, {
    method: "GET",
    credentials: "include", // Include cookies in the request
  });
  if (!response.ok) {
    throw new Error("Failed to fetch notifications");
  }
  return response.json();
}

export async function readNotification(notificationData) {
  const response = await fetch("/frontend-api/readNotification", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(notificationData),
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.message || "Failed to make notification read");
  }
  return data;
}