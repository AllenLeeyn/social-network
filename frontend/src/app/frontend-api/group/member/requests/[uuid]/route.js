import { proxyToBackend } from "../../../../proxyToBackend";

export async function GET(req, { params }) {
    const { uuid } = await params;
    return proxyToBackend(req, `/api/group/member/requests/${uuid}`, 'GET');
}
