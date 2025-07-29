import { apiClient } from './apiClient';
import { PurchaseTransaction, SalesTransaction } from './transactionService';

export interface VehicleProfitability {
	vehicle_id: number;
	vehicle_code: string;
	brand: string;
	model: string;
	year: number;
	purchase_price: number;
	total_repair_cost: number;
	final_selling_price: number;
	profit: number;
	sold_at: string;
}

export const reportService = {
	async getVehicleProfitability(startDate: string, endDate: string): Promise<VehicleProfitability[]> {
		const response = await apiClient.get('/reports/profitability', {
			params: { start_date: startDate, end_date: endDate },
		});
		return response.data.data;
	},

	async getSalesReport(startDate: string, endDate: string): Promise<SalesTransaction[]> {
		const response = await apiClient.get('/reports/sales', {
			params: { start_date: startDate, end_date: endDate },
		});
		return response.data.data;
	},

	async getPurchaseReport(startDate: string, endDate: string): Promise<PurchaseTransaction[]> {
		const response = await apiClient.get('/reports/purchases', {
			params: { start_date: startDate, end_date: endDate },
		});
		return response.data.data;
	},
};
