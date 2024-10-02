import { useEffect, useState } from "react";
import { transactions } from "../Dummy";
import { FaCheck } from "react-icons/fa6";
import { LiaTimesSolid } from "react-icons/lia";
import { useSnackbar } from "notistack";
import { useDispatch, useSelector } from "react-redux"; // Import useSelector
import {
  setTotalRevenue,
  setAllOrders,
  setPendingOrders,
  setCompletedOrders,
  resetOrdersState,
} from "../../../adminSlices/amountSlice/amounts";
import { RootState } from "../../../app/store"; // Import RootState for type safety

type TabType = "All" | "Pending" | "Completed";

const OrdersTable = () => {
  const [activeTab, setActiveTab] = useState<TabType>("All");
  const [newData, setNewData] = useState(transactions);
  const { enqueueSnackbar } = useSnackbar();
  const dispatch = useDispatch();

  const searchQuery = useSelector((state: RootState) => state.search.query); // Get the search query from

  const tabs: { label: string; value: TabType }[] = [
    { label: "All", value: "All" },
    { label: "Pending", value: "Pending" },
    { label: "Completed", value: "Completed" },
  ];

  // Calculate total amounts
  const totalPendingAmount = newData.reduce((total, item) => {
    if (item.status === "Pending") {
      return total + Number(item.amountNaira);
    }
    return total;
  }, 0);

  const totalRevenueAmount = newData.reduce((total, item) => {
    if (item.status === "Confirmed") {
      return total + Number(item.amountNaira);
    }
    return total;
  }, 0);

  const allTransactionsAmount = newData.reduce((total, item) => {
    if (item.status === "Confirmed" || item.status === "Pending") {
      return total + Number(item.amountNaira);
    }
    return total;
  }, 0);

  useEffect(() => {
    dispatch(setTotalRevenue(totalRevenueAmount * (1 / 100)));
    dispatch(setAllOrders(allTransactionsAmount));
    dispatch(setPendingOrders(totalPendingAmount));
    dispatch(setCompletedOrders(totalRevenueAmount));

    return () => {
      dispatch(resetOrdersState());
    };
  }, [
    dispatch,
    totalPendingAmount,
    totalRevenueAmount,
    allTransactionsAmount,
    newData,
  ]);

  // Filter the data based on activeTab and search query
  const filteredTransactions = newData.filter((transaction) => {
    const matchesSearch = transaction.name
      .toLowerCase()
      .includes(searchQuery.toLowerCase()); // Adjust to the field you want to search
    const matchesTab =
      activeTab === "All" ||
      (activeTab === "Pending" && transaction.status === "Pending") ||
      (activeTab === "Completed" &&
        (transaction.status === "Confirmed" ||
          transaction.status === "Cancelled"));
    return matchesSearch && matchesTab;
  });

  const confirmOrder = (accountNo: number | string, orderStatus: string) => {
    const updatedData = newData.map((item) => {
      if (item.accountNo === accountNo) {
        return {
          ...item,
          status: orderStatus,
        };
      }
      return item;
    });

    enqueueSnackbar("Transaction Confirmed!", { variant: "success" });
    setNewData(updatedData);
  };

  const cancelOrder = (accountNo: number | string, orderStatus: string) => {
    const updatedData = newData.map((item) => {
      if (item.accountNo === accountNo) {
        return {
          ...item,
          status: orderStatus,
        };
      }
      return item;
    });
    enqueueSnackbar("Cancelled! Maybe a Scammer", { variant: "error" });
    setNewData(updatedData);
  };

  return (
    <div className="w-full bg-white py-8 px-10 rounded-2xl flex flex-col gap-4">
      <p className="text-2xl font-bold ml-3">Orders</p>

      {/* Tab Buttons */}
      <div className="flex space-x-4 mb-4">
        {tabs.map((tab) => (
          <button
            key={tab.value}
            className={`px-4 py-2 rounded-[8px] ${
              activeTab === tab.value
                ? "bg-orange-1 text-white"
                : "bg-[#FDF9F6]"
            }`}
            onClick={() => setActiveTab(tab.value)}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Transaction Table */}
      <div className="w-full overflow-x-auto">
        <table className="min-w-full border table-fixed border-gray-300 ">
          <thead className="bg-[#F5F5F5]">
            <tr>
              <th className="border p-2">S/N</th>
              <th className="border p-2">Name</th>
              <th className="border p-2">Account No</th>
              <th className="border p-2">Amount (₦)</th>
              <th className="border p-2">Amount ($)</th>
              <th className="border p-2">Date</th>
              <th className="border p-2">Status</th>
              <th className="border p-2">Ref</th>
              {activeTab === "Pending" && (
                <th className="border p-2" colSpan={2}>
                  Action
                </th>
              )}
            </tr>
          </thead>
          <tbody>
            {filteredTransactions.map((transaction) => (
              <tr key={transaction.serialNo}>
                <td className="border-b p-2">{transaction.serialNo}</td>
                <td className="border-b p-2">{transaction.name}</td>
                <td className="border-b p-2">{transaction.accountNo} </td>
                <td className="border-b p-2">
                  {transaction.amountNaira} Naira
                </td>
                <td className="border-b p-2">
                  {transaction.amountDollar} Dollars
                </td>
                <td className="border-b p-2">{transaction.date}</td>

                <td className="border-b p-2">
                  <span
                    className={`bg-[#FDF9F6] ${
                      transaction.status === "Pending"
                        ? " text-orange-1"
                        : "text-gray-700"
                    } shadowbtn px-3 py-1 w-[150px]`}
                  >
                    {transaction.status}
                  </span>
                </td>
                <td className="border-b p-1">{transaction.ref}</td>

                {activeTab === "Pending" && (
                  <td className="border-b p-1" colSpan={2}>
                    {transaction.status === "Pending" && (
                      <div className="flex gap-2">
                        <button
                          onClick={() =>
                            confirmOrder(transaction.accountNo, "Confirmed")
                          }
                          className="bg-orange-1 text-white px-2 py-2 rounded flex items-center gap-1"
                        >
                          Confirm
                          <FaCheck size={20} />
                        </button>
                        <button
                          onClick={() =>
                            cancelOrder(transaction.accountNo, "Cancelled")
                          }
                          className="bg-[#FDF9F6] text-black px-2 py-2 rounded flex items-center gap-1"
                        >
                          Cancel
                          <LiaTimesSolid size={20} />
                        </button>
                      </div>
                    )}
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default OrdersTable;
