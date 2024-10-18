// src/features/orders/ordersSlice.ts
import { createSlice, PayloadAction } from "@reduxjs/toolkit";

// Define the structure for monthly totals
interface MonthlyTotal {
  month: string;
  total: number;
}

// Define the slice state type
interface graphDataState {
  adminLogin: MonthlyTotal[];
}

// Set the initial state
const initialState: graphDataState = {
  adminLogin: [], // Empty array to hold monthly totals
};

// Create the slice
const graphData: any = createSlice({
  name: "adminLogin",
  initialState,
  reducers: {
    // Action to set the adminLogin with the monthly totals array
    setGraphData: (state, action: PayloadAction<MonthlyTotal[]>) => {
      state.adminLogin = action.payload; // Store the monthly totals array in the state
    },
  },
});

// Export the actions and reducer
export const { setGraphData } = graphData.actions;

export default graphData.reducer;
