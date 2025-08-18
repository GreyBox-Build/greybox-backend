import { Search, Filter, ArrowUpRight, ArrowDownLeft } from 'lucide-react'
import AppLayout from './AppLayout'

interface Transaction {
  id: string
  name: string
  type: 'debit' | 'credit'
  amount: number
  date: string
  status: 'success' | 'pending' | 'failed'
  category?: string
}
const transactions: Transaction[] = [
  {
    id: '1',
    name: 'OWOLABI SAMUEL AWOLU',
    type: 'debit',
    amount: 105,
    date: '11 July 2025',
    status: 'success',
  },
  {
    id: '2',
    name: 'OWOLABI SAMUEL AWOLU',
    type: 'debit',
    amount: 93.04,
    date: '11 July 2025',
    status: 'success',
  },
  {
    id: '3',
    name: 'PAYPAL',
    type: 'credit',
    amount: 95.54,
    date: '11 July 2025',
    status: 'success',
  },
  {
    id: '4',
    name: 'OWOLABI SAMUEL AWOLU',
    type: 'debit',
    amount: 20.71,
    date: '11 July 2025',
    status: 'success',
  },
  {
    id: '5',
    name: 'OWOLABI SAMUEL AWOLU',
    type: 'debit',
    amount: 103.0,
    date: '10 July 2025',
    status: 'success',
  },
  {
    id: '6',
    name: 'PAYPAL',
    type: 'credit',
    amount: 105.71,
    date: '10 July 2025',
    status: 'success',
  },
  {
    id: '7',
    name: 'OWOLABI SAMUEL AWOLU',
    type: 'debit',
    amount: 104.94,
    date: '10 July 2025',
    status: 'success',
  },
  {
    id: '8',
    name: 'PAYPAL',
    type: 'credit',
    amount: 106.94,
    date: '10 July 2025',
    status: 'success',
  },
]
export const TransactionHistory = () => (
  <AppLayout
    child={
      <div className="p-6 bg-white shadow-md w-full max-w-xl mx-auto min-h-screen">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Transactions</h1>

        {/* Search and Filter */}
        <div className="flex space-x-3 mb-6">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search here"
              className="w-full pl-10 pr-4 py-3 bg-white rounded-xl border border-gray-200 focus:outline-none focus:border-orange-500"
            />
          </div>
          <button className="p-3 bg-white rounded-xl border border-gray-200">
            <Filter className="w-5 h-5 text-gray-600" />
          </button>
        </div>

        {/* Filter Tabs */}
        <div className="mb-6">
          <button className="text-gray-900 font-medium border-b-2 border-orange-500 pb-2">
            All Transactions
          </button>
        </div>

        {/* Transactions List */}
        <div className="space-y-1">
          {transactions.map((transaction) => (
            <div
              key={transaction.id}
              className="bg-white p-4 rounded-xl shadow-sm"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                    {transaction.type === 'credit' ? (
                      <ArrowDownLeft className="w-5 h-5 text-orange-1" />
                    ) : (
                      <ArrowUpRight className="w-5 h-5 text-orange-1" />
                    )}
                  </div>
                  <div>
                    <div className="font-medium text-gray-900">
                      {transaction.name}
                    </div>
                    <div className="text-sm text-gray-500">
                      Bank account • {transaction.date}
                    </div>
                  </div>
                </div>
                <div className="text-right">
                  <div
                    className={`font-semibold ${
                      transaction.type === 'credit'
                        ? 'text-green-600'
                        : 'text-gray-900'
                    }`}
                  >
                    {transaction.type === 'credit' ? '+' : '-'}£
                    {transaction.amount.toFixed(2)}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    }
  />
)
