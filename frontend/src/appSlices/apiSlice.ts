// Import the RTK Query methods from the React-specific entry point
import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
// Define our single API slice object
// baseUrl: "https://apis.greyboxpay.com/api",
export const apiSlice = createApi({
  reducerPath: "api",
  refetchOnReconnect: true,
  baseQuery: fetchBaseQuery({
    baseUrl: "http://localhost:8080/api",
    prepareHeaders: (headers, { endpoint }) => {
      const token = localStorage.getItem("access_token");

      if (
        token &&
        endpoint !== "createUser" &&
        endpoint !== "obtainToken" &&
        endpoint !== "forgetPassword" &&
        endpoint !== "resetPassword" &&
        endpoint !== "getChains"
      ) {
        headers.set("Authorization", `Bearer ${token}`);
      }
      return headers;
    },
  }),
  tagTypes: ["User"],
  endpoints: (builder) => ({
    createUser: builder.mutation({
      query: (user) => ({
        url: "/v2/user/register",
        method: "POST",
        body: user,
      }),
    }),
    getChains: builder.query({
      query: () => ({
        url: "/v1/chains",
      }),
    }),
    obtainToken: builder.mutation({
      query: (user) => ({
        url: "/v1/user/login",
        method: "POST",
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
        url: "/v2/transaction/off-ramp",
        method: "POST",
        body: details,
      }),
      invalidatesTags: ["User"],
    }),
    signUrl: builder.mutation({
      query: (url) => ({
        url: "/v1/transaction/sign-url",
        method: "POST",
        body: url,
      }),
    }),
    getAuthUser: builder.query({
      query: () => ({
        url: "/v1/auth/user",
      }),
      providesTags: ["User"],
    }),
    onramp: builder.mutation({
      query: (details) => ({
        url: "/v2/transaction/on-ramp",
        method: "POST",
        body: details,
      }),
      invalidatesTags: ["User"],
    }),
    getTransaction: builder.query({
      query: (chain) => ({
        url: `/v1/transaction?chain=${chain}`,
      }),
    }),
    forgetPassword: builder.mutation({
      query: (user) => ({
        url: "/v1/token/forget-password",
        method: "POST",
        body: user,
      }),
    }),
    resetPassword: builder.mutation({
      query: (user) => ({
        url: "/v1/token/reset-password",
        method: "POST",
        body: user,
      }),
    }),
  }),
});

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
} = apiSlice;
