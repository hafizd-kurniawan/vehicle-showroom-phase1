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
import { transactionService, CreatePurchaseTransactionRequest } from '../services/transactionService';
import { customerService, Customer } from '../services/customerService';
import { vehicleService, Vehicle } from '../services/vehicleService';

export default function PurchaseFormPage() {
  const [formData, setFormData] = useState<CreatePurchaseTransactionRequest>({
    vehicle_id: 0,
    customer_id: 0,
    vehicle_price: 0,
    tax_amount: 0,
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
      // Load all vehicles for purchase transactions
      const response = await vehicleService.list(1, 100);
      setVehicles(response.vehicles);
    } catch (error) {
      console.error('Failed to load vehicles:', error);
    }
  };

  const calculateTotal = () => {
    return formData.vehicle_price + formData.tax_amount;
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

      await transactionService.createPurchase(submitData);
      toast({
        title: "Success",
        description: "Purchase transaction created successfully",
      });
      
      navigate('/transactions');
    } catch (error) {
      console.error('Failed to create purchase transaction:', error);
      toast({
        title: "Error",
        description: "Failed to create purchase transaction",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (field: keyof CreatePurchaseTransactionRequest, value: string | number) => {
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
                New Purchase Transaction
              </h1>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Card>
          <CardHeader>
            <CardTitle>Create Purchase Transaction</CardTitle>
            <CardDescription>
              Record a vehicle purchase from a customer
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
                    <SelectValue placeholder="Select vehicle to purchase" />
                  </SelectTrigger>
                  <SelectContent>
                    {vehicles.map((vehicle) => (
                      <SelectItem key={vehicle.id} value={vehicle.id.toString()}>
                        {vehicle.brand} {vehicle.model} {vehicle.variant} - {vehicle.vehicle_code}
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
                    <p className="text-sm text-blue-600">
                      Status: {selectedVehicle.status}
                    </p>
                  </div>
                )}
              </div>

              {/* Customer Selection */}
              <div className="space-y-2">
                <Label htmlFor="customer_id">Customer (Seller) *</Label>
                <Select
                  value={formData.customer_id.toString()}
                  onValueChange={(value) => handleInputChange('customer_id', parseInt(value))}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select customer selling the vehicle" />
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
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <Label htmlFor="vehicle_price">Purchase Price (IDR) *</Label>
                  <Input
                    id="vehicle_price"
                    type="number"
                    min="0"
                    placeholder="Enter purchase price"
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
              </div>

              {/* Total Calculation */}
              <Card className="bg-red-50 border-red-200">
                <CardContent className="p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Calculator className="h-5 w-5 text-red-600" />
                    <h3 className="font-semibold text-red-800">Total Calculation</h3>
                  </div>
                  <div className="space-y-1 text-sm">
                    <div className="flex justify-between">
                      <span>Purchase Price:</span>
                      <span>{formatPrice(formData.vehicle_price)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Tax Amount:</span>
                      <span>{formatPrice(formData.tax_amount)}</span>
                    </div>
                    <hr className="border-red-300" />
                    <div className="flex justify-between font-bold text-lg text-red-800">
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
                    onValueChange={(value: 'cash' | 'transfer' | 'check') => handleInputChange('payment_method', value)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select payment method" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="cash">Cash</SelectItem>
                      <SelectItem value="transfer">Bank Transfer</SelectItem>
                      <SelectItem value="check">Check</SelectItem>
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
                      Create Purchase Transaction
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
