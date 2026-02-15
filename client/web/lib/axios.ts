"use client";
import axios from "axios";
import { getSession } from "next-auth/react";

// Client-side API URL
// In browser on HTTPS (e.g. Vercel), we must not call HTTP APIs (mixed content blocked).
// Use same-origin proxy /api-backend when no HTTPS base URL is set.
const getApiUrl = () => {
  const envUrl = process.env.NEXT_PUBLIC_BASE_URL;

  if (envUrl) {
    if (typeof window !== 'undefined' && (envUrl.includes('backend:') || envUrl.includes('backend/'))) {
      const port = envUrl.includes(':8080') ? '8090' : '8090';
      return `http://localhost:${port}`;
    }
    return envUrl;
  }

  if (typeof window !== 'undefined') {
    return '/api-backend';
  }

  return process.env.BACKEND_API_URL ? `${process.env.BACKEND_API_URL}/api` : 'http://185.222.240.66/api';
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

// Client-side interceptor that gets token from localStorage (basic auth) or next-auth session
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
        
        // Skip auth for login endpoint
        if (config.url?.includes('/api/v1/auth/login')) {
            return config;
        }
        
        try {
            if (typeof window !== 'undefined') {
                // Use saved_email and saved_password for Basic Auth (matching mobile app)
                const savedEmail = localStorage.getItem('saved_email');
                const savedPassword = localStorage.getItem('saved_password');
                
                if (savedEmail && savedPassword) {
                    // Create Basic Auth header: base64(email:password)
                    const credentials = `${savedEmail}:${savedPassword}`;
                    const encodedCredentials = btoa(credentials);
                    config.headers.Authorization = `Basic ${encodedCredentials}`;
                } else {
                    // Fallback to Bearer token if Basic Auth credentials not available
                    const token = localStorage.getItem('access_token');
                    if (token) {
                        config.headers.Authorization = `Bearer ${token}`;
                    } else {
                        // Try NextAuth session (for backward compatibility)
                        try {
                            const session = await getSession();
                            // @ts-ignore - accessToken is added to session in auth.ts
                            if (session && session.accessToken) {
                                // @ts-ignore
                                config.headers.Authorization = `Bearer ${session.accessToken}`;
                            } else {
                                // Log when no auth is available
                                if (process.env.NODE_ENV === 'development') {
                                    console.log('ℹ️ No auth credentials available (this is OK for public endpoints)');
                                }
                            }
                        } catch (sessionError) {
                            // If session access fails, continue without auth header
                            if (process.env.NODE_ENV === 'development') {
                                console.log('ℹ️ No auth credentials available (this is OK for public endpoints)');
                            }
                        }
                    }
                }
            }
        } catch (error) {
            // If auth setup fails, continue without auth header
            console.warn('Failed to get auth credentials:', error);
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
    async (error) => {
        // Handle 401 Unauthorized - token might be expired
        if (error.response?.status === 401) {
            // Clear stored tokens on unauthorized
            if (typeof window !== 'undefined') {
                localStorage.removeItem('access_token');
                localStorage.removeItem('refresh_token');
                // Optionally redirect to login
                // window.location.href = '/login';
            }
        }
        
        // Enhanced error logging - only log specific error types with meaningful information
        // Skip general error logging to avoid empty object logs
        
        if (error?.code === 'ECONNABORTED') {
            console.error('⏱️ Request timeout:', error?.config?.url || 'Unknown URL');
        } else if (error?.code === 'ERR_NETWORK') {
            const networkErrorInfo: any = {};
            if (error?.config?.url) networkErrorInfo.url = error.config.url;
            networkErrorInfo.baseURL = apiUrl;
            if (error?.config?.baseURL && error?.config?.url) {
                networkErrorInfo.fullURL = error.config.baseURL + error.config.url;
            }
            if (error?.message) networkErrorInfo.message = error.message;
            
            console.error('🌐 Network error:', networkErrorInfo);
            
            // Check if it's a CORS issue
            if (error?.message && (error.message.includes('CORS') || error.message.includes('cross-origin'))) {
                console.error('🚫 CORS error detected! Backend may not be allowing cross-origin requests.');
            }
        } else if (error?.response) {
            const apiErrorInfo: any = {};
            let hasApiErrorData = false;
            
            if (error.response?.status) {
                apiErrorInfo.status = error.response.status;
                hasApiErrorData = true;
            }
            if (error.response?.statusText) {
                apiErrorInfo.statusText = error.response.statusText;
                hasApiErrorData = true;
            }
            if (error?.config?.url) {
                apiErrorInfo.url = error.config.url;
                hasApiErrorData = true;
            }
            // Only add data if it's not an empty object and has meaningful content
            if (error.response?.data) {
                const data = error.response.data;
                // Check if data is an empty object or has meaningful content
                if (typeof data === 'object' && data !== null) {
                    // Include if it's an array (even if empty, might be meaningful) or has keys
                    if (Array.isArray(data) || Object.keys(data).length > 0) {
                        apiErrorInfo.data = data;
                        hasApiErrorData = true;
                    }
                } else if (data !== '') {
                    // Include non-empty strings, numbers, booleans, etc.
                    apiErrorInfo.data = data;
                    hasApiErrorData = true;
                }
            }
            
            // Only log if we have meaningful data
            if (hasApiErrorData && Object.keys(apiErrorInfo).length > 0) {
                // Filter out falsy values and empty objects
                const filteredApiError = Object.fromEntries(
                    Object.entries(apiErrorInfo).filter(([_, val]) => {
                        if (val === undefined || val === null || val === '') {
                            return false;
                        }
                        // Filter out empty objects
                        if (typeof val === 'object' && !Array.isArray(val) && Object.keys(val).length === 0) {
                            return false;
                        }
                        return true;
                    })
                );
                
                // // Only log if we have meaningful data after filtering
                // if (Object.keys(filteredApiError).length > 0) {
                //     console.error('📡 API error response:', filteredApiError);
                // }
            }
        } else if (error?.request) {
            const requestErrorInfo: any = {};
            let hasRequestErrorData = false;
            
            if (error?.config?.url) {
                requestErrorInfo.url = error.config.url;
                hasRequestErrorData = true;
            }
            if (error?.request) {
                requestErrorInfo.request = error.request;
                hasRequestErrorData = true;
            }
            
            // Only log if we have meaningful data
            if (hasRequestErrorData && Object.keys(requestErrorInfo).length > 0) {
                const filteredRequestError = Object.fromEntries(
                    Object.entries(requestErrorInfo).filter(([_, val]) => 
                        val !== undefined && val !== null && val !== ''
                    )
                );
                
                if (Object.keys(filteredRequestError).length > 0) {
                    console.error('📭 No response received:', filteredRequestError);
                }
            }
        } else if (error?.message) {
            console.error('⚠️ Request setup error:', error.message, error?.config?.url || 'Unknown URL');
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
