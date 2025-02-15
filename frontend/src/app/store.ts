// app/store.ts
import { configureStore } from "@reduxjs/toolkit";
import { apiSlice } from "../appSlices/apiSlice";
import ordersReducer from "../adminSlices/amountSlice/amounts";
import searchReducer from "../adminSlices/searchSlice";
import graphDataReducer from "../adminSlices/graphDataSlice";

export const store = configureStore({
  reducer: {
    [apiSlice.reducerPath]: apiSlice.reducer,
    orders: ordersReducer,
    search: searchReducer,
    graphData: graphDataReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(apiSlice.middleware),
});
export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
