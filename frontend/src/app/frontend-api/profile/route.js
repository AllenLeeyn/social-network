import { proxyToBackend } from "../proxyToBackend";

export async function GET(req) {
  const url = new URL(req.url);
  const uuid = url.searchParams.get("uuid");
  const backendUrl = uuid ? `/api/user/${uuid}` : `/api/user`;
  return proxyToBackend(req, backendUrl, "GET");
}
