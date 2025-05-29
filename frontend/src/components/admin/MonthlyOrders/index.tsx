import React from "react";
import BarChart from "../BarChart";

const MonthlyOrders = () => {
  return (
    <div className="w-full my-4">
      <div className="w-full flex-col lg:flex-row flex gap-3">
        <div className="overflow-x-auto lg:w-3/4 md:2/3  w-full p-3 rounded-2xl bg-white pl-10 overflow-hidden">
          <BarChart />
        </div>
        <div className="bg-orange-1 p-4 lg:w-1/4 md:1/3 w-full flex flex-col gap-3 text-white rounded-2xl">
          <h2 className="text-base font-bold text-center mt-3">
            Lorem ipsum is a dummy{" "}
          </h2>
          <p>
            Lorem ipsum is a dummy text, Lorem ipsum is a dummy text, Lorem
            ipsum is a dummy text,Lorem{" "}
          </p>
          <button className="bg-white mt-4 !text-orange-1 font-medium block w-full rounded-full py-2 ">
            Get Started
          </button>
        </div>
      </div>
    </div>
  );
};

export default MonthlyOrders;
