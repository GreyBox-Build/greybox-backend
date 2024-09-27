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
import { useOnrampMobileMutation } from "../../appSlices/apiSlice";
import { enqueueSnackbar } from "notistack";

type MobileDepositForm = z.infer<typeof mobileDepositSchema>;

const MobileDeposit = () => {
  const navigate = useNavigate();
  const [openCountry, setOpenCountry] = useState<boolean>(false);
  const [openNetwork, setOpenNetwork] = useState(false);
  const [selectedCountry, setSelectedCountry] = useState<any>(null);
  const [networks, setNetworks] = useState<string[]>([]);
  const [mobileCode, setMobileCode] = useState<string>("");
  const [phoneNumber, setPhoneNumber] = useState<string>(""); // State for phone number

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
  } = useForm<MobileDepositForm>({
    resolver: zodResolver(mobileDepositSchema),
  });

  const { currentData: userData } = useGetAuthUserQuery({});
  const walletInfo = userData?.data?.wallet_details;
  const personal_details = userData?.data?.personal_details;

  const { data: network } = useGetNetworksQuery({});

  const [onrampMobile, { isLoading }] = useOnrampMobileMutation();

  const onSubmit = async (data: MobileDepositForm) => {
    // Convert phoneNumber and amount to numbers before submitting

    const requestData = {
      collection: {
        customerName: `${personal_details.first_name} ${personal_details.last_name}`,
        customerEmail: personal_details.email,
        phoneNumber: data.phoneNumber,
        countryCode: data.countryCode,
        network: data.network,
        amount: Number(data.amount),
      },
      transfer: {
        digitalNetwork: personal_details.crypto_currency,
        digitalAsset: "cUSD",
        walletAddress: personal_details.account_address,
      },
    };

    setSelectedCountry(requestData);

    try {
      const response = await onrampMobile(requestData).unwrap();
      enqueueSnackbar(response?.status, { variant: "success" });
      setTimeout(() => {
        navigate("/dashboard");
      }, 3000);

      console.log(response?.status, requestData);
    } catch (error: any) {
      console.log(error);
      enqueueSnackbar(error?.data?.error, { variant: "success" });
    }
  };

  const handleCountrySelect = (code: string) => {
    const selected = network.data.find(
      (country: any) => country.countryCode === code
    );
    if (selected) {
      // setSelectedCountry(selected);
      setNetworks(selected.networks); // Set networks based on selected country
      setValue("countryCode", selected.countryCode); // Set the country code in the form
      setMobileCode(selected.mobileCode);
    }
    setOpenCountry(false);
  };

  const handleNetworkSelect = (network: any) => {
    setValue("network", network);
    setOpenNetwork(false);
  };

  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative">
            <span className="absolute left-[24px]" onClick={() => navigate(-1)}>
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
                <label
                  htmlFor="countryCode"
                  className="block text-sm font-medium"
                >
                  Country Code
                </label>
                <input
                  type="text"
                  {...register("countryCode")}
                  readOnly
                  placeholder="Select Country Code"
                  onClick={() => setOpenCountry(true)}
                  className="mt-1 block w-full px-3 py-3 rounded-xl bg-transparent border-[1px] border-gray-400 outline-none"
                />
                {errors.countryCode && (
                  <p className="text-red-500 text-sm">
                    {errors.countryCode.message || "Invalid phone number"}
                  </p>
                )}
              </div>

              {/* Phone Number */}
              <div className="flex items-center">
                <div className="flex-1">
                  <label
                    htmlFor="phoneNumber"
                    className="block text-sm font-medium"
                  >
                    Phone Number
                  </label>
                  <input
                    type="text"
                    {...register("phoneNumber")}
                    placeholder="Enter phone number"
                    value={phoneNumber}
                    className="mt-1 block w-full px-3 py-3 rounded-xl bg-transparent border-[1px] border-gray-400 outline-none"
                    onChange={(e) => {
                      if (e.target.value.startsWith(mobileCode)) {
                        setPhoneNumber(e.target.value);
                      } else {
                        setPhoneNumber(mobileCode + e.target.value); // Append if necessary
                      }
                    }}
                  />
                  {errors.phoneNumber && (
                    <p className="text-red-500 text-sm">
                      {errors.phoneNumber.message || "Invalid phone number"}
                    </p>
                  )}
                </div>
              </div>

              {/* Network */}
              <div>
                <label htmlFor="network" className="block text-sm font-medium">
                  Network
                </label>
                <input
                  type="text"
                  {...register("network")}
                  readOnly
                  placeholder="Select Network"
                  onClick={() => setOpenNetwork(true)}
                  className="mt-1 block w-full px-3 py-3 rounded-xl bg-transparent border-[1px] border-gray-400 outline-none"
                />
                {errors.network && (
                  <p className="text-red-500 text-sm">
                    {errors.network.message || "Invalid phone number"}
                  </p>
                )}
              </div>

              {/* Amount */}
              <div>
                <label htmlFor="amount" className="block text-sm font-medium">
                  Amount
                </label>
                <input
                  type="number"
                  {...register("amount")}
                  placeholder="Enter amount"
                  className="mt-1 block w-full px-3 py-3 rounded-xl bg-transparent border-[1px] border-gray-400 outline-none"
                />
                {errors.amount && (
                  <p className="text-red-500 text-sm">
                    {errors.amount.message || "Invalid phone number"}
                  </p>
                )}
              </div>

              {/* Submit Button */}
              <FormButton label="Submit" type="submit" loading={isLoading} />

              {JSON.stringify(selectedCountry)}
            </form>

            {/* Country Select Modal */}
            {openCountry && (
              <div className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50">
                <div className="bg-white p-4 rounded-md shadow-md">
                  <h2 className="text-lg mb-2">Select Country Code</h2>
                  <ul className="space-y-2">
                    {network.data.map((country: any) => (
                      <li
                        key={country.countryCode}
                        onClick={() => handleCountrySelect(country.countryCode)}
                        className="cursor-pointer hover:bg-gray-200 p-2 rounded-md"
                      >
                        {country.countryName} ({country.countryCode})
                      </li>
                    ))}
                  </ul>
                  <button
                    onClick={() => setOpenCountry(false)}
                    className="mt-4 text-slate-500"
                  >
                    Close
                  </button>
                </div>
              </div>
            )}

            {/* Network Select Modal */}
            {openNetwork && (
              <div className="absolute w-full inset-0 flex items-center justify-center bg-black bg-opacity-50">
                <div className="bg-white p-4 rounded-md shadow-md w-1/2">
                  <h2 className="text-lg mb-2">Select Network</h2>
                  <ul className="space-y-2">
                    {networks.map((network) => (
                      <li
                        key={network}
                        onClick={() => handleNetworkSelect(network)}
                        className="cursor-pointer hover:bg-gray-200 p-2 rounded-md"
                      >
                        {network}
                      </li>
                    ))}
                  </ul>

                  <button
                    onClick={() => setOpenNetwork(false)}
                    className="mt-4 text-slate-500"
                  >
                    Close
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      }
    />
  );
};

export default MobileDeposit;
