import React, { useState, ChangeEvent } from 'react'
import {
  ArrowUpRight,
  ArrowDownLeft,
  CreditCard,
  Building,
  Wallet,
  DollarSign,
  Shield,
  CheckCircle,
  AlertCircle,
  Sun,
  Moon,
} from 'lucide-react'
import AppLayout from './AppLayout'
import { CancelIcon } from 'src/components/icons/Icons'

// Type definitions
interface OffRampData {
  currency: string
  bank_id: number
  account_holder_name: string
  amount: string
  payment_purpose: string
  account_number: string
  account_type: string
  asset: string
  master_wallet: string
  account_id: string
  method: string
}

interface OnRampData {
  account_id: string
  amount: string
  fiat: string
  asset: string
  country: string
  method: string
}

interface ValidationResult {
  isValid: boolean
  errors: Record<string, string>
}

interface SelectOption {
  value: string
  label: string
}

interface InputFieldProps {
  label: string
  type?: string
  value: string
  onChange: (e: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => void
  error?: string
  placeholder?: string
  required?: boolean
  options?: SelectOption[]
}

type FormType = 'onramp' | 'offramp' | null

// Validation functions
const validateOffRamp = (data: OffRampData): ValidationResult => {
  const errors: Record<string, string> = {}

  if (!data.account_holder_name.trim()) {
    errors.account_holder_name = 'Account holder name is required'
  }
  if (!data.amount || parseFloat(data.amount) <= 0) {
    errors.amount = 'Valid amount is required'
  }
  if (!data.payment_purpose.trim()) {
    errors.payment_purpose = 'Payment purpose is required'
  }
  if (!data.account_number.trim()) {
    errors.account_number = 'Account number is required'
  }
  if (!data.account_type) {
    errors.account_type = 'Account type is required'
  }
  if (!data.asset) {
    errors.asset = 'Asset is required'
  }
  if (!data.master_wallet.trim()) {
    errors.master_wallet = 'Master wallet is required'
  }
  if (!data.account_id.trim()) {
    errors.account_id = 'Account ID is required'
  }
  if (!data.method) {
    errors.method = 'Payment method is required'
  }

  return { isValid: Object.keys(errors).length === 0, errors }
}

const validateOnRamp = (data: OnRampData): ValidationResult => {
  const errors: Record<string, string> = {}

  if (!data.account_id.trim()) {
    errors.account_id = 'Account ID is required'
  }
  if (!data.amount || parseFloat(data.amount) <= 0) {
    errors.amount = 'Valid amount is required'
  }
  if (!data.fiat) {
    errors.fiat = 'Fiat currency is required'
  }
  if (!data.asset) {
    errors.asset = 'Asset is required'
  }
  if (!data.country) {
    errors.country = 'Country is required'
  }
  if (!data.method) {
    errors.method = 'Payment method is required'
  }

  return { isValid: Object.keys(errors).length === 0, errors }
}

const CryptoRampInterface: React.FC = () => {
  const [isDarkMode, setIsDarkMode] = useState<boolean>(true)
  const [showForm, setShowForm] = useState<FormType>(null)
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false)
  const [submitSuccess, setSubmitSuccess] = useState<boolean>(false)

  const [offRampData, setOffRampData] = useState<OffRampData>({
    currency: 'USD',
    bank_id: 1,
    account_holder_name: '',
    amount: '',
    payment_purpose: '',
    account_number: '',
    account_type: '',
    asset: '',
    master_wallet: '',
    account_id: '',
    method: '',
  })

  const [onRampData, setOnRampData] = useState<OnRampData>({
    account_id: '',
    amount: '',
    fiat: 'USD',
    asset: 'USDC_POLYGON',
    country: 'US',
    method: 'Wire',
  })

  const [errors, setErrors] = useState<Record<string, string>>({})

  const handleOffRampSubmit = async (): Promise<void> => {
    const validation = validateOffRamp(offRampData)

    if (!validation.isValid) {
      setErrors(validation.errors)
      return
    }

    setIsSubmitting(true)
    setErrors({})

    // Simulate API call
    await new Promise((resolve) => setTimeout(resolve, 2000))

    setIsSubmitting(false)
    setSubmitSuccess(true)
    setTimeout(() => setSubmitSuccess(false), 3000)
  }

