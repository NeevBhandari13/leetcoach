import axios from 'axios';

// FILE TO DEFINE AXIOS INSTANCE
const axiosInstance = axios.create({
    baseURL: process.env.BACKEND_API_BASE_URL,
    timeout: 10000,
});

export default axiosInstance
