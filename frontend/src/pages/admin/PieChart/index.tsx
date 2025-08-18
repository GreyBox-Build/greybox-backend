import React from "react";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import { Pie } from "react-chartjs-2";
import { useGetHurupayStatsQuery } from "../../../appSlices/apiSlice";

ChartJS.register(ArcElement, Tooltip, Legend);

interface HurupayStatsResponse {
  data: {
    created_onramp: number;
    failed_off_ramp: number;
    failed_on_ramp: number;
    pending_off_ramp: number;
    pending_on_ramp: number;
    successful_off_ramp: number;
    successful_on_ramp: number;
    total_requests: number;
  };
  errors: boolean;
  status: string;
}

const PieChart = () => {
  const { data: HurustatsResponse } = useGetHurupayStatsQuery({});

  const generateChartData = (
    statsResponse: HurupayStatsResponse | undefined
  ): any => {
    if (!statsResponse || !statsResponse.data) {
      return {
        labels: [],
        datasets: [
          { data: [], backgroundColor: [], borderColor: [], borderWidth: 1 },
        ],
      };
    }

    const stats = statsResponse.data;

    const labels = [
      "Created Onramp",
      "Failed Offramp",
      "Failed Onramp",
      "Pending Offramp",
      "Pending Onramp",
      "Successful Offramp",
      "Successful Onramp",
    ];

    const data = [
      stats.created_onramp,
      stats.failed_off_ramp,
      stats.failed_on_ramp,
      stats.pending_off_ramp,
      stats.pending_on_ramp,
      stats.successful_off_ramp,
      stats.successful_on_ramp,
    ];

    const backgroundColors = [
      "rgba(54, 162, 235, 0.9)", // Blue
      "rgba(255, 99, 132, 0.8)", // Red
      "rgba(255, 99, 132, 0.5)", // Red
      "rgba(255, 206, 86, 0.8)", // Yellow
      "rgba(255, 206, 86, 0.5)", // Yellow
      "rgba(75, 192, 192, 0.9)", // Teal
      "rgba(75, 192, 192, 0.5)", // Teal
    ];

    const borderColors = [
      "rgba(54, 162, 235, 1)",
      "rgba(255, 99, 132, 1)",
      "rgba(255, 99, 132, 1)",
      "rgba(255, 206, 86, 1)",
      "rgba(255, 206, 86, 1)",
      "rgba(75, 192, 192, 1)",
      "rgba(75, 192, 192, 1)",
    ];

    return {
      labels,
      datasets: [
        {
          label: "Requests",
          data,
          backgroundColor: backgroundColors,
          borderColor: borderColors,
          borderWidth: 1,
        },
      ],
    };
  };

  const chartData = generateChartData(HurustatsResponse);

  return <Pie data={chartData} />;
};

export default PieChart;
