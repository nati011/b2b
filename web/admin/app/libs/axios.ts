import { getSession } from "@/app/actions/getSession";
import axios, { AxiosError, AxiosResponse } from "axios";
import { signOut } from "next-auth/react";


interface ExtendedSession {
  user?: {
    id: string;
    name: string;
    email: string;
    username: string;
    roles: string[];
  } | null;
  accessToken?: string;
  error?: string;
  expires: string;
}

const API_BASE_URL = process.env.NEXT_PUBLIC_BASE_URL || "https://b2b-67gk.onrender.com/api/v1";

const axiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

axiosInstance.interceptors.request.use(
  async (config) => {
    try {
      const session = await getSession() as ExtendedSession;
      if (session?.accessToken) {
        config.headers.Authorization = `Bearer ${session.accessToken}`;
      }
      return config;
    } catch (error) {
      console.error("Error in request interceptor:", error);
      return config;
    }
  },
  (error) => {
    console.error("Request interceptor error:", error);
    return Promise.reject(error);
  }
);

axiosInstance.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  async (error: AxiosError) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && originalRequest) {
      const retryCount = (originalRequest as any)._retryCount || 0;
      const MAX_RETRY_ATTEMPTS = 2;

      if (retryCount < MAX_RETRY_ATTEMPTS) {
        try {
          (originalRequest as any)._retryCount = retryCount + 1;
          
          const session = await getSession() as ExtendedSession;
          if (session?.accessToken) {
            originalRequest.headers.Authorization = `Bearer ${session.accessToken}`;
            return axiosInstance(originalRequest);
          }
        } catch (refreshError) {
          console.error("Token refresh failed:", refreshError);
        }
      } else {
        console.error("Maximum retry attempts exceeded for 401 error");
      }

      if (typeof window !== 'undefined') {
        try {
          await signOut({
            callbackUrl: '/auth/signin',
            redirect: true
          });
        } catch (signOutError) {
          console.error("SignOut failed, using fallback redirect:", signOutError);
          window.location.replace('/auth/signin');
        }
      }
    }

    if (error.response?.status === 403) {
      console.error("Forbidden: User doesn't have permission to access this resource");
      
      const forbiddenError = new Error("FORBIDDEN_ACCESS");
      forbiddenError.name = "ForbiddenError";
      (forbiddenError as any).redirectTo = "/forbidden";
      (forbiddenError as any).statusCode = 403;
      
      return Promise.reject(forbiddenError);
    }

    if (error.response?.status === 404) {
      console.error("Resource not found");
    }

    if (error.response?.status && error.response.status >= 500) {
      console.error("Server error:", error.response.status);
    }

    return Promise.reject(error);
  }
);

export default axiosInstance;
export const getApiBaseUrl = () => API_BASE_URL;