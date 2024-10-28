import { useNavigate } from "react-router-dom";
import { CancelIcon, DropDown } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
// import { InputLabel, TextInput } from "../../components/inputs/TextInput";
import { FormButton } from "../../components/buttons/FormButton";
// import SelectBox from "../../components/modals/SelectBox";
import { useState } from "react";
import { useForm } from "react-hook-form";
import {
  useGetNetworksQuery,
  useGetAuthUserQuery,
} from "../../appSlices/apiSlice";
import { mobileDepositSchema } from "../../utils/Validations";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useOfframpMobileMutation } from "../../appSlices/apiSlice";
import { enqueueSnackbar } from "notistack";
import {
  TextInput,
  InputInfoLabel,
  InputLabel,
} from "../../components/inputs/TextInput";
import SelectBox from "../../components/modals/SelectBox";
import { PhoneInput } from "../../components/inputs/PhoneInput";

type MobileDepositForm = z.infer<typeof mobileDepositSchema>;

const SendMobileMoney = () => {
  const navigate = useNavigate();
  const [openCountry, setOpenCountry] = useState<boolean>(false);
  const [openNetwork, setOpenNetwork] = useState(false);
  const [datame, setDatame] = useState<any>([]);

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
  console.log(personal_details);
  const { data: network } = useGetNetworksQuery({});

  const [offrampMobile, { isLoading }] = useOfframpMobileMutation();

  const onSubmit = async (data: MobileDepositForm) => {
    // Convert phoneNumber and amount to numbers before submitting

    const sentData = {
      sendingAddress: personal_details.account_address,
      amountSending: data.amount?.replace(/,/g, ""),
      network: personal_details.crypto_currency,
      token: personal_details.crypto_currency === "CELO" ? "cUSD" : "USDC",
      customerName: `${personal_details.first_name} ${personal_details.last_name}`,
      phoneNumber: selectCountryCodeByCountry(network?.data) + data.phoneNumber,
      countryCode: data.countryCode,
      mobileProvider: data.network,
    };

    setDatame(sentData);
    try {
      const response = await offrampMobile(sentData).unwrap();
      enqueueSnackbar(response?.status, { variant: "success" });
      setTimeout(() => {
        navigate("/dashboard");
      }, 3000);
    } catch (error: any) {
      enqueueSnackbar(error?.data?.error, { variant: "warning" });
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
              Send Via Mobile Money
            </h2>
          </div>
          <p className="text-black-3 text-[0.875rem] text-center">
            (Bal ${walletInfo?.balance || 0})
          </p>

          {/* Use handleSubmit to process form data on submission */}
          <div className="w-full  mx-auto p-4">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              {/* Country Code */}
              <div>
                <InputLabel text="Select Country" />
                <TextInput
                  control={control}
                  name="country"
                  placeholder="Country"
                  readOnly
                  type="text"
                  onClick={() => {
                    setOpenCountry(true);
                  }}
                  img={<DropDown />}
                />
              </div>

              {/* Phone Number */}
              <div>
                <InputLabel text="Your Number" />
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

              {/* Network */}
              <div>
                <InputLabel text="Select Network Provider" />
                <TextInput
                  control={control}
                  name="network"
                  placeholder="Network"
                  readOnly
                  type="text"
                  onClick={() => {
                    setOpenNetwork(true);
                  }}
                  img={<DropDown />}
                />
              </div>

              {/* Amount */}
              <div>
                <InputLabel text="Amount to send" />
                <TextInput
                  name="amount"
                  control={control}
                  placeholder="Amount"
                  localType="figure"
                />
                <InputInfoLabel title="Exchange Rate" value="1cUSD = 1USD" />
              </div>

              {/* Submit Button */}
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
                clearErrors("country");
              }}
              onClose={() => setOpenCountry(false)}
            />

            {/* Network Select Modal */}
            {/* Country Select Modal */}
            <SelectBox
              state={openNetwork}
              title="Select Network"
              placeholder="Search Network"
              childList={
                network?.data ? groupNetworkByCountry(network?.data) : []
              }
              type="network"
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

export default SendMobileMoney;
