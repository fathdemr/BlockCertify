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

api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error?.response?.status === 401) {
            console.warn("Session expired -> redirecting to login");

            alert("Oturumunuz zaman aşımına uğradı. Tekrar giriş yapın.")

            localStorage.removeItem('blockcertify_user');

            window.location.href = '/login';

            return Promise.reject(error);
        }
    }
)

export const diplomaService = {
    upload: async (formData: FormData) => {
        try {
            const response = await api.post('/v1/diploma/upload', formData, {
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
            const response = await api.post('/v1/diploma/verify', { DiplomaID: diplomaId });
            return response.data;
        } catch (error) {
            console.error('Error verifying diploma:', error);
            throw error;
        }
    },

    getAllDiplomas: async () => {
        try{
            const response = await api.get('/v1/diploma/records');
            return response.data;
        }catch (error) {
            console.error('Error retrieving diploma:', error);
            throw error;
        }
    },

    getDiplomaFile: async (diplomaId: string) => {
        try{
            const response = await api.get(`v1/diploma/records/${diplomaId}`, {
                responseType: 'blob',
            });

            const blob = new Blob([response.data], {type: 'application/pdf'});
            const url = window.URL.createObjectURL(blob);
            window.open(url, '_blank');
        }catch (error: any) {
            console.error("Diploma fetch error:", error);

            if(error?.response?.status === 401) {
                localStorage.removeItem('blockcertify_user');
                window.location.href = '/login';
                return;
            }
            alert("Diploma yüklenemedi veya bulunamadı.");
        }
    }
};

export default api;
