import AuthLayout from './AuthLayout'
import { BackArrow, LockOpen, Mail } from '../../components/icons/Icons'
import { TextInput } from '../../components/inputs/TextInput'
import { Link, useNavigate } from 'react-router-dom'
import { FormButton } from '../../components/buttons/FormButton'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { obtainTokenSchema } from '../../utils/Validations'
import { useObtainTokenMutation } from '../../appSlices/apiSlice'
import { useSnackbar } from 'notistack'
import InputTextRem from '../../components/InputTextRem'
import { useState } from 'react'

const Login = () => {
  const navigate = useNavigate()
  const [obtainToken, { isLoading }] = useObtainTokenMutation()
  const { enqueueSnackbar } = useSnackbar()
  const [showPassword, setShowPassword] = useState(false)
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm({
    defaultValues: {
      email: '',
      password: '',
    },
    resolver: zodResolver(obtainTokenSchema),
  })
  const handleObtainToken = async (data: any) => {
    try {
      const response: any = await obtainToken(data).unwrap()
      if (response.status) {
        localStorage.setItem('access_token', response?.data?.access_token)
        navigate('/dashboard')
      }
    } catch (error: any) {
      enqueueSnackbar(
        error?.data?.error ? error?.data?.error : 'Connction failed!',
        { variant: 'success' }
      )
    }
  }
  return (
    <AuthLayout
      child={
        <div className="bg-white/95 backdrop-blur-xl rounded-3xl shadow-2xl p-8 w-full max-w-xl relative z-10 animate-fade-in-up">
          <BackArrow />
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-orange-1 to-orange-600 rounded-2xl mb-4 shadow-lg animate-pulse">
              <img src="/images/greybox_logo.png" alt="logo" />
            </div>
            <h1 className="text-3xl font-bold text-gray-800 mb-2">
              GreyBox Wallet
            </h1>
            <p className="text-gray-600">
              Fill in the details below, to login to your account.
            </p>
          </div>

          <form
            className="mt-[24px]"
            onSubmit={handleSubmit(handleObtainToken)}
          >
            <section className="flex flex-col gap-y-[32px]">
              <div className="space-y-2">
                <label className="block text-sm font-semibold text-gray-700">
                  Email Address
                </label>
                <input
                  type="email"
                  {...control.register('email')}
                  className={`w-full px-4 py-3 bg-gray-50 border-2 rounded-xl focus:bg-white focus:outline-none transition-all duration-300 hover:border-gray-300 ${
                    errors.email
                      ? 'border-orange-1 focus:border-orange-1'
                      : 'border-orange-1/40 focus:border-orange-1/20'
                  }`}
                  placeholder="john@example.com"
                />
                {errors.email && (
                  <p className="text-red-500 text-xs mt-1 animate-pulse">
                    {errors.email.message}
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
                        ? 'border-orange-1 focus:border-orange-1'
                        : 'border-orange-1/90 focus:border-orange-1/90'
                    }`}
                    placeholder="Create a strong password"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500 hover:text-orange-1 transition-colors duration-200 focus:outline-none"
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
                  <p className="text-orange-1 text-xs mt-1 animate-pulse">
                    {errors.password.message}
                  </p>
                )}
                <div className="text-xs text-gray-500 mt-1">
                  Password must be at least 8 characters with uppercase,
                  lowercase, and number
                </div>
              </div>
            </section>
            <Link
              to={'/forgot-password'}
              className="text-[0.875rem] text-black-3 leading-[12px] mt-[12px]"
            >
              Forgot Password?
            </Link>
            {/* Submit Button */}
            <button
              type="submit"
              disabled={isLoading}
              className={`w-full py-4 mt-10 px-6 rounded-xl font-semibold text-white transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-4 focus:ring-orange-500/30 ${
                isLoading
                  ? 'bg-orange-1/80 cursor-not-allowed'
                  : 'bg-gradient-to-r from-orange-1 to-orange-600 hover:from-orange-1 hover:to-orange-700 shadow-lg hover:shadow-xl'
              }`}
            >
              {isLoading ? (
                <div className="flex items-center justify-center space-x-2">
                  <div className="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
                  <span>Login you in...</span>
                </div>
              ) : (
                'Login'
              )}
            </button>

            <p className="text-center text-sm text-gray-600 mt-6">
              Don't have an account?{' '}
              <button
                type="button"
                onClick={() => navigate('/sign-up')}
                className="text-orange-1 hover:text-orange-1 font-semibold transition-colors bg-transparent border-none cursor-pointer"
              >
                Sign Up
              </button>
            </p>
          </form>
        </div>
      }
    />
  )
}

export default Login
