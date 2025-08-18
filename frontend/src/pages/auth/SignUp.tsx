import { useState } from 'react'
import { useSnackbar } from 'notistack'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from 'react-router-dom'
import { currencyDataT, BorderlessCountry } from '../../utils/Dummies'
import {
  useCreateUserMutation,
  useGetChainsQuery,
} from '../../appSlices/apiSlice'
import { createUserSchema } from 'src/utils/Validations'
import SelectBoxT from '../../components/modals/SelectBoxT'

const SignupForm = () => {
  const navigate = useNavigate()
  const { enqueueSnackbar } = useSnackbar()

  // State for modal controls
  const [openCurrency, setOpenCurrency] = useState(false)
  const [openCountry, setOpenCountry] = useState(false)
  const [openChain, setOpenChain] = useState(false)
  const [showPassword, setShowPassword] = useState(false)

  // RTK Query hooks
  const [createUser, { isLoading }] = useCreateUserMutation()
  const { data: chains } = useGetChainsQuery({})

  // React Hook Form
  const {
    control,
    handleSubmit,
    clearErrors,
    setValue,
    getValues,
    formState: { errors },
  } = useForm({
    defaultValues: {
      first_name: '',
      last_name: '',
      email: '',
      password: '',
      currency: '',
      country: '',
      chain: '',
      country_code: '',
    },
    resolver: zodResolver(createUserSchema),
  })

  const handleCreateUser = async (data: any) => {
    const updatedData = { ...data, country_code: getValues('country_code') }
    try {
      const response: any = await createUser(updatedData).unwrap()
      enqueueSnackbar(response?.message, { variant: 'success' })
      setTimeout(() => {
        navigate('/sign-in')
      }, 5000)
    } catch (error: any) {
      enqueueSnackbar(
        error?.data?.error ? error?.data?.error : 'Connection failed!',
        { variant: 'error' }
      )
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-orange-50 to-orange-100 flex items-center justify-center p-4 relative overflow-hidden">
      {/* Floating particles */}
      <div className="absolute inset-0 pointer-events-none">
        {[...Array(12)].map((_, i) => (
          <div
            key={i}
            className="absolute w-1 h-1 bg-orange-400 rounded-full opacity-30 animate-pulse"
            style={{
              left: `${Math.random() * 100}%`,
              top: `${Math.random() * 100}%`,
              animationDelay: `${i * 0.5}s`,
              animationDuration: `${3 + Math.random() * 4}s`,
            }}
          />
        ))}
      </div>

      {/* Floating geometric shapes */}
      <div
        className="absolute top-20 left-20 w-20 h-20 border border-orange-500/20 rounded-full animate-spin"
        style={{ animationDuration: '20s' }}
      />
      <div
        className="absolute bottom-20 right-20 w-16 h-16 border border-orange-400/30 rounded-lg animate-bounce"
        style={{ animationDuration: '3s' }}
      />
      <div className="absolute top-1/2 left-10 w-12 h-12 bg-orange-500/10 rounded-full animate-pulse" />

      <div className="bg-white/95 backdrop-blur-xl rounded-3xl shadow-2xl p-8 w-full max-w-xl relative z-10 animate-fade-in-up">
        {/* Logo Section */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-orange-500 to-orange-600 rounded-2xl mb-4 shadow-lg animate-pulse">
            <img src="/images/greybox_logo.png" alt="logo" />
          </div>
          <h1 className="text-3xl font-bold text-gray-800 mb-2">
            GreyBox Wallet
          </h1>
          <p className="text-gray-600">Create your secure wallet account</p>
        </div>

        {/* Form */}
        <form className="space-y-6" onSubmit={handleSubmit(handleCreateUser)}>
          {/* Name Row */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-semibold text-gray-700">
                First Name
              </label>
              <input
                type="text"
                {...control.register('first_name')}
                className={`w-full px-4 py-3 bg-gray-50 border-2 rounded-xl focus:bg-white focus:outline-none transition-all duration-300 hover:border-gray-300 ${
                  errors.first_name
                    ? 'border-red-500 focus:border-red-500'
                    : 'border-gray-200 focus:border-orange-500'
                }`}
                placeholder="John"
              />
              {errors.first_name && (
                <p className="text-red-500 text-xs mt-1 animate-pulse">
                  {errors.first_name.message}
                </p>
              )}
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-semibold text-gray-700">
                Last Name
              </label>
              <input
                type="text"
                {...control.register('last_name')}
                className={`w-full px-4 py-3 bg-gray-50 border-2 rounded-xl focus:bg-white focus:outline-none transition-all duration-300 hover:border-gray-300 ${
                  errors.last_name
                    ? 'border-red-500 focus:border-red-500'
                    : 'border-gray-200 focus:border-orange-500'
                }`}
                placeholder="Doe"
              />
              {errors.last_name && (
                <p className="text-red-500 text-xs mt-1 animate-pulse">
                  {errors.last_name.message}
                </p>
              )}
            </div>
          </div>

          {/* Email */}
          <div className="space-y-2">
            <label className="block text-sm font-semibold text-gray-700">
              Email Address
            </label>
            <input
              type="email"
              {...control.register('email')}
              className={`w-full px-4 py-3 bg-gray-50 border-2 rounded-xl focus:bg-white focus:outline-none transition-all duration-300 hover:border-gray-300 ${
                errors.email
                  ? 'border-red-500 focus:border-red-500'
                  : 'border-gray-200 focus:border-orange-500'
              }`}
              placeholder="john@example.com"
            />
            {errors.email && (
              <p className="text-red-500 text-xs mt-1 animate-pulse">
                {errors.email.message}
              </p>
            )}
          </div>

          {/* Country and Currency Row */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-semibold text-gray-700">
                Country
              </label>
              <input
                type="text"
                {...control.register('country')}
                readOnly
                onClick={() => setOpenCountry(true)}
                className={`w-full px-4 py-3 bg-gray-50 border-2 rounded-xl focus:bg-white focus:outline-none transition-all duration-300 hover:border-gray-300 ${
                  errors.country
                    ? 'border-red-500 focus:border-red-500'
                    : 'border-gray-200 focus:border-orange-500'
                }`}
                placeholder="Select Country"
              />
              {errors.country && (
                <p className="text-red-500 text-xs mt-1 animate-pulse">
                  {errors.country.message}
                </p>
              )}
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-semibold text-gray-700">
                Currency
              </label>
              <input
                type="text"
                {...control.register('currency')}
                readOnly
                onClick={() => setOpenCurrency(true)}
                className={`w-full px-4 py-3 bg-gray-50 border-2 rounded-xl focus:bg-white focus:outline-none transition-all duration-300 hover:border-gray-300 ${
                  errors.currency
                    ? 'border-red-500 focus:border-red-500'
                    : 'border-gray-200 focus:border-orange-500'
                }`}
                placeholder="Select Currency"
              />
              {errors.currency && (
                <p className="text-red-500 text-xs mt-1 animate-pulse">
                  {errors.currency.message}
                </p>
              )}
            </div>
          </div>

          {/* Chain */}
          <div className="space-y-2">
            <label className="block text-sm font-semibold text-gray-700">
              Preferred Chain
            </label>
            <input
              type="text"
              {...control.register('chain')}
              readOnly
              onClick={() => setOpenChain(true)}
              className={`w-full px-4 py-3 bg-gray-50 border-2 rounded-xl focus:bg-white focus:outline-none transition-all duration-300 hover:border-gray-300 ${
                errors.chain
                  ? 'border-red-500 focus:border-red-500'
                  : 'border-gray-200 focus:border-orange-500'
              }`}
              placeholder="Select Chain"
            />
            {errors.chain && (
              <p className="text-red-500 text-xs mt-1 animate-pulse">
                {errors.chain.message}
              </p>
            )}
          </div>

          {/* Password */}
          <div className="space-y-2">
            <label className="block text-sm font-semibold text-gray-700">
              Password
            </label>
            <div className="relative">
              <input
                type={showPassword ? 'text' : 'password'}
                {...control.register('password')}
                className={`w-full px-4 py-3 pr-12 bg-gray-50 border-2 rounded-xl focus:bg-white focus:outline-none transition-all duration-300 hover:border-gray-300 ${
                  errors.password
                    ? 'border-red-500 focus:border-red-500'
                    : 'border-gray-200 focus:border-orange-500'
                }`}
                placeholder="Create a strong password"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500 hover:text-orange-500 transition-colors duration-200 focus:outline-none"
              >
                {showPassword ? (
                  <svg
                    className="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"
                    />
                  </svg>
                ) : (
                  <svg
                    className="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                    />
                  </svg>
                )}
              </button>
            </div>
            {errors.password && (
              <p className="text-red-500 text-xs mt-1 animate-pulse">
                {errors.password.message}
              </p>
            )}
            <div className="text-xs text-gray-500 mt-1">
              Password must be at least 8 characters with uppercase, lowercase,
              and number
            </div>
          </div>

          {/* Submit Button */}
          <button
            type="submit"
            disabled={isLoading}
            className={`w-full py-4 px-6 rounded-xl font-semibold text-white transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-4 focus:ring-orange-500/30 ${
              isLoading
                ? 'bg-gray-400 cursor-not-allowed'
                : 'bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700 shadow-lg hover:shadow-xl'
            }`}
          >
            {isLoading ? (
              <div className="flex items-center justify-center space-x-2">
                <div className="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
                <span>Creating Account...</span>
              </div>
            ) : (
              'Create Account'
            )}
          </button>
        </form>

        {/* Footer */}
        <p className="text-center text-sm text-gray-600 mt-6">
          Already have an account?{' '}
          <button
            type="button"
            onClick={() => navigate('/sign-in')}
            className="text-orange-500 hover:text-orange-600 font-semibold transition-colors bg-transparent border-none cursor-pointer"
          >
            Sign In
          </button>
        </p>
      </div>

      {/* SelectBoxT Modals */}
      <SelectBoxT
        state={openCurrency}
        title="Select Currency"
        placeholder="Search Currency"
        childList={currencyDataT}
        onPickChild={(list) => {
          setValue('currency', list?.code!)
          clearErrors('currency')
        }}
        onClose={() => setOpenCurrency(false)}
      />
      <SelectBoxT
        state={openCountry}
        title="Select Country"
        placeholder="Search Country"
        childList={BorderlessCountry}
        onPickChild={(list) => {
          setValue('country', list?.name)
          setValue('country_code', list?.code!)
          clearErrors('country')
        }}
        onClose={() => setOpenCountry(false)}
      />
      <SelectBoxT
        state={openChain}
        title="Select Chain"
        placeholder="Search Chain"
        type="chain"
        childList={chains?.data === undefined ? [] : chains?.data}
        onPickChild={(list: any) => {
          setValue('chain', list?.chain)
          clearErrors('chain')
        }}
        onClose={() => setOpenChain(false)}
      />

      <style>{`
        @keyframes fade-in-up {
          from {
            opacity: 0;
            transform: translateY(30px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }

        .animate-fade-in-up {
          animation: fade-in-up 0.8s ease-out;
        }
      `}</style>
    </div>
  )
}

export default SignupForm
