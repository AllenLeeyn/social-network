import { cookies } from 'next/headers';

export async function GET() {
    const cookieStore = cookies();
    const sessionId = cookieStore.get('session-id')?.value;
    return Response.json({ sessionId });
}