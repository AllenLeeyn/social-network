import { proxyToBackend } from "../../proxyToBackend";

export async function GET(req, context) {
  const { id } = await context.params;
  if (!id) {
    return new Response(JSON.stringify({ message: "Missing post id" }), {
      status: 400,
    });
  }
  return proxyToBackend(req, `/api/post/${id}`, "GET");
}
