import { Link, useNavigate } from "react-router-dom";
import { CancelIcon, CopyWhite } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import { useGetAuthUserQuery } from "../../appSlices/apiSlice";
import QRCode from "qrcode.react";
import { useCopyTextToClipboard } from "../../utils/Copy";

const ExchangeDeposit = () => {
  const navigate = useNavigate();
  const { currentData: userData } = useGetAuthUserQuery({});
  const copyText = useCopyTextToClipboard();
  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative ">
            <span
              className="absolute left-[24px] top-1"
              onClick={() => navigate(-1)}
            >
              <CancelIcon />
            </span>
            <h2 className=" text-black text-[1.5rem] font-[600] mb-[91px]">
              Deposit Via Exchange
            </h2>
          </div>

          <div className="box m-[auto] relative">
            <QRCode value={userData?.data?.personal_details?.account_address} />
            <img
              src="/images/scanMe.svg"
              alt=""
              className="absolute left-[16px] bottom-[-45px]"
            />
          </div>
          <p className=" text-orange-1 text-[1rem] font-[600] leading-[22px] text-center mt-[88px]">
            Deposit Address
          </p>
          <p className="text-black text-[0.875rem] leading-[22px] text-center mt-[22px]">
            {userData?.data?.personal_details?.account_address}
          </p>

          <button
            className=" bg-grey-2 p-[17px_91px] rounded-[48px] flex items-center gap-x-[10px] text-white text-[0.875rem] font-[700] m-[0_auto] mt-[41px]"
            onClick={() =>
              copyText(
                userData?.data?.personal_details?.account_address,
                "Wallet address copied"
              )
            }
          >
            Copy Address <CopyWhite />
          </button>

          <p className="px-[10px] text-center text-[1rem] leading-[22px] mt-[127px] max-w-[387px] m-[0_auto]">
            Only use exchanges or{" "}
            <Link to={""} className=" text-orange-1">
              assets compatible with networks supported by Greybox
            </Link>{" "}
            to prevent loss of funds
          </p>
        </div>
      }
    />
  );
};

export default ExchangeDeposit;
