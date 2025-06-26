import { create } from "zustand";
import axiosIns from "@/app/libs/axios";
import { Category } from "@/app/libs/types";
import { Create, Delete, GetAll, Update } from "@/app/actions/category";

interface CategoryStore {
    success: string;
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;
    categories: Category[];
    categoriesLoading: boolean;
    categoriesError: string | null;

    fetchCategories: () => Promise<void>;
    createCategory: (name: string) => Promise<void>;
    deleteCategory: (id: number) => Promise<void>;
    editCategory: (id: number, name: string) => Promise<void>;
}

const useCategoryStore = create<CategoryStore>((set) => ({
    success: "",
    loading: false,
    error: null,
    next: null,
    previous: null,
    categories: [],
    categoriesLoading: false,
    categoriesError: null,


    fetchCategories: async () => {
        set({ categoriesLoading: true, categoriesError: null });
        try {
            const response = await GetAll()
            set({
                categories: response,
                categoriesLoading: false,
            });
        } catch (error) {
            set({
                categoriesError: "Failed to fetch categories",
                categoriesLoading: false,
            });
        }
    },
    createCategory: async (name: string) => {
        set({ categoriesLoading: true, categoriesError: null });
        try {
            const response = await Create(name)
            await useCategoryStore.getState().fetchCategories();
            set({ categoriesLoading: false, success: response });
        } catch (error: any) {
            set({ categoriesLoading: false, categoriesError: error.message });
        }
    },
    deleteCategory: async (id: number) => {
        set({ categoriesLoading: true, categoriesError: null });
        try {
            await Delete(id)
            set({ categoriesLoading: false });
            await useCategoryStore.getState().fetchCategories();
        } catch (error: any) {
            set({ categoriesLoading: false, categoriesError: error.message });
        } 0
    },
    editCategory: async (id: number, name: string) => {
        set({ categoriesLoading: true, categoriesError: null });
        try {
            await Update(id, name)
            set({ categoriesLoading: false });
            await useCategoryStore.getState().fetchCategories();
        } catch (error: any) {
            set({ categoriesLoading: false, categoriesError: error.message });
        }
    }
}));

export default useCategoryStore;
