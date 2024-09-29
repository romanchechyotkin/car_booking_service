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

axiosInstance.interceptors.response.use(
    (response) => {
        return response; 
    }, 
    async function (error) {
        const originalRequest = error.config;

        if ((error.response.status === 401 || error.response.status === 403) && !originalRequest._retry) {
            originalRequest._retry = true;

            try {
                const res = await refreshAccessToken(); 
                
                if (res.status === 200) {
                    const newAccessToken = res.data.access_token;
                    const newRefreshToken = res.data.refresh_token;

                    localStorage.setItem('access_token', JSON.stringify(newAccessToken));
                    localStorage.setItem('refresh_token', JSON.stringify(newRefreshToken));

                    originalRequest.headers['Authorization'] = 'Bearer ' + newAccessToken;
                    
                    return axiosInstance(originalRequest);
                }
            } catch (err) {
                console.error('Token refresh failed', err);
            }
        }

        return Promise.reject(error); 
    }
);

const refreshAccessToken = async () => {
    const refreshToken = JSON.parse(localStorage.getItem('refresh_token')); 

    if (!refreshToken) {
        throw new Error("No refresh token available");
    }

    try {
        const response = await axiosInstance.post("/auth/refresh", {
            "refresh_token": refreshToken 
        });
        
        return response;
    } catch (error) {
        console.error("Error refreshing access token", error); 
        throw error;
    }
};
