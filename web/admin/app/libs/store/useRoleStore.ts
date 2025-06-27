import { create } from "zustand";
import axiosIns from "@/app/libs/axios";
import { Role } from "@/app/libs/types";
import { Create, Delete, GetAll, Update } from "@/app/actions/roles";

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
            const response = await GetAll();
            console.log(response.data)
            set({
                roles: response,
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
            const response = await Create(data)
            await useRoleStore.getState().fetchRoles();
            set({ rolesLoading: false, success: response });
        } catch (error: any) {
            set({ rolesLoading: false, rolesError: error.message });
        }
    },
    deleteRole: async (id: number) => {
        set({ rolesLoading: true, rolesError: null });
        try {
            await Delete(id)
            set({ rolesLoading: false, success:"Role deleted successfully" });
            await useRoleStore.getState().fetchRoles();
        } catch (error: any) {
            set({ rolesLoading: false, rolesError: error.message });
        } 
    },
    editRole: async (data: Partial<Role>) => {
        set({ rolesLoading: true, rolesError: null });
        try {
            await Update(data)
            set({ rolesLoading: false,  success:"Role updated successfully" });
            await useRoleStore.getState().fetchRoles();
        } catch (error: any) {
            set({ rolesLoading: false, rolesError: error.message });
        }
    }
}));

export default useRoleStore;
