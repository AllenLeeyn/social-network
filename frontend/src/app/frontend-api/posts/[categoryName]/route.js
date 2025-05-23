import { proxyToBackend } from "../../proxyToBackend";

export async function GET(req, { params }) {
  const { categoryName } = params;
  return proxyToBackend(
    req,
    `http://localhost:8080/api/posts/${encodeURIComponent(categoryName)}`,
    "GET"
  );
}
