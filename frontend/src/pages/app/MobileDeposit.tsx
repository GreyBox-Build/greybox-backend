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
} from "../../appSlices/apiSlice";
import { mobileDepositSchema } from "../../utils/Validations";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useOnrampMobileMutation } from "../../appSlices/apiSlice";
import { enqueueSnackbar } from "notistack";
import { PhoneInput } from "../../components/inputs/PhoneInput";
import ProcessingOverlay from "../../components/Processing";

type MobileDepositForm = z.infer<typeof mobileDepositSchema>;

const MobileDeposit = () => {
  const navigate = useNavigate();
  const [openCountry, setOpenCountry] = useState<boolean>(false);
  const [openNetwork, setOpenNetwork] = useState(false);
  const [isProcessing, setIsProcessing] = useState<boolean>(false);
  const [selectedCountry, setSelectedCountry] = useState<any>(null); // State to track selected country

  const { control, handleSubmit, clearErrors, setValue, watch } = useForm({
    defaultValues: {
      country: "",
      countryCode: "",
      phoneNumber: "",
      network: "",
      amount: "",
    },
    resolver: zodResolver(mobileDepositSchema),
  });

  const { currentData: userData } = useGetAuthUserQuery({});
  const walletInfo = userData?.data?.wallet_details;
  const personal_details = userData?.data?.personal_details;

  const { data: network } = useGetNetworksQuery({});
  const [onrampMobile, { isLoading }] = useOnrampMobileMutation();

  const onSubmit = async (data: MobileDepositForm) => {
    const requestData = {
      collection: {
        customerName: `${personal_details.first_name} ${personal_details.last_name}`,
        customerEmail: personal_details.email,
        phoneNumber:
          selectCountryCodeByCountry(network?.data) + data.phoneNumber,
        countryCode: data.countryCode,
        network: data.network,
        amount: Number(data.amount?.replace(/,/g, "")),
      },
      transfer: {
        digitalNetwork: personal_details.crypto_currency,
        digitalAsset: "cUSD",
        walletAddress: personal_details.account_address,
      },
    };

    try {
      const response = await onrampMobile(requestData).unwrap();
      enqueueSnackbar(response?.status, { variant: "success" });
      setTimeout(() => {
        navigate("/dashboard");
      }, 3000);
    } catch (error: any) {
      enqueueSnackbar(error?.data?.error, { variant: "error" });
    }
  };

  const country_code = watch("countryCode");

  const groupNetworkByCountry = (data: any) => {
    const target = data?.find(
      (entry: any) => entry?.countryCode === country_code
    );
    if (target) {
      return target.networks;
    }
  };

  const selectCountryCodeByCountry = (data: any) => {
    const target = data?.find(
      (entry: any) => entry?.countryCode === country_code
    );
    if (target) {
      return target.mobileCode;
    }
  };

  const handleStartProcessing = (): void => {
    setIsProcessing(true);
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

          <div className="w-full  mx-auto p-4">
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
                <InputLabel text="Phone Number" />
                <PhoneInput
                  name="phoneNumber"
                  control={control}
                  placeholder="0000000000"
                  localType="number"
                  countryCode={
                    selectCountryCodeByCountry(network?.data) || "+123"
                  }
                />
              </div>

              <div>
                <InputLabel text="Select Network Provider" />
                <TextInput
                  control={control}
                  name="network"
                  placeholder="Network"
                  readOnly
                  type="text"
                  onClick={() => setOpenNetwork(true)}
                  img={<DropDown />}
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

            {/* Country Select Modal */}
            <SelectBox
              state={openCountry}
              title="Select Country"
              placeholder="Search Country"
              childList={network?.data ? network?.data : []}
              type="countryName"
              onPickChild={(list: any) => {
                setValue("country", list?.countryName);
                setValue("countryCode", list?.countryCode);
                setValue("network", "");
                setSelectedCountry(list); // Update selected country state
                clearErrors("country");
              }}
              onClose={() => setOpenCountry(false)}
            />

            {/* Network Select Modal */}

            <SelectBox
              state={openNetwork}
              title="Select Network"
              placeholder="Search Network"
              childList={
                network?.data ? groupNetworkByCountry(network?.data) : []
              }
              type="network"
              selectedCountry={selectedCountry} // Pass the selected country here
              onPickChild={(list: any) => {
                setValue("network", list);
                clearErrors("network");
              }}
              onClose={() => setOpenNetwork(false)}
            />
          </div>
        </div>
      }
    />
  );
};

export default MobileDeposit;
