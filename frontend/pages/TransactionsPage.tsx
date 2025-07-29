import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useToast } from '@/components/ui/use-toast';
import { Search, Plus, ArrowLeft, TrendingUp, TrendingDown, Eye } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { transactionService, PurchaseTransaction, SalesTransaction } from '../services/transactionService';
import { useAuth } from '../contexts/AuthContext';

export default function TransactionsPage() {
  const [purchases, setPurchases] = useState<PurchaseTransaction[]>([]);
  const [sales, setSales] = useState<SalesTransaction[]>([]);
  const [purchaseLoading, setPurchaseLoading] = useState(true);
  const [salesLoading, setSalesLoading] = useState(true);
  const [purchaseSearch, setPurchaseSearch] = useState('');
  const [salesSearch, setSalesSearch] = useState('');
  const [purchasePage, setPurchasePage] = useState(1);
  const [salesPage, setSalesPage] = useState(1);
  const [purchaseTotal, setPurchaseTotal] = useState(0);
  const [salesTotal, setSalesTotal] = useState(0);
  const navigate = useNavigate();
  const { toast } = useToast();
  const { user } = useAuth();

  const canCreateTransaction = user?.role === 'admin' || user?.role === 'cashier';
  const limit = 10;

  useEffect(() => {
    loadPurchases();
  }, [purchasePage, purchaseSearch]);

  useEffect(() => {
    loadSales();
  }, [salesPage, salesSearch]);

  const loadPurchases = async () => {
    try {
      setPurchaseLoading(true);
      const response = await transactionService.listPurchases(purchasePage, limit, purchaseSearch);
      setPurchases(response.transactions as PurchaseTransaction[]);
      setPurchaseTotal(response.total);
    } catch (error) {
      console.error('Failed to load purchases:', error);
      toast({
        title: "Error",
        description: "Failed to load purchase transactions",
        variant: "destructive",
      });
    } finally {
      setPurchaseLoading(false);
    }
  };

  const loadSales = async () => {
    try {
      setSalesLoading(true);
      const response = await transactionService.listSales(salesPage, limit, salesSearch);
      setSales(response.transactions as SalesTransaction[]);
      setSalesTotal(response.total);
    } catch (error) {
      console.error('Failed to load sales:', error);
      toast({
        title: "Error",
        description: "Failed to load sales transactions",
        variant: "destructive",
      });
    } finally {
      setSalesLoading(false);
    }
  };

  const handlePurchaseSearch = (value: string) => {
    setPurchaseSearch(value);
    setPurchasePage(1);
  };

  const handleSalesSearch = (value: string) => {
    setSalesSearch(value);
    setSalesPage(1);
  };

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
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const purchaseTotalPages = Math.ceil(purchaseTotal / limit);
  const salesTotalPages = Math.ceil(salesTotal / limit);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <Button
                variant="ghost"
                onClick={() => navigate('/dashboard')}
                className="mr-4"
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <h1 className="text-xl font-semibold text-gray-900">
                Transaction Management
              </h1>
            </div>
            {canCreateTransaction && (
              <div className="flex gap-2">
                <Button onClick={() => navigate('/transactions/purchase/new')} variant="outline">
                  <TrendingDown className="h-4 w-4 mr-2" />
                  New Purchase
                </Button>
                <Button onClick={() => navigate('/transactions/sales/new')}>
                  <TrendingUp className="h-4 w-4 mr-2" />
                  New Sale
                </Button>
              </div>
            )}
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Tabs defaultValue="purchases" className="space-y-6">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="purchases">Purchase Transactions</TabsTrigger>
            <TabsTrigger value="sales">Sales Transactions</TabsTrigger>
          </TabsList>

          {/* Purchase Transactions Tab */}
          <TabsContent value="purchases" className="space-y-6">
            {/* Search */}
            <Card>
              <CardHeader>
                <CardTitle>Search Purchase Transactions</CardTitle>
                <CardDescription>Find transactions by number, invoice, customer, or vehicle</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                  <Input
                    placeholder="Search purchases..."
                    value={purchaseSearch}
                    onChange={(e) => handlePurchaseSearch(e.target.value)}
                    className="pl-10"
                  />
                </div>
              </CardContent>
            </Card>

            {/* Purchase List */}
            <div className="space-y-4">
              {purchaseLoading ? (
                <div className="text-center py-8">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                  <p className="mt-2 text-gray-500">Loading purchases...</p>
                </div>
              ) : purchases.length === 0 ? (
                <Card>
                  <CardContent className="text-center py-8">
                    <p className="text-gray-500">No purchase transactions found</p>
                    {canCreateTransaction && (
                      <Button 
                        onClick={() => navigate('/transactions/purchase/new')}
                        className="mt-4"
                      >
                        <Plus className="h-4 w-4 mr-2" />
                        Create First Purchase
                      </Button>
                    )}
                  </CardContent>
                </Card>
              ) : (
                purchases.map((purchase) => (
                  <Card key={purchase.id} className="hover:shadow-md transition-shadow">
                    <CardContent className="p-6">
                      <div className="flex justify-between items-start">
                        <div className="flex-1">
                          <div className="flex items-center gap-3 mb-3">
                            <TrendingDown className="h-5 w-5 text-red-600" />
                            <h3 className="text-lg font-semibold">{purchase.transaction_number}</h3>
                            <Badge variant="outline">{purchase.invoice_number}</Badge>
                            <Badge className="bg-red-100 text-red-800">{purchase.status}</Badge>
                          </div>
                          
                          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-sm text-gray-600 mb-4">
                            <div>
                              <span className="text-gray-500">Customer:</span>
                              <p className="font-medium">{purchase.customer?.name}</p>
                            </div>
                            <div>
                              <span className="text-gray-500">Vehicle:</span>
                              <p className="font-medium">
                                {purchase.vehicle?.brand} {purchase.vehicle?.model}
                              </p>
                            </div>
                            <div>
                              <span className="text-gray-500">Payment:</span>
                              <p className="font-medium capitalize">{purchase.payment_method}</p>
                            </div>
                            <div>
                              <span className="text-gray-500">Date:</span>
                              <p className="font-medium">{formatDate(purchase.transaction_date)}</p>
                            </div>
                          </div>

                          <div className="flex items-center justify-between">
                            <div>
                              <span className="text-gray-500">Total Amount:</span>
                              <p className="text-lg font-bold text-red-600">
                                {formatPrice(purchase.total_amount)}
                              </p>
                            </div>
                            <div className="text-right text-sm text-gray-500">
                              <p>Cashier: {purchase.cashier?.full_name}</p>
                            </div>
                          </div>
                        </div>
                        
                        <div className="ml-4">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => navigate(`/transactions/purchase/${purchase.id}`)}
                          >
                            <Eye className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))
              )}
            </div>

            {/* Purchase Pagination */}
            {purchaseTotalPages > 1 && (
              <div className="flex justify-center gap-2 mt-8">
                <Button
                  variant="outline"
                  onClick={() => setPurchasePage(purchasePage - 1)}
                  disabled={purchasePage === 1}
                >
                  Previous
                </Button>
                <span className="flex items-center px-4 text-sm text-gray-600">
                  Page {purchasePage} of {purchaseTotalPages}
                </span>
                <Button
                  variant="outline"
                  onClick={() => setPurchasePage(purchasePage + 1)}
                  disabled={purchasePage === purchaseTotalPages}
                >
                  Next
                </Button>
              </div>
            )}
          </TabsContent>

          {/* Sales Transactions Tab */}
          <TabsContent value="sales" className="space-y-6">
            {/* Search */}
            <Card>
              <CardHeader>
                <CardTitle>Search Sales Transactions</CardTitle>
                <CardDescription>Find transactions by number, invoice, customer, or vehicle</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                  <Input
                    placeholder="Search sales..."
                    value={salesSearch}
                    onChange={(e) => handleSalesSearch(e.target.value)}
                    className="pl-10"
                  />
                </div>
              </CardContent>
            </Card>

            {/* Sales List */}
            <div className="space-y-4">
              {salesLoading ? (
                <div className="text-center py-8">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                  <p className="mt-2 text-gray-500">Loading sales...</p>
                </div>
              ) : sales.length === 0 ? (
                <Card>
                  <CardContent className="text-center py-8">
                    <p className="text-gray-500">No sales transactions found</p>
                    {canCreateTransaction && (
                      <Button 
                        onClick={() => navigate('/transactions/sales/new')}
                        className="mt-4"
                      >
                        <Plus className="h-4 w-4 mr-2" />
                        Create First Sale
                      </Button>
                    )}
                  </CardContent>
                </Card>
              ) : (
                sales.map((sale) => (
                  <Card key={sale.id} className="hover:shadow-md transition-shadow">
                    <CardContent className="p-6">
                      <div className="flex justify-between items-start">
                        <div className="flex-1">
                          <div className="flex items-center gap-3 mb-3">
                            <TrendingUp className="h-5 w-5 text-green-600" />
                            <h3 className="text-lg font-semibold">{sale.transaction_number}</h3>
                            <Badge variant="outline">{sale.invoice_number}</Badge>
                            <Badge className="bg-green-100 text-green-800">{sale.status}</Badge>
                          </div>
                          
                          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-sm text-gray-600 mb-4">
                            <div>
                              <span className="text-gray-500">Customer:</span>
                              <p className="font-medium">{sale.customer?.name}</p>
                            </div>
                            <div>
                              <span className="text-gray-500">Vehicle:</span>
                              <p className="font-medium">
                                {sale.vehicle?.brand} {sale.vehicle?.model}
                              </p>
                            </div>
                            <div>
                              <span className="text-gray-500">Payment:</span>
                              <p className="font-medium capitalize">{sale.payment_method}</p>
                            </div>
                            <div>
                              <span className="text-gray-500">Date:</span>
                              <p className="font-medium">{formatDate(sale.transaction_date)}</p>
                            </div>
                          </div>

                          <div className="flex items-center justify-between">
                            <div>
                              <span className="text-gray-500">Total Amount:</span>
                              <p className="text-lg font-bold text-green-600">
                                {formatPrice(sale.total_amount)}
                              </p>
                            </div>
                            <div className="text-right text-sm text-gray-500">
                              <p>Cashier: {sale.cashier?.full_name}</p>
                            </div>
                          </div>
                        </div>
                        
                        <div className="ml-4">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => navigate(`/transactions/sales/${sale.id}`)}
                          >
                            <Eye className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))
              )}
            </div>

            {/* Sales Pagination */}
            {salesTotalPages > 1 && (
              <div className="flex justify-center gap-2 mt-8">
                <Button
                  variant="outline"
                  onClick={() => setSalesPage(salesPage - 1)}
                  disabled={salesPage === 1}
                >
                  Previous
                </Button>
                <span className="flex items-center px-4 text-sm text-gray-600">
                  Page {salesPage} of {salesTotalPages}
                </span>
                <Button
                  variant="outline"
                  onClick={() => setSalesPage(salesPage + 1)}
                  disabled={salesPage === salesTotalPages}
                >
                  Next
                </Button>
              </div>
            )}
          </TabsContent>
        </Tabs>
      </main>
    </div>
  );
}
