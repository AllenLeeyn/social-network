import { proxyToBackend } from "../proxyToBackend";

export async function GET(req) {
  return proxyToBackend(req, "/api/categories", "GET");
}