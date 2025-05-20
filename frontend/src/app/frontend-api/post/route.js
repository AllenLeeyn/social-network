import { cookies } from "next/headers";

export async function GET(req) {
  const backendUrl = "http://localhost:8080/api/post";
  try {
    // Get the post id from the query string
    const { searchParams } = new URL(req.url);
    const id = searchParams.get("id");
    console.log("Fetching post with id:", id);

    if (!id) {
      return new Response(JSON.stringify({ message: "Missing post id" }), {
        status: 400,
      });
    }

    // Retrieve the session-id cookie
    const cookieStore = await cookies();
    const sessionId = cookieStore.get("session-id")?.value;

    if (!sessionId) {
      return new Response(
        JSON.stringify({ message: "Unauthorized: No session-id cookie found" }),
        { status: 401 }
      );
    }

    // Forward the session-id cookie and id (as path param) to the backend
    const response = await fetch(`${backendUrl}/${id}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Cookie: `session-id=${sessionId}`,
      },
    });

    if (!response.ok) {
      return new Response(JSON.stringify({ message: "Failed to fetch post" }), {
        status: response.status,
      });
    }

    const data = await response.json();
    console.log("data is: ",data)
    return new Response(JSON.stringify(data), { status: 200 });
  } catch (error) {
    console.error("Error fetching post:", error);
    return new Response(JSON.stringify({ message: "Internal Server Error" }), {
      status: 500,
    });
  }
}
