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
        }
        return Promise.reject(error);
    }
)

export interface PrepareUploadResponse {
    diplomaHash: string;
    arweaveTxID: string;
    arweaveUrl: string;
}

export interface ConfirmUploadPayload {
    diplomaHash: string;
    arweaveTxID: string;
    polygonTxHash: string;
    blockNumber: number;
    firstName: string;
    lastName: string;
    email: string;
    university: string;
    faculty: string;
    department: string;
    graduationYear: number;
    studentNumber: string;
    nationality: string;
}

export interface ConfirmUploadResponse {
    success: boolean;
    diplomaHash: string;
    arweaveTxID: string;
    arweaveUrl: string;
    polygonTxHash: string;
    blockNumber: number;
}

export const diplomaService = {
    /**
     * Phase 1: Upload PDF to Arweave, get back diploma hash + arweave tx ID.
     * The frontend then uses these to sign the Polygon transaction via MetaMask.
     */
    prepare: async (formData: FormData): Promise<PrepareUploadResponse> => {
        const response = await api.post('/v1/diploma/prepare', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
        });
        return response.data;
    },

    /**
     * Phase 2: After MetaMask signed the Polygon tx, tell the backend to save the record.
     */
    confirm: async (payload: ConfirmUploadPayload): Promise<ConfirmUploadResponse> => {
        const response = await api.post('/v1/diploma/confirm', payload);
        return response.data;
    },

    verify: async (diplomaId: string) => {
        const response = await api.post('/v1/diploma/verify', { DiplomaID: diplomaId });
        return response.data;
    },

    getAllDiplomas: async () => {
        const response = await api.get('/v1/diploma/records');
        return response.data;
    },

    getDiplomaFile: async (diplomaId: string) => {
        try {
            const response = await api.get(`v1/diploma/records/${diplomaId}`, {
                responseType: 'blob',
            });

            const blob = new Blob([response.data], { type: 'application/pdf' });
            const url = window.URL.createObjectURL(blob);
            window.open(url, '_blank');
        } catch (error: any) {
            console.error("Diploma fetch error:", error);

            if (error?.response?.status === 401) {
                localStorage.removeItem('blockcertify_user');
                window.location.href = '/login';
                return;
            }
            alert("Diploma yüklenemedi veya bulunamadı.");
        }
    }
};

export default api;
