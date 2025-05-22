import { getSession } from "@/actions/getSession";
import axios from "axios";

const apiUrl = "http://localhost:3000/api/v1";

const axiosIns = axios.create({
  baseURL: apiUrl,
  headers: {
    "Content-Type": "application/json",
  },
});

let isRefreshing = false;
axiosIns.interceptors.request.use(
  async (config) => {
    console.log("Config called");
    const session = await getSession()
    console.log("*****************************")
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

axiosIns.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config;
    let retryLimit = 2;
    console.log(error, "________________________________________________");
    if (error.response.status === 401 && retryLimit > 0) {
      console.log(retryLimit);
      retryLimit -= 1;

      try {
        const newToken = await refreshAccessToken();
        if (newToken) {
          axios.defaults.headers.common[
            "Authorization"
          ] = `Bearer ${newToken.access}`;
          originalRequest.headers[
            "Authorization"
          ] = `Bearer ${newToken.access}`;
          return axiosIns(originalRequest);
        }
      } catch (err) {
        return Promise.reject(err);
      }
    }

    return Promise.reject(error);
  }
);

async function refreshAccessToken() {
  try {
    isRefreshing = true;
    console.log("Refreshing");

    try {
      const sessionData = await getSession();
      console.log("Session Data_________________________")
      console.log(sessionData)

      const tokens = (
        await axiosIns.post(`/auth/refresh`, {
          // @ts-ignore
          refresh: sessionData.refreshToken,
        })
      ).data;

      isRefreshing = false;
      return tokens;
    } catch (refreshError) {
      isRefreshing = false;
    }
  } catch (e) {
    console.log(e, "Error");
    window.location.href = "/login";
  }
}
export default axiosIns;
