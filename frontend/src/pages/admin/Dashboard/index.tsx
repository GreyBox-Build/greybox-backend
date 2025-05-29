import Sidebar from "../../../components/admin/Sidebar";
import Header from "../../../components/admin/Header";

import { Outlet } from "react-router-dom";

import { useSnackbar } from "notistack";

import { useAdminOnRampRetrieveWithParamsQuery } from "../../../appSlices/apiSlice";
import Dashboard from "../../app/Dashboard";

const AdminDashboard = () => {
  const { data, isError, isLoading } = useAdminOnRampRetrieveWithParamsQuery(
    {}
  );

  const { enqueueSnackbar } = useSnackbar();

  if (isError) {
    enqueueSnackbar("You are not allowed as the admin", { variant: "success" });
    return <Dashboard />;
  }

  return (
    <>
      {!isLoading && (
        <main className="w-full flex relative bg-[#f2f2f2] min-h-screen">
          <aside className="w-[266px] py-4 px-5 h-screen hidden md:block fixed left-0 top-0 border-r-2">
            <Sidebar />
          </aside>
          <main className="md:w-[calc(100%-266px)] w-full md:left-[266px] relative py-4 pl-4 pr-4 md:pr-14">
            <Header />
            <Outlet />
          </main>
        </main>
      )}
    </>
  );
};

export default AdminDashboard;
