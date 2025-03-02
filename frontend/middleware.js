import { NextResponse } from "next/server";
import { jwtVerify } from "jose";

// // Secret key (same as the backend)
// const SECRET_KEY = new TextEncoder().encode(process.env.JWT_SECRET); // Store in .env.local

// ✅ Function to verify JWT token
async function verifyToken(token) {
  try {
    const { payload } = await jwtVerify(token);
    return payload; // Return decoded payload if valid
  } catch (error) {
    return null; // Invalid token
  }
}

export async function middleware(req) {
  const { pathname } = req.nextUrl;
  const authToken = req.cookies.get("auth_token")?.value;

  //for dev purpose only
  return NextResponse.next();

  if (pathname.startsWith("/login")) {
    if (!authToken) return NextResponse.next();
    return NextResponse.redirect(new URL("/codePreview", req.url));
  }

  // If no token → Redirect to login
  if (!authToken) {
    return NextResponse.redirect(new URL("/login", req.url));
  }

  // ✅ Token is valid → Allow access
  return NextResponse.next();
}

// ✅ Protect specific routes
export const config = {
  matcher: ["/((?!api|_next|favicon.ico).*)"], // Protects everything except specified paths
};