  const handleOnRampSubmit = async (): Promise<void> => {
    const validation = validateOnRamp(onRampData)

    if (!validation.isValid) {
      setErrors(validation.errors)
      return
    }

    setIsSubmitting(true)
    setErrors({})

    // Simulate API call
    await new Promise((resolve) => setTimeout(resolve, 2000))

    setIsSubmitting(false)
    setSubmitSuccess(true)
    setTimeout(() => setSubmitSuccess(false), 3000)
  }

  const toggleTheme = (): void => {
    setIsDarkMode(!isDarkMode)
  }

  const InputField: React.FC<InputFieldProps> = ({
    label,
    type = 'text',
    value,
    onChange,
    error,
    placeholder,
    required = true,
    options = null,
  }) => (
    <div className="space-y-2">
      <label
        className={`block text-sm font-medium ${
          isDarkMode ? 'text-gray-300' : 'text-gray-700'
        }`}
      >
        {label} {required && <span className="text-red-500">*</span>}
      </label>
      {options ? (
        <select
          value={value}
          onChange={onChange}
          className={`w-full px-4 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-orange-500 transition-all duration-200 ${
            isDarkMode
              ? `bg-gray-800/50 border-gray-600 hover:border-gray-500 text-white ${
                  error ? 'border-red-500 bg-red-900/20' : ''
                }`
              : `bg-white border-gray-300 hover:border-gray-400 text-gray-900 ${
                  error ? 'border-red-500 bg-red-50' : ''
                }`
          }`}
        >
          <option value="">Select {label}</option>
          {options.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      ) : (
        <input
          type={type}
          value={value}
          onChange={onChange}
          placeholder={placeholder}
          className={`w-full px-4 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-orange-500 transition-all duration-200 ${
            isDarkMode
              ? `bg-gray-800/50 border-gray-600 hover:border-gray-500 text-white placeholder-gray-400 ${
                  error ? 'border-red-500 bg-red-900/20' : ''
                }`
              : `bg-white border-gray-300 hover:border-gray-400 text-gray-900 placeholder-gray-500 ${
                  error ? 'border-red-500 bg-red-50' : ''
                }`
          }`}
        />
      )}
      {error && (
        <div className="flex items-center gap-2 text-red-500 text-sm">
          <AlertCircle size={16} />
          {error}
        </div>
      )}
    </div>
  )

  const themeClasses = isDarkMode
    ? 'bg-gradient-to-br from-gray-900 via-orange-900 to-red-900'
    : 'bg-gradient-to-br from-orange-50 via-orange-100 to-red-50'

  const cardClasses = isDarkMode
    ? 'bg-gray-800/30 border-gray-700'
    : 'bg-white/90 border-gray-200'

  // Options for dropdowns
  const fiatOptions: SelectOption[] = [
    { value: 'USD', label: 'USD - US Dollar' },
    { value: 'EUR', label: 'EUR - Euro' },
    { value: 'GBP', label: 'GBP - British Pound' },
  ]

  const assetOptions: SelectOption[] = [
    { value: 'USDC_POLYGON', label: 'USDC (Polygon)' },
    { value: 'USDT_ETHEREUM', label: 'USDT (Ethereum)' },
    { value: 'BTC', label: 'Bitcoin' },
  ]

  const countryOptions: SelectOption[] = [
    { value: 'US', label: 'United States' },
    { value: 'GB', label: 'United Kingdom' },
    { value: 'CA', label: 'Canada' },
  ]

  const onRampMethodOptions: SelectOption[] = [
    { value: 'Wire', label: 'Wire Transfer' },
    { value: 'ACH', label: 'ACH Transfer' },
    { value: 'Card', label: 'Credit/Debit Card' },
  ]

  const offRampMethodOptions: SelectOption[] = [
    { value: 'Wire', label: 'Wire Transfer' },
    { value: 'ACH', label: 'ACH Transfer' },
    { value: 'SEPA', label: 'SEPA Transfer' },
  ]

  const accountTypeOptions: SelectOption[] = [
    { value: 'checking', label: 'Checking' },
    { value: 'savings', label: 'Savings' },
    { value: 'business', label: 'Business' },
  ]

  return (
    <AppLayout
      child={
        <div
          className={`min-h-screen transition-all duration-500 ${themeClasses} p-4`}
        >
          <div className="max-w-4xl mx-auto">
            {/* Header with Theme Toggle */}
            <div className="flex justify-between items-center mb-8">
              <div className="text-center flex-1">
                <h1
                  className={`text-4xl font-bold mb-2 ${
                    isDarkMode ? 'text-white' : 'text-gray-900'
                  }`}
                >
                  Withdraw or Deposit{' '}
                  <span className="text-orange-500">Funds</span>
                </h1>
                <p
                  className={`${
                    isDarkMode ? 'text-gray-300' : 'text-gray-600'
                  }`}
                >
                  Seamless crypto-to-fiat conversions
                </p>
              </div>
              <button
                onClick={toggleTheme}
                className={`p-3 rounded-xl transition-all duration-300 ${
                  isDarkMode
                    ? 'bg-gray-800/50 hover:bg-gray-700/50 text-yellow-400'
                    : 'bg-white/50 hover:bg-white/70 text-orange-600'
                }`}
              >
                {isDarkMode ? <Sun size={24} /> : <Moon size={24} />}
              </button>
            </div>

            {/* Success Message */}
            {submitSuccess && (
              <div
                className={`mb-6 p-4 border rounded-xl flex items-center gap-3 ${
                  isDarkMode
                    ? 'bg-green-900/50 border-green-500 text-green-400'
                    : 'bg-green-50 border-green-200 text-green-700'
                }`}
              >
                <CheckCircle size={20} />
                Transaction submitted successfully!
              </div>
            )}

            {/* Main Action Buttons */}
            {!showForm && (
              <div
                className={`backdrop-blur-sm rounded-2xl p-8 shadow-2xl border ${cardClasses}`}
              >
                <div className="text-center mb-8">
                  <h2
                    className={`text-2xl font-semibold mb-4 ${
                      isDarkMode ? 'text-white' : 'text-gray-900'
                    }`}
                  >
                    Choose Your Transaction
                  </h2>
                  <p
                    className={`${
                      isDarkMode ? 'text-gray-400' : 'text-gray-600'
                    }`}
                  >
                    Select whether you want to withdraw or deposit funds
                  </p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {/* OnRamp Button */}
                  <button
                    onClick={() => setShowForm('onramp')}
                    className="group relative p-8 bg-gradient-to-br from-orange-500 to-red-600 hover:from-orange-600 hover:to-red-700 text-white rounded-2xl transition-all duration-300 shadow-lg shadow-orange-500/25 hover:shadow-orange-500/40 transform hover:scale-105"
                  >
                    <div className="flex flex-col items-center gap-4">
                      <div className="p-4 bg-white/20 rounded-full group-hover:bg-white/30 transition-all duration-300">
                        <ArrowUpRight size={32} />
                      </div>
                      <div className="text-center">
                        <h3 className="text-xl font-bold mb-2">Deposit Fund</h3>
                        <p className="text-orange-100">
                          Deposit fiat to cryptocurrency
                        </p>
                      </div>
                    </div>
                    <div className="absolute inset-0 bg-white/10 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
                  </button>

                  {/* OffRamp Button */}
                  <button
                    onClick={() => setShowForm('offramp')}
                    className="group relative p-8 bg-gradient-to-br from-red-500 to-orange-600 hover:from-red-600 hover:to-orange-700 text-white rounded-2xl transition-all duration-300 shadow-lg shadow-red-500/25 hover:shadow-red-500/40 transform hover:scale-105"
                  >
                    <div className="flex flex-col items-center gap-4">
                      <div className="p-4 bg-white/20 rounded-full group-hover:bg-white/30 transition-all duration-300">
                        <ArrowDownLeft size={32} />
                      </div>
                      <div className="text-center">
                        <h3 className="text-xl font-bold mb-2">
                          Withdraw Fund
                        </h3>
                        <p className="text-red-100">
                          Convert cryptocurrency to fiat
                        </p>
                      </div>
                    </div>
                    <div className="absolute inset-0 bg-white/10 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
                  </button>
                </div>
              </div>
            )}

            {/* Forms */}
            {showForm && (
              <div
                className={`backdrop-blur-sm rounded-2xl p-8 shadow-2xl border ${cardClasses}`}
              >
                {/* Back Button */}
                <button
                  onClick={() => {
                    setShowForm(null)
                    setErrors({})
                  }}
                  className={`mb-6 px-4 py-2 rounded-lg transition-all duration-300 ${
                    isDarkMode
                      ? 'bg-gray-700 hover:bg-gray-600 text-gray-300'
                      : 'bg-gray-200 hover:bg-gray-300 text-gray-700'
                  }`}
                >
                  ← Back to Options
                </button>

                {showForm === 'onramp' ? (
                  <div className="space-y-6">
                    <div className="flex items-center gap-3 mb-6">
                      <div className="p-3 bg-orange-500/20 rounded-xl">
                        <Wallet className="text-orange-500" size={24} />
                      </div>
                      <div>
                        <h2
                          className={`text-xl font-semibold ${
                            isDarkMode ? 'text-white' : 'text-gray-900'
                          }`}
                        >
                          Buy Crypto
                        </h2>
                        <p
                          className={`${
                            isDarkMode ? 'text-gray-400' : 'text-gray-600'
                          }`}
                        >
                          Convert fiat to cryptocurrency
                        </p>
                      </div>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <InputField
                        label="Account ID"
                        value={onRampData.account_id}
                        onChange={(e) =>
                          setOnRampData({
                            ...onRampData,
                            account_id: e.target.value,
                          })
                        }
                        error={errors.account_id}
                        placeholder="Enter your account ID"
                      />

                      <InputField
                        label="Amount"
                        type="number"
                        value={onRampData.amount}
                        onChange={(e) =>
                          setOnRampData({
                            ...onRampData,
                            amount: e.target.value,
                          })
                        }
                        error={errors.amount}
                        placeholder="0.00"
                      />

                      <InputField
                        label="Fiat Currency"
                        value={onRampData.fiat}
                        onChange={(e) =>
                          setOnRampData({ ...onRampData, fiat: e.target.value })
                        }
                        error={errors.fiat}
                        options={fiatOptions}
                      />

                      <InputField
                        label="Asset"
                        value={onRampData.asset}
                        onChange={(e) =>
                          setOnRampData({
                            ...onRampData,
                            asset: e.target.value,
                          })
                        }
                        error={errors.asset}
                        options={assetOptions}
                      />

                      <InputField
                        label="Country"
                        value={onRampData.country}
                        onChange={(e) =>
                          setOnRampData({
                            ...onRampData,
                            country: e.target.value,
                          })
                        }
                        error={errors.country}
                        options={countryOptions}
                      />

                      <InputField
                        label="Payment Method"
                        value={onRampData.method}
                        onChange={(e) =>
                          setOnRampData({
                            ...onRampData,
                            method: e.target.value,
                          })
                        }
                        error={errors.method}
                        options={onRampMethodOptions}
                      />
                    </div>

                    <button
                      onClick={handleOnRampSubmit}
                      disabled={isSubmitting}
                      className="w-full bg-gradient-to-r from-orange-500 to-red-600 hover:from-orange-600 hover:to-red-700 text-white font-medium py-4 px-8 rounded-xl transition-all duration-300 shadow-lg shadow-orange-500/25 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                    >
                      {isSubmitting ? (
                        <>
                          <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                          Processing...
                        </>
                      ) : (
                        <>
                          <CreditCard size={20} />
                          Deposit Fund
                        </>
                      )}
                    </button>
                  </div>
                ) : (
                  <div className="space-y-6">
                    <div className="flex items-center gap-3 mb-6">
                      <div className="p-3 bg-red-500/20 rounded-xl">
                        <Building className="text-red-500" size={24} />
                      </div>
                      <div>
                        <h2
                          className={`text-xl font-semibold ${
                            isDarkMode ? 'text-white' : 'text-gray-900'
                          }`}
                        >
                          Sell Crypto
                        </h2>
                        <p
                          className={`${
                            isDarkMode ? 'text-gray-400' : 'text-gray-600'
                          }`}
                        >
                          Convert cryptocurrency to fiat
                        </p>
                      </div>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <InputField
                        label="Account Holder Name"
                        value={offRampData.account_holder_name}
                        onChange={(e) =>
                          setOffRampData({
                            ...offRampData,
                            account_holder_name: e.target.value,
                          })
                        }
                        error={errors.account_holder_name}
                        placeholder="Enter account holder name"
                      />

                      <InputField
                        label="Amount"
                        type="number"
                        value={offRampData.amount}
                        onChange={(e) =>
                          setOffRampData({
                            ...offRampData,
                            amount: e.target.value,
                          })
                        }
                        error={errors.amount}
                        placeholder="0.00"
                      />

                      <InputField
                        label="Payment Purpose"
                        value={offRampData.payment_purpose}
                        onChange={(e) =>
                          setOffRampData({
                            ...offRampData,
                            payment_purpose: e.target.value,
                          })
                        }
                        error={errors.payment_purpose}
                        placeholder="e.g., Investment withdrawal"
                      />

                      <InputField
                        label="Account Number"
                        value={offRampData.account_number}
                        onChange={(e) =>
                          setOffRampData({
                            ...offRampData,
                            account_number: e.target.value,
                          })
                        }
                        error={errors.account_number}
                        placeholder="Enter bank account number"
                      />

                      <InputField
                        label="Account Type"
                        value={offRampData.account_type}
                        onChange={(e) =>
                          setOffRampData({
                            ...offRampData,
                            account_type: e.target.value,
                          })
                        }
                        error={errors.account_type}
                        options={accountTypeOptions}
                      />

                      <InputField
                        label="Asset"
                        value={offRampData.asset}
                        onChange={(e) =>
                          setOffRampData({
                            ...offRampData,
                            asset: e.target.value,
                          })
                        }
                        error={errors.asset}
                        options={assetOptions}
                      />

                      <InputField
                        label="Master Wallet"
                        value={offRampData.master_wallet}
                        onChange={(e) =>
                          setOffRampData({
                            ...offRampData,
                            master_wallet: e.target.value,
                          })
                        }
                        error={errors.master_wallet}
                        placeholder="Enter wallet address"
                      />

                      <InputField
                        label="Account ID"
                        value={offRampData.account_id}
                        onChange={(e) =>
                          setOffRampData({
                            ...offRampData,
                            account_id: e.target.value,
                          })
                        }
                        error={errors.account_id}
                        placeholder="Enter your account ID"
                      />

                      <InputField
                        label="Payment Method"
                        value={offRampData.method}
                        onChange={(e) =>
                          setOffRampData({
                            ...offRampData,
                            method: e.target.value,
                          })
                        }
                        error={errors.method}
                        options={offRampMethodOptions}
                      />
                    </div>

                    <button
                      onClick={handleOffRampSubmit}
                      disabled={isSubmitting}
                      className="w-full bg-gradient-to-r from-red-500 to-orange-600 hover:from-red-600 hover:to-orange-700 text-white font-medium py-4 px-8 rounded-xl transition-all duration-300 shadow-lg shadow-red-500/25 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                    >
                      {isSubmitting ? (
                        <>
                          <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                          Processing...
                        </>
                      ) : (
                        <>
                          <DollarSign size={20} />
                          Withdraw Fund
                        </>
                      )}
                    </button>
                  </div>
                )}
              </div>
            )}

            {/* Security Notice */}
            <div
              className={`mt-8 p-4 border rounded-xl flex items-center gap-3 ${
                isDarkMode
                  ? 'bg-gray-800/30 border-gray-700 text-gray-300'
                  : 'bg-white/50 border-gray-300 text-gray-700'
              }`}
            >
              <Shield className="text-orange-500 flex-shrink-0" size={20} />
              <p className="text-sm">
                Your transactions are secured with bank-level encryption. All
                data is processed securely and never stored permanently.
              </p>
            </div>
          </div>
        </div>
      }
    />
  )
}

export default CryptoRampInterface
