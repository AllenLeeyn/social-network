import { cookies } from "next/headers";

export async function GET() {
  const backendUrl = "http://localhost:8080/posts"; // Backend endpoint

  try {
    // Retrieve the session-id cookie
    const cookieStore = cookies();
    const sessionId = cookieStore.get("session-id")?.value;

    if (!sessionId) {
      return new Response(
        JSON.stringify({ message: "Unauthorized: No session-id cookie found" }),
        { status: 401 }
      );
    }

    // Forward the session-id cookie to the backend
    const response = await fetch(backendUrl, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Cookie: `session-id=${sessionId}`, // Include the session-id cookie
      },
    });

    if (!response.ok) {
      return new Response(
        JSON.stringify({ message: "Failed to fetch posts" }),
        {
          status: response.status,
        }
      );
    }

    const data = await response.json();
    return new Response(JSON.stringify(data), { status: 200 });
  } catch (error) {
    console.error("Error fetching posts:", error);
    return new Response(JSON.stringify({ message: "Internal Server Error" }), {
      status: 500,
    });
  }
}
