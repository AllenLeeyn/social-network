import { NextResponse } from "next/server";

const baseURL = process.env.NEXT_PUBLIC_BACKEND_URL;

export async function GET(_req, { params }) {
  const { filename } = params;

  const backendUrl = `${baseURL}/uploads/${filename}`;

  console.log(backendUrl)
  try {
    const response = await fetch(backendUrl, {
      method: "GET",
    });

    if (!response.ok) {
      return new NextResponse("Image not found", { status: 404 });
    }

    const contentType = response.headers.get("Content-Type");
    const buffer = await response.arrayBuffer();

    return new NextResponse(buffer, {
      status: 200,
      headers: {
        "Content-Type": contentType || "image/jpeg",
        "Cache-Control": "public, max-age=86400",
      },
    });
  } catch (err) {
    console.error("Error proxying image:", err);
    return new NextResponse("Server error", { status: 500 });
  }
}
