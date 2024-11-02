import {
  Deposit,
  Notification,
  SendQ,
  UserLogo,
  Withdraw,
} from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import { useState } from "react";
import { DateHead, DetailsCard, QuickLink } from "../../components/Cards";
import { useNavigate } from "react-router-dom";
import {
  useGetAuthUserQuery,
  useGetTransactionQuery,
} from "../../appSlices/apiSlice";
import moment from "moment";
import { Oval } from "react-loader-spinner";
import { findSubArray, groupByDate } from "../../utils/Helpers";

const Dashboard = () => {
  const navigate = useNavigate();

  const [tab, setTab] = useState<"deposits" | "withdrawals" | "history">(
    "deposits"
  );

  const { currentData: userData, isFetching } = useGetAuthUserQuery({});

  const { currentData: transactions, isFetching: isFetchingTransactions } =
    useGetTransactionQuery(
      userData?.data?.personal_details?.crypto_currency.toLowerCase(),
      {
        refetchOnMountOrArgChange: true,
      }
    );

  console.log(transactions);

  const depositArray = transactions?.data?.filter(
    (transaction: any) => transaction?.transaction_sub_type === "Deposit"
  );
  const withdrawalArray = transactions?.data?.filter(
    (transaction: any) => transaction?.transaction_sub_type === "Withdrawal"
  );

  const transactionArray = tab === "deposits" ? depositArray : withdrawalArray;

  const personInfo = userData?.data?.personal_details;
  const walletInfo = userData?.data?.wallet_details;

  return (
    <AppLayout
      child={
        <div className="w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-2 px-[25px] pt-[51px] pb-[130px]">
          <section className="w-full flex justify-between items-center">
            <Notification
              current={0}
              onClick={() => navigate("/notifications")}
            />
            <div>
              <p className="text-center text-[0.875rem] text-grey-1 leading-[18px]">
                Welcome,
              </p>
              <p className="text-center text-white text-[1.25rem] font-[600]">
                {isFetching
                  ? null
                  : `${personInfo?.first_name} ${personInfo?.last_name}`}
              </p>
            </div>
            <UserLogo
              firstCharacterOfName={personInfo?.first_name?.substring(0, 1)}
              onClick={() => navigate("/settings")}
            />
          </section>
          <section className="w-full pb-[21px] bg-grey-3 rounded-[40px] mt-[59px]">
            <div className="w-full flex flex-col items-center py-[18px] gap-y-[10px] text-white bg-orange-1 rounded-[40px]">
              <p className="text-[0.875rem] text-center leading-[12px]">
                Total Balance
              </p>
              <p className="text-[2.125rem] font-[700]">
                cUSD
                {walletInfo?.balance !== undefined
                  ? parseFloat(walletInfo?.balance)
                      .toFixed(2)
                      .toString()
                      .replace(/\B(?=(\d{3})+(?!\d))/g, ",")
                  : "0.00"}
              </p>
              <p className="text-[0.625rem] text-grey-1">
                Last Updated {moment(personInfo?.UpdatedAt).fromNow()}
              </p>
            </div>

            <section className="w-full flex justify-center px-[51px] gap-x-[55px] mt-[15px]">
              <QuickLink
                icon={<Deposit />}
                label="Deposit"
                onClick={() => navigate("/deposit-options")}
              />
              <QuickLink
                icon={<SendQ />}
                label="Send"
                onClick={() => navigate("/send-options")}
              />

              <QuickLink
                icon={<Withdraw />}
                label="Withdraw"
                onClick={() => navigate("/withdrawal-options")}
              />
            </section>
          </section>
          <section className="bg-grey-1 rounded-[40px_40px_0px_0px] mt-[24px] h-[400px]">
            <div className=" pt-[34px] mx-[23px] border-b-[2px] border-b-grey-4 flex items-center justify-between">
              <div className="flex items-center gap-x-[10px] md:gap-x-[24px]">
                <button
                  className={`${
                    tab === "deposits"
                      ? "border-b-orange-1 border-b-[2px]  font-[700] text-black-2  "
                      : "font-[400] text-black-3"
                  } text-[0.875rem] mb-[-2px] pr-[5px] leading-[18px] transition-all duration-300 ease-in-out  `}
                  onClick={() => setTab("deposits")}
                >
                  Deposits
                </button>
                <button
                  className={`${
                    tab === "withdrawals"
                      ? "border-b-orange-1 border-b-[2px]  font-[700] text-black-2  "
                      : "font-[400] text-black-3"
                  } text-[0.875rem] mb-[-2px] pr-[5px] leading-[18px] transition-all duration-300 ease-in-out  `}
                  onClick={() => setTab("withdrawals")}
                >
                  Withdrawals
                </button>
              </div>
              <button
                className="text-[0.875rem] text-orange-1 leading-[18px] min-w-fit"
                onClick={() => navigate("/all-transactions")}
              >
                View all
              </button>
            </div>
            <section className="h-[400px] overflow-y-auto scrollbar-thin">
              {!isFetchingTransactions &&
                groupByDate(transactionArray)[0]?.date?.map(
                  (date: string, index) => {
                    return (
                      <div key={index}>
                        <DateHead date={date} />

                        {findSubArray(
                          groupByDate(transactionArray)[0]?.transactions,
                          date
                        )
                          ?.reverse()
                          ?.map((details: any, index: number) => {
                            return (
                              <DetailsCard
                                key={index}
                                label={details?.transaction_sub_type}
                                time={moment(details?.timestamp).format(
                                  "hh:mm A"
                                )}
                                amount={`${details?.amount}${details?.asset}`}
                                index={index}
                                length={date?.length}
                                onClick={() => {}}
                              />
                            );
                          })}
                      </div>
                    );
                  }
                )}

              {isFetchingTransactions ||
                (isFetching && (
                  <div className=" w-full flex justify-center items-center p-[20px_0]">
                    <Oval
                      height={50}
                      width={50}
                      color="#fff"
                      wrapperStyle={{}}
                      wrapperClass=""
                      visible={true}
                      ariaLabel="oval-loading"
                      secondaryColor="#22262B"
                      strokeWidth={2}
                      strokeWidthSecondary={2}
                    />
                  </div>
                ))}
              {transactionArray?.length === 0 &&
                userData?.data?.personal_details?.crypto_currency && (
                  <p className="text-[0.875rem] text-center leading-[12px] mt-[24px]">
                    No transaction here yet!
                  </p>
                )}
            </section>
          </section>
        </div>
      }
    />
  );
};

export default Dashboard;
