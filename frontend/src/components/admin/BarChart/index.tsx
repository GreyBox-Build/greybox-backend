import React from "react";
import { Bar } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";

// Register the necessary components from Chart.js
ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

// Define the types for the monthly orders data
interface MonthlyOrder {
  month: string;
  amount: number;
}

const BarChart: React.FC = () => {
  // Monthly orders data
  const monthlyOrders: MonthlyOrder[] = [
    { month: "Jan", amount: 120 },
    { month: "Feb", amount: 150 },
    { month: "Mar", amount: 90 },
    { month: "Apr", amount: 200 },
    { month: "May", amount: 170 },
    { month: "Jun", amount: 250 },
    { month: "Jul", amount: 300 },
    { month: "Aug", amount: 220 },
    { month: "Sep", amount: 180 },
    { month: "Oct", amount: 210 },
    { month: "Nov", amount: 160 },
    { month: "Dec", amount: 240 },
  ];

  const currentMonth: number = new Date().getMonth();

  // Bar chart data
  const data = {
    labels: monthlyOrders.map((order) => order.month),
    datasets: [
      {
        label: "Monthly Orders",
        data: monthlyOrders.map((order) => order.amount),
        backgroundColor: monthlyOrders.map((_, index) =>
          index === currentMonth ? "#CD5928" : "#FCF2EE"
        ),

        borderWidth: 0,
        borderRadius: 8,
      },
    ],
  };

  // Bar chart options
  const options = {
    responsive: true,
    maintainAspectRatio: false, // To control the height manually
    plugins: {
      legend: {
        display: false, // Hide the legend
      },
      tooltip: {
        callbacks: {
          label: (tooltipItem: any) => `Amount: $${tooltipItem.raw}`,
        },
      },
    },
    scales: {
      x: {
        grid: {
          display: false, // Hide grid lines on the x-axis
        },
        title: {
          display: false,
          text: "Months",
        },
        border: {
          display: false, // Hide the x-axis line
        },
      },
      y: {
        min: 0,
        grid: {
          display: false, // Hide grid lines on the x-axis
        },
        display: false, // Hides the Y-axis
        // Limits the Y-axis value to 100 per scale
      },
    },
  };
  return (
    <>
      <div>
        <p className="font-medium">Monthly Transactions:</p>
        <p className="font-bold text-2xl">
          ${monthlyOrders[currentMonth].amount}
        </p>
      </div>
      <div style={{ height: "130px" }} className="mt-4">
        <Bar data={data} options={options} />
      </div>
    </>
  );
};

export default BarChart;
