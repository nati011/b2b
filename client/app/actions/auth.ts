import axiosIns from "@/lib/axios";
import { SupplierRequest, RegisterRequest, User } from "@/lib/types";
import { jwtDecode } from "jwt-decode";

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token?: string;
}

export interface LoginResponseWrapper {
  body: LoginResponse;
}

export const RegisterCustomer = async (profile: RegisterRequest) => {
  try {
    const response = await axiosIns.post("/api/v1/customer", profile);
    console.log(response.data);
    return response.data.message;
  } catch (error: any) {
    if (error.response) {
      throw new Error(
        error.response.data.message ||
        "An error has occured while creating account"
      );
    }
    throw new Error("An error has occured while creating account");
  }
};

export const RegisterSupplier = async(profile: SupplierRequest) => {
  try {
    // Filter out fields that backend doesn't expect
    const requestPayload = {
      tin: profile.tin,
      latitude: profile.latitude,
      longitude: profile.longitude,
      general_zone: profile.general_zone,
      region: profile.region,
      woreda: profile.woreda,
      licence_url: profile.licence_url || "",
      first_name: profile.first_name,
      last_name: profile.last_name,
      email: profile.email,
      phone: profile.phone,
      password: profile.password,
      // username will be set by backend from phone
    };
    const response = await axiosIns.post("/api/v1/supplier", requestPayload);
    console.log(response.data);
    // Backend returns user_id as supplier id
    return response.data.body.supplier;
  } catch (error: any) {
    if (error.response) {
      throw new Error(
        error.response.data.message ||
        "An error has occured while creating account"
      );
    }
    throw new Error("An error has occured while creating account");
  }
}

export const InitResetPassword = async (email: string) => {
  try {
    console.log(email);
    const response = await axiosIns.post("/user/init_reset", {
      email: email,
    });
    return response.data.message;
  } catch (error: any) {
    if (error.response) {
      throw new Error(
        error.response.data.message ||
        "An error has occured while reseting the password"
      );
    }
    throw new Error("An error has occured while reseting the password");
  }
};

export const ResetPassword = async (token: string, password: string) => {
  try {
    const response = await axiosIns.post(`/auth/reset/${token}`, {
      password: password,
    });
    return response.data.message;
  } catch (error: any) {
    console.log(error);
    if (error.response) {
      throw new Error(
        error.response.data.message ||
        "An error has occured while creating the product"
      );
    }
    throw new Error("An error has occured while creating the product");
  }
};

export const UpdateProfile = async (data: Partial<User>) => {
  try {
    const response = await axiosIns.patch("/user", data);
    console.log(response.data);
    return response.data.message;
  } catch (error: any) {
    console.log(error);
    if (error.response) {
      throw new Error(
        error.response.data.message ||
        "An error has occured while updating your profile"
      );
    }
    throw new Error("An error has occured while updating your profile");
  }
};

/**
 * Login with email and password using basic auth
 * Stores the access token in localStorage for use in axios interceptor
 */
export const Login = async (credentials: LoginRequest): Promise<LoginResponse> => {
  try {
    // Make request without auth header (login endpoint is public)
    const response = await axiosIns.post<LoginResponseWrapper>(
      "/api/v1/auth/login",
      credentials
    );

    if (response.status !== 202) {
      throw new Error(response.data?.body?.access_token ? "Login successful but unexpected status" : "Authentication failed");
    }

    const loginResponse = response.data.body;

    if (!loginResponse.access_token) {
      throw new Error("No access token received");
    }

    // Store tokens and user info in localStorage
    if (typeof window !== 'undefined') {
      localStorage.setItem('access_token', loginResponse.access_token);
      if (loginResponse.refresh_token) {
        localStorage.setItem('refresh_token', loginResponse.refresh_token);
      }
      // Store credentials for auto-fill on next login
      localStorage.setItem('saved_email', credentials.email);
      localStorage.setItem('saved_password', credentials.password);
      
      // Decode JWT token to extract user information
      try {
        const decoded = jwtDecode<{
          sub: string;
          name?: string;
          email?: string;
          preferred_username?: string;
          realm_access?: { roles?: string[] };
        }>(loginResponse.access_token);
        
        // Store user information for session use
        if (decoded.name) {
          localStorage.setItem('user_name', decoded.name);
        }
        if (decoded.email) {
          localStorage.setItem('user_email', decoded.email);
        }
        if (decoded.sub) {
          localStorage.setItem('user_id', decoded.sub);
        }
        if (decoded.preferred_username) {
          localStorage.setItem('user_username', decoded.preferred_username);
        }
        if (decoded.realm_access?.roles) {
          localStorage.setItem('user_roles', JSON.stringify(decoded.realm_access.roles));
        }
      } catch (error) {
        console.error('Error decoding JWT token:', error);
        // If decoding fails, still store basic info from credentials
        localStorage.setItem('user_email', credentials.email);
      }
      
      // Dispatch custom event to notify AuthContext of login
      window.dispatchEvent(new Event('auth-state-changed'));
    }

    return loginResponse;
  } catch (error: any) {
    if (error.response) {
      const errorMessage = error.response.data?.message || 
                          error.response.data?.error || 
                          "Invalid email or password";
      throw new Error(errorMessage);
    }
    if (error.request) {
      throw new Error("Network error: Unable to reach the server");
    }
    throw new Error(error.message || "An error occurred during login");
  }
};

/**
 * Logout - removes stored tokens and optionally credentials
 */
export const Logout = (clearCredentials: boolean = false) => {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user_name');
    localStorage.removeItem('user_email');
    localStorage.removeItem('user_id');
    localStorage.removeItem('user_username');
    localStorage.removeItem('user_roles');
    if (clearCredentials) {
      localStorage.removeItem('saved_email');
      localStorage.removeItem('saved_password');
    }
    // Dispatch custom event to notify AuthContext of logout
    window.dispatchEvent(new Event('auth-state-changed'));
  }
};

/**
 * Get saved credentials from localStorage
 */
export const getSavedCredentials = (): { email: string; password: string } | null => {
  if (typeof window !== 'undefined') {
    const email = localStorage.getItem('saved_email');
    const password = localStorage.getItem('saved_password');
    if (email && password) {
      return { email, password };
    }
  }
  return null;
};

/**
 * Clear saved credentials from localStorage
 */
export const clearSavedCredentials = () => {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('saved_email');
    localStorage.removeItem('saved_password');
  }
};

/**
 * Get stored access token
 */
export const getAccessToken = (): string | null => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('access_token');
  }
  return null;
};

// Subscription payment is no longer handled client-side.