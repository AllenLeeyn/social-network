// gets members for a group

import { proxyToBackend } from "../../../proxyToBackend";

export async function GET(req, context) {
    const { params } = await context;
    const { uuid } = await params;
    return proxyToBackend(
        req,
        `/api/group/members/${uuid}`, 
        "GET"
    );
}
