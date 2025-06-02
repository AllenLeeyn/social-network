import { proxyToBackend } from "../proxyToBackend";

export async function GET(req) {
  console.log("Fetching user profile data");
  return proxyToBackend(req, `/api/user`, "GET");
}
