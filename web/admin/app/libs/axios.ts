import { getSession } from "@/app/actions/getSession";
import axios, { AxiosError, AxiosResponse } from "axios";
import { signOut } from "next-auth/react";

// Extend the Session type to include our custom properties
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

// API configuration
const API_BASE_URL = process.env.NEXT_BASE_URL || "https://b2b-67gk.onrender.com/api/v1";

const axiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000, // 30 seconds timeout
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add authentication token
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

// Response interceptor to handle token refresh and errors
axiosInstance.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  async (error: AxiosError) => {
    const originalRequest = error.config;

    // Handle 401 Unauthorized errors
    if (error.response?.status === 401 && originalRequest) {
      // Check if this request has already been retried
      const retryCount = (originalRequest as any)._retryCount || 0;
      const MAX_RETRY_ATTEMPTS = 2; // Limit to 2 retry attempts

      if (retryCount < MAX_RETRY_ATTEMPTS) {
        try {
          // Mark this request as retried
          (originalRequest as any)._retryCount = retryCount + 1;

          // Try to refresh the token
          const session = await getSession() as ExtendedSession;

          if (session?.accessToken) {
            // Update the authorization header with the new token
            originalRequest.headers.Authorization = `Bearer ${session.accessToken}`;

            // Retry the original request
            return axiosInstance(originalRequest);
          }
        } catch (refreshError) {
          console.error("Token refresh failed:", refreshError);
        }
      } else {
        console.error("Maximum retry attempts exceeded for 401 error");
      }

      // If refresh fails or max retries exceeded, sign out the user
      await signOut({
        callbackUrl: '/auth/signin',
        redirect: true
      });
    }

    // Handle other errors
    if (error.response?.status === 403) {
      console.error("Forbidden: User doesn't have permission to access this resource");
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

// Export the configured instance
export default axiosInstance;

// Export a function to get the base URL for external use
export const getApiBaseUrl = () => API_BASE_URL;
