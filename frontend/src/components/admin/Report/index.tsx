import React from "react";
import Card from "../Card";
import { useSelector } from "react-redux";

import { RootState } from "../../../app/store";

const Report = () => {
  const { totalRevenue, allOrders, pendingOrders, completedOrders } =
    useSelector((state: RootState) => state.orders);

  const data = [
    {
      name: "Total Revenew",
      amount: totalRevenue,
      interest: (totalRevenue * 2.4) / 100, // 2.4% of 100k
    },
    {
      name: "All orders",
      amount: allOrders,
      interest: (allOrders * 2.4) / 100,
    },
    {
      name: "Pending orders",
      amount: pendingOrders,
      interest: (pendingOrders * 2.4) / 100,
    },
    {
      name: "Completed orders",
      amount: completedOrders,
      interest: (completedOrders * 2.4) / 100,
    },
  ];

  return (
    <div
      style={{
        backgroundImage:
          "linear-gradient(rgba(0,0,0,0.5), rgba(0,0,0,0.5)), url('images/dashboardbg.jpg')",
      }}
      className="w-full bg-cover bg-center bg-no-repeat rounded-[16px] px-10 py-8"
    >
      <div>
        <div className="mb-12">
          {" "}
          <p className="text-2xl text-white">Dashboard</p>
          <p className="text-sm text-white">Welcome Admin to lorem ipsum</p>
        </div>
        <div className="gridIT">
          {data.map((item) => (
            <Card
              key={item.name}
              name={item.name}
              amount={item.amount}
              interest={item.interest}
            />
          ))}
        </div>
      </div>
    </div>
  );
};

export default Report;
