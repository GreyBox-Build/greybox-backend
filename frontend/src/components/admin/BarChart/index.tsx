// src/components/BarChart.tsx
import React, { useEffect, useState } from "react";
import { Bar } from "react-chartjs-2";
import { useSelector } from "react-redux";
import { RootState } from "../../../app/store"; // Import the RootState type for typing the selector

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
  total: number;
}

interface Data {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    backgroundColor: ("#CD5928" | "#FCF2EE")[];
    borderWidth: number;
    borderRadius: number;
  }[];
}

// Define the structure of the state slice for graphData
interface GraphDataState {
  adminLogin: MonthlyOrder[]; // This should match your MonthlyTotal type
}

const BarChart: React.FC = () => {
  // Monthly orders data
  const monthlyTotals = useSelector(
    (state: RootState) => state.graphData
  ) as GraphDataState;

  const currentMonth: number = new Date().getMonth();

  // Initialize chartData with a default structure
  const [chartData, setChartData] = useState<Data>({
    labels: [],
    datasets: [
      {
        label: "Monthly Orders",
        data: [],
        backgroundColor: [],
        borderWidth: 0,
        borderRadius: 8,
      },
    ],
  });

  // Update chartData when monthlyTotals or currentMonth changes
  useEffect(() => {
    const data = {
      labels: monthlyTotals.adminLogin.map((order) => order.month.slice(0, 3)),
      datasets: [
        {
          label: "Monthly Orders",
          data: monthlyTotals.adminLogin.map((order) => order.total),
          backgroundColor: monthlyTotals.adminLogin.map((_, index) =>
            index === currentMonth ? "#CD5928" : "#FCF2EE"
          ),
          borderWidth: 0,
          borderRadius: 8,
        },
      ],
    };

    setChartData(data);
  }, [monthlyTotals, currentMonth]);

  // Bar chart options
  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
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
          display: false,
        },
        title: {
          display: false,
          text: "Months",
        },
        border: {
          display: false,
        },
      },
      y: {
        min: 0,
        grid: {
          display: false,
        },
        display: false,
      },
    },
  };

  return (
    <>
      <div>
        <p className="font-medium">Monthly Transactions:</p>
        <p className="font-bold text-2xl">
          {monthlyTotals.adminLogin[currentMonth]?.total || 0}
        </p>
      </div>
      <div style={{ height: "130px" }} className="mt-4">
        <Bar data={chartData} options={options} />
      </div>
    </>
  );
};

export default BarChart;
