import axios from 'axios';


const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
    headers: {
        'Content-Type': 'application/json',
    },
});

api.interceptors.request.use((config) => {
    const savedUser = localStorage.getItem('blockcertify_user');
    if (savedUser) {
        const { token } = JSON.parse(savedUser);
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
    }
    return config;
});

export const diplomaService = {
    upload: async (formData: FormData) => {
        try {
            const response = await api.post('/upload-diploma', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });
            return response.data;
        } catch (error) {
            console.error('Error uploading diploma:', error);
            throw error;
        }
    },

    verify: async (diplomaId: string) => {
        try {
            const response = await api.post('/verify-diploma', { DiplomaID: diplomaId });
            return response.data;
        } catch (error) {
            console.error('Error verifying diploma:', error);
            throw error;
        }
    },

    getDiplomaFile: (diplomaId: string) => {
        const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';
        return `${baseURL}/diploma/${diplomaId}`;
    },
};

export default api;
