import axios from "axios";

// const apiUrl = "http://localhost:8084/api/v1";

const axiosIns = axios.create({
    // baseURL: apiUrl,
    headers: {
        "Content-Type": "application/json",
    },
});

let isRefreshing = false;
let refreshSubscribers: ((token: string) => void)[] = [];

function subscribeTokenRefresh(cb: (token: string) => void) {
    refreshSubscribers.push(cb);
}

function onRefreshed(token: string) {
    refreshSubscribers.forEach(cb => cb(token));
    refreshSubscribers = [];
}

axiosIns.interceptors.request.use(
    (config) => {
        const accessToken = localStorage.getItem('accessToken');
        if (accessToken) {
            config.headers.Authorization = `Bearer ${accessToken}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

axiosIns.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;

            if (isRefreshing) {
                return new Promise((resolve) => {
                    subscribeTokenRefresh((newToken) => {
                        originalRequest.headers.Authorization = `Bearer ${newToken}`;
                        resolve(axiosIns(originalRequest));
                    });
                });
            }

            isRefreshing = true;

            try {
                const newTokens = await refreshAccessToken();
                localStorage.setItem('accessToken', newTokens.access);

                if (newTokens.refresh) {
                    localStorage.setItem('refreshToken', newTokens.refresh);
                }

                axiosIns.defaults.headers.common['Authorization'] = `Bearer ${newTokens.access}`;
                onRefreshed(newTokens.access);

                return axiosIns(originalRequest);
            } catch (refreshError) {
                localStorage.removeItem('accessToken');
                localStorage.removeItem('refreshToken');
                window.location.href = '/login';
                return Promise.reject(refreshError);
            } finally {
                isRefreshing = false;
            }
        }

        return Promise.reject(error);
    }
);

async function refreshAccessToken() {
    const refreshToken = localStorage.getItem('refreshToken');
    if (!refreshToken) {
        throw new Error('No refresh token available');
    }

    try {
        const response = await axiosIns.post("/api/auth/refresh/", {
            refresh: refreshToken,
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}

export default axiosIns;