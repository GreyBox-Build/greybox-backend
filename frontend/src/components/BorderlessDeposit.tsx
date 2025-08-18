import { ArrowLeft } from 'lucide-react'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

// Import your actual RTK Query hook
import { useOnrampBorderlessPaymentMutation } from 'src/appSlices/apiSlice'

export default function OnRampDepositForm() {
  const [formData, setFormData] = useState({
    fiat: 'USD',
    country: 'US',
    asset: 'USDT_ETHEREUM',
    amount: '10',
    accountId: 'cm8r9xplu389tph251q56h7li',
    paymentMethod: 'Wire',
  })

  const [onrampBorderlessPayment, { isLoading }] =
    useOnrampBorderlessPaymentMutation()
  const [showPreview, setShowPreview] = useState(false)
  const [submitStatus, setSubmitStatus] = useState({ type: '', message: '' })

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }))
    // Clear any previous status when form is modified
    if (submitStatus.type) {
      setSubmitStatus({ type: '', message: '' })
    }
  }

  const handleSubmit = async () => {
    setSubmitStatus({ type: '', message: '' })

    try {
      const response = await onrampBorderlessPayment(formData).unwrap()
      setSubmitStatus({
        type: 'success',
        message: `Deposit successful! Transaction ID: ${
          response.transactionId || 'Processing'
        }`,
      })
    } catch (err) {
      setSubmitStatus({
        type: 'error',
        message: 'An error occurred while processing your deposit',
      })
    }
  }

  const navigate = useNavigate()

  const isFormValid = () => {
    return (
      formData.amount &&
      parseFloat(formData.amount) > 0 &&
      formData.accountId.trim() !== ''
    )
  }

  const fiatOptions = ['USD', 'EUR', 'GBP', 'CAD', 'AUD']
  const countryOptions = [
    { value: 'US', label: 'United States' },
    { value: 'GB', label: 'United Kingdom' },
    { value: 'CA', label: 'Canada' },
    { value: 'AU', label: 'Australia' },
    { value: 'DE', label: 'Germany' },
    { value: 'FR', label: 'France' },
  ]
  const assetOptions = [
    { value: 'USDT_ETHEREUM', label: 'USDT (Ethereum)' },
    { value: 'USDC_ETHEREUM', label: 'USDC (Ethereum)' },
    { value: 'BTC', label: 'Bitcoin' },
    { value: 'ETH', label: 'Ethereum' },
  ]
  const paymentMethodOptions = [
    'Wire',
    'Bank Transfer',
    'Credit Card',
    'Debit Card',
  ]

  return (
    <div className="min-h-screen bg-gray-50 p-2 md:p-4">
      <div className="max-w-2xl pt-4 rounded-2xl shadow-lg shadow-gray-200 border border-gray-100 mx-auto">
        {/* Header */}
        <div className="mb-8 ">
          <ArrowLeft onClick={() => navigate(-1)} className="cursor-pointer" />
          <div className="flex flex-col items-center">
            <h1 className="text-3xl text-center font-bold text-orange-1 mb-2">
              Deposit Funds
            </h1>
          </div>

          <p className="text-gray-600 text-center">
            Deposit funds via digital assets
          </p>
        </div>

        {/* Status Message */}
        {submitStatus.type && (
          <div
            className={`mb-6 p-4 rounded-lg ${
              submitStatus.type === 'success'
                ? 'bg-green-50 text-green-800 border border-green-200'
                : 'bg-red-50 text-red-800 border border-red-200'
            }`}
          >
            <div className="flex items-center">
              <div
                className={`w-2 h-2 rounded-full mr-3 ${
                  submitStatus.type === 'success'
                    ? 'bg-green-600'
                    : 'bg-red-600'
                }`}
              ></div>
              {submitStatus.message}
            </div>
          </div>
        )}

        {/* Form Card */}
        <div className="bg-white  p-4 md:p-8">
          <div className="space-y-6">
            {/* Fiat Currency */}
            <div>
              <label
                htmlFor="fiat"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Fiat Currency
              </label>
              <select
                id="fiat"
                name="fiat"
                value={formData.fiat}
                onChange={handleInputChange}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
              >
                {fiatOptions.map((option) => (
                  <option key={option} value={option}>
                    {option}
                  </option>
                ))}
              </select>
            </div>

            {/* Country */}
            <div>
              <label
                htmlFor="country"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Country
              </label>
              <select
                id="country"
                name="country"
                value={formData.country}
                onChange={handleInputChange}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
              >
                {countryOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </div>

            {/* Asset */}
            <div>
              <label
                htmlFor="asset"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Asset
              </label>
              <select
                id="asset"
                name="asset"
                value={formData.asset}
                onChange={handleInputChange}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
              >
                {assetOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </div>

            {/* Amount */}
            <div>
              <label
                htmlFor="amount"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Amount ({formData.fiat})
              </label>
              <input
                type="number"
                id="amount"
                name="amount"
                value={formData.amount}
                onChange={handleInputChange}
                min="1"
                step="0.01"
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-1 focus:border-orange-1 transition-colors"
                placeholder="Enter amount"
              />
            </div>

            {/* Account ID */}
            <div>
              <label
                htmlFor="accountId"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Account ID
              </label>
              <input
                type="text"
                id="accountId"
                name="accountId"
                value={formData.accountId}
                onChange={handleInputChange}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
                placeholder="Enter account ID"
              />
            </div>

            {/* Payment Method */}
            <div>
              <label
                htmlFor="paymentMethod"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Payment Method
              </label>
              <select
                id="paymentMethod"
                name="paymentMethod"
                value={formData.paymentMethod}
                onChange={handleInputChange}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
              >
                {paymentMethodOptions.map((option) => (
                  <option key={option} value={option}>
                    {option}
                  </option>
                ))}
              </select>
            </div>

            {/* Preview Section */}
            {showPreview && (
              <div className="bg-gray-50 rounded-lg p-4 border border-gray-200">
                <h4 className="font-medium text-gray-900 mb-3">
                  Transaction Preview
                </h4>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-600">Amount:</span>
                    <span className="font-medium">
                      {formData.amount} {formData.fiat}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Asset:</span>
                    <span className="font-medium">
                      {
                        assetOptions.find((opt) => opt.value === formData.asset)
                          ?.label
                      }
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Payment Method:</span>
                    <span className="font-medium">
                      {formData.paymentMethod}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Country:</span>
                    <span className="font-medium">
                      {
                        countryOptions.find(
                          (opt) => opt.value === formData.country
                        )?.label
                      }
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Account ID:</span>
                    <span className="font-medium font-mono text-xs">
                      {formData.accountId}
                    </span>
                  </div>
                </div>
              </div>
            )}

            {/* Buttons */}
            <div className="flex gap-4 pt-4">
              <button
                type="button"
                onClick={() => setShowPreview(!showPreview)}
                className="flex-1 bg-gray-700 text-white font-medium py-3 px-6 rounded-lg hover:bg-gray-200 transition-colors"
              >
                {showPreview ? 'Hide Preview' : 'Preview'}
              </button>

              <button
                type="button"
                onClick={handleSubmit}
                disabled={isLoading || !isFormValid()}
                className="flex-1 bg-orange-1 hover:bg-orange-2 disabled:bg-gray-400 disabled:cursor-not-allowed text-white font-medium py-3 px-6 rounded-lg transition-colors flex items-center justify-center"
              >
                {isLoading ? (
                  <div className="flex items-center">
                    <div className="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent mr-2"></div>
                    Processing...
                  </div>
                ) : (
                  'Deposit'
                )}
              </button>
            </div>
          </div>
        </div>

        {/* Info Section */}
        <div className="mt-8 bg-orange-50 rounded-lg p-6 border border-orange-200">
          <h3 className="text-lg font-semibold text-orange-900 mb-3">
            How it works
          </h3>
          <div className="space-y-2 text-sm text-orange-800">
            <div className="flex items-start">
              <div className="w-2 h-2 bg-orange-600 rounded-full mt-2 mr-3 flex-shrink-0"></div>
              <p>Select your preferred fiat currency and payment method</p>
            </div>
            <div className="flex items-start">
              <div className="w-2 h-2 bg-orange-600 rounded-full mt-2 mr-3 flex-shrink-0"></div>
              <p>Choose the digital asset you want to purchase</p>
            </div>
            <div className="flex items-start">
              <div className="w-2 h-2 bg-orange-600 rounded-full mt-2 mr-3 flex-shrink-0"></div>
              <p>Enter the amount and submit your deposit request</p>
            </div>
            <div className="flex items-start">
              <div className="w-2 h-2 bg-orange-600 rounded-full mt-2 mr-3 flex-shrink-0"></div>
              <p>
                Your assets will be credited to your account once payment is
                processed
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
