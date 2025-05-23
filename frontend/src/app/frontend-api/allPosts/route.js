import { proxyToBackend } from "../proxyToBackend";

export async function GET(req) {
  return proxyToBackend(req, "http://localhost:8080/api/allPosts", "GET");
}