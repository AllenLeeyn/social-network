import { proxyToBackend } from "../../proxyToBackend";

export async function POST(req) {
    return proxyToBackend(req, "/api/group/event/create", "POST");
}
