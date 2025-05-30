const baseURL = process.env.NEXT_PUBLIC_BACKEND_URL;

export async function POST(req) {
  const backendUrl = `${baseURL}/api/uploadFile`;

  try {
    const headers = {};
    req.headers.forEach((value, key) => {
      if (key.toLowerCase() === 'content-type' ||
            key.toLowerCase() === 'cookie' ||
            key.toLowerCase() === 'authorization') {
        headers[key] = value;
      }
    });

    const response = await fetch(backendUrl, {
      method: "POST",
      headers,
      body: req.body,
      duplex: 'half',
      credentials: "include",
    });

    const responseBody = await response.text();
  console.log(responseBody)

    return new Response(responseBody, {
      status: response.status,
      headers: { "Content-Type": "application/json" },
    });
  } catch (error) {
    console.error("Error proxying upload:", error);
    return new Response(JSON.stringify({ message: "Internal Server Error" }), {
      status: 500,
      headers: { "Content-Type": "application/json" },
    });
  }
}