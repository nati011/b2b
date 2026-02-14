"use client";
import axios from "@/lib/axios";

const API_BASE = "/api/v1";

export interface User {
  id: string;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  status: string;
  role_ids: string[];
  created_at: string;
  updated_at: string;
}

export interface Role {
  id: string;
  name: string;
  description: string;
  permission_ids: string[];
  created_at: string;
  updated_at: string;
}

export interface Permission {
  id: string;
  resource_code: string;
  action: string;
  description: string;
  created_at: string;
  updated_at: string;
}

// Users
// Note: Backend automatically filters users by supplier when a supplier user is authenticated.
// The backend detects the supplier from the authenticated user's email and filters accordingly.
// No need to pass supplier_id explicitly - backend handles it securely.
export interface UserListResponse {
  items: User[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

export async function getUsers(page: number = 1, limit: number = 10): Promise<UserListResponse> {
  try {
    const response = await axios.get(`${API_BASE}/users?page=${page}&limit=${limit}`);
    const data = response.data;
    
    // Handle different response formats
    if (data.items && Array.isArray(data.items)) {
      return {
        items: data.items,
        total: data.total || data.items.length,
        page: data.page || page,
        limit: data.limit || limit,
        total_pages: data.total_pages || Math.ceil((data.total || data.items.length) / (data.limit || limit)),
        has_next: data.has_next || false,
        has_prev: data.has_prev || false,
      };
    } else if (data.body?.items && Array.isArray(data.body.items)) {
      return {
        items: data.body.items,
        total: data.body.total || data.body.items.length,
        page: data.body.page || page,
        limit: data.body.limit || limit,
        total_pages: data.body.total_pages || Math.ceil((data.body.total || data.body.items.length) / (data.body.limit || limit)),
        has_next: data.body.has_next || false,
        has_prev: data.body.has_prev || false,
      };
    } else if (Array.isArray(data.body)) {
      // Fallback for non-paginated response
      return {
        items: data.body,
        total: data.body.length,
        page: 1,
        limit: limit,
        total_pages: Math.ceil(data.body.length / limit),
        has_next: false,
        has_prev: false,
      };
    }
    
    return {
      items: [],
      total: 0,
      page: 1,
      limit: limit,
      total_pages: 0,
      has_next: false,
      has_prev: false,
    };
  } catch (error: any) {
    console.error("Error fetching users:", error);
    throw new Error(error.response?.data?.error || "Failed to fetch users");
  }
}

export async function getUserById(id: string): Promise<User> {
  try {
    const response = await axios.get(`${API_BASE}/users/${id}`);
    return response.data.body;
  } catch (error: any) {
    console.error("Error fetching user:", error);
    throw new Error(error.response?.data?.error || "Failed to fetch user");
  }
}

export async function createUser(data: {
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  password: string;
}): Promise<User> {
  try {
    const response = await axios.post(`${API_BASE}/users`, data);
    return response.data.body;
  } catch (error: any) {
    console.error("Error creating user:", error);
    throw new Error(error.response?.data?.error || "Failed to create user");
  }
}

export async function updateUser(id: string, data: Partial<User>): Promise<User> {
  try {
    const response = await axios.put(`${API_BASE}/users/${id}`, data);
    return response.data.body;
  } catch (error: any) {
    console.error("Error updating user:", error);
    throw new Error(error.response?.data?.error || "Failed to update user");
  }
}

export async function deleteUser(id: string): Promise<void> {
  try {
    await axios.delete(`${API_BASE}/users/${id}`);
  } catch (error: any) {
    console.error("Error deleting user:", error);
    throw new Error(error.response?.data?.error || "Failed to delete user");
  }
}

export async function assignRolesToUser(userId: string, roleIds: string[]): Promise<User> {
  try {
    const response = await axios.post(`${API_BASE}/users/${userId}/roles`, {
      role_ids: roleIds,
    });
    return response.data.body;
  } catch (error: any) {
    console.error("Error assigning roles:", error);
    throw new Error(error.response?.data?.error || "Failed to assign roles");
  }
}

export async function activateUser(id: string): Promise<User> {
  try {
    const response = await axios.post(`${API_BASE}/users/${id}/activate`);
    return response.data.body;
  } catch (error: any) {
    console.error("Error activating user:", error);
    throw new Error(error.response?.data?.error || "Failed to activate user");
  }
}

export async function deactivateUser(id: string): Promise<User> {
  try {
    const response = await axios.post(`${API_BASE}/users/${id}/deactivate`);
    return response.data.body;
  } catch (error: any) {
    console.error("Error deactivating user:", error);
    throw new Error(error.response?.data?.error || "Failed to deactivate user");
  }
}

export async function suspendUser(id: string): Promise<User> {
  try {
    const response = await axios.post(`${API_BASE}/users/${id}/suspend`);
    return response.data.body;
  } catch (error: any) {
    console.error("Error suspending user:", error);
    throw new Error(error.response?.data?.error || "Failed to suspend user");
  }
}

// Roles
export async function getRoles(): Promise<Role[]> {
  try {
    const response = await axios.get(`${API_BASE}/roles`);
    return response.data.body || [];
  } catch (error: any) {
    console.error("Error fetching roles:", error);
    throw new Error(error.response?.data?.error || "Failed to fetch roles");
  }
}

export async function getRoleById(id: string): Promise<Role> {
  try {
    const response = await axios.get(`${API_BASE}/roles/${id}`);
    return response.data.body;
  } catch (error: any) {
    console.error("Error fetching role:", error);
    throw new Error(error.response?.data?.error || "Failed to fetch role");
  }
}

export async function createRole(data: {
  name: string;
  description: string;
  permission_ids: string[];
}): Promise<Role> {
  try {
    const response = await axios.post(`${API_BASE}/roles`, data);
    return response.data.body;
  } catch (error: any) {
    console.error("Error creating role:", error);
    throw new Error(error.response?.data?.error || "Failed to create role");
  }
}

export async function updateRole(id: string, data: Partial<Role>): Promise<Role> {
  try {
    const response = await axios.put(`${API_BASE}/roles/${id}`, data);
    return response.data.body;
  } catch (error: any) {
    console.error("Error updating role:", error);
    throw new Error(error.response?.data?.error || "Failed to update role");
  }
}

export async function deleteRole(id: string): Promise<void> {
  try {
    await axios.delete(`${API_BASE}/roles/${id}`);
  } catch (error: any) {
    console.error("Error deleting role:", error);
    throw new Error(error.response?.data?.error || "Failed to delete role");
  }
}

// Permissions
export async function getPermissions(): Promise<Permission[]> {
  try {
    const response = await axios.get(`${API_BASE}/permissions`);
    return response.data.body || [];
  } catch (error: any) {
    console.error("Error fetching permissions:", error);
    throw new Error(error.response?.data?.error || "Failed to fetch permissions");
  }
}

export async function getPermissionById(id: string): Promise<Permission> {
  try {
    const response = await axios.get(`${API_BASE}/permissions/${id}`);
    return response.data.body;
  } catch (error: any) {
    console.error("Error fetching permission:", error);
    throw new Error(error.response?.data?.error || "Failed to fetch permission");
  }
}

export async function createPermission(data: {
  resource_code: string;
  action: string;
  description: string;
}): Promise<Permission> {
  try {
    const response = await axios.post(`${API_BASE}/permissions`, data);
    return response.data.body;
  } catch (error: any) {
    console.error("Error creating permission:", error);
    throw new Error(error.response?.data?.error || "Failed to create permission");
  }
}

export async function updatePermission(id: string, data: Partial<Permission>): Promise<Permission> {
  try {
    const response = await axios.put(`${API_BASE}/permissions/${id}`, data);
    return response.data.body;
  } catch (error: any) {
    console.error("Error updating permission:", error);
    throw new Error(error.response?.data?.error || "Failed to update permission");
  }
}

export async function deletePermission(id: string): Promise<void> {
  try {
    await axios.delete(`${API_BASE}/permissions/${id}`);
  } catch (error: any) {
    console.error("Error deleting permission:", error);
    throw new Error(error.response?.data?.error || "Failed to delete permission");
  }
}

