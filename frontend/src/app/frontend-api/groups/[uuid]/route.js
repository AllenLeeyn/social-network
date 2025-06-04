// Single group detail

import { proxyToBackend } from "../../proxyToBackend";

export async function GET(req, context) {
    const { params } = await context;
    const { uuid } = await params;

    const url = new URL(req.url);
    const search = url.search || "";

    return proxyToBackend(
        req,
        `/api/groups/${encodeURIComponent(uuid)}${search}`,
        "GET"
    );
}
