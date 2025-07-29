import { apiClient } from './apiClient';

interface LoginRequest {
  username: string;
  password: string;
}

interface LoginResponse {
  token: string;
  user: {
    id: number;
    username: string;
    email: string;
    full_name: string;
    phone?: string;
    role: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
  };
}

interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  phone?: string;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export const authService = {
  async login(username: string, password: string): Promise<LoginResponse> {
    const response = await apiClient.post('/auth/login', { username, password });
    return response.data.data;
  },

  async logout(): Promise<void> {
    await apiClient.post('/auth/logout');
  },

  async getProfile(): Promise<User> {
    const response = await apiClient.get('/auth/me');
    return response.data.data;
  },
};
