import { create } from 'zustand';
import { Catalogue } from '../types';

interface CatalogueStore {
    success: string | null;
    catalogue: Catalogue | null;
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    setProduct: (product: Catalogue) => void;
}

const useCatalogueStore = create<CatalogueStore>((set) => ({
    loading: false,
    success: null,
    error: null,
    next: null,
    previous: null,
    catalogue: null,
    setProduct: (product: Catalogue) => {
        set({ catalogue: product })
    },
}));

export default useCatalogueStore;