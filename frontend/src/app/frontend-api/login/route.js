export async function POST(req) {
  const backendUrl = "http://localhost:8080/api/login"; // Backend login endpoint

  try {
    const body = await req.json(); // Parse the request body
    const response = await fetch(backendUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    });

    if (!response.ok) {
      const errorData = await response.json();
      return new Response(JSON.stringify(errorData), {
        status: response.status,
      });
    }

    // Forward cookies from backend to client
    const cookies = response.headers.get("set-cookie");
    if (cookies) {
      return new Response(await response.text(), {
        status: response.status,
        headers: {
          "Set-Cookie": cookies,
          "Access-Control-Allow-Credentials": "true", // Allow credentials
        },
      });
    }

    return new Response(await response.text(), { status: response.status });
  } catch (error) {
    console.error("Error during login:", error);
    return new Response(JSON.stringify({ message: "Internal Server Error" }), {
      status: 500,
    });
  }
}
