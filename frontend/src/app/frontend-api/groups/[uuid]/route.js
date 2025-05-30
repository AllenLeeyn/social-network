import { proxyToBackend } from "../../proxyToBackend";

export async function GET(req, { params }) {
    const { uuid } = params;
    return proxyToBackend(
        req,
        `/api/group/${encodeURIComponent(uuid)}`,
        "GET"
    );
}
