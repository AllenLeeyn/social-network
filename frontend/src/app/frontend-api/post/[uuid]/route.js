import { proxyToBackend } from "../../proxyToBackend";

export async function GET(req, context) {
  const { uuid } = await context.params;
  if (!uuid) {
    return new Response(JSON.stringify({ message: "Missing post id" }), {
      status: 400,
    });
  }
  return proxyToBackend(req, `/api/post/${uuid}`, "GET");
}
