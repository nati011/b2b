import axiosIns from '../axios';

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token?: string;
}

export interface AuthResponse {
  body: LoginResponse;
}

// JWT payload interface
interface JWTPayload {
  sub: string;
  name: string;
  email: string;
  preferred_username: string;
  realm_access?: {
    roles: string[];
  };
  exp?: number;
  iat?: number;
}

// Decode JWT token (simple base64 decode, no verification)
const decodeJWT = (token: string): JWTPayload | null => {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) {
      return null;
    }
    const payload = parts[1];
    // Add padding if needed
    const paddedPayload = payload + '='.repeat((4 - (payload.length % 4)) % 4);
    const decoded = atob(paddedPayload.replace(/-/g, '+').replace(/_/g, '/'));
    return JSON.parse(decoded) as JWTPayload;
  } catch (error) {
    console.error('Failed to decode JWT:', error);
    return null;
  }
};

export const Login = async (credentials: LoginRequest): Promise<LoginResponse> => {
  try {
    const response = await axiosIns.post<AuthResponse>('/api/v1/auth/login', credentials);
    
    if (response.status !== 202) {
      throw new Error('Authentication failed');
    }

    // Extract token from response body
    const token = response.data.body?.access_token;
    if (!token) {
      throw new Error('No access token received');
    }

    // Decode JWT to extract user information
    const payload = decodeJWT(token);
    if (payload) {
      // Store user ID (sub) as customer_id for compatibility
      if (payload.sub) {
        localStorage.setItem('user_id', payload.sub);
        // Try to extract numeric customer ID if it's a UUID
        // For now, we'll store the UUID and let the backend handle it
      }
      
      // Store user email and name for quick access
      if (payload.email) {
        localStorage.setItem('user_email', payload.email);
      }
      if (payload.name) {
        localStorage.setItem('user_name', payload.name);
      }
    }

    // Store token in localStorage
    localStorage.setItem('auth_token', token);
    
    // Store refresh token if available
    if (response.data.body?.refresh_token) {
      localStorage.setItem('refresh_token', response.data.body.refresh_token);
    }

    return response.data.body;
  } catch (error: any) {
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while logging in';
      throw new Error(errorMessage);
    }
    throw new Error('An error occurred while logging in');
  }
};

export interface GoogleSSORequest {
  token: string;
  first_name: string;
  last_name: string;
  email: string;
}

export interface SSOResponse {
  body: LoginResponse;
}

export const GoogleSSO = async (ssoData: GoogleSSORequest): Promise<LoginResponse> => {
  try {
    const response = await axiosIns.post<SSOResponse>('/api/v1/auth/sso', ssoData);
    
    if (response.status !== 202) {
      throw new Error('SSO authentication failed');
    }

    // Extract token from response body
    const token = response.data.body?.access_token;
    if (!token) {
      throw new Error('No access token received');
    }

    // Store token in localStorage
    localStorage.setItem('auth_token', token);
    
    // Store refresh token if available
    if (response.data.body?.refresh_token) {
      localStorage.setItem('refresh_token', response.data.body.refresh_token);
    }

    return response.data.body;
  } catch (error: any) {
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while signing in with Google';
      throw new Error(errorMessage);
    }
    throw new Error('An error occurred while signing in with Google');
  }
};

export const Logout = () => {
  localStorage.removeItem('auth_token');
  localStorage.removeItem('refresh_token');
  localStorage.removeItem('customer_id');
  localStorage.removeItem('customer_data'); // Clear customer data
  localStorage.removeItem('user_id');
  localStorage.removeItem('user_email');
  localStorage.removeItem('user_password'); // Clear password
  localStorage.removeItem('user_name');
};

