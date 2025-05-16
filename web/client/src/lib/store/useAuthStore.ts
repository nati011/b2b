import { create } from "zustand";
import axiosIns from "@/lib/axios";
import { Retailer, AuthModel, Register, User } from "@/lib/types";
import axios from "axios";

interface ProductsStore {
    success: string | null;
    retailer: Retailer;
    loading: boolean;
    error: string | null;
    isLoggedIn: boolean

    fetchLoggedInUser: () => Promise<void>;
    login: (user: AuthModel) => Promise<void>;
    signup: (user: Register) => Promise<void>;
}

const useAuthStore = create<ProductsStore>((set) => ({
    success: null,
    isLoggedIn: true,
    retailer: {
        id: 0,
        name: "",
        tin: "",
        latitude: "",
        longitude: "",
        general_zone: "",
        region: "",
        woreda: "",
        user: {
            id: 0,
            first_name: "",
            last_name: "",
            email: "",
            phone: "",
            username: "",
            dob: "",
            external_id: ""
        }
    },
    loading: false,
    error: null,


    fetchLoggedInUser: async () => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get("/api/v1/retailer");
            const email = localStorage.getItem("email")
            const retailer = response.data.retailers.list.map((r: Retailer) => {
                if (r.user.email == email) {
                    console.log(email)
                    console.log(r.user.email)
                    set({
                        retailer: r,
                        loading: false,
                    });
                }
                return r
            })
            console.log(retailer)

        } catch (error) {
            set({ error: "Failed to fetch configurableProducts", loading: false });
        }
    },
    login: async (user: AuthModel) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.post("/api/v1/auth/login", user);
            localStorage.setItem('accessToken', response.data.body.access_token);
            localStorage.setItem('refreshToken', response.data.body.refresh_token);
            if (response.status == 202) {
                localStorage.setItem('email', user.email);
            }
            set({
                success: "Login successful",
                loading: false,
            });
            window.location.href = '/product'
            await useAuthStore.getState().fetchLoggedInUser();
        } catch (error) {
            set({ error: "Failed to fetch configurableProducts", loading: false });
        }
    },
    signup: async (user: Register) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.post("/api/v1/retailer", user);
            set({
                success: "Sign up successful",
                loading: false,
            });
            window.location.href = '/login'
        } catch (error) {
            set({ error: "Failed to fetch configurableProducts", loading: false });
        }
    },
    // updateUserData: async(user:User) =>{
    //     set({ loading: true, error: null });
    //     try {
    //         const response = await axiosIns.put(`/api/v1/retailer/${re}`, user);
    //         localStorage.setItem('accessToken', response.data.body.access_token);
    //         localStorage.setItem('refreshToken', response.data.body.refresh_token);
    //         if (response.status == 202) {
    //             localStorage.setItem('email', user.email);
    //         }
    //         set({
    //             success: "Login successful",
    //             loading: false,
    //         });
    //         window.location.href = '/product'
    //         await useAuthStore.getState().fetchLoggedInUser();
    //     } catch (error) {
    //         set({ error: "Failed to fetch configurableProducts", loading: false });
    //     }
    // }

}));

export default useAuthStore;
