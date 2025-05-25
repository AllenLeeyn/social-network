import { cookies } from "next/headers";

export async function POST() {

  const backendUrl = "http://localhost:8080/api/logout";

  try {
    const cookieStore = await cookies();
    const sessionId = cookieStore.get("session-id")?.value;

    if (!sessionId) {
      return new Response(
        JSON.stringify({ message: "No session cookie found" }),
        { status: 400 }
      );
    }


    // Forward the request to the backend
    const response = await fetch(backendUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Cookie: `session-id=${sessionId}`,
      },
    });

    // Forward the backend's response as-is
    const responseBody = await response.text(); // Read the response body
    return new Response(responseBody, {
      status: response.status,
      headers: { "Content-Type": response.headers.get("Content-Type") },
    });
  } catch (error) {
    console.error("Error during logout:", error);
    return new Response(JSON.stringify({ message: "Internal Server Error" }), {
      status: 500,
    });
  }
}
