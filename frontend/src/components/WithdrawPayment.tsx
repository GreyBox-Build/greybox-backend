import { useState } from "react";
import { InputLabel, TextInput } from "./inputs/TextInput";
import { FormButton } from "./buttons/FormButton";
import { useForm } from "react-hook-form";
import { useOfframpMutation, useOnrampQuery } from "../appSlices/apiSlice";
import { zodResolver } from "@hookform/resolvers/zod";
import { withdrawPaymentSchema } from "../utils/Validations";
import { enqueueSnackbar } from "notistack";
import { CopyWhite } from "./icons/Icons";
import { useCopyTextToClipboard } from "../utils/Copy";

const WithdrawPayment = () => {
  const { currentData: offrampInfo } = useOnrampQuery({});
  const [transactionHash, setTransactionHash] = useState<string | null>(null);
  const copyText = useCopyTextToClipboard();
  const { control, handleSubmit } = useForm({
    defaultValues: {
      amount: "",
      account_address: "",
      chain: "celo",
    },
    resolver: zodResolver(withdrawPaymentSchema),
  });

  const [offramp, { isLoading }] = useOfframpMutation();

  const handleWithdrawPayment = async (data: any) => {
    try {
      const response = await offramp(data).unwrap();
      console.log(response);
      if (response?.data?.transaction_hash) {
        setTransactionHash(response?.data?.transaction_hash);
      }
    } catch (error: any) {
      console.log(error);
      if (error?.status === 403) {
        enqueueSnackbar("Insufficient fund", { variant: "success" });
      }
      if (error?.status === 400) {
        enqueueSnackbar(error?.data?.message, { variant: "success" });
      }
    }
  };
  return (
    <form
      className="mt-[29px] px-[24px] pb-[80px]"
      onSubmit={handleSubmit(handleWithdrawPayment)}
    >
      <section className="flex flex-col gap-y-[32px]">
        <div>
          <InputLabel text={`Enter amount you want to withdraw in cUSD`} />
          <TextInput
            name="amount"
            control={control}
            placeholder="0"
            type="number"
          />
        </div>
        <div>
          <InputLabel text={`Amount will be transferred from:`} />
          <TextInput
            name="from_address"
            control={control}
            value={offrampInfo?.data?.wallet_address}
            readOnly
          />
        </div>
        <div>
          <InputLabel text={`Enter recipient's wallet address`} />
          <TextInput name="account_address" control={control} type="text" />
        </div>
      </section>

      {transactionHash !== null ? (
        <section>
          <button
            className="w-full bg-grey-2 p-[17px_91px] rounded-[48px] flex items-center justify-center gap-x-[10px] text-white text-[0.875rem] font-[700] m-[0_auto] mt-[41px]"
            onClick={() =>
              copyText(transactionHash, "Transaction hash copied to clipboard")
            }
            type="button"
          >
            Copy Transaction Hash <CopyWhite />
          </button>
        </section>
      ) : (
        <FormButton
          label="Submit Request"
          extraClass="mt-[50px]"
          loading={isLoading}
        />
      )}
    </form>
  );
};

export default WithdrawPayment;
