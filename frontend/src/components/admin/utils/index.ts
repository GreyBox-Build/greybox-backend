// export const getMonthName = (dateString: string) => {
//   const date = new Date(dateString);
//   return date.toLocaleString("default", { month: "long" });
// };

// // Function to filter and sum transactions by month, excluding NaN values
// export const calculateMonthlyTotals = (transactions: any[]) => {
//   return transactions.reduce((acc, transaction) => {
//     const month = getMonthName(transaction.CreatedAt); // Get month name from CreatedAt date
//     const amount = parseFloat(transaction.asset_equivalent); // Convert fiat_amount to a number

//     // Check if amount is a valid number and not NaN
//     if (!isNaN(amount)) {
//       // If the month already exists in the accumulator, add the current amount
//       if (acc[month]) {
//         acc[month] += amount;
//       } else {
//         // Otherwise, initialize it with the current amount
//         acc[month] = amount;
//       }
//     }

//     return acc;
//   }, {}); // Initial value is an empty object
// };
// Function to get month name from a date string
const getMonthName = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleString("default", { month: "long" });
};

// Initialize all months with 0 total
const initializeMonthlyTotals = (): Record<string, number> => {
  return {
    January: 0,
    February: 0,
    March: 0,
    April: 0,
    May: 0,
    June: 0,
    July: 0,
    August: 0,
    September: 0,
    October: 0,
    November: 0,
    December: 0,
  };
};

// Function to filter and sum transactions by month, including empty months and excluding NaN values
export const calculateMonthlyTotals = (transactions: any[]) => {
  const monthlyTotals = initializeMonthlyTotals(); // Now it's scoped inside this function

  // Aggregate transactions into the monthly totals object
  transactions.forEach((transaction) => {
    const month = getMonthName(transaction.CreatedAt); // Get month name from CreatedAt date
    const amount = parseFloat(transaction.asset_equivalent); // Convert fiat_amount to a number

    // Check if amount is a valid number and not NaN
    if (!isNaN(amount)) {
      monthlyTotals[month] += amount; // Update the total for the specific month
    }
  });

  // Convert the monthly totals object into an array of objects
  return Object.entries(monthlyTotals).map(([month, total]) => ({
    month,
    total,
  }));
};
