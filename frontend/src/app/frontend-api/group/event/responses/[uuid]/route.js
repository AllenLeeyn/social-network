import { proxyToBackend } from "../../../../proxyToBackend";

export async function GET(req, context) {
  const { params } = await context;
  const { uuid } = await params;

  return proxyToBackend(
    req,
    `/api/group/event/responses/${encodeURIComponent(uuid)}`,
    "GET"
  );
}
