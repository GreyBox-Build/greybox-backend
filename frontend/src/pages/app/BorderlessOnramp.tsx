import { useNavigate } from "react-router-dom";
import { CancelIcon, DropDown } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import { InputLabel, TextInput } from "../../components/inputs/TextInput";
import { FormButton } from "../../components/buttons/FormButton";
import SelectBox from "../../components/modals/SelectBox";
import { useState } from "react";
import { useForm } from "react-hook-form";
import {
  useGetNetworksQuery,
  useGetAuthUserQuery,
  useOnrampBorderlessPaymentMutation,
} from "../../appSlices/apiSlice";
import { borderlessSchema } from "../../utils/Validations";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { enqueueSnackbar } from "notistack";
import ProcessingOverlay from "../../components/Processing";

interface OnrampBorderlessPaymentRequest {
  amount: string;
  asset: string;
  country: string;
  fiat: string;
}

type MobileDepositForm = z.infer<typeof borderlessSchema>;

const BorderlessOnramp = () => {
  const navigate = useNavigate();
  const [openCountry, setOpenCountry] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const { control, handleSubmit, setValue } = useForm<MobileDepositForm>({
    defaultValues: {
      amount: "",
      asset: "",
      country: "",
      fiat: "USD", // Default to USD, but allow change
    },
    resolver: zodResolver(borderlessSchema),
  });

  const { currentData: userData } = useGetAuthUserQuery({});
  const walletInfo = userData?.data?.wallet_details;

  const { data: network } = useGetNetworksQuery({});
  const [onrampBorderlessPayment] = useOnrampBorderlessPaymentMutation();

  const onSubmit = async (data: MobileDepositForm) => {
    setIsLoading(true);

    const requestData: OnrampBorderlessPaymentRequest = {
      amount: data.amount.replace(/,/g, ""),
      asset: data.asset, // Use the selected asset from the form
      country: data.country,
      fiat: data.fiat, // Use the selected fiat from the form
    };
    const countryChange = {
      ...requestData,
      country:
        data.country === "Ghana"
          ? "GH"
          : data.country === "Tanzania"
          ? "TZN"
          : "KEN",
    };

    try {
      const response = await onrampBorderlessPayment(countryChange).unwrap();

      if (response.status === "success") {
        enqueueSnackbar(response.message || "Payment initiated successfully", {
          variant: "success",
        });
        setTimeout(() => navigate("/dashboard"), 3000);
      } else {
        enqueueSnackbar(response.message || "Payment initiation failed", {
          variant: "error",
        });
      }
    } catch (error: any) {
      console.error("API Error:", error);
      const errorMessage =
        error?.data?.message || error?.message || "An error occurred.";
      enqueueSnackbar(errorMessage, { variant: "error" });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative">
            <span
              className="absolute left-[15px] md:left-[24px]"
              onClick={() => navigate(-1)}
            >
              <CancelIcon />
            </span>
            <h2 className="text-black text-[1.5rem] font-[600]">
              Deposit Via Mobile Money
            </h2>
          </div>
          <p className="text-black-3 text-[0.875rem] text-center">
            (Bal ${walletInfo?.balance || 0})
          </p>
          <div className="w-full mx-auto p-4">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div>
                <InputLabel text="Select Country" />
                <TextInput
                  control={control}
                  name="country"
                  placeholder="Country"
                  readOnly
                  type="text"
                  onClick={() => setOpenCountry(true)}
                  img={<DropDown />}
                />
              </div>
              <div>
                <InputLabel text="Select Asset" />
                <TextInput
                  control={control}
                  name="asset"
                  placeholder="Asset (e.g., USDC_POLYGON)" // Provide example
                  type="text"
                />
              </div>
              <div>
                <InputLabel text="Select Fiat" />
                <TextInput
                  control={control}
                  name="fiat"
                  placeholder="Fiat (USD or USDC)" // Provide options
                  type="text"
                />
              </div>

              <div>
                <InputLabel text="Enter Amount" />
                <TextInput
                  name="amount"
                  control={control}
                  placeholder="Amount"
                  localType="figure"
                />
              </div>

              <FormButton label="Submit" type="submit" loading={isLoading} />
            </form>

            <SelectBox
              state={openCountry}
              title="Select Country"
              placeholder="Search Country"
              childList={network?.data || []}
              type="countryName"
              onPickChild={(list: any) => {
                setValue("country", list?.countryName);
                setOpenCountry(false);
              }}
              onClose={() => setOpenCountry(false)}
            />
          </div>
        </div>
      }
    />
  );
};

export default BorderlessOnramp;
