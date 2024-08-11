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
import { withdrawViaBankSchema } from "../../utils/Validations";
import {
  useGetAuthUserQuery,
  useGetEquivalentAmountQuery,
  useGetExchangeRateQuery,
  useOfframpMutation,
} from "../../appSlices/apiSlice";
import {
  assignLocalError,
  removeLocalError,
  returnAsset,
} from "../../utils/Helpers";
import { ZodIssue } from "zod";
import { useState } from "react";
import { useSnackbar } from "notistack";

const BankWithdrawal = () => {
  const navigate = useNavigate();
  const { currentData: user } = useGetAuthUserQuery({});
  const [localErrors, setLocalErrors] = useState<ZodIssue[]>([]);

  const { enqueueSnackbar } = useSnackbar();

  const { control, handleSubmit, watch } = useForm({
    defaultValues: {
      cryptoAmount: "0",
    },
    resolver: zodResolver(withdrawViaBankSchema),
  });

  const cryptoAmount = watch("cryptoAmount");
  const userData = user?.data?.personal_details;

  const { currentData: rate } = useGetExchangeRateQuery({
    fiat: userData?.currency,
    asset: returnAsset(userData?.crypto_currency)?.toLocaleUpperCase(),
  });

  const [offramp, { isLoading }] = useOfframpMutation();

  const {
    currentData: equivalent,
    isError: isEquivalentError,
    error: equivalentError,
  }: any = useGetEquivalentAmountQuery({
    amount: cryptoAmount?.replace(/,/g, ""),
    currency: userData?.currency,
    cryptoAsset: returnAsset(userData?.crypto_currency)?.toLocaleUpperCase(),
    type: "off-ramp",
  });

  const handleDepositViaBankTransfer = async (data: any) => {
    const { cryptoAmount, bankName, accountNumber, accountName } = data;

    const details = {
      cryptoAmount: cryptoAmount?.replace(/,/g, ""),
      asset: userData?.crypto_currency,
      fiatEquivalent: equivalent?.data?.amount,
      chain: userData?.crypto_currency,
      currencyCode: userData?.currency_code,
      bankName,
      accountNumber,
      accountName,
    };

    try {
      const response = await offramp(details).unwrap();
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
            <span className="absolute left-[24px]" onClick={() => navigate(-1)}>
              <CancelIcon />
            </span>
            <h2 className=" text-black text-[1.5rem] font-[600]">
              Withdraw Via Bank
            </h2>
          </div>
          <form
            className="mt-[29px] px-[24px] pb-[80px]"
            onSubmit={handleSubmit(handleDepositViaBankTransfer)}
          >
            <section className="flex flex-col gap-y-[32px]">
              <div>
                <InputLabel text={`Enter amount in cUSD`} />
                <TextInput
                  name="cryptoAmount"
                  control={control}
                  placeholder="0"
                  // localType="figure"
                  onLocalChange={() => {
                    parseFloat(cryptoAmount) < 1
                      ? assignLocalError("cryptoAmount", localErrors)
                      : removeLocalError(
                          "cryptoAmount",
                          localErrors,
                          setLocalErrors
                        );
                  }}
                />
                <InputInfoLabel
                  title="Buying Rate"
                  value={`1${returnAsset(
                    userData?.crypto_currency
                  )} = ${parseFloat(rate?.data)?.toFixed(2)}${
                    userData?.currency ? userData?.currency : ""
                  }`}
                />
                {isEquivalentError && (
                  <p className=" text-red-500 text-[10px]">
                    {equivalentError?.data?.error}
                  </p>
                )}
              </div>
              <div>
                <InputLabel text={`You will recieve`} />
                <TextInput
                  name=""
                  control={control}
                  readOnly
                  value={`${
                    equivalent?.data?.asset ? equivalent?.data?.asset : "-"
                  } ${
                    equivalent?.data?.amount
                      ? equivalent?.data?.amount?.replace(
                          /\B(?=(\d{3})+(?!\d))/g,
                          ","
                        )
                      : ""
                  }`}
                />
              </div>

              <div>
                <InputLabel text={`Bank Name`} />
                <TextInput
                  name="bankName"
                  control={control}
                  placeholder="Enter Bank Name"
                />
              </div>
              <div>
                <InputLabel text={`Account Number`} />
                <TextInput
                  name="accountNumber"
                  control={control}
                  placeholder="Enter Account Number"
                  localType="number"
                />
              </div>
              <div>
                <InputLabel text={`Account Name`} />
                <TextInput
                  name="accountName"
                  control={control}
                  placeholder="Enter Account Name"
                />
              </div>
            </section>

            <FormButton
              label="Submit Request"
              extraClass="mt-[50px]"
              loading={isLoading}
            />
          </form>
        </div>
      }
    />
  );
};

export default BankWithdrawal;
