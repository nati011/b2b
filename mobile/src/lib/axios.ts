import axios from 'axios';

const getApiUrl = () => {
  const envUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8082';
  return envUrl;
};

const apiUrl = getApiUrl();

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
    // Get token from localStorage if available
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
axiosIns.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized - clear token and redirect to login
      localStorage.removeItem('auth_token');
    }
    return Promise.reject(error);
  }
);

export default axiosIns;

