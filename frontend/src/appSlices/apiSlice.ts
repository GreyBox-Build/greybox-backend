// Import the RTK Query methods from the React-specific entry point
import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
// Define our single API slice object
// baseUrl: "https://apis.greyboxpay.com/api",
// baseUrl: "http://localhost:8080/api",

interface RejectKycPayload {
  id: string | number | undefined // The KYC ID from the URL
  email: string
  rejection_reason: string
}

// Define the shape of the success response (adjust as per your API's actual response)
interface RejectKycSuccess {
  message: string
  // Add any other data your API returns on success, e.g., the updated KYC status
}
export const apiSlice = createApi({
  reducerPath: 'api',
  refetchOnReconnect: true,

  baseQuery: fetchBaseQuery({
    baseUrl:
      'http://greygo-env.eba-g4cj2keb.us-east-1.elasticbeanstalk.com/api',
    prepareHeaders: (headers, { endpoint }) => {
      const token = localStorage.getItem('access_token')

      if (
        token &&
        endpoint !== 'createUser' &&
        endpoint !== 'obtainToken' &&
        endpoint !== 'forgetPassword' &&
        endpoint !== 'resetPassword' &&
        endpoint !== 'getChains'
      ) {
        headers.set('Authorization', `Bearer ${token}`)
      }
      return headers
    },
  }),
  tagTypes: ['User', 'Kyc', 'Virtual'],
  endpoints: (builder) => ({
    createUser: builder.mutation({
      query: (user) => ({
        url: '/v2/user/register',
        method: 'POST',
        body: user,
      }),
    }),
    getChains: builder.query({
      query: () => ({
        url: '/v1/chains',
      }),
    }),
    obtainToken: builder.mutation({
      query: (user) => ({
        url: '/v1/user/login',
        method: 'POST',
        body: user,
      }),
    }),
    getEquivalentAmount: builder.query({
      query: ({ amount, currency, cryptoAsset, type }) => ({
        url: `/v2/transaction/equivalent-amount?amount=${amount}&currency=${currency}&cryptoAsset=${cryptoAsset}&type=${type}`,
      }),
    }),
    getBankDetails: builder.query({
      query: (countryCode) => ({
        url: `/v2/transaction/destination-bank?countryCode=${countryCode}`,
      }),
    }),
    getTransactionReference: builder.query({
      query: () => ({
        url: `/v2/transaction/reference`,
      }),
    }),
    getExchangeRate: builder.query({
      query: ({ fiat, asset }) => ({
        url: `/v1/exchange-rate?fiat_currency=${fiat}&asset=${asset}`,
      }),
    }),
    offramp: builder.mutation({
      query: (details) => ({
        url: '/v2/transaction/off-ramp',
        method: 'POST',
        body: details,
      }),
      invalidatesTags: ['User'],
    }),
    signUrl: builder.mutation({
      query: (url) => ({
        url: '/v1/transaction/sign-url',
        method: 'POST',
        body: url,
      }),
    }),
    getAuthUser: builder.query({
      query: () => ({
        url: '/v1/auth/user',
      }),
      providesTags: ['User'],
    }),
    onramp: builder.mutation({
      query: (details) => ({
        url: '/v2/transaction/on-ramp',
        method: 'POST',
        body: details,
      }),
      invalidatesTags: ['User'],
    }),
    getTransaction: builder.query({
      query: (chain) => ({
        url: `/v1/transaction?chain=${chain}`,
      }),
    }),
    forgetPassword: builder.mutation({
      query: (user) => ({
        url: '/v1/user/forget-password',
        method: 'POST',
        body: user,
      }),
    }),
    resetPassword: builder.mutation({
      query: (user) => ({
        url: '/v1/token/reset-password',
        method: 'POST',
        body: user,
      }),
    }),

    getMobileEquivAmount: builder.query({
      query: ({ amount, currency, cryptoAsset, type }) => ({
        url: `/v2/transaction/on-ramp/mobile/equivalent-amount?amount=${amount}&currency=${currency}&type=${type}&cryptoAsset=${cryptoAsset}`,
      }),
    }),

    getNetworks: builder.query({
      query: () => ({
        url: '/v1/networks',
      }),
    }),
    onrampMobile: builder.mutation({
      query: (requestData) => ({
        url: '/v2/transaction/on-ramp/mobile',
        method: 'POST',
        body: requestData,
      }),
      invalidatesTags: ['User'],
    }),

    kyc: builder.mutation({
      query: (data) => ({
        url: '/v2/kyc',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['Kyc'],
    }),

    createVirtualAccount: builder.mutation({
      query: (data) => ({
        url: '/v2/user/account',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['Virtual'],
    }),

    getKyc: builder.query({
      query: () => ({
        url: '/v2/kyc',
      }),
      providesTags: ['Kyc'],
    }),
    approveKyc: builder.mutation({
      query: (id) => ({
        url: `/v2/kyc/${id}/approve`,
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({}),
      }),
      invalidatesTags: ['Kyc'],
    }),

    offrampMobile: builder.mutation({
      query: (sentData) => ({
        url: '/v2/transaction/off-ramp/mobile',
        method: 'POST',
        body: sentData,
      }),
      invalidatesTags: ['User'],
    }),

    adminGetOffRampWithdrawalReq: builder.query({
      query: () => ({
        url: `/v1/requests/off-ramp?chain&hash&address&account_number&status`,
      }),
    }),
    getTransHistory: builder.query({
      query: () => ({
        url: `/v1/transaction?chain=CELO`,
      }),
    }),
    rejectKyc: builder.mutation<RejectKycSuccess, RejectKycPayload>({
      query: ({ id, email, rejection_reason }) => ({
        url: `kyc/${id}/reject`, // Endpoint path relative to baseUrl
        method: 'POST',
        body: { email, rejection_reason },
      }),
      invalidatesTags: [{ type: 'Kyc' }],
    }),
    adminOffRampRetrieveData: builder.query({
      query: (id) => ({
        url: `v1/requests/off-ramp/:${id}`,
      }),
    }),
    adminOnRampRetrieveDataReq: builder.query({
      query: (id) => ({
        url: `v1/requests/on-ramp/:${id}`,
      }),
    }),
    adminOnRampRetrieveWithParams: builder.query({
      query: () => ({
        url: `v1/requests/on-ramp?$ref&$cu`,
      }),
    }),

    getHurupayStats: builder.query({
      query: () => ({
        url: `v1/requests/hurupay-requests/stats`,
      }),
    }),
    getMyKyc: builder.query({
      query: () => ({
        url: `/v2/user/account`,
      }),
    }),
    // adminOnRampRetrieveWithParams: builder.query({
    //   query: ({ ref, cu }) => ({
    //     url: `v1/requests/on-ramp?${ref}&${cu}`,
    //   }),
    // }),
    getVirtualAccountStatus: builder.query({
      query: () => ({
        url: 'v2/user/virtual-account-status',
      }),
      providesTags: ['Virtual'],
    }),
    getUserAccountdetails: builder.query({
      query: () => ({
        url: 'v2/user/account',
      }),
      providesTags: ['Virtual'],
    }),
    adminVerifyOnRampReqWithId: builder.mutation({
      query: ({ id, action }) => ({
        url: `v1/requests/on-ramp/${id}/verify`, // Make sure the `id` is used directly in the URL
        method: 'POST',
        body: action, // Send the action payload
      }),
    }),

    adminVerifyOffRampReqWithId: builder.mutation({
      query: ({ id, actionBankRef }) => ({
        url: `v1/requests/off-ramp/${id}/verify`,
        method: 'POST',
        body: actionBankRef,
      }),
    }),
    onrampBorderlessPayment: builder.mutation({
      query: ({ countryChange }) => ({
        url: `v1/payments/borderless-onramp/variable`,
        method: 'POST',
        body: countryChange,
      }),
    }),
    offrampBorderlessPayment: builder.mutation({
      query: ({ countryChange }) => ({
        url: `v1/payments/borderless-oframp/variable`,
        method: 'POST',
        body: countryChange,
      }),
    }),
  }),
})

export const {
  useCreateUserMutation,
  useGetChainsQuery,
  useObtainTokenMutation,
  useGetAuthUserQuery,
  useGetBankDetailsQuery,
  useGetEquivalentAmountQuery,
  useGetTransactionReferenceQuery,
  useOnrampMutation,
  useForgetPasswordMutation,
  useResetPasswordMutation,
  useOfframpMutation,
  useGetTransactionQuery,
  useGetExchangeRateQuery,
  useSignUrlMutation,
  useGetMobileEquivAmountQuery,
  useGetNetworksQuery,
  useOnrampMobileMutation,
  useOfframpMobileMutation,
  useGetTransHistoryQuery,
  useAdminGetOffRampWithdrawalReqQuery,
  useAdminOnRampRetrieveDataReqQuery,
  useAdminOnRampRetrieveWithParamsQuery,
  useAdminVerifyOnRampReqWithIdMutation,
  useAdminVerifyOffRampReqWithIdMutation,
  useGetHurupayStatsQuery,
  useOnrampBorderlessPaymentMutation,
  useKycMutation,
  useGetKycQuery,
  useCreateVirtualAccountMutation,
  useApproveKycMutation,
  useGetMyKycQuery,
  useRejectKycMutation,
  useGetVirtualAccountStatusQuery,
  useGetUserAccountdetailsQuery,
  useOfframpBorderlessPaymentMutation,
} = apiSlice
