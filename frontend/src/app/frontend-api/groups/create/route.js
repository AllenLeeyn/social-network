import { proxyToBackend } from '../../proxyToBackend';


export async function POST(req) {
    return proxyToBackend(
        req,
        "/api/group/create", // Your Go backend endpoint
        "POST"
    );
}
