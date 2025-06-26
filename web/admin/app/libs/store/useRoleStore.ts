import { create } from "zustand";
import axiosIns from "@/app/libs/axios";
import { Role } from "@/app/libs/types";

interface RoleStore {
    success: string;
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;
    roles: Role[];
    rolesLoading: boolean;
    rolesError: string | null;

    fetchRoles: () => Promise<void>;
    createRole: (data: Partial<Role>) => Promise<void>;
    deleteRole: (id: number) => Promise<void>;
    editRole: (data: Partial<Role>) => Promise<void>;
}

const useRoleStore = create<RoleStore>((set) => ({
    success: "",
    loading: false,
    error: null,
    next: null,
    previous: null,
    roles: [],
    rolesLoading: false,
    rolesError: null,


    fetchRoles: async () => {
        set({ rolesLoading: true, rolesError: null });
        try {
            const response = await axiosIns.get("/role");
            set({
                roles: response.data.body.roles,
                rolesLoading: false,
            });
        } catch (error) {
            set({
                rolesError: "Failed to fetch roles",
                rolesLoading: false,
            });
        }
    },
    createRole: async (data: Partial<Role>) => {
        set({ rolesLoading: true, rolesError: null });
        try {
            const response = await axiosIns.post(
                "/api/role",
                data,
            );
            await useRoleStore.getState().fetchRoles();
            set({ rolesLoading: false });
        } catch (error: any) {
            set({ rolesLoading: false, rolesError: error.message });
        }
    },
    deleteRole: async (id: number) => {
        set({ rolesLoading: true, rolesError: null });
        try {
            await axiosIns.delete(`/api/role?id=${id}`);
            set({ rolesLoading: false });
            await useRoleStore.getState().fetchRoles();
        } catch (error: any) {
            set({ rolesLoading: false, rolesError: error.message });
        } 0
    },
    editRole: async (data: Partial<Role>) => {
        set({ rolesLoading: true, rolesError: null });
        try {
            console.log(data)
            await axiosIns.put(`/api/role`, data);
            set({ rolesLoading: false });
            await useRoleStore.getState().fetchRoles();
        } catch (error: any) {
            set({ rolesLoading: false, rolesError: error.message });
        }
    }
}));

export default useRoleStore;
