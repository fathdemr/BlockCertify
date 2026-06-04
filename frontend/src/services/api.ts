import axios from 'axios';
import type { Diploma, University, VerifyResult } from '../types';

// Public routes: /exapi/...
const exapi = axios.create({
  baseURL: '/exapi',
  withCredentials: true,
});

// Private routes: /api/...
const api = axios.create({
  baseURL: '/api',
  withCredentials: true,
});

export const authApi = {
  login: (email: string, password: string) =>
    exapi.post('/user/login', { email, password }),

  register: (name: string, email: string, password: string) =>
    exapi.post('/user/register/admin', { name, email, password }),

  logout: () =>
    exapi.post('/user/logout'),
};

export const diplomaApi = {
  upload: (formData: FormData) =>
    api.post<Diploma>('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),

  verify: (txHash: string) =>
    api.post<VerifyResult>('/verify', { tx_hash: txHash }),

  getAll: () =>
    api.get<Diploma[]>('/records'),

  getById: (id: number) =>
    api.get<Diploma>(`/records/${id}`),
};

export const universityApi = {
  getAll: () =>
    exapi.get<University[]>('/universities'),
};

export const facultyApi = {
  getByUniversity: (universityId: string) =>
    exapi.get('/faculties', { params: { university_id: universityId } }),
};

export const departmentApi = {
  getByFaculty: (facultyId: string) =>
    exapi.get('/departments', { params: { faculty_id: facultyId } }),
};

export const walletApi = {
  getStatus: () =>
    api.get('/wallet/status'),
};

export default api;
