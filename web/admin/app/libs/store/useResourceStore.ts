import { create } from "zustand";
import { Resource } from "@/app/libs/types";
import { toast } from "sonner";
import { 
    fetchResources, 
    createResource, 
    deleteResource, 
    updateResource, 
    fetchResourceById,
    addResourceToRole,
    removeResourceFromRole,
    addMultipleResourcesToRole,
    removeMultipleResourcesFromRole
} from "@/app/actions/resource";

interface ResourceStore {
    success: string;
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;
    resources: Resource[];
    resourcesLoading: boolean;
    resourcesError: string | null;
    selectedResource: Resource | null;
    roleResources: Resource[];
    
    roleResourcesLoading: boolean;
    roleResourcesError: string | null;
    roleOperationLoading: boolean;
    roleOperationError: string | null;

    fetchResources: () => Promise<void>;
    fetchResourceById: (id: number) => Promise<void>;
    createResource: (data: Partial<Resource>) => Promise<void>;
    deleteResource: (id: number) => Promise<void>;
    editResource: (data: Partial<Resource>) => Promise<void>;
    setSelectedResource: (resource: Resource | null) => void;
    
    addResourceToRole: (roleId: number, resourceId: number) => Promise<void>;
    removeResourceFromRole: (roleId: number, resourceId: number) => Promise<void>;
    addMultipleResourcesToRole: (roleId: number, resourceIds: number[]) => Promise<void>;
    removeMultipleResourcesFromRole: (roleId: number, resourceIds: number[]) => Promise<void>;
    setRoleResources: (resources: Resource[]) => void;
    clearRoleResources: () => void;
}

const useResourceStore = create<ResourceStore>((set, get) => ({
    success: "",
    loading: false,
    error: null,
    next: null,
    previous: null,
    resources: [],
    resourcesLoading: false,
    resourcesError: null,
    selectedResource: null,
    roleResources: [],
    
    roleResourcesLoading: false,
    roleResourcesError: null,
    roleOperationLoading: false,
    roleOperationError: null,

    fetchResources: async () => {
        set({ resourcesLoading: true, resourcesError: null });
        try {
            const response = await fetchResources();
            set({
                resources: response.body?.resources || response.body || [],
                resourcesLoading: false,
            });
        } catch (error: any) {
            set({
                resourcesError: error.message || "Failed to fetch resources",
                resourcesLoading: false,
            });
            toast.error(error.message || "Failed to fetch resources");
        }
    },

    fetchResourceById: async (id: number) => {
        set({ loading: true, error: null });
        try {
            const response = await fetchResourceById(id);
            set({
                selectedResource: response.body || response,
                loading: false,
            });
        } catch (error: any) {
            set({
                error: error.message || "Failed to fetch resource",
                loading: false,
            });
            toast.error(error.message || "Failed to fetch resource");
        }
    },

    createResource: async (data: Partial<Resource>) => {
        set({ resourcesLoading: true, resourcesError: null });
        try {
            await createResource(data);
            await get().fetchResources();
            set({ resourcesLoading: false });
            toast.success("Resource created successfully");
        } catch (error: any) {
            set({ resourcesLoading: false, resourcesError: error.message });
            toast.error(error.message || "Failed to create resource");
        }
    },

    deleteResource: async (id: number) => {
        set({ resourcesLoading: true, resourcesError: null });
        try {
            await deleteResource(id);
            await get().fetchResources();
            set({ resourcesLoading: false });
            toast.success("Resource deleted successfully");
        } catch (error: any) {
            set({ resourcesLoading: false, resourcesError: error.message });
            toast.error(error.message || "Failed to delete resource");
        }
    },

    editResource: async (data: Partial<Resource>) => {
        set({ resourcesLoading: true, resourcesError: null });
        try {
            await updateResource(data);
            await get().fetchResources();
            set({ resourcesLoading: false });
            toast.success("Resource updated successfully");
        } catch (error: any) {
            set({ resourcesLoading: false, resourcesError: error.message });
            toast.error(error.message || "Failed to update resource");
        }
    },

    setSelectedResource: (resource: Resource | null) => {
        set({ selectedResource: resource });
    },

    addResourceToRole: async (roleId: number, resourceId: number) => {
        set({ roleOperationLoading: true, roleOperationError: null });
        try {
            const result = await addResourceToRole(roleId, resourceId);
            
            // Update roleResources state optimistically
            const currentRoleResources = get().roleResources;
            const resourceToAdd = get().resources.find(r => r.id === resourceId);
            
            if (resourceToAdd && !currentRoleResources.some(r => r.id === resourceId)) {
                set({ 
                    roleResources: [...currentRoleResources, resourceToAdd],
                    roleOperationLoading: false 
                });
            } else {
                set({ roleOperationLoading: false });
            }
            
            toast.success(result.message || "Resource added to role successfully");
        } catch (error: any) {
            set({ 
                roleOperationLoading: false, 
                roleOperationError: error.message 
            });
            toast.error(error.message || "Failed to add resource to role");
        }
    },

    removeResourceFromRole: async (roleId: number, resourceId: number) => {
        set({ roleOperationLoading: true, roleOperationError: null });
        try {
            const result = await removeResourceFromRole(roleId, resourceId);

            const currentRoleResources = get().roleResources;
            set({ 
                roleResources: currentRoleResources.filter(r => r.id !== resourceId),
                roleOperationLoading: false 
            });
            
            toast.success(result.message || "Resource removed from role successfully");
        } catch (error: any) {
            set({ 
                roleOperationLoading: false, 
                roleOperationError: error.message 
            });
            toast.error(error.message || "Failed to remove resource from role");
        }
    },

    addMultipleResourcesToRole: async (roleId: number, resourceIds: number[]) => {
        set({ roleOperationLoading: true, roleOperationError: null });
        try {
            await addMultipleResourcesToRole(roleId, resourceIds);
            
            const currentRoleResources = get().roleResources;
            const resourcesToAdd = get().resources.filter(r => 
                resourceIds.includes(r.id) && !currentRoleResources.some(rr => rr.id === r.id)
            );
            
            set({ 
                roleResources: [...currentRoleResources, ...resourcesToAdd],
                roleOperationLoading: false 
            });
            
        } catch (error: any) {
            set({ 
                roleOperationLoading: false, 
                roleOperationError: error.message 
            });
            toast.error(error.message || "Failed to add resources to role");
        }
    },

    removeMultipleResourcesFromRole: async (roleId: number, resourceIds: number[]) => {
        set({ roleOperationLoading: true, roleOperationError: null });
        try {
            const results = await removeMultipleResourcesFromRole(roleId, resourceIds);

            const currentRoleResources = get().roleResources;
            set({ 
                roleResources: currentRoleResources.filter(r => !resourceIds.includes(r.id)),
                roleOperationLoading: false 
            });
            
            toast.success(`Successfully removed ${results.length} resources from role`);
        } catch (error: any) {
            set({ 
                roleOperationLoading: false, 
                roleOperationError: error.message 
            });
            toast.error(error.message || "Failed to remove resources from role");
        }
    },

    setRoleResources: (resources: Resource[]) => {
        set({ roleResources: resources });
    },

    clearRoleResources: () => {
        set({ roleResources: [] });
    }
}));

export default useResourceStore;