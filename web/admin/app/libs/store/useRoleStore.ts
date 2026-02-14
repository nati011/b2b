import { create } from "zustand";
import { Role } from "@/app/libs/types";
import { toast } from "sonner";
import { 
    fetchRoles, 
    createRole, 
    deleteRole, 
    updateRole, 
    fetchRoleById,
    updateRoleStatus 
} from "@/app/actions/role";

interface RoleStore {
    success: string;
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;
    roles: Role[];
    rolesLoading: boolean;
    rolesError: string | null;
    selectedRole: Role | null;

    fetchRoles: () => Promise<void>;
    fetchRoleById: (id: number) => Promise<void>;
    createRole: (data: Partial<Role>) => Promise<void>;
    deleteRole: (id: number) => Promise<void>;
    editRole: (data: Partial<Role>) => Promise<void>;
    updateRoleStatus: (id: number, command: string) => Promise<void>;
    setSelectedRole: (role: Role | null) => void;
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
    selectedRole: null,

    fetchRoles: async () => {
        set({ rolesLoading: true, rolesError: null });
        try {
            const response = await fetchRoles();
            set({
                roles: response.body?.roles || response.body || [],
                rolesLoading: false,
            });
        } catch (error: any) {
            set({
                rolesError: error.message || "Failed to fetch roles",
                rolesLoading: false,
            });
            toast.error(error.message || "Failed to fetch roles");
        }
    },

    fetchRoleById: async (id: number) => {
        set({ loading: true, error: null });
        try {
            const response = await fetchRoleById(id);
            set({
                selectedRole: response.body || response,
                loading: false,
            });
        } catch (error: any) {
            set({
                error: error.message || "Failed to fetch role",
                loading: false,
            });
            toast.error(error.message || "Failed to fetch role");
        }
    },

    createRole: async (data: Partial<Role>) => {
        set({ rolesLoading: true, rolesError: null });
        try {
            await createRole(data);
            await useRoleStore.getState().fetchRoles();
            set({ rolesLoading: false });
            toast.success("Role created successfully");
        } catch (error: any) {
            set({ rolesLoading: false, rolesError: error.message });
            toast.error(error.message || "Failed to create role");
        }
    },

    deleteRole: async (id: number) => {
        set({ rolesLoading: true, rolesError: null });
        try {
            await deleteRole(id);
            await useRoleStore.getState().fetchRoles();
            set({ rolesLoading: false });
            toast.success("Role deleted successfully");
        } catch (error: any) {
            set({ rolesLoading: false, rolesError: error.message });
            toast.error(error.message || "Failed to delete role");
        }
    },

    editRole: async (data: Partial<Role>) => {
        set({ rolesLoading: true, rolesError: null });
        try {
            await updateRole(data);
            await useRoleStore.getState().fetchRoles();
            set({ rolesLoading: false });
            toast.success("Role updated successfully");
        } catch (error: any) {
            set({ rolesLoading: false, rolesError: error.message });
            toast.error(error.message || "Failed to update role");
        }
    },

    updateRoleStatus: async (id: number, command: string) => {
        set({ rolesLoading: true, rolesError: null });
        try {
            await updateRoleStatus(id, command);
            await useRoleStore.getState().fetchRoles();
            set({ rolesLoading: false });
            toast.success(`Role ${command === 'activate' ? 'activated' : 'deactivated'} successfully`);
        } catch (error: any) {
            set({ rolesLoading: false, rolesError: error.message });
            toast.error(error.message || "Failed to update role status");
        }
    },

    setSelectedRole: (role: Role | null) => {
        set({ selectedRole: role });
    }
}));

export default useRoleStore;
