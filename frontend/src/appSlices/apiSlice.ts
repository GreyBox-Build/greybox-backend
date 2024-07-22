// Import the RTK Query methods from the React-specific entry point
import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
// Define our single API slice object
export const apiSlice = createApi({
  reducerPath: "api",
  refetchOnReconnect: true,
  baseQuery: fetchBaseQuery({
    baseUrl: "https://apis.greyboxpay.com/api",
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
  tagTypes: ["Staffs", "Zones", "Areas", "Branches"],
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
    offramp: builder.mutation({
      query: (details) => ({
        url: "/v1/transaction/off-ramp",
        method: "POST",
        body: details,
      }),
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
    }),
    onramp: builder.query({
      query: () => ({
        url: "/v1/transaction/on-ramp",
      }),
    }),
    getTransaction: builder.query({
      query: ({ category, pageSize }) => ({
        url: `/v1/transaction?category=${category}&pageSize=${pageSize}`,
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
  useOnrampQuery,
  useForgetPasswordMutation,
  useResetPasswordMutation,
  useOfframpMutation,
  useGetTransactionQuery,
  useSignUrlMutation,
} = apiSlice;
