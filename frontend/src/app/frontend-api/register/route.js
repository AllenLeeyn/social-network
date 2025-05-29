    const baseURL = process.env.NEXT_PUBLIC_BACKEND_URL;
    
    export async function POST(req) {
    const backendUrl = `${baseURL}/api/register`;
  
    try {
      const body = await req.json(); 
      const response = await fetch(backendUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body), // Forward the request body
      });
  
      if (!response.ok) {
        const errorData = await response.json();
        return new Response(JSON.stringify(errorData), 
        { status: response.status });
      }
  
      // Forward cookies from backend to client
      const cookies = response.headers.get("set-cookie");
      if (cookies) {
        return new Response(await response.text(), {
          status: response.status,
          headers: {
            "Set-Cookie": cookies,
            "Access-Control-Allow-Credentials": "true",
          },
        });
      }

      return new Response(await response.text(), { status: response.status });
    } catch (error) {
      console.error("Error in signup proxy:", error);
      return new Response(JSON.stringify({ message: "Internal Server Error" }), {
        status: 500,
      });
    }
  }