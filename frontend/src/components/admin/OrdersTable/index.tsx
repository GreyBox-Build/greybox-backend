import { useEffect, useState } from "react";
import { FaCheck } from "react-icons/fa6";
import { LiaTimesSolid } from "react-icons/lia";
import { useSnackbar } from "notistack";
import { useDispatch, useSelector } from "react-redux";
import {
  setTotalRevenue,
  setAllOrders,
  setPendingOrders,
  setCompletedOrders,
} from "../../../adminSlices/amountSlice/amounts";
import { RootState } from "../../../app/store";
import {
  useAdminOnRampRetrieveWithParamsQuery,
  useAdminVerifyOnRampReqWithIdMutation,
  useAdminGetOffRampWithdrawalReqQuery,
  useAdminVerifyOffRampReqWithIdMutation,
} from "../../../appSlices/apiSlice";

import { setGraphData } from "../../../adminSlices/graphDataSlice";
import { calculateMonthlyTotals } from "../utils";

type TabType = "All" | "pending" | "completed";

const OrdersTable = () => {
  const [activeTab, setActiveTab] = useState<TabType>("All");
  const [anoData, setAnoData] = useState<any[]>([]);
  const [withdrawalData, setWithdrawalData] = useState<any[]>([]);
  const [type, setType] = useState("Deposit");
  const { enqueueSnackbar } = useSnackbar();
  const dispatch = useDispatch();
  const searchQuery = useSelector((state: RootState) => state.search.query);

  // Retrieve transactions data from the backend
  const { data, isError } = useAdminOnRampRetrieveWithParamsQuery({});
  const { data: offrampData } = useAdminGetOffRampWithdrawalReqQuery({});

  // Mutation to verify the transaction
  const [adminVerifyOnRampReqWithId] = useAdminVerifyOnRampReqWithIdMutation();
  const [adminVerifyOffRampReqWithId] =
    useAdminVerifyOffRampReqWithIdMutation();

  // UseEffect to load data

  useEffect(() => {
    if (offrampData?.data) {
      setWithdrawalData(offrampData.data);
    }
  }, [offrampData, withdrawalData]);

  useEffect(() => {
    if (data && data.data) {
      const fetchedData = data.data;
      setAnoData(fetchedData); // set fetched data to anoData

      // Calculate amounts for different statuses
      const totalPendingAmount = fetchedData.reduce(
        (total: number, item: any) => {
          const amount = Number(item.asset_equivalent);
          if (item.status === "pending" && !isNaN(amount)) {
            return total + amount;
          }
          return total;
        },
        0
      );

      const totalRevenueAmount = fetchedData.reduce(
        (total: number, item: any) => {
          const amount = Number(item.asset_equivalent);
          if (item.status === "Approve" && !isNaN(amount)) {
            return total + amount;
          }
          return total;
        },
        0
      );

      const allTransactionsAmount = fetchedData.reduce(
        (total: number, item: any) => {
          const amount = Number(item.asset_equivalent);
          if (
            (item.status === "Approve" || item.status === "pending") &&
            !isNaN(amount)
          ) {
            return total + amount;
          }
          return total;
        },
        0
      );

      // Dispatch the calculated totals
      dispatch(setTotalRevenue(totalRevenueAmount * (1 / 100)));
      dispatch(setAllOrders(allTransactionsAmount));
      dispatch(setPendingOrders(totalPendingAmount));
      dispatch(setCompletedOrders(totalRevenueAmount));
    }
  }, [data, dispatch]);

  // Tab labels with explicit TabType typing
  const tabs: { label: string; value: TabType }[] = [
    { label: "All", value: "All" },
    { label: "pending", value: "pending" },
    { label: "completed", value: "completed" },
  ];

  //====================================== Deposit transactions ======================================
  // Confirm order handler
  const confirmOrder = async (id: string | number) => {
    try {
      const actionPayload = { action: "Approve" };
      const response = await adminVerifyOnRampReqWithId({
        id: id,
        action: actionPayload,
      });

      if ("error" in response) {
        enqueueSnackbar("Error confirming this transaction.", {
          variant: "error",
        });
      } else if ("data" in response) {
        const updatedData = anoData.map((item) => {
          if (item.ID === id) {
            return { ...item, status: "Approved" };
          }
          return item;
        });

        setAnoData([...updatedData]); // Use the spread operator to ensure state update triggers a re-render
        enqueueSnackbar("Transaction Confirmed!", { variant: "info" });
      }
    } catch (error) {
      console.error(error); // Log the error for debugging
      enqueueSnackbar("Error confirming transaction.", { variant: "error" });
    }
  };

  // Cancel order handler
  const cancelOrder = async (id: string | number) => {
    try {
      const actionPayload = { action: "Reject" };
      await adminVerifyOnRampReqWithId({
        id: id,
        action: actionPayload,
      });

      const updatedData = anoData.map((item) => {
        if (item.ID === id) {
          return { ...item, status: "Rejected" };
        }
        return item;
      });

      setAnoData([...updatedData]); // Ensure state update triggers a re-render
      enqueueSnackbar("Transaction Cancelled.", { variant: "error" });
    } catch (error) {
      enqueueSnackbar("Error cancelling transaction.", { variant: "error" });
    }
  };

  //================================= Withdrawal Functionalities =================================

  const confirmWithdrawalOrder = async (
    id: string | number,
    bankRef: string | number
  ) => {
    try {
      const actionPayload = { action: "Verified", bankRef: bankRef };
      const response = await adminVerifyOffRampReqWithId({
        id: id,
        action: actionPayload,
      });

      if ("error" in response) {
        enqueueSnackbar("Error confirming this transaction.", {
          variant: "error",
        });
      } else if ("data" in response) {
        const updatedData = withdrawalData.map((item) => {
          if (item.ID === id) {
            return { ...item, status: "Completed" };
          }
          return item;
        });

        setAnoData([...updatedData]); // Use the spread operator to ensure state update triggers a re-render
        enqueueSnackbar("Transaction Confirmed!", { variant: "info" });
      }
    } catch (error) {
      console.error(error); // Log the error for debugging
      enqueueSnackbar("Error confirming transaction.", { variant: "error" });
    }
  };

  const cancelWithdrawalOrder = async (
    id: string | number,
    bankRef: string | number
  ) => {
    try {
      const actionPayload = { action: "Reject", bankRef: bankRef };
      const response = await adminVerifyOffRampReqWithId({
        id: id,
        action: actionPayload,
      });

      if ("error" in response) {
        enqueueSnackbar("Error cancelling this transaction.", {
          variant: "error",
        });
      } else if ("data" in response) {
        const updatedData = withdrawalData.map((item) => {
          if (item.ID === id) {
            return { ...item, status: "Rejected" };
          }
          return item;
        });

        setAnoData([...updatedData]); // Ensure state update triggers a re-render
        enqueueSnackbar("Transaction Cancelled.", { variant: "error" });
      }
    } catch (error) {
      console.error(error); // Log the error for debugging
      enqueueSnackbar("Error cancelling transaction.", { variant: "error" });
    }
  };

  // Filter transactions based on active tab and search query
  const filteredTransactions = anoData.filter((transaction) => {
    const matchesSearch = transaction?.account_name
      ?.toLowerCase()
      .includes(searchQuery?.toLowerCase());
    const matchesSearchTwo = transaction.ref
      ?.toLowerCase()
      .includes(searchQuery?.toLowerCase());
    const matchesTab =
      activeTab === "All" ||
      (activeTab === "pending" && transaction.status === "pending") ||
      (activeTab === "completed" &&
        (transaction.status === "Approved" ||
          transaction.status === "Rejected"));

    return matchesSearch && matchesTab && matchesSearchTwo;
  });

  const filteredWithTransactions =
    withdrawalData &&
    withdrawalData.filter((transaction) => {
      const matchesSearch = transaction.account_name
        .toLowerCase()
        .includes(searchQuery.toLowerCase());

      const matchesTab =
        activeTab === "All" ||
        (activeTab === "pending" &&
          transaction.status === "Awaiting Payment") ||
        (activeTab === "completed" &&
          (transaction.status === "Completed" ||
            transaction.status === "Rejected"));

      return matchesSearch && matchesTab;
    });

  // Function to get month name from a date string

  const monthlyTotalsDep = calculateMonthlyTotals(anoData);
  const monthlyTotalsWith = calculateMonthlyTotals(withdrawalData);

  useEffect(() => {
    dispatch(setGraphData(monthlyTotalsDep));
  }, [data, anoData, dispatch, monthlyTotalsDep]);

  return (
    <div className="w-full bg-white py-8 px-10 rounded-2xl flex flex-col gap-4">
      <p className="text-2xl font-bold ml-3">Transactions</p>

      {/* Tab buttons */}
      <div className="flex space-x-4 mb-4">
        {tabs.map((tab) => (
          <button
            key={tab.value}
            className={`px-4 py-2 capitalize rounded-[8px] ${
              activeTab === tab.value
                ? "bg-orange-1 text-white"
                : "bg-[#FDF9F6]"
            }`}
            onClick={() => setActiveTab(tab.value)}
          >
            {tab.label}
          </button>
        ))}

        <select
          className="outline-none border-none bg-[#FDF9F6]"
          onChange={(e) => {
            setType(e.target.value);
            console.log(type);
          }}
        >
          <option value="Deposit">Deposit</option>
          <option value="Withdraw">Withdrawal</option>
        </select>
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
              {activeTab === "pending" && (
                <th className="border p-2" colSpan={2}>
                  Action
                </th>
              )}
            </tr>
          </thead>

          {/* table body  */}

          {type === "Deposit" ? (
            <tbody>
              {filteredTransactions.map((transaction, index) => (
                <tr key={transaction.ID}>
                  <td className="border-b p-2">{index + 1}</td>
                  <td className="border-b p-2">{transaction.account_name}</td>
                  <td className="border-b p-2">{transaction.account_number}</td>
                  <td className="border-b p-2">₦{transaction.fiat_amount}</td>
                  <td className="border-b p-2">
                    ${transaction.asset_equivalent}
                  </td>
                  <td className="border-b p-2">
                    {new Date(transaction.CreatedAt).toLocaleDateString()}
                  </td>
                  <td
                    className={`border-b  ${
                      transaction.status === "pending"
                        ? " text-orange-1"
                        : "text-gray-700"
                    }`}
                  >
                    <span
                      className={`bg-[#FDF9F6]  px-2 py-1 inline-block rounded-full shadow-sm `}
                    >
                      {" "}
                      {transaction.status}
                    </span>
                  </td>
                  <td className="border-b p-2">{transaction.ref}</td>
                  {activeTab === "pending" && (
                    <td className="border-b p-2">
                      <div className="flex items-center gap-2">
                        {" "}
                        <button
                          onClick={() => confirmOrder(transaction.ID)}
                          className={`bg-[#FDF9F6] flex items-center gap-2 ${
                            transaction.status === "pending"
                              ? " text-orange-1"
                              : "text-gray-700"
                          } shadowbtn px-3 py-2`}
                        >
                          Confirm
                          <FaCheck size={24} />
                        </button>
                        <button
                          onClick={() => cancelOrder(transaction.ID)}
                          className={`bg-[#FDF9F6] flex items-center gap-2  ${
                            transaction.status !== "pending"
                              ? " text-orange-1"
                              : "text-gray-700"
                          } shadowbtn px-3 py-2`}
                        >
                          Cancel
                          <LiaTimesSolid size={24} />
                        </button>
                      </div>
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          ) : (
            <tbody>
              {filteredWithTransactions?.map((transaction, index) => (
                <tr key={transaction.ID}>
                  <td className="border-b p-2">{index + 1}</td>
                  <td className="border-b p-2">{transaction.account_name}</td>
                  <td className="border-b p-2">{transaction.account_number}</td>
                  <td className="border-b p-2">
                    ₦{transaction.equivalent_fiat}
                  </td>
                  <td className="border-b p-2">${transaction.crypto_amount}</td>
                  <td className="border-b p-2">
                    {new Date(transaction.CreatedAt).toLocaleDateString()}
                  </td>
                  <td
                    className={`border-b  ${
                      transaction.status === "Awaiting Payment"
                        ? " text-orange-1"
                        : "text-gray-700"
                    }`}
                  >
                    <span
                      className={`bg-[#FDF9F6]  px-2 py-1 inline-block rounded-full shadow-sm `}
                    >
                      {" "}
                      {transaction.status}
                    </span>
                  </td>
                  <td className="border-b p-2">{transaction.bank_ref}</td>
                  {activeTab === "pending" &&
                    transaction.status === "Awaiting Payment" && (
                      <td className="border-b p-2">
                        <div className="flex items-center gap-2">
                          {" "}
                          <button
                            onClick={() =>
                              confirmWithdrawalOrder(
                                transaction.ID,
                                transaction.bank_ref
                              )
                            }
                            className={`bg-[#FDF9F6] flex items-center gap-2 ${
                              transaction.status === "Awaiting Payment"
                                ? " text-orange-1"
                                : "text-gray-700"
                            } shadowbtn px-3 py-2`}
                          >
                            Confirm
                            <FaCheck size={24} />
                          </button>
                          <button
                            onClick={() =>
                              cancelWithdrawalOrder(
                                transaction.ID,
                                transaction.bank_ref
                              )
                            }
                            className={`bg-[#FDF9F6] flex items-center gap-2  ${
                              transaction.status !== "Awaiting Payment"
                                ? " text-orange-1"
                                : "text-gray-700"
                            } shadowbtn px-3 py-2`}
                          >
                            Cancel
                            <LiaTimesSolid size={24} />
                          </button>
                        </div>
                      </td>
                    )}
                </tr>
              ))}
            </tbody>
          )}
        </table>
      </div>
    </div>
  );
};

export default OrdersTable;
