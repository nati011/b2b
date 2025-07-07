import { authOptions } from "@/lib/auth";
import NextAuth, { AuthOptions, Session, TokenSet } from "next-auth";

const handler = NextAuth(authOptions);

export { handler as GET, handler as POST };
