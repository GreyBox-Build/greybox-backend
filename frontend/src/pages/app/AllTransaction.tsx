import { useNavigate } from "react-router-dom";
import { DateHead, TransactionHistoryCard } from "../../components/Cards";
import { CancelIcon } from "../../components/icons/Icons";
import { allTransactionData } from "../../utils/Dummies";
import AppLayout from "./AppLayout";

const AllTransaction = () => {
  const navigate = useNavigate();
  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative ">
            <span className="absolute left-[24px]" onClick={() => navigate(-1)}>
              <CancelIcon />
            </span>{" "}
            <h2 className=" text-black text-[1.5rem] font-[600]">
              Transaction History
            </h2>
          </div>
          {allTransactionData?.map((deposit) => {
            return (
              <div key={deposit.id}>
                <DateHead date={deposit.date} />
                {deposit.details.map((details, index) => {
                  return (
                    <TransactionHistoryCard
                      label={details.label}
                      status={details.status}
                      channel={details.channel}
                      time={details.time}
                      amount={details.amount}
                      index={index}
                      length={deposit.details.length}
                      onClick={() => {}}
                    />
                  );
                })}
              </div>
            );
          })}
        </div>
      }
    />
  );
};

export default AllTransaction;
