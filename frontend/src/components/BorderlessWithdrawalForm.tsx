import { useState } from 'react'
import { useOfframpBorderlessPaymentMutation } from 'src/appSlices/apiSlice'

export default function OffRampWithdrawalForm() {
  const [formData, setFormData] = useState({
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

  const [offrampBorderlessPayment, { isLoading }] =
    useOfframpBorderlessPaymentMutation()
  const [showPreview, setShowPreview] = useState(false)
  const [submitStatus, setSubmitStatus] = useState({ type: '', message: '' })

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value, type } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: type === 'number' ? parseInt(value) || '' : value,
    }))
    // Clear any previous status when form is modified
    if (submitStatus.type) {
      setSubmitStatus({ type: '', message: '' })
    }
  }

  const handleSubmit = async () => {
    setSubmitStatus({ type: '', message: '' })

    try {
      const response = await offrampBorderlessPayment(formData).unwrap()
      setSubmitStatus({
        type: 'success',
        message: `Withdrawal successful! Transaction ID: ${
          response.transactionId || 'Processing'
        }`,
      })
    } catch (err) {
      setSubmitStatus({
        type: 'error',
        message: 'An error occurred while processing your withdrawal',
      })
    }
  }

  const isFormValid = () => {
    return (
      formData.amount &&
      parseFloat(formData.amount) > 0 &&
      formData.account_holder_name.trim() !== '' &&
      formData.asset !== '' &&
      formData.account_number.trim() !== '' &&
      formData.account_type !== '' &&
      formData.payment_purpose !== '' &&
      formData.method !== '' &&
      formData.master_wallet.trim() !== '' &&
      formData.account_id.trim() !== ''
    )
  }

  const currencyOptions = ['USD', 'EUR', 'GBP', 'CAD', 'AUD']
  const assetOptions = [
    { value: 'USDT_ETHEREUM', label: 'USDT (Ethereum)' },
    { value: 'USDC_ETHEREUM', label: 'USDC (Ethereum)' },
    { value: 'BTC', label: 'Bitcoin' },
    { value: 'ETH', label: 'Ethereum' },
    { value: 'USDC_POLYGON', label: 'USDC (Polygon)' },
    { value: 'DAI', label: 'DAI' },
  ]
  const accountTypeOptions = ['Checking', 'Savings', 'Business']
  const methodOptions = ['Wire Transfer', 'ACH', 'SWIFT', 'Local Bank Transfer']
  const paymentPurposeOptions = [
    'Personal Investment',
    'Business Operation',
    'Trading Profit',
    'Salary Payment',
    'Investment Return',
    'Other',
  ]

  return (
    <div className="min-h-screen bg-gray-50 p-4">
      <div className="max-w-2xl mx-auto">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Crypto Off-Ramp
          </h1>
          <p className="text-gray-600">
            Convert your digital assets to fiat currency
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
        <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-8">
          <div className="space-y-6">
            {/* Currency & Bank ID Row */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label
                  htmlFor="currency"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Currency
                </label>
                <select
                  id="currency"
                  name="currency"
                  value={formData.currency}
                  onChange={handleInputChange}
                  className="w-full px-4 py-3 border border-gray-300 outline-none  rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
                >
                  {currencyOptions.map((option) => (
                    <option key={option} value={option}>
                      {option}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label
                  htmlFor="bank_id"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Bank ID
                </label>
                <input
                  type="number"
                  id="bank_id"
                  name="bank_id"
                  value={formData.bank_id}
                  onChange={handleInputChange}
                  className="w-full px-4 py-3 border border-gray-300 outline-none  rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
                  placeholder="Enter bank ID"
                />
              </div>
            </div>

            {/* Account Holder Name */}
            <div>
              <label
                htmlFor="account_holder_name"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Account Holder Name *
              </label>
              <input
                type="text"
                id="account_holder_name"
                name="account_holder_name"
                value={formData.account_holder_name}
                onChange={handleInputChange}
                className="w-full px-4 py-3 border border-gray-300 outline-none  rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
                placeholder="Enter full name as it appears on bank account"
                required
              />
            </div>

            {/* Amount & Asset Row */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label
                  htmlFor="amount"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Amount ({formData.currency}) *
                </label>
                <input
                  type="number"
                  id="amount"
                  name="amount"
                  value={formData.amount}
                  onChange={handleInputChange}
                  min="1"
                  step="0.01"
                  className="w-full px-4 py-3 border border-gray-300 outline-none  rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
                  placeholder="Enter withdrawal amount"
                  required
                />
              </div>

              <div>
                <label
                  htmlFor="asset"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Asset to Convert *
                </label>
                <select
                  id="asset"
                  name="asset"
                  value={formData.asset}
                  onChange={handleInputChange}
                  className="w-full px-4 py-3 border border-gray-300 outline-none  rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
                  required
                >
                  <option value="">Select asset</option>
                  {assetOptions.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            {/* Payment Purpose */}
            <div>
              <label
                htmlFor="payment_purpose"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Payment Purpose *
              </label>
              <select
                id="payment_purpose"
                name="payment_purpose"
                value={formData.payment_purpose}
                onChange={handleInputChange}
                className="w-full px-4 py-3 border border-gray-300 outline-none  rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
                required
              >
                <option value="">Select purpose</option>
                {paymentPurposeOptions.map((option) => (
                  <option key={option} value={option}>
                    {option}
                  </option>
                ))}
              </select>
            </div>

            {/* Account Details Row */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label
                  htmlFor="account_number"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Account Number *
                </label>
                <input
                  type="text"
                  id="account_number"
                  name="account_number"
                  value={formData.account_number}
                  onChange={handleInputChange}
                  className="w-full px-4 py-3 border border-gray-300 outline-none  rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
                  placeholder="Enter bank account number"
                  required
                />
              </div>

              <div>
                <label
                  htmlFor="account_type"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Account Type *
                </label>
                <select
                  id="account_type"
                  name="account_type"
                  value={formData.account_type}
                  onChange={handleInputChange}
                  className="w-full px-4 py-3 border border-gray-300 outline-none  rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
                  required
                >
                  <option value="">Select account type</option>
                  {accountTypeOptions.map((option) => (
                    <option key={option} value={option}>
                      {option}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            {/* Master Wallet */}
            <div>
              <label
                htmlFor="master_wallet"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Master Wallet Address *
              </label>
              <input
                type="text"
                id="master_wallet"
                name="master_wallet"
                value={formData.master_wallet}
                onChange={handleInputChange}
                className="w-full px-4 py-3 border border-gray-300 outline-none  rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors font-mono text-sm"
                placeholder="0x..."
                required
              />
            </div>

            {/* Account ID & Method Row */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label
                  htmlFor="account_id"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Account ID *
                </label>
                <input
                  type="text"
                  id="account_id"
                  name="account_id"
                  value={formData.account_id}
                  onChange={handleInputChange}
                  className="w-full px-4 py-3 border border-gray-300 outline-none  rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
                  placeholder="Enter account ID"
                  required
                />
              </div>

              <div>
                <label
                  htmlFor="method"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Transfer Method *
                </label>
                <select
                  id="method"
                  name="method"
                  value={formData.method}
                  onChange={handleInputChange}
                  className="w-full px-4 py-3 border border-gray-300 outline-none  rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors"
                  required
                >
                  <option value="">Select method</option>
                  {methodOptions.map((option) => (
                    <option key={option} value={option}>
                      {option}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            {/* Preview Section */}
            {showPreview && (
              <div className="bg-gray-50 rounded-lg p-4 border border-gray-200">
                <h4 className="font-medium text-gray-900 mb-3">
                  Withdrawal Preview
                </h4>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-600">Amount:</span>
                    <span className="font-medium">
                      {formData.amount} {formData.currency}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Asset:</span>
                    <span className="font-medium">
                      {assetOptions.find((opt) => opt.value === formData.asset)
                        ?.label || formData.asset}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Account Holder:</span>
                    <span className="font-medium">
                      {formData.account_holder_name}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Account Number:</span>
                    <span className="font-medium font-mono">
                      {formData.account_number
                        ? `****${formData.account_number.slice(-4)}`
                        : ''}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Account Type:</span>
                    <span className="font-medium">{formData.account_type}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Transfer Method:</span>
                    <span className="font-medium">{formData.method}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Payment Purpose:</span>
                    <span className="font-medium">
                      {formData.payment_purpose}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Bank ID:</span>
                    <span className="font-medium">{formData.bank_id}</span>
                  </div>
                </div>
              </div>
            )}

            {/* Buttons */}
            <div className="flex gap-4 pt-4">
              <button
                type="button"
                onClick={() => setShowPreview(!showPreview)}
                className="flex-1 bg-gray-100 text-gray-700 font-medium py-3 px-6 rounded-lg hover:bg-gray-200 transition-colors"
              >
                {showPreview ? 'Hide Preview' : 'Preview'}
              </button>

              <button
                type="button"
                onClick={handleSubmit}
                disabled={isLoading || !isFormValid()}
                className="flex-1 bg-orange-600 hover:bg-orange-700 disabled:bg-gray-400 disabled:cursor-not-allowed text-white font-medium py-3 px-6 rounded-lg transition-colors flex items-center justify-center"
              >
                {isLoading ? (
                  <div className="flex items-center">
                    <div className="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent mr-2"></div>
                    Processing...
                  </div>
                ) : (
                  'Withdraw'
                )}
              </button>
            </div>
          </div>
        </div>

        {/* Info Section */}
        <div className="mt-8 bg-orange-50 rounded-lg p-6 border border-orange-200">
          <h3 className="text-lg font-semibold text-orange-900 mb-3">
            Withdrawal Process
          </h3>
          <div className="space-y-2 text-sm text-orange-800">
            <div className="flex items-start">
              <div className="w-2 h-2 bg-orange-600 rounded-full mt-2 mr-3 flex-shrink-0"></div>
              <p>Select the cryptocurrency asset you want to convert to fiat</p>
            </div>
            <div className="flex items-start">
              <div className="w-2 h-2 bg-orange-600 rounded-full mt-2 mr-3 flex-shrink-0"></div>
              <p>Provide your bank account details for the fiat deposit</p>
            </div>
            <div className="flex items-start">
              <div className="w-2 h-2 bg-orange-600 rounded-full mt-2 mr-3 flex-shrink-0"></div>
              <p>Specify the withdrawal amount and payment purpose</p>
            </div>
            <div className="flex items-start">
              <div className="w-2 h-2 bg-orange-600 rounded-full mt-2 mr-3 flex-shrink-0"></div>
              <p>
                Your crypto will be converted and deposited to your bank account
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
