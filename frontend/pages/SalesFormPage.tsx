import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useToast } from '@/components/ui/use-toast';
import { ArrowLeft, Save, Calculator } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { transactionService, CreateSalesTransactionRequest } from '../services/transactionService';
import { customerService, Customer } from '../services/customerService';
import { vehicleService, Vehicle } from '../services/vehicleService';

export default function SalesFormPage() {
  const [formData, setFormData] = useState<CreateSalesTransactionRequest>({
    vehicle_id: 0,
    customer_id: 0,
    vehicle_price: 0,
    tax_amount: 0,
    discount_amount: 0,
    payment_method: 'cash',
    payment_reference: '',
    notes: '',
  });
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [selectedVehicle, setSelectedVehicle] = useState<Vehicle | null>(null);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { toast } = useToast();

  useEffect(() => {
    loadCustomers();
    loadVehicles();
  }, []);

  useEffect(() => {
    if (formData.vehicle_id) {
      const vehicle = vehicles.find(v => v.id === formData.vehicle_id);
      setSelectedVehicle(vehicle || null);
      if (vehicle && vehicle.approved_selling_price) {
        setFormData(prev => ({
          ...prev,
          vehicle_price: vehicle.approved_selling_price || 0
        }));
      }
    }
  }, [formData.vehicle_id, vehicles]);

  const loadCustomers = async () => {
    try {
      const response = await customerService.list(1, 100);
      setCustomers(response.customers);
    } catch (error) {
      console.error('Failed to load customers:', error);
    }
  };

  const loadVehicles = async () => {
    try {
      const response = await vehicleService.list(1, 100, '', 'ready_to_sell');
      setVehicles(response.vehicles);
    } catch (error) {
      console.error('Failed to load vehicles:', error);
    }
  };

  const calculateTotal = () => {
    return formData.vehicle_price + formData.tax_amount - formData.discount_amount;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.vehicle_id || !formData.customer_id || !formData.vehicle_price) {
      toast({
        title: "Error",
        description: "Please fill in all required fields",
        variant: "destructive",
      });
      return;
    }

    try {
      setLoading(true);
      
      const submitData = {
        ...formData,
        payment_reference: formData.payment_reference || undefined,
        notes: formData.notes || undefined,
      };

      await transactionService.createSales(submitData);
      toast({
        title: "Success",
        description: "Sales transaction created successfully",
      });
      
      navigate('/transactions');
    } catch (error) {
      console.error('Failed to create sales transaction:', error);
      toast({
        title: "Error",
        description: "Failed to create sales transaction",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (field: keyof CreateSalesTransactionRequest, value: string | number) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(price);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <Button
                variant="ghost"
                onClick={() => navigate('/transactions')}
                className="mr-4"
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <h1 className="text-xl font-semibold text-gray-900">
                New Sales Transaction
              </h1>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Card>
          <CardHeader>
            <CardTitle>Create Sales Transaction</CardTitle>
            <CardDescription>
              Record a vehicle sale to a customer
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Vehicle Selection */}
              <div className="space-y-2">
                <Label htmlFor="vehicle_id">Vehicle *</Label>
                <Select
                  value={formData.vehicle_id.toString()}
                  onValueChange={(value) => handleInputChange('vehicle_id', parseInt(value))}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select vehicle to sell" />
                  </SelectTrigger>
                  <SelectContent>
                    {vehicles.map((vehicle) => (
                      <SelectItem key={vehicle.id} value={vehicle.id.toString()}>
                        {vehicle.brand} {vehicle.model} {vehicle.variant} - {vehicle.vehicle_code}
                        {vehicle.approved_selling_price && (
                          <span className="text-green-600 ml-2">
                            ({formatPrice(vehicle.approved_selling_price)})
                          </span>
                        )}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {selectedVehicle && (
                  <div className="mt-2 p-3 bg-blue-50 rounded-lg">
                    <p className="text-sm text-blue-800">
                      <strong>Selected Vehicle:</strong> {selectedVehicle.brand} {selectedVehicle.model} {selectedVehicle.variant}
                    </p>
                    <p className="text-sm text-blue-600">
                      Chassis: {selectedVehicle.chassis_number} | Year: {selectedVehicle.year}
                    </p>
                    {selectedVehicle.approved_selling_price && (
                      <p className="text-sm text-green-600">
                        Approved Price: {formatPrice(selectedVehicle.approved_selling_price)}
                      </p>
                    )}
                  </div>
                )}
              </div>

              {/* Customer Selection */}
              <div className="space-y-2">
                <Label htmlFor="customer_id">Customer *</Label>
                <Select
                  value={formData.customer_id.toString()}
                  onValueChange={(value) => handleInputChange('customer_id', parseInt(value))}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select customer" />
                  </SelectTrigger>
                  <SelectContent>
                    {customers.map((customer) => (
                      <SelectItem key={customer.id} value={customer.id.toString()}>
                        {customer.name} ({customer.customer_code}) - {customer.type}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              {/* Pricing */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="space-y-2">
                  <Label htmlFor="vehicle_price">Vehicle Price (IDR) *</Label>
                  <Input
                    id="vehicle_price"
                    type="number"
                    min="0"
                    placeholder="Enter vehicle price"
                    value={formData.vehicle_price}
                    onChange={(e) => handleInputChange('vehicle_price', parseFloat(e.target.value) || 0)}
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="tax_amount">Tax Amount (IDR)</Label>
                  <Input
                    id="tax_amount"
                    type="number"
                    min="0"
                    placeholder="Enter tax amount"
                    value={formData.tax_amount}
                    onChange={(e) => handleInputChange('tax_amount', parseFloat(e.target.value) || 0)}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="discount_amount">Discount Amount (IDR)</Label>
                  <Input
                    id="discount_amount"
                    type="number"
                    min="0"
                    placeholder="Enter discount amount"
                    value={formData.discount_amount}
                    onChange={(e) => handleInputChange('discount_amount', parseFloat(e.target.value) || 0)}
                  />
                </div>
              </div>

              {/* Total Calculation */}
              <Card className="bg-green-50 border-green-200">
                <CardContent className="p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Calculator className="h-5 w-5 text-green-600" />
                    <h3 className="font-semibold text-green-800">Total Calculation</h3>
                  </div>
                  <div className="space-y-1 text-sm">
                    <div className="flex justify-between">
                      <span>Vehicle Price:</span>
                      <span>{formatPrice(formData.vehicle_price)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Tax Amount:</span>
                      <span>{formatPrice(formData.tax_amount)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Discount:</span>
                      <span>-{formatPrice(formData.discount_amount)}</span>
                    </div>
                    <hr className="border-green-300" />
                    <div className="flex justify-between font-bold text-lg text-green-800">
                      <span>Total Amount:</span>
                      <span>{formatPrice(calculateTotal())}</span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Payment Information */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <Label htmlFor="payment_method">Payment Method *</Label>
                  <Select
                    value={formData.payment_method}
                    onValueChange={(value: 'cash' | 'transfer' | 'check' | 'credit') => handleInputChange('payment_method', value)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select payment method" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="cash">Cash</SelectItem>
                      <SelectItem value="transfer">Bank Transfer</SelectItem>
                      <SelectItem value="check">Check</SelectItem>
                      <SelectItem value="credit">Credit</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="payment_reference">Payment Reference</Label>
                  <Input
                    id="payment_reference"
                    type="text"
                    placeholder="Enter payment reference (optional)"
                    value={formData.payment_reference}
                    onChange={(e) => handleInputChange('payment_reference', e.target.value)}
                  />
                </div>
              </div>

              {/* Notes */}
              <div className="space-y-2">
                <Label htmlFor="notes">Notes</Label>
                <Textarea
                  id="notes"
                  placeholder="Enter additional notes (optional)"
                  value={formData.notes}
                  onChange={(e) => handleInputChange('notes', e.target.value)}
                  rows={3}
                />
              </div>

              {/* Submit Button */}
              <div className="flex gap-4 pt-4">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => navigate('/transactions')}
                  className="flex-1"
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={loading}
                  className="flex-1"
                >
                  {loading ? (
                    "Creating..."
                  ) : (
                    <>
                      <Save className="h-4 w-4 mr-2" />
                      Create Sales Transaction
                    </>
                  )}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </main>
    </div>
  );
}
