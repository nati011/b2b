import { type NextRequest, NextResponse } from "next/server";

const BACKEND_BASE = process.env.BACKEND_API_URL || "https://api.efoyetastore.com";

export async function GET(request: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  return proxy(request, await params);
}

export async function POST(request: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  return proxy(request, await params);
}

export async function PUT(request: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  return proxy(request, await params);
}

export async function PATCH(request: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  return proxy(request, await params);
}

export async function DELETE(request: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  return proxy(request, await params);
}

export async function OPTIONS() {
  return new NextResponse(null, { status: 204 });
}

async function proxy(request: NextRequest, params: { path: string[] }) {
  const path = params.path?.length ? params.path.join("/") : "";
  const search = request.nextUrl.searchParams.toString();
  const url = `${BACKEND_BASE}/api/${path}${search ? `?${search}` : ""}`;

  const headers = new Headers();
  const forwardHeaders = ["authorization", "content-type", "accept", "x-requested-with", "x-user-id"];
  request.headers.forEach((value, key) => {
    if (forwardHeaders.includes(key.toLowerCase())) {
      headers.set(key, value);
    }
  });

  let body: string | undefined;
  try {
    body = await request.text();
  } catch {
    // no body
  }

  try {
    const res = await fetch(url, {
      method: request.method,
      headers,
      body: body || undefined,
      cache: "no-store",
    });

    const data = await res.text();
    const response = new NextResponse(data, {
      status: res.status,
      statusText: res.statusText,
    });
    res.headers.forEach((value, key) => {
      const lower = key.toLowerCase();
      if (lower.startsWith("access-control-") || lower === "content-type") {
        response.headers.set(key, value);
      }
    });
    return response;
  } catch (err) {
    console.error("[backend proxy]", url, err);
    return NextResponse.json(
      { error: "Backend unreachable", details: err instanceof Error ? err.message : String(err) },
      { status: 502 }
    );
  }
}
