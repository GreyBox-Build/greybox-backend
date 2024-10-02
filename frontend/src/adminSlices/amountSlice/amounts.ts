// src/features/orders/ordersSlice.ts
import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface OrdersState {
  totalRevenue: number;
  allOrders: number;
  pendingOrders: number;
  completedOrders: number;
}

const initialState: OrdersState = {
  totalRevenue: 0,
  allOrders: 0,
  pendingOrders: 0,
  completedOrders: 0,
};

const ordersSlice = createSlice({
  name: "orders",
  initialState,
  reducers: {
    setTotalRevenue: (state, action: PayloadAction<number>) => {
      state.totalRevenue = action.payload;
    },
    setAllOrders: (state, action: PayloadAction<number>) => {
      state.allOrders = action.payload;
    },
    setPendingOrders: (state, action: PayloadAction<number>) => {
      state.pendingOrders = action.payload;
    },
    setCompletedOrders: (state, action: PayloadAction<number>) => {
      state.completedOrders = action.payload;
    },
    resetOrdersState: (state) => {
      state.totalRevenue = 0;
      state.allOrders = 0;
      state.pendingOrders = 0;
      state.completedOrders = 0;
    },
  },
});

export const {
  setTotalRevenue,
  setAllOrders,
  setPendingOrders,
  setCompletedOrders,
  resetOrdersState,
} = ordersSlice.actions;

export default ordersSlice.reducer;
