import { useNavigate } from "react-router-dom";
import { CancelIcon, PlusWhiteIcon } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import { PaymentDetailsCard } from "../../components/Cards";

const PaymentDetails = () => {
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
              Payment Details
            </h2>
          </div>
          <section className="px-[24px] mt-[59px] flex flex-col gap-y-[24px]">
            <PaymentDetailsCard
              text="GTB- 0903849729 "
              subText="Lorem Ipsum"
              deleteFunc={() => {}}
            />
            <PaymentDetailsCard
              text="GTB- 0903849729 "
              subText="Lorem Ipsum"
              deleteFunc={() => {}}
            />
            <button className="w-full flex flex-col items-center gap-y-[9px] p-[21px_16px] rounded-[4px] border-[1px] border-[#999999] bg-grey-2 text-white ">
              Add Another Bank
              <PlusWhiteIcon />
            </button>
          </section>
        </div>
      }
    />
  );
};

export default PaymentDetails;
