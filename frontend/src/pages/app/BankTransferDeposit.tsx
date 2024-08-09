import { useNavigate } from "react-router-dom";
import { CancelIcon } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import {
  InputInfoLabel,
  InputLabel,
  TextInput,
} from "../../components/inputs/TextInput";
import { FormButton } from "../../components/buttons/FormButton";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { depositViaMobileSchema } from "../../utils/Validations";
import {
  useGetAuthUserQuery,
  useGetEquivalentAmountQuery,
} from "../../appSlices/apiSlice";

const BankTransferDeposit = () => {
  const navigate = useNavigate();
  const { currentData: user } = useGetAuthUserQuery({});

  const { control, handleSubmit, watch } = useForm({
    defaultValues: {
      amount: "0",
    },
    resolver: zodResolver(depositViaMobileSchema),
  });

  const handleDepositViaMobileMoney = (data: any) => {
    const { amount } = data;
  };

  const amount = watch("amount");
  const userData = user?.data?.personal_details;

  const { currentData: equivalent } = useGetEquivalentAmountQuery({
    amount,
    currency: userData?.currency,
    cryptoAsset: userData?.crypto_currency,
    type: "on-ramp",
  });

  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative ">
            <span className="absolute left-[24px]" onClick={() => navigate(-1)}>
              <CancelIcon />
            </span>
            <h2 className=" text-black text-[1.5rem] font-[600]">
              Deposit Via Bank Transfer
            </h2>
          </div>
          <form
            className="mt-[29px] px-[24px] pb-[80px]"
            onSubmit={handleSubmit(handleDepositViaMobileMoney)}
          >
            <section className="flex flex-col gap-y-[32px]">
              <div>
                <InputLabel text={`Enter amount you want to receive in cUSD`} />
                <TextInput
                  name="amount"
                  control={control}
                  placeholder="0"
                  localType="figure"
                />
                <InputInfoLabel title="Buying Rate" value="1cUSD = 1USD" />
              </div>
              <div>
                <InputLabel text={`You will recieve`} />
                <TextInput name="" control={control} readOnly value={0} />
              </div>
            </section>
            <FormButton
              label="Submit Request"
              extraClass="mt-[50px]"
              onClick={() => {}}
            />
          </form>
        </div>
      }
    />
  );
};

export default BankTransferDeposit;
