"use client";
import axios from "axios";
import { getSession } from "next-auth/react";

// Client-side API URL
// When running in browser, we need to use localhost with external port
// Docker service names (like 'backend:8080') don't work in browser
// Check if we're in browser and if URL contains Docker service name, use localhost instead
const getApiUrl = () => {
  const envUrl = process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:8082';
  
  // If running in browser and URL contains Docker service name, use localhost
  if (typeof window !== 'undefined') {
    // Browser environment - can't use Docker service names
    if (envUrl.includes('backend:') || envUrl.includes('backend/')) {
      // Extract port from env URL or use default 8082
      const port = envUrl.includes(':8080') ? '8082' : '8082';
      return `http://localhost:${port}`;
    }
  }
  
  return envUrl;
};

const apiUrl = getApiUrl();

const axiosIns = axios.create({
    baseURL: apiUrl,
    timeout: 15000, // Increased timeout to 15 seconds
    headers: {
        'Content-Type': 'application/json',
    },
});

// Log the API URL being used (only in development)
if (process.env.NODE_ENV === 'development') {
    console.log('Axios configured with baseURL:', apiUrl);
}

// Client-side interceptor that gets token from next-auth session
axiosIns.interceptors.request.use(
    async (config) => {
        // Log request details in development
        if (process.env.NODE_ENV === 'development') {
            const fullURL = config.baseURL ? `${config.baseURL}${config.url}` : config.url;
            console.log('📤 Axios request:', {
                method: config.method,
                url: config.url,
                baseURL: config.baseURL,
                fullURL: fullURL,
                headers: config.headers,
            });
        }
        
        try {
            // Use next-auth's getSession which works client-side
            const session = await getSession();
            // @ts-ignore - accessToken is added to session in auth.ts
            if (session && session.accessToken) {
                // @ts-ignore
                config.headers.Authorization = `Bearer ${session.accessToken}`;
            } else {
                // Log when no auth token is available (this endpoint doesn't require auth)
                if (process.env.NODE_ENV === 'development') {
                    console.log('ℹ️ No auth token available (this is OK for /api/v1/plan)');
                }
            }
        } catch (error) {
            // If session access fails, continue without auth header
            console.warn('Failed to get auth token:', error);
        }

        return config;
    },
    (error) => {
        console.error('❌ Request interceptor error:', error);
        return Promise.reject(error);
    }
);

// Response interceptor for better error handling and logging
axiosIns.interceptors.response.use(
    (response) => {
        // Log successful responses in development
        if (process.env.NODE_ENV === 'development') {
            console.log('✅ Axios response received:', {
                url: response.config?.url,
                status: response.status,
                statusText: response.statusText,
                data: response.data,
            });
        }
        return response;
    },
    (error) => {
        // Enhanced error logging
        console.error('❌ Axios request failed:', {
            url: error.config?.url,
            method: error.config?.method,
            baseURL: error.config?.baseURL,
            code: error.code,
            message: error.message,
        });
        
        if (error.code === 'ECONNABORTED') {
            console.error('⏱️ Request timeout:', error.config?.url);
        } else if (error.code === 'ERR_NETWORK') {
            console.error('🌐 Network error:', {
                url: error.config?.url,
                baseURL: apiUrl,
                fullURL: error.config?.baseURL + error.config?.url,
                message: error.message,
            });
            // Check if it's a CORS issue
            if (error.message?.includes('CORS') || error.message?.includes('cross-origin')) {
                console.error('🚫 CORS error detected! Backend may not be allowing cross-origin requests.');
            }
        } else if (error.response) {
            console.error('📡 API error response:', {
                status: error.response.status,
                statusText: error.response.statusText,
                url: error.config?.url,
                data: error.response.data,
                headers: error.response.headers,
            });
        } else if (error.request) {
            console.error('📭 No response received:', {
                url: error.config?.url,
                request: error.request,
            });
        } else {
            console.error('⚠️ Request setup error:', error.message, error.config?.url);
        }
        return Promise.reject(error);
    }
);

/**
 * Create an axios instance with auth token
 * Use this in components where you have access to useSession hook
 */
export function createAuthenticatedAxios(token: string) {
    const instance = axios.create({
        baseURL: apiUrl,
        timeout: 15000, // Increased timeout to 15 seconds
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
    });
    return instance;
}

export default axiosIns;
