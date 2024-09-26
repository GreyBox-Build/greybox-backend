import { useNavigate } from "react-router-dom";
import { CancelIcon, DropDown } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import { InputLabel, TextInput } from "../../components/inputs/TextInput";
import { FormButton } from "../../components/buttons/FormButton";
import SelectBox from "../../components/modals/SelectBox";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { useGetNetworksQuery } from "../../appSlices/apiSlice";

const MobileDeposit = () => {
  const navigate = useNavigate();
  const [selectedCountryCode, setSelectedCountryCode] = useState<string>("");
  const [networks, setNetworks] = useState<string[]>([]);
  const [mobileCode, setMobileCode] = useState<string>("");
  const [momo, setMomo] = useState<string>("");
  const [openCountry, setOpenCountry] = useState<boolean>(false);
  const [openMomo, setOpenMomo] = useState<boolean>(false);
  const { control } = useForm();

  // Yet to use this fetched network
  const { data: network } = useGetNetworksQuery({});
  console.log(network.data);

  const handleCountrySelect = (countryCode: string) => {
    const selected = network?.data?.find(
      (country: any) => country.countryCode === countryCode
    );
    if (selected) {
      setSelectedCountryCode(selected.countryCode);
      setNetworks(selected.networks);
      setMobileCode(selected.mobileCode);
    }
    setOpenCountry(false);
  };

  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative ">
            <span className="absolute left-[24px]" onClick={() => navigate(-1)}>
              <CancelIcon />
            </span>{" "}
            <h2 className=" text-black text-[1.5rem] font-[600]">
              Send Via Mobile Money
            </h2>
          </div>
          <p className="text-black-3 text-[0.875rem] text-center">
            (Bal $0.00)
          </p>
          <form className="mt-[29px] px-[24px] pb-[80px]">
            <section className="flex flex-col gap-y-[32px]">
              <div>
                <InputLabel text="Select Country" />
                <TextInput
                  name="country"
                  control={control}
                  placeholder="Select Country"
                  readOnly
                  type="text"
                  value={selectedCountryCode || ""}
                  onClick={() => setOpenCountry(true)}
                  img={<DropDown />}
                />
              </div>

              <div>
                <InputLabel text="Mobile Code" />
                <TextInput
                  name="mobile-code"
                  control={control}
                  placeholder="Mobile Code"
                  readOnly
                  type="text"
                  value={mobileCode}
                />
              </div>

              <div>
                <InputLabel text="Select Network" />
                <TextInput
                  name="momo"
                  control={control}
                  placeholder="Select Network"
                  readOnly
                  type="text"
                  value={momo}
                  onClick={() => setOpenMomo(true)}
                  img={<DropDown />}
                />
              </div>

              <div>
                <InputLabel text="Mobile Number" />
                <TextInput
                  name="mobile-number"
                  control={control}
                  placeholder="Enter Mobile number"
                  type="text"
                  onChange={() => {}}
                />
              </div>
              <div>
                <InputLabel text="Amount" />
                <TextInput
                  name="amount"
                  control={control}
                  placeholder="Enter Amount"
                  type="number"
                  onChange={() => {}}
                />
              </div>
            </section>

            <FormButton
              label="Send"
              extraClass="mt-[80px]"
              onClick={() => {}}
            />

            <SelectBox
              state={openCountry}
              title="Select Country"
              placeholder="Search Country"
              childList={network?.data.map((country: any) => ({
                name: country.countryCode,
              }))}
              onPickChild={(list) => handleCountrySelect(list?.name)}
              onClose={() => setOpenCountry(false)}
            />

            <SelectBox
              state={openMomo}
              title="Select Network"
              placeholder="Search Network"
              childList={networks.map((network) => ({ name: network }))}
              onPickChild={(list) => setMomo(list?.name)}
              onClose={() => setOpenMomo(false)}
            />
          </form>
        </div>
      }
    />
  );
};

export default MobileDeposit;
