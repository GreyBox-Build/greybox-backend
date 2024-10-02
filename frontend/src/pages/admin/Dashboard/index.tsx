import React from "react";
import Sidebar from "../../../components/admin/Sidebar";
import Header from "../../../components/admin/Header";
import Report from "../../../components/admin/Report";
import MonthlyOrders from "../../../components/admin/MonthlyOrders";
import OrdersTable from "../../../components/admin/OrdersTable";

const AdminDashboard = () => {
  return (
    <main className="w-full flex relative bg-[#f2f2f2] min-h-screen">
      <aside className="w-[266px] py-4 px-5 h-screen hidden md:block fixed left-0 top-0 border-r-2">
        <Sidebar />
      </aside>
      <main className="md:w-[calc(100%-266px)] w-full md:left-[266px] relative py-4 pl-4 pr-4 md:pr-14">
        <Header />
        <Report />
        <MonthlyOrders />
        <OrdersTable />
      </main>
    </main>
  );
};

export default AdminDashboard;
