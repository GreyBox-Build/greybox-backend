import { useNavigate } from "react-router-dom";
import { CancelIcon, CopyBlack } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import {
  InputInfoLabel,
  InputLabel,
  TextInput,
} from "../../components/inputs/TextInput";
import { FormButton } from "../../components/buttons/FormButton";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { depositViaBankTransferSchema } from "../../utils/Validations";
import {
  useGetAuthUserQuery,
  useGetBankDetailsQuery,
  useGetEquivalentAmountQuery,
  useGetExchangeRateQuery,
  useGetTransactionReferenceQuery,
  useOnrampMutation,
} from "../../appSlices/apiSlice";
import { returnAsset } from "../../utils/Helpers";
import { useCopyTextToClipboard } from "../../utils/Copy";
import { useSnackbar } from "notistack";

const BankTransferDeposit = () => {
  const navigate = useNavigate();
  const { currentData: user } = useGetAuthUserQuery({});

  const userData = user?.data?.personal_details;
  const copyText = useCopyTextToClipboard();

  const { enqueueSnackbar } = useSnackbar();

  const { control, handleSubmit, watch } = useForm({
    defaultValues: {
      amount: "0",
    },
    resolver: zodResolver(
      depositViaBankTransferSchema({ currency: userData?.currency })
    ),
  });

  const amount = watch("amount");

  const { currentData: bank } = useGetBankDetailsQuery(userData?.country_code);

  const bankDetails = bank?.data;

  const { currentData: reference } = useGetTransactionReferenceQuery({});

  const { currentData: rate } = useGetExchangeRateQuery({
    fiat: userData?.currency,
    asset: returnAsset(userData?.crypto_currency)?.toLocaleUpperCase(),
  });

  const [onramp, { isLoading }] = useOnrampMutation();

  const {
    currentData: equivalent,
    isError: isEquivalentError,
    error: equivalentError,
  }: any = useGetEquivalentAmountQuery({
    amount: amount?.replace(/,/g, ""),
    currency: userData?.currency,
    cryptoAsset: returnAsset(userData?.crypto_currency)?.toLocaleUpperCase(),
    type: "on-ramp",
  });

  const DetailsCard = ({ text, value }: { text: string; value: string }) => (
    <div
      onClick={() => copyText(value, `Copied ${text}`)}
      className="rounded-[8px]  flex items-center justify-between bg-grey-1 p-[8px_22px] text-black-2 text-[0.875rem] leading-[18px] border-[#99999961] border-[1px] gap-x-[5px] shadow-shadow-1"
    >
      <div className="flex flex-col gap-y-[5px]">
        <p>{text}</p>
        <p className="font-[700]">{value}</p>
      </div>
      <CopyBlack />
    </div>
  );
  const handleDepositViaBankTransfer = async (data: any) => {
    const { amount } = data;

    const details = {
      amount: amount?.replace(/,/g, ""),
      asset: userData?.crypto_currency,
      countryCode: userData?.country_code,
      ref: reference?.data?.reference,
      bankName: bankDetails?.BankName,
      accountNumber: bankDetails?.AccountNumber,
      accountName: bankDetails?.AccountName,
      currency: userData?.currency,
      assetAmount: equivalent?.data?.amount,
    };

    try {
      const response = await onramp(details).unwrap();
      enqueueSnackbar(response?.status, { variant: "success" });
      setTimeout(() => {
        navigate("/dashboard");
      }, 3000);
    } catch (error: any) {
      console.log(error);
      enqueueSnackbar(error?.data?.error, { variant: "success" });
    }
  };

  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative ">
            <span
              className="absolute left-[15px] md:left-[24px]"
              onClick={() => navigate(-1)}
            >
              <CancelIcon />
            </span>
            <h2 className=" text-black text-[1.5rem] font-[600]">
              Deposit Via Bank Transfer
            </h2>
          </div>
          <form
            className="mt-[29px] px-[24px] pb-[80px]"
            onSubmit={handleSubmit(handleDepositViaBankTransfer)}
          >
            <section className="flex flex-col gap-y-[32px]">
              <div>
                <InputLabel
                  text={`Enter amount in ${
                    userData?.currency ? userData?.currency : ""
                  }`}
                />
                <TextInput
                  name="amount"
                  control={control}
                  placeholder="0"
                  localType="figure"
                />
                <InputInfoLabel
                  title="Buying Rate"
                  value={`1${returnAsset(userData?.crypto_currency)} = ${
                    !isNaN(parseFloat(rate?.data))
                      ? parseFloat(rate?.data)?.toFixed(2)
                      : "0.0"
                  }${userData?.currency ? userData?.currency : ""}`}
                />
                {isEquivalentError && (
                  <p className=" text-red-500 text-[10px]">
                    {equivalentError?.data?.error}
                  </p>
                )}
              </div>
              <div>
                <InputLabel text={`You will recieve (minus 1% service fee)`} />
                <TextInput
                  name=""
                  control={control}
                  readOnly
                  value={` ${
                    equivalent?.data?.amount
                      ? equivalent?.data?.amount?.replace(
                          /\B(?=(\d{3})+(?!\d))/g,
                          ","
                        )
                      : ""
                  }${
                    userData?.crypto_currency
                      ? returnAsset(userData?.crypto_currency)
                      : ""
                  }`}
                />
              </div>
            </section>

            <section className="mt-[24px]">
              <InputLabel text="Make your transfer with the account details below. Add reference to the description." />
              <div className="mt-[15px] flex flex-col gap-y-[10px]">
                <DetailsCard
                  text="Account Number"
                  value={bankDetails?.AccountNumber}
                />
                <DetailsCard text="Bank Name" value={bankDetails?.BankName} />
                <DetailsCard
                  text="Account Name"
                  value={bankDetails?.AccountName}
                />
                <DetailsCard
                  text="Reference"
                  value={reference?.data?.reference}
                />
              </div>
            </section>
            <FormButton
              label="I Have Paid"
              extraClass="mt-[50px]"
              loading={isLoading}
            />
          </form>
        </div>
      }
    />
  );
};

export default BankTransferDeposit;
