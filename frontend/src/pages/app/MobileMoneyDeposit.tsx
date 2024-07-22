import { useNavigate } from "react-router-dom";
import { CancelIcon } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import { InputLabel, TextInput } from "../../components/inputs/TextInput";
import { FormButton } from "../../components/buttons/FormButton";
import { useForm } from "react-hook-form";
import {
  useGetAuthUserQuery,
  useOnrampQuery,
  // useSignUrlMutation,
} from "../../appSlices/apiSlice";
import { zodResolver } from "@hookform/resolvers/zod";
import { depositViaMobileSchema } from "../../utils/Validations";
import { MoonPayBuyWidget } from "@moonpay/moonpay-react";

const MobileMoneyDeposit = () => {
  const navigate = useNavigate();

  const { currentData: onrampInfo } = useOnrampQuery({});
  const { currentData: user } = useGetAuthUserQuery({});
  // const [signUrl] = useSignUrlMutation();
  const { control, handleSubmit } = useForm({
    defaultValues: {
      amount: "0",
    },
    resolver: zodResolver(depositViaMobileSchema),
  });

  const handleDepositViaMobileMoney = (data: any) => {
    const { amount } = data;
    window.open(
      `https://pay.fonbnk.com/?source=${process.env.REACT_APP_SOURCE_PARAM}&asset=${onrampInfo?.data?.asset}&country=${onrampInfo?.data?.country}&provider=mobile_money&amount=${amount}&network=${onrampInfo?.data?.network}&address=${onrampInfo?.data?.wallet_address}&freezeWallet=1&freezeNetwork=1&freezeAmount=1&redirectUrl=${process.env.REACT_APP_FRONTEND_BASE_URL}/dashboard/{status}`,
      "_self"
    );
  };

  // const handleGetSignature = async (url: string): Promise<any> => {
  //   console.log(url);
  //   const signature = await signUrl({
  //     url,
  //   }).unwrap();
  //   console.log(signature?.signedUrl);
  //   return signature?.signedUrl;
  // };

  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          {onrampInfo?.data?.network === "XLM" ? (
            <MoonPayBuyWidget
              variant="overlay"
              baseCurrencyCode={user?.data?.personal_details?.currency?.toLowerCase()}
              baseCurrencyAmount="100"
              defaultCurrencyCode={onrampInfo?.data?.network?.toLowerCase()}
              // onUrlSignatureRequested={handleGetSignature}
              // walletAddress={onrampInfo?.data?.wallet_address}
              onCloseOverlay={() => navigate(-1)}
              visible
            />
          ) : (
            <>
              <div className="flex items-center justify-center relative ">
                <span
                  className="absolute left-[24px]"
                  onClick={() => navigate(-1)}
                >
                  <CancelIcon />
                </span>
                <h2 className=" text-black text-[1.5rem] font-[600]">
                  Deposit Via Mobile Money
                </h2>
              </div>
              <form
                className="mt-[29px] px-[24px] pb-[80px]"
                onSubmit={handleSubmit(handleDepositViaMobileMoney)}
              >
                <section className="flex flex-col gap-y-[32px]">
                  <div>
                    <InputLabel
                      text={`Enter amount you want to receive in cUSD`}
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
                <FormButton
                  label="Submit Request"
                  extraClass="mt-[50px]"
                  onClick={() => {}}
                />
              </form>
            </>
          )}
        </div>
      }
    />
  );
};

export default MobileMoneyDeposit;
