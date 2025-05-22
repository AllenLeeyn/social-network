import { cookies } from "next/headers";

export async function POST(req) {
  const backendUrl = "http://localhost:8080/api/submitComment";
  try {
    // Get session-id cookie
    const cookieStore = await cookies();
    const sessionId = cookieStore.get("session-id")?.value;

    if (!sessionId) {
      return new Response(
        JSON.stringify({ message: "Unauthorized: No session-id cookie found" }),
        { status: 401 }
      );
    }

    // Get comment data from request body
    const body = await req.json();

    // Forward to backend
    const response = await fetch(backendUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Cookie: `session-id=${sessionId}`,
      },
      body: JSON.stringify(body),
    });

    const data = await response.json();
    return new Response(JSON.stringify(data), { status: response.status });
  } catch (error) {
    console.error("Error submitting comment:", error);
    return new Response(JSON.stringify({ message: "Internal Server Error" }), {
      status: 500,
    });
  }
}
