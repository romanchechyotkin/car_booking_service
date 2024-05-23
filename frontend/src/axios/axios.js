import axios from "axios";

export const STATIC = "http://localhost:9000/images/"

export const axiosInstance = axios.create({
    baseURL: "http://localhost:8000",
    headers: {
        "Content-Type": ['application/x-www-form-urlencoded', "application/json"],
    },
})

axiosInstance.interceptors.request.use(
    async config => {
        config.headers = {
            'Authorization': `Bearer ${JSON.parse(localStorage.getItem("access_token"))}`,
            'Accept': 'application/json',
            'Content-Type': 'application/x-www-form-urlencoded',
        }
        return config;
    },
    error => {
        Promise.reject(error)
});

axiosInstance.interceptors.response.use((response) => {
    return response
}, async function (error) {
    const originalRequest = error.config;
    if (error.response.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true;
        const res = await refreshAccessToken();
        localStorage.setItem('access_token', JSON.stringify(res.data.access_token))
        localStorage.setItem('refresh_token', JSON.stringify(res.data.refresh_token))
        axios.defaults.headers.common['Authorization'] = 'Bearer ' + res.data.access_token;
        return axiosInstance(originalRequest);
    }
    return Promise.reject(error);
});

const refreshAccessToken = async () => {
    return await axiosInstance.post("/auth/refresh", JSON.stringify({
        "refresh_token": JSON.parse(localStorage.getItem("refresh_token"))
    }))
}
