import { proxyToBackend } from "../../proxyToBackend";

export async function GET(req, { params }) {
  const { categoryName } = params;
  return proxyToBackend(
    req,
    `/api/posts/${encodeURIComponent(categoryName)}`,
    "GET"
  );
}
