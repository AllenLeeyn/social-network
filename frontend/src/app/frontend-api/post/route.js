import { proxyToBackend } from "../proxyToBackend";

export async function GET(req) {
  const { searchParams } = new URL(req.url);
  const id = searchParams.get("id");
  if (!id) {
    return new Response(JSON.stringify({ message: "Missing post id" }), {
      status: 400,
    });
  }
  return proxyToBackend(req, `/api/post/${id}`, "GET");
}