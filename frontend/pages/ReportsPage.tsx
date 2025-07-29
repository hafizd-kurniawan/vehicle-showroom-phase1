import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { ArrowLeft, BarChart3, Calendar as CalendarIcon } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { DateRange } from 'react-day-picker';
import { addDays, format } from 'date-fns';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { Calendar } from '@/components/ui/calendar';
import { cn } from '@/lib/utils';
import { reportService, VehicleProfitability } from '../services/reportService';
import { useToast } from '@/components/ui/use-toast';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { SalesTransaction, PurchaseTransaction } from '../services/transactionService';

const formatPrice = (price: number) => {
	return new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR',
		minimumFractionDigits: 0,
	}).format(price);
};

const formatDate = (dateString: string) => {
	return new Date(dateString).toLocaleDateString('id-ID', {
		year: 'numeric',
		month: 'long',
		day: 'numeric',
	});
};

function ProfitabilityReportTab() {
	const [date, setDate] = useState<DateRange | undefined>({
		from: addDays(new Date(), -30),
		to: new Date(),
	});
	const [data, setData] = useState<VehicleProfitability[]>([]);
	const [loading, setLoading] = useState(false);
	const { toast } = useToast();

	const handleGenerate = async () => {
		if (!date?.from || !date?.to) {
			toast({ title: 'Error', description: 'Please select a date range.', variant: 'destructive' });
			return;
		}
		setLoading(true);
		try {
			const result = await reportService.getVehicleProfitability(
				format(date.from, 'yyyy-MM-dd'),
				format(date.to, 'yyyy-MM-dd')
			);
			setData(result);
		} catch (error) {
			toast({ title: 'Error', description: 'Failed to generate report.', variant: 'destructive' });
		} finally {
			setLoading(false);
		}
	};

	return (
		<Card>
			<CardHeader>
				<CardTitle>Vehicle Profitability Report</CardTitle>
				<CardDescription>Analyze profit for each vehicle sold within a date range.</CardDescription>
			</CardHeader>
			<CardContent className="space-y-4">
				<div className="flex items-center gap-4">
					<DateRangePicker date={date} setDate={setDate} />
					<Button onClick={handleGenerate} disabled={loading}>
						{loading ? 'Generating...' : 'Generate Report'}
					</Button>
				</div>
				<div className="border rounded-md">
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead>Vehicle</TableHead>
								<TableHead className="text-right">Purchase Price</TableHead>
								<TableHead className="text-right">Repair Cost</TableHead>
								<TableHead className="text-right">Selling Price</TableHead>
								<TableHead className="text-right font-bold">Profit</TableHead>
								<TableHead>Sold Date</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							{data.length > 0 ? (
								data.map((item) => (
									<TableRow key={item.vehicle_id}>
										<TableCell>
											<div className="font-medium">{`${item.brand} ${item.model}`}</div>
											<div className="text-sm text-muted-foreground">{item.vehicle_code}</div>
										</TableCell>
										<TableCell className="text-right">{formatPrice(item.purchase_price)}</TableCell>
										<TableCell className="text-right">{formatPrice(item.total_repair_cost)}</TableCell>
										<TableCell className="text-right">{formatPrice(item.final_selling_price)}</TableCell>
										<TableCell className={`text-right font-bold ${item.profit >= 0 ? 'text-green-600' : 'text-red-600'}`}>
											{formatPrice(item.profit)}
										</TableCell>
										<TableCell>{formatDate(item.sold_at)}</TableCell>
									</TableRow>
								))
							) : (
								<TableRow>
									<TableCell colSpan={6} className="h-24 text-center">
										No data to display. Generate a report to see results.
									</TableCell>
								</TableRow>
							)}
						</TableBody>
					</Table>
				</div>
			</CardContent>
		</Card>
	);
}

// Similar components for SalesReportTab and PurchaseReportTab can be created
// For brevity, I'll create a generic report tab component
interface ReportTabProps<T> {
	title: string;
	description: string;
	fetchData: (startDate: string, endDate: string) => Promise<T[]>;
	columns: { header: string; accessor: (item: T) => React.ReactNode; className?: string }[];
}

function GenericReportTab<T>({ title, description, fetchData, columns }: ReportTabProps<T>) {
	const [date, setDate] = useState<DateRange | undefined>({
		from: addDays(new Date(), -30),
		to: new Date(),
	});
	const [data, setData] = useState<T[]>([]);
	const [loading, setLoading] = useState(false);
	const { toast } = useToast();

	const handleGenerate = async () => {
		if (!date?.from || !date?.to) {
			toast({ title: 'Error', description: 'Please select a date range.', variant: 'destructive' });
			return;
		}
		setLoading(true);
		try {
			const result = await fetchData(format(date.from, 'yyyy-MM-dd'), format(date.to, 'yyyy-MM-dd'));
			setData(result);
		} catch (error) {
			toast({ title: 'Error', description: 'Failed to generate report.', variant: 'destructive' });
		} finally {
			setLoading(false);
		}
	};

	return (
		<Card>
			<CardHeader>
				<CardTitle>{title}</CardTitle>
				<CardDescription>{description}</CardDescription>
			</CardHeader>
			<CardContent className="space-y-4">
				<div className="flex items-center gap-4">
					<DateRangePicker date={date} setDate={setDate} />
					<Button onClick={handleGenerate} disabled={loading}>
						{loading ? 'Generating...' : 'Generate Report'}
					</Button>
				</div>
				<div className="border rounded-md">
					<Table>
						<TableHeader>
							<TableRow>
								{columns.map((col, i) => (
									<TableHead key={i} className={col.className}>{col.header}</TableHead>
								))}
							</TableRow>
						</TableHeader>
						<TableBody>
							{data.length > 0 ? (
								data.map((item, rowIndex) => (
									<TableRow key={rowIndex}>
										{columns.map((col, colIndex) => (
											<TableCell key={colIndex} className={col.className}>{col.accessor(item)}</TableCell>
										))}
									</TableRow>
								))
							) : (
								<TableRow>
									<TableCell colSpan={columns.length} className="h-24 text-center">
										No data to display. Generate a report to see results.
									</TableCell>
								</TableRow>
							)}
						</TableBody>
					</Table>
				</div>
			</CardContent>
		</Card>
	);
}

export default function ReportsPage() {
	const navigate = useNavigate();

	const salesColumns = [
		{ header: 'Invoice #', accessor: (item: SalesTransaction) => item.invoice_number },
		{ header: 'Vehicle', accessor: (item: SalesTransaction) => `${item.vehicle?.brand} ${item.vehicle?.model}` },
		{ header: 'Customer', accessor: (item: SalesTransaction) => item.customer?.name },
		{ header: 'Date', accessor: (item: SalesTransaction) => formatDate(item.transaction_date) },
		{ header: 'Total Amount', accessor: (item: SalesTransaction) => formatPrice(item.total_amount), className: 'text-right' },
	];

	const purchaseColumns = [
		{ header: 'Invoice #', accessor: (item: PurchaseTransaction) => item.invoice_number },
		{ header: 'Vehicle', accessor: (item: PurchaseTransaction) => `${item.vehicle?.brand} ${item.vehicle?.model}` },
		{ header: 'Customer', accessor: (item: PurchaseTransaction) => item.customer?.name },
		{ header: 'Date', accessor: (item: PurchaseTransaction) => formatDate(item.transaction_date) },
		{ header: 'Total Amount', accessor: (item: PurchaseTransaction) => formatPrice(item.total_amount), className: 'text-right' },
	];

	return (
		<div className="min-h-screen bg-gray-50">
			<header className="bg-white shadow-sm border-b">
				<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
					<div className="flex items-center h-16">
						<Button variant="ghost" onClick={() => navigate('/dashboard')} className="mr-4">
							<ArrowLeft className="h-4 w-4" />
						</Button>
						<h1 className="text-xl font-semibold text-gray-900">Reports & Analytics</h1>
					</div>
				</div>
			</header>

			<main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
				<Tabs defaultValue="profitability" className="space-y-6">
					<TabsList className="grid w-full grid-cols-3">
						<TabsTrigger value="profitability">Vehicle Profitability</TabsTrigger>
						<TabsTrigger value="sales">Sales Report</TabsTrigger>
						<TabsTrigger value="purchases">Purchase Report</TabsTrigger>
					</TabsList>

					<TabsContent value="profitability">
						<ProfitabilityReportTab />
					</TabsContent>
					<TabsContent value="sales">
						<GenericReportTab
							title="Sales Report"
							description="View all sales transactions within a specific date range."
							fetchData={reportService.getSalesReport}
							columns={salesColumns}
						/>
					</TabsContent>
					<TabsContent value="purchases">
						<GenericReportTab
							title="Purchase Report"
							description="View all purchase transactions within a specific date range."
							fetchData={reportService.getPurchaseReport}
							columns={purchaseColumns}
						/>
					</TabsContent>
				</Tabs>
			</main>
		</div>
	);
}

function DateRangePicker({ date, setDate }: { date?: DateRange; setDate: (date?: DateRange) => void }) {
	return (
		<Popover>
			<PopoverTrigger asChild>
				<Button
					id="date"
					variant={'outline'}
					className={cn('w-[300px] justify-start text-left font-normal', !date && 'text-muted-foreground')}
				>
					<CalendarIcon className="mr-2 h-4 w-4" />
					{date?.from ? (
						date.to ? (
							<>
								{format(date.from, 'LLL dd, y')} - {format(date.to, 'LLL dd, y')}
							</>
						) : (
							format(date.from, 'LLL dd, y')
						)
					) : (
						<span>Pick a date</span>
					)}
				</Button>
			</PopoverTrigger>
			<PopoverContent className="w-auto p-0" align="start">
				<Calendar
					initialFocus
					mode="range"
					defaultMonth={date?.from}
					selected={date}
					onSelect={setDate}
					numberOfMonths={2}
				/>
			</PopoverContent>
		</Popover>
	);
}
