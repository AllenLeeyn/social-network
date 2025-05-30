export async function submitGroupRequestOrInviteResponse(responseData) {
  const response = await fetch("/frontend-api/group/member/response", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(responseData),
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.message || "Failed to submit group request response");
  }
  return data;
}
