import axios from 'axios';

const getApiUrl = () => {
  const envUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8090';
  return envUrl;
};

const apiUrl = getApiUrl();

// Log API URL in development
if (import.meta.env.DEV) {
  console.log('📱 Mobile API Base URL:', apiUrl);
}

const axiosIns = axios.create({
  baseURL: apiUrl,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for auth token
axiosIns.interceptors.request.use(
  async (config) => {
    // Get user credentials from localStorage for Basic Auth
    const userEmail = localStorage.getItem('user_email');
    const userPassword = localStorage.getItem('user_password');
    
    // If we have email and password, use Basic Auth
    if (userEmail && userPassword) {
      // Create Basic Auth header: base64(email:password)
      const credentials = `${userEmail}:${userPassword}`;
      const encodedCredentials = btoa(credentials);
      config.headers.Authorization = `Basic ${encodedCredentials}`;
    } else {
      // Fallback to Bearer token if Basic Auth credentials not available
      const token = localStorage.getItem('auth_token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    
    // Log request in development
    if (import.meta.env.DEV) {
      const fullURL = config.baseURL ? `${config.baseURL}${config.url}` : config.url;
      console.log('📤 Mobile API Request:', {
        method: config.method?.toUpperCase(),
        url: config.url,
        baseURL: config.baseURL,
        fullURL: fullURL,
        hasBasicAuth: !!(userEmail && userPassword),
        hasBearerToken: !!localStorage.getItem('auth_token'),
      });
    }
    
    return config;
  },
  (error) => {
    console.error('❌ Mobile API Request Error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor
axiosIns.interceptors.response.use(
  (response) => {
    // Log successful responses in development
    if (import.meta.env.DEV) {
      console.log('✅ Mobile API Response:', {
        url: response.config?.url,
        status: response.status,
        statusText: response.statusText,
      });
    }
    return response;
  },
  (error) => {
    // Enhanced error logging
    if (error.response) {
      console.error('❌ Mobile API Error Response:', {
        url: error.config?.url,
        status: error.response.status,
        statusText: error.response.statusText,
        data: error.response.data,
      });
    } else if (error.request) {
      console.error('❌ Mobile API No Response:', {
        url: error.config?.url,
        baseURL: error.config?.baseURL,
        fullURL: error.config?.baseURL ? `${error.config.baseURL}${error.config.url}` : error.config?.url,
        message: error.message,
        code: error.code,
      });
      
      // Check for network/CORS issues
      if (error.code === 'ERR_NETWORK' || error.message?.includes('Network Error')) {
        console.error('🌐 Network Error - Possible causes:');
        console.error('  1. Backend server is not running');
        console.error('  2. CORS is not configured on backend');
        console.error('  3. Wrong API URL:', apiUrl);
        console.error('  4. Firewall blocking the connection');
      }
    } else {
      console.error('❌ Mobile API Request Setup Error:', error.message);
    }
    
    // Don't automatically logout on 401 - let components handle it
    // Credentials are persisted and will be reused
    return Promise.reject(error);
  }
);

export default axiosIns;

