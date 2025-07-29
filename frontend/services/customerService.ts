import { apiClient } from './apiClient';

export interface Customer {
  id: number;
  customer_code: string;
  name: string;
  phone?: string;
  email?: string;
  address?: string;
  id_card_number?: string;
  type: 'individual' | 'corporate';
  created_at: string;
  updated_at: string;
  created_by?: number;
  is_active: boolean;
}

export interface CreateCustomerRequest {
  name: string;
  phone?: string;
  email?: string;
  address?: string;
  id_card_number?: string;
  type: 'individual' | 'corporate';
}

export interface UpdateCustomerRequest {
  name: string;
  phone?: string;
  email?: string;
  address?: string;
  id_card_number?: string;
  type: 'individual' | 'corporate';
}

export interface CustomerListResponse {
  customers: Customer[];
  total: number;
  page: number;
  limit: number;
}

export const customerService = {
  async list(page = 1, limit = 10, search = ''): Promise<CustomerListResponse> {
    const response = await apiClient.get('/customers', {
      params: { page, limit, search }
    });
    return response.data.data;
  },

  async getById(id: number): Promise<Customer> {
    const response = await apiClient.get(`/customers/${id}`);
    return response.data.data;
  },

  async create(data: CreateCustomerRequest): Promise<Customer> {
    const response = await apiClient.post('/customers', data);
    return response.data.data;
  },

  async update(id: number, data: UpdateCustomerRequest): Promise<Customer> {
    const response = await apiClient.put(`/customers/${id}`, data);
    return response.data.data;
  },

  async delete(id: number): Promise<void> {
    await apiClient.delete(`/customers/${id}`);
  },
};
