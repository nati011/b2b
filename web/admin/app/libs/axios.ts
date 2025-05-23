import { getSession } from "@/actions/getSession";
import axios from "axios";
import { signIn } from "next-auth/react";

const apiUrl = "http://localhost:3000/api/v1";

const axiosIns = axios.create({
  baseURL: apiUrl
});

axiosIns.interceptors.request.use(
  async (config) => {
    const session = await getSession()
    // @ts-ignore
    if (session?.error === "RefreshAccessTokenError") {
      window.location.href = "/auth/signin"
    }

    // @ts-ignore
    if (session && session?.accessToken) {
      // @ts-ignore
      config.headers.Authorization = `Bearer ${session.accessToken}`;
    }

    // console.log(config.baseURL, 'Config Here')
    return config;
  },
  (error) => {
    console.log(error, "Errorin interceptor");
    Promise.reject(error);
  }
);

export default axiosIns;
