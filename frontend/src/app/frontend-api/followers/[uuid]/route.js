import { proxyToBackend } from "../../proxyToBackend";

export async function GET(req, context) {
    const { uuid } = await context.params;
    return proxyToBackend(
        req,
        `/api/followers/${uuid}`,
        "GET"
    );
}
