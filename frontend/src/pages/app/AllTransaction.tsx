import { useNavigate } from "react-router-dom";
import { DateHead, TransactionHistoryCard } from "../../components/Cards";
import { CancelIcon } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import {
  useGetAuthUserQuery,
  useGetTransactionQuery,
} from "../../appSlices/apiSlice";
import { findSubArray, groupByDate } from "../../utils/Helpers";
import moment from "moment";
import { Oval } from "react-loader-spinner";

const AllTransaction = () => {
  const navigate = useNavigate();

  const { currentData: userData } = useGetAuthUserQuery({});

  const { currentData: transactions, isFetching: isFetchingTransactions } =
    useGetTransactionQuery(userData?.data?.personal_details?.crypto_currency, {
      refetchOnMountOrArgChange: true,
    });

  const transactionArray = transactions?.data;

  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative ">
            <span className="absolute left-[24px]" onClick={() => navigate(-1)}>
              <CancelIcon />
            </span>
            <h2 className=" text-black text-[1.5rem] font-[600]">
              Transaction History
            </h2>
          </div>

          {!isFetchingTransactions &&
            groupByDate(transactionArray)[0]?.date?.map(
              (date: string, index) => {
                return (
                  <div key={index}>
                    <DateHead date={date} />

                    {findSubArray(
                      groupByDate(transactionArray)[0]?.transactions,
                      date
                    )?.map((details: any, index: number) => {
                      return (
                        <TransactionHistoryCard
                          key={index}
                          label={details?.transaction_sub_type}
                          status={details.status}
                          channel={details.description}
                          time={moment(details?.CreatedAt).format("h:mm A")}
                          amount={`${details.amount}${details?.asset}`}
                          index={index}
                          length={transactionArray?.length}
                          onClick={() => {}}
                        />
                      );
                    })}
                  </div>
                );
              }
            )}
          {isFetchingTransactions && (
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
          )}
          {transactionArray?.length === 0 &&
            userData?.data?.personal_details?.crypto_currency && (
              <p className="text-[0.875rem] text-center leading-[12px] mt-[24px]">
                No transaction here yet!
              </p>
            )}
        </div>
      }
    />
  );
};

export default AllTransaction;
