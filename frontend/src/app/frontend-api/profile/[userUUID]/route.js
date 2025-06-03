import { proxyToBackend } from "../../proxyToBackend";

export async function GET(req, context) {
  const { userUUID } = await context.params;
  const url = `/api/user/${userUUID}`;

  return proxyToBackend(req, url, "GET");
}
