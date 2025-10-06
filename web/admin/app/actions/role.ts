'use server'
import axiosIns from "@/app/libs/axios";
import { Role } from "@/app/libs/types";

export const fetchRoles = async () => {
    try {
        const response = await axiosIns.get("/role");
        console.log(response.data)
        return response.data;
    } catch (error) {
        throw error;
    }
}

export const fetchRoleById = async (id: number) => {
    try {
        const response = await axiosIns.get(`/role?id=${id}`);
        return response.data;
    } catch (error) {
        throw error;
    }
}

export const createRole = async (roleData: Partial<Role>) => {
    try {
        const response = await axiosIns.post("/role", roleData);
        return response.data;
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "An error occurred while creating the role";
        }
        throw "An error occurred while creating the role";
    }
}

export const updateRole = async (roleData: Partial<Role>) => {
    try {
        const response = await axiosIns.put("/role", roleData);
        return response.data;
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "An error occurred while updating the role";
        }
        throw "An error occurred while updating the role";
    }
}

export const deleteRole = async (id: number) => {
    try {
        const response = await axiosIns.delete(`/role?id=${id}`);
        return response.data;
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "An error occurred while deleting the role";
        }
        throw "An error occurred while deleting the role";
    }
}

export const updateRoleStatus = async (id: number, command: string) => {
    try {
        const response = await axiosIns.patch(`/role/${id}/status?command=${command}`);
        return response.data;
    } catch (error) {
        throw error;
    }
} 