import { cookies } from "next/headers";

export async function proxyToBackend(req, backendUrl, method = "POST") {
  try {
    const cookieStore = await cookies();
    const sessionId = cookieStore.get("session-id")?.value;

    if (!sessionId) {
      return new Response(
        JSON.stringify({ message: "Unauthorized: No session-id cookie found" }),
        { status: 401 }
      );
    }

    let body;
    if (method !== "GET") {
      body = await req.json();
    }

    const response = await fetch(backendUrl, {
      method,
      headers: {
        "Content-Type": "application/json",
        Cookie: `session-id=${sessionId}`,
      },
      ...(body ? { body: JSON.stringify(body) } : {}),
    });

    const data = await response.json();
    return new Response(JSON.stringify(data), { status: response.status });
  } catch (error) {
    console.error("Proxy error:", error);
    return new Response(JSON.stringify({ message: "Internal Server Error" }), {
      status: 500,
    });
  }
}
