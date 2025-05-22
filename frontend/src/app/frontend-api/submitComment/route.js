import { proxyToBackend } from "../proxyToBackend";

export async function POST(req) {
  return proxyToBackend(req, "http://localhost:8080/api/submitComment");
}