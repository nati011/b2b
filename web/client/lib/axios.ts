import { getSession } from "@/app/actions/getSession";
import axios from "axios";
import { NextResponse } from "next/server";

// For server-side requests, use internal Docker network URL
// For client-side requests, use external URL
const apiUrl = process.env.NEXT_PUBLIC_BASE_URL || 
               (typeof window === 'undefined' ? 'http://backend:8080' : 'http://localhost:8082');

const axiosIns = axios.create({
    baseURL: apiUrl,
    timeout: 5000, // Reduced to 5 seconds for faster failure
    headers: {
        'Content-Type': 'application/json',
    },
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
