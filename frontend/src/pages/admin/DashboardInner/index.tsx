import MonthlyOrders from "../../../components/admin/MonthlyOrders";
import OrdersTable from "../../../components/admin/OrdersTable";
import Report from "../../../components/admin/Report";

const DashboardInner = () => {
  return (
    <>
      <Report />
      <MonthlyOrders />
      <OrdersTable />
    </>
  );
};

export default DashboardInner;
