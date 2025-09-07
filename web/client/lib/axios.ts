import { getSession } from "@/app/actions/getSession";
import axios from "axios";
import { NextResponse } from "next/server";

const apiUrl = "https://b2b-67gk.onrender.com/api/v1";

const axiosIns = axios.create({
    baseURL: apiUrl
});

axiosIns.interceptors.request.use(
    async (config) => {
        const session = await getSession()
        // @ts-ignore
        if (session && session?.accessToken) {
            // @ts-ignore
            config.headers.Authorization = `Bearer ${session.accessToken}`;
        }

        return config;
    },
    (error) => {
        Promise.reject(error);
    }
);

export default axiosIns;
