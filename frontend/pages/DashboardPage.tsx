import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel, 
  DropdownMenuSeparator, 
  DropdownMenuTrigger 
} from '@/components/ui/dropdown-menu';
import { Car, Users, Receipt, Wrench, LogOut, Plus, TrendingUp, TrendingDown, DollarSign, BarChart3, FileText } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { transactionService, DashboardStats } from '../services/transactionService';

export default function DashboardPage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);

  const canCreateVehicle = user?.role === 'admin' || user?.role === 'cashier';
  const canCreateCustomer = user?.role === 'admin' || user?.role === 'cashier';
  const canCreateTransaction = user?.role === 'admin' || user?.role === 'cashier';
  const canViewDashboardStats = user?.role === 'admin';
  const canViewReports = user?.role === 'admin';
  const canAccessRepairs = user?.role === 'admin' || user?.role === 'mechanic';

  useEffect(() => {
    if (canViewDashboardStats) {
      loadDashboardStats();
    } else {
      setLoading(false);
    }
  }, [canViewDashboardStats]);

  const loadDashboardStats = async () => {
    try {
      const dashboardStats = await transactionService.getDashboardStats();
      setStats(dashboardStats);
    } catch (error) {
      console.error('Failed to load dashboard stats:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    await logout();
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(price);
  };

  const dashboardCards = [
    {
      title: 'Vehicles',
      subtitle: 'Manage vehicle inventory',
      icon: Car,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
      onClick: () => navigate('/vehicles'),
      visible: true,
    },
    {
      title: 'Customers',
      subtitle: 'Manage customer database',
      icon: Users,
      color: 'text-green-600',
      bgColor: 'bg-green-50',
      onClick: () => navigate('/customers'),
      visible: true,
    },
    {
      title: 'Transactions',
      subtitle: 'Purchase & Sales',
      icon: Receipt,
      color: 'text-orange-600',
      bgColor: 'bg-orange-50',
      onClick: () => navigate('/transactions'),
      visible: canCreateTransaction,
    },
    {
      title: 'Repairs & Parts',
      subtitle: 'Workshop Management',
      icon: Wrench,
      color: 'text-purple-600',
      bgColor: 'bg-purple-50',
      onClick: () => navigate('/repairs'),
      visible: canAccessRepairs,
    },
    {
      title: 'Reports',
      subtitle: 'Analytics & Insights',
      icon: FileText,
      color: 'text-red-600',
      bgColor: 'bg-red-50',
      onClick: () => navigate('/reports'),
      visible: canViewReports,
    },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <Car className="h-8 w-8 text-blue-600 mr-3" />
              <h1 className="text-xl font-semibold text-gray-900">
                Vehicle Showroom Dashboard
              </h1>
            </div>
            
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                  <Avatar className="h-8 w-8">
                    <AvatarFallback className="bg-blue-600 text-white">
                      {user?.full_name?.[0]?.toUpperCase() || 'U'}
                    </AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-56" align="end" forceMount>
                <DropdownMenuLabel className="font-normal">
                  <div className="flex flex-col space-y-1">
                    <p className="text-sm font-medium leading-none">{user?.full_name}</p>
                    <p className="text-xs leading-none text-muted-foreground">
                      {user?.role?.toUpperCase()}
                    </p>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={handleLogout}>
                  <LogOut className="mr-2 h-4 w-4" />
                  <span>Log out</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h2 className="text-3xl font-bold text-gray-900 mb-2">
            Welcome, {user?.full_name}!
          </h2>
          <p className="text-lg text-green-600 font-semibold">
            🎉 Phase 6: Repair & Parts Management - COMPLETED! ✅
          </p>
        </div>

        {/* Statistics Cards */}
        {canViewDashboardStats && (
          loading ? (
            <div className="text-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
              <p className="mt-2 text-gray-500">Loading dashboard...</p>
            </div>
          ) : stats && (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Total Vehicles</CardTitle>
                  <Car className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{stats.total_vehicles}</div>
                  <p className="text-xs text-muted-foreground">
                    {stats.vehicles_for_sale} ready to sell
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Total Customers</CardTitle>
                  <Users className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{stats.total_customers}</div>
                  <p className="text-xs text-muted-foreground">
                    Active customers
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Today's Revenue</CardTitle>
                  <DollarSign className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{formatPrice(stats.today_revenue)}</div>
                  <p className="text-xs text-muted-foreground">
                    {stats.today_sales} sales today
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Total Profit</CardTitle>
                  <BarChart3 className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{formatPrice(stats.total_profit)}</div>
                  <p className="text-xs text-muted-foreground">
                    All time profit
                  </p>
                </CardContent>
              </Card>
            </div>
          )
        )}

        {/* Quick Actions */}
        <div className="mb-8">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
          <div className="flex flex-wrap gap-4">
            {canCreateVehicle && (
              <Button onClick={() => navigate('/vehicles/new')} className="flex items-center gap-2">
                <Plus className="h-4 w-4" />
                Add Vehicle
              </Button>
            )}
            {canCreateCustomer && (
              <Button onClick={() => navigate('/customers/new')} variant="outline" className="flex items-center gap-2">
                <Plus className="h-4 w-4" />
                Add Customer
              </Button>
            )}
            {canCreateTransaction && (
              <Button onClick={() => navigate('/transactions/purchase/new')} variant="outline" className="flex items-center gap-2">
                <TrendingDown className="h-4 w-4" />
                New Purchase
              </Button>
            )}
            {canCreateTransaction && (
              <Button onClick={() => navigate('/transactions/sales/new')} className="flex items-center gap-2">
                <TrendingUp className="h-4 w-4" />
                New Sale
              </Button>
            )}
            {canAccessRepairs && (
              <Button onClick={() => navigate('/repairs/new')} variant="outline" className="flex items-center gap-2">
                <Wrench className="h-4 w-4" />
                New Repair
              </Button>
            )}
          </div>
        </div>

        {/* Dashboard Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {dashboardCards.filter(card => card.visible).map((card, index) => {
            const IconComponent = card.icon;
            return (
              <Card 
                key={index} 
                className="hover:shadow-lg transition-shadow cursor-pointer"
                onClick={card.onClick}
              >
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    {card.title}
                  </CardTitle>
                  <div className={`p-2 rounded-lg ${card.bgColor}`}>
                    <IconComponent className={`h-6 w-6 ${card.color}`} />
                  </div>
                </CardHeader>
                <CardContent>
                  <CardDescription className="text-xs text-muted-foreground">
                    {card.subtitle}
                  </CardDescription>
                </CardContent>
              </Card>
            );
          })}
        </div>
      </main>
    </div>
  );
}
