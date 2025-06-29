import getCurrentUser from "@/app/actions/getCurrentUser";
import {User, UserIdentity} from "@/lib/types"
import { create } from "zustand";


interface UserStore {
    user: UserIdentity | null,
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    fetchUser: () => Promise<void>;
    updateProfile:(data: Partial<User>)=>Promise<void>
}

export const useUserStore = create<UserStore>((set) => ({
    user: null,
    loading: false,
    error: null,
    next: null,
    previous: null,
    fetchUser: async () => {
        set({loading: true})
        try{
            const user = await getCurrentUser()
            set({loading:false, user: user})
            
        } catch (error: any) {
            set({ error: error.message, loading: false });
        }
    },
    updateProfile: async(data: Partial<User>)=>{
        
    }
}))

