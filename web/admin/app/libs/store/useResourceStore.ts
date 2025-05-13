import { create } from "zustand";
import axiosIns from "@/app/libs/axios";
import { Resource } from "@/app/libs/types";

interface ResourceStore {
    success: string;
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;
    resource: Resource[];
    resourceLoading: boolean;
    resourceError: string | null;

    fetchResources: () => Promise<void>;
    createResource: (data: Partial<Resource>) => Promise<void>;
    deleteResource: (id: number) => Promise<void>;
    editResource: (data: Partial<Resource>) => Promise<void>;
}

const useResourceStore = create<ResourceStore>((set) => ({
    success: "",
    loading: false,
    error: null,
    next: null,
    previous: null,
    resource: [],
    resourceLoading: false,
    resourceError: null,


    fetchResources: async () => {
        set({ resourceLoading: true, resourceError: null });
        try {
            const response = await axiosIns.get("/api/role");
            set({
                resource: response.data.body.resource,
                resourceLoading: false,
            });
        } catch (error) {
            set({
                resourceError: "Failed to fetch resource",
                resourceLoading: false,
            });
        }
    },
    createResource: async (data: Partial<Resource>) => {
        set({ resourceLoading: true, resourceError: null });
        try {
            const response = await axiosIns.post(
                "/api/role",
                data,
            );
            await useResourceStore.getState().fetchResources();
            set({ resourceLoading: false });
        } catch (error: any) {
            set({ resourceLoading: false, resourceError: error.message });
        }
    },
    deleteResource: async (id: number) => {
        set({ resourceLoading: true, resourceError: null });
        try {
            await axiosIns.delete(`/api/role?id=${id}`);
            set({ resourceLoading: false });
            await useResourceStore.getState().fetchResources();
        } catch (error: any) {
            set({ resourceLoading: false, resourceError: error.message });
        } 0
    },
    editResource: async (data: Partial<Resource>) => {
        set({ resourceLoading: true, resourceError: null });
        try {
            console.log(data)
            await axiosIns.put(`/api/role`, data);
            set({ resourceLoading: false });
            await useResourceStore.getState().fetchResources();
        } catch (error: any) {
            set({ resourceLoading: false, resourceError: error.message });
        }
    }
}));

export default useResourceStore;
