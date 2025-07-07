import getCurrentUser from "@/app/actions/getCurrentUser";
import {UpdateProfile} from "@/app/actions/auth"
import { UserIdentity} from "@/lib/types"
import { create } from "zustand";


interface UserStore {
    success: string | null
    user: UserIdentity | null,
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    fetchUser: () => Promise<void>;
    updateProfile:(data: Partial<UserIdentity>)=>Promise<void>
}

export const useUserStore = create<UserStore>((set) => ({
    success: null,
    user: null,
    loading: false,
    error: null,
    next: null,
    previous: null,
    fetchUser: async () => {
        set({loading: true})
        try{
            const user = await getCurrentUser()
            console.log(user)
            set({loading:false, user: user})
            
        } catch (error: any) {
            set({ error: error.message, loading: false });
        }
    },
    updateProfile: async(data: Partial<UserIdentity>)=>{
        set({loading: true})
        try{
            await UpdateProfile(data)
            await useUserStore.getState().fetchUser()
            set({loading:false, success: "Profile updating successfully."})
        } catch (error: any) {
            set({ error: error.message, loading: false });
        }
    }
}))

