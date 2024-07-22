import { useNavigate } from "react-router-dom";
import { CancelIcon } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import { InputLabel, TextInput } from "../../components/inputs/TextInput";
import { FormButton } from "../../components/buttons/FormButton";
import { useForm } from "react-hook-form";
import { useGetAuthUserQuery, useOnrampQuery } from "../../appSlices/apiSlice";
import { zodResolver } from "@hookform/resolvers/zod";
import { withdrawViaMobileSchema } from "../../utils/Validations";
import { useState } from "react";
import { Modal } from "@mui/material";
import { HomeButton } from "../../components/buttons/HomeButton";
import WithdrawPayment from "../../components/WithdrawPayment";

const MobileMoneyWithdrawal = () => {
  const navigate = useNavigate();
  const { currentData: userInfo } = useGetAuthUserQuery({});
  const [currentPage, setCurrentPage] = useState<number>(0);
  const [openNote, setOpenNote] = useState(false);

  const { currentData: onrampInfo } = useOnrampQuery({});
  const { control, handleSubmit } = useForm({
    defaultValues: {
      amount: "0",
    },
    resolver: zodResolver(withdrawViaMobileSchema),
  });

  const handleDepositViaMobileMoney = (data: any) => {
    const { amount } = data;
    window.open(
      `https://pay.fonbnk.com/offramp/?source=${process.env.REACT_APP_SOURCE_PARAM}&asset=${onrampInfo?.data?.asset}&offrampCurrency=local&country=${onrampInfo?.data?.country}&provider=mobile_money&amount=${amount}&network=${onrampInfo?.data?.network}&fromAddress=${onrampInfo?.data?.wallet_address}&freezeWallet=1&freezeNetwork=1&freezeAmount=1`,
      "_blank"
    );
    setOpenNote(false);
    setCurrentPage(1);
  };
  const balance =
    userInfo?.data?.wallet_details?.balance === undefined
      ? "0.00"
      : parseFloat(userInfo?.data?.wallet_details?.balance)?.toFixed(2);
  const currency = userInfo?.data?.personal_details?.currency;
  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative ">
            <span className="absolute left-[24px]" onClick={() => navigate(-1)}>
              <CancelIcon />
            </span>
            <h2 className=" text-black text-[1.5rem] font-[600]">Withdraw</h2>
          </div>
          <p className="text-black-3 text-[0.875rem] text-center">{`(Bal cUSD${balance})`}</p>
          {currentPage === 0 ? (
            <form
              className="mt-[29px] px-[24px] pb-[80px]"
              onSubmit={handleSubmit(() => setOpenNote(true))}
            >
              <section className="flex flex-col gap-y-[32px]">
                <div>
                  <InputLabel
                    text={`Enter amount you want to receive in ${currency}`}
                  />
                  <TextInput
                    name="amount"
                    control={control}
                    placeholder="0"
                    type="number"
                    step={1}
                  />
                </div>
              </section>
              <FormButton label="Submit Request" extraClass="mt-[50px]" />
            </form>
          ) : (
            <WithdrawPayment />
          )}

          <Modal
            open={openNote}
            onClose={() => setOpenNote(false)}
            className="flex justify-center"
          >
            <section className="max-w-[350px] h-fit bg-grey-1 p-[18px_20px] mt-[8rem] rounded-[8px]">
              <p>
                You will now be redirected to Fonbnk to complete your
                transaction.
              </p>
              <p>
                Make sure you come back to this page to obtain your transaction
                hash which is required to comfirm your transaction.
              </p>
              <div className="mt-[20px] flex items-center gap-x-[10px]">
                <HomeButton
                  label="Cancel"
                  onClick={() => setOpenNote(false)}
                  extraClass="bg-white text-orange-1 w-[50%]"
                />
                <HomeButton
                  label="Continue"
                  onClick={handleSubmit(handleDepositViaMobileMoney)}
                  extraClass="bg-orange-1 text-white w-[50%]"
                />
              </div>
            </section>
          </Modal>
        </div>
      }
    />
  );
};

export default MobileMoneyWithdrawal;
