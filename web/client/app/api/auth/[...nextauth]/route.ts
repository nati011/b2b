import NextAuth, { AuthOptions, Session, TokenSet } from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";
import GoogleProvider from "next-auth/providers/google";
import { jwtDecode } from "jwt-decode";
import axios from "axios";

const baseURL = process.env.NEXT_BASE_URL || "https://b2b-67gk.onrender.com";

interface KeycloakJWT {
  exp: number;
  iat: number;
  sub: string;
  email: string;
  name: string;
  preferred_username: string;
  given_name: string;
  family_name: string;
  realm_access: {
    roles: string[];
  };
}

interface UserToken {
  id: string;
  name: string;
  email: string;
  username: string;
  roles: string[];
}

interface AppToken extends TokenSet {
  accessToken: string;
  refreshToken: string;
  accessTokenExpires: number;
  user: UserToken;
  error?: string;
}

async function linkAccount(profile: any) {
  try {
    return {};
  } catch (error) {
    throw new Error("User not found");
  }
}

async function refreshAccessToken(token: AppToken): Promise<AppToken> {
  try {
    const response = await axios.post(`${baseURL}/api/v1/auth/refresh`, {
      refresh_token: token.refreshToken,
    });

    const decoded = jwtDecode<KeycloakJWT>(response.data.body.access_token);

    return {
      ...token,
      accessToken: response.data.body.access_token,
      refreshToken: response.data.body.refresh_token,
      accessTokenExpires: decoded.exp * 1000,
      user: {
        id: decoded.sub,
        name: decoded.name,
        email: decoded.email,
        username: decoded.preferred_username,
        roles: decoded.realm_access?.roles || [],
      },
      error: undefined,
    };
  } catch (error) {
    console.error("Refresh token error:", error);

    if (
      axios.isAxiosError(error) &&
      (error.response?.status === 400 || error.response?.status === 401)
    ) {
      return {
        ...token,
        error: "RefreshAccessTokenError",
        accessToken: "",
        refreshToken: "",
        accessTokenExpires: 0,
        user: {
          id: "",
          name: "",
          email: "",
          username: "",
          roles: [],
        },
      };
    }

    return {
      ...token,
      error: "RefreshAccessTokenError",
    };
  }
}

export const authOptions: AuthOptions = {
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID as string,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET as string,
      // @ts-ignore
      profile: (profile, tokens) => {
        console.log("_________Profile_______________");
        console.log(profile);
        console.log("_________Tokens_______________");
        console.log(tokens);
        return {};
      },
    }),
    CredentialsProvider({
      name: "Credentials",
      credentials: {
        email: { label: "Email", type: "email" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        try {
          const response = await axios.post(
            `${baseURL}/api/v1/auth/login`,
            credentials
          );

          if (response.status !== 202) {
            throw new Error(response.data?.message || "Authentication failed");
          }

          const decoded = jwtDecode<KeycloakJWT>(
            response.data.body.access_token
          );

          return {
            ...response.data,
            user: {
              id: decoded.sub,
              name: decoded.name,
              email: decoded.email,
              username: decoded.preferred_username,
              roles: decoded.realm_access?.roles || [],
            },
          };
        } catch (error) {
          if (axios.isAxiosError(error)) {
            console.error("Authentication error:", error.response?.data);
            throw new Error(
              error.response?.data?.message || "Authentication failed"
            );
          }
          console.error("Authentication error:", error);
          return null;
        }
      },
    }),
  ],
  secret: process.env.NEXTAUTH_SECRET,
  pages: {
    signIn: "/login",
    error: "/login",
  },
  callbacks: {
    async redirect({ url, baseUrl }) {
      if (url.startsWith("/")) return `${baseUrl}${url}`;
      else if (new URL(url).origin === baseUrl) return url;
      return baseUrl;
    },
    async jwt({ token, user, account }) {
      if (user && account) {
        return {
          accessToken: (user as any).body.access_token,
          refreshToken: (user as any).body.refresh_token,
          accessTokenExpires:
            jwtDecode<KeycloakJWT>((user as any).body.access_token).exp * 1000,
          user: (user as any).user,
        };
      }

      if (
        token.error != "RefreshAccessTokenError" &&
        Date.now() > (token as AppToken).accessTokenExpires
      ) {
        const refreshedToken = await refreshAccessToken(token as AppToken);
        if (
          refreshedToken.error === "RefreshAccessTokenError" &&
          !refreshedToken.accessToken
        ) {
          return {
            ...refreshedToken,
            user: {
              id: "",
              name: "",
              email: "",
              username: "",
              roles: [],
            },
          };
        }
        return refreshedToken;
      }

      return token;
    },
    async session({ session, token }) {
      console.log("______________________");
      console.log(session, token);
      if (
        (token as AppToken).error === "RefreshAccessTokenError" &&
        !(token as AppToken).accessToken
      ) {
        // @ts-ignore
        session.user = null;
        session.expires = new Date(0).toISOString();
      } else {
        session.user = (token as AppToken).user;
        // @ts-ignore
        session.accessToken = (token as AppToken).accessToken;
      }

      // @ts-ignore
      session.error = (token as AppToken).error;
      return session;
    },
  },
  session: {
    strategy: "jwt",
    maxAge: 72 * 60 * 60,
  },
  debug: process.env.NODE_ENV !== "production",
};

const handler = NextAuth(authOptions);

export { handler as GET, handler as POST };
