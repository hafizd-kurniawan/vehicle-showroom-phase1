import { apiClient } from './apiClient';

export interface SparePart {
  id: number;
  part_code: string;
  name: string;
  description?: string;
  brand?: string;
  cost_price: number;
  selling_price: number;
  stock_quantity: number;
  min_stock_level: number;
  unit_measure?: string;
  created_at: string;
  updated_at: string;
  is_active: boolean;
}

export interface CreateSparePartRequest {
  name: string;
  description?: string;
  brand?: string;
  cost_price: number;
  selling_price: number;
  stock_quantity: number;
  min_stock_level: number;
  unit_measure?: string;
}

export interface UpdateSparePartRequest {
  name: string;
  description?: string;
  brand?: string;
  cost_price: number;
  selling_price: number;
  min_stock_level: number;
  unit_measure?: string;
  is_active?: boolean;
}

export interface SparePartListResponse {
  spare_parts: SparePart[];
  total: number;
  page: number;
  limit: number;
}

export const sparePartService = {
  async list(page = 1, limit = 10, search = ''): Promise<SparePartListResponse> {
    const response = await apiClient.get('/spare-parts', {
      params: { page, limit, search }
    });
    return response.data.data;
  },

  async getById(id: number): Promise<SparePart> {
    const response = await apiClient.get(`/spare-parts/${id}`);
    return response.data.data;
  },

  async create(data: CreateSparePartRequest): Promise<SparePart> {
    const response = await apiClient.post('/spare-parts', data);
    return response.data.data;
  },

  async update(id: number, data: UpdateSparePartRequest): Promise<SparePart> {
    const response = await apiClient.put(`/spare-parts/${id}`, data);
    return response.data.data;
  },

  async delete(id: number): Promise<void> {
    await apiClient.delete(`/spare-parts/${id}`);
  },
};
