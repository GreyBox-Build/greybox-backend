import React from 'react'
import { Plus, Info, Snowflake, MoreHorizontal } from 'lucide-react'
import AppLayout from './AppLayout'

interface CardData {
  number: string
  balance: string
  cardType: 'mastercard' | 'visa'
}

const VirtualCard: React.FC = () => {
  const cardData: CardData = {
    number: 'Coming Soon...',
    balance: '*****',
    cardType: 'mastercard',
  }

  const Cards = () => (
    <div className="p-6 bg-gray-50 min-h-screen">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Cards</h1>

      {/* Card Display */}
      <div className="relative mb-8">
        <div className="bg-gradient-to-br from-orange-600 via-orange-700 to-orange-900 rounded-2xl p-6 text-white shadow-xl">
          <div className="flex justify-between items-start mb-8">
            <div className="w-8 h-8 bg-white/20 rounded-full flex items-center justify-center">
              <div className="w-4 h-4 bg-white rounded-full"></div>
            </div>
          </div>

          <div className="mb-8">
            <div className="text-xl font-mono tracking-wider">
              {cardData.number}
            </div>
          </div>

          <div className="flex justify-between items-end">
            <div>
              <div className="text-2xl font-bold">${cardData.balance}</div>
            </div>
            <div className="flex space-x-1">
              <div className="w-8 h-8 bg-red-500 rounded-full opacity-90"></div>
              <div className="w-8 h-8 bg-yellow-500 rounded-full opacity-90 -ml-2"></div>
            </div>
          </div>
        </div>
      </div>

      {/* Action Buttons */}
      <div className="grid grid-cols-4 gap-4 mb-8">
        <button className="flex flex-col items-center p-4 bg-white rounded-2xl shadow-sm">
          <div className="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center mb-2">
            <Info className="w-6 h-6 text-orange-600" />
          </div>
          <span className="text-sm text-gray-700">Details</span>
        </button>

        <button className="flex flex-col items-center p-4 bg-white rounded-2xl shadow-sm">
          <div className="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center mb-2">
            <Plus className="w-6 h-6 text-orange-600" />
          </div>
          <span className="text-sm text-gray-700">Add Money</span>
        </button>

        <button className="flex flex-col items-center p-4 bg-white rounded-2xl shadow-sm">
          <div className="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center mb-2">
            <Snowflake className="w-6 h-6 text-orange-600" />
          </div>
          <span className="text-sm text-gray-700">Freeze</span>
        </button>

        <button className="flex flex-col items-center p-4 bg-white rounded-2xl shadow-sm">
          <div className="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center mb-2">
            <MoreHorizontal className="w-6 h-6 text-orange-600" />
          </div>
          <span className="text-sm text-gray-700">More</span>
        </button>
      </div>

      {/* Transaction History
      <div className="bg-white rounded-2xl p-6 shadow-sm">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-semibold text-gray-900">
            Transaction history
          </h3>
          <button className="text-orange-600 text-sm font-medium">
            See all
          </button>
        </div>

        <div className="space-y-4">
          {recentTransactions.map((transaction, index) => (
            <div key={index} className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center">
                  <div className="w-6 h-6 bg-gray-400 rounded"></div>
                </div>
                <div>
                  <div className="font-medium text-gray-900">
                    {transaction.name}
                  </div>
                  <div className="text-sm text-gray-500">
                    {transaction.date}
                  </div>
                </div>
              </div>
              <div className="text-right">
                <div
                  className={`font-semibold ${
                    transaction.amount >= 0 ? 'text-green-600' : 'text-gray-900'
                  }`}
                >
                  {transaction.amount >= 0 ? '+' : '-'}$
                  {Math.abs(transaction.amount).toFixed(2)}
                </div>
                <div className="text-sm text-green-600">• Success</div>
              </div>
            </div>
          ))}
        </div>
      </div> */}
    </div>
  )

  return (
    <AppLayout
      child={
        <div className="max-w-xl mx-auto bg-white min-h-screen">
          {/* Main Content */}
          <div className="flex-1">
            {' '}
            <Cards />
          </div>
        </div>
      }
    />
  )
}

export default VirtualCard
