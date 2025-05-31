import { proxyToBackend } from "../proxyToBackend";

export async function GET(req) {
    const url = new URL(req.url);
    const search = url.search || "";

    return proxyToBackend(
        req,
        `/api/groups/${search}`,
        "GET"
    );
}
