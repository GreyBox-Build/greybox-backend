import { useNavigate } from "react-router-dom";
import { CancelIcon, DropDown } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import {
  InputInfoLabel,
  InputLabel,
  TextInput,
} from "../../components/inputs/TextInput";
import { FormButton } from "../../components/buttons/FormButton";
import SelectBox from "../../components/modals/SelectBox";
import { momoData, currencyData } from "../../utils/Dummies";
import { useState } from "react";
import { useForm } from "react-hook-form";

const SendMobile = () => {
  const navigate = useNavigate();

  const [momo, setMomo] = useState<string>("");
  const [openMomo, setOpenMomo] = useState<boolean>(false);

  const [currency, setCurrency] = useState<string>("");
  const [openCurrency, setOpenCurrency] = useState<boolean>(false);
  const { control } = useForm();

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
                <InputLabel text="Select MoMo" />
                <TextInput
                  name="momo"
                  control={control}
                  placeholder="Select Momo"
                  readOnly
                  type="text"
                  value={momo}
                  onClick={() => {
                    setOpenMomo(true);
                  }}
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
                <InputLabel text="Currency recipient will get" />
                <TextInput
                  name="currency"
                  control={control}
                  placeholder="Select currency"
                  readOnly
                  type="text"
                  value={currency}
                  onClick={() => {
                    setOpenCurrency(true);
                  }}
                  img={<DropDown />}
                />
              </div>
              <div>
                <InputLabel text="Amount to send" />
                <TextInput
                  name="amount-to-send"
                  control={control}
                  placeholder="0"
                  type="text"
                  onChange={() => {}}
                />
                <InputInfoLabel title="Exchange Rate" value="1cUSD = USD" />
              </div>
              <div>
                <InputLabel text="Recipient will get" />
                <TextInput
                  name="recipient-will-get"
                  control={control}
                  placeholder="₦ Amount to send"
                  type="text"
                  onChange={() => {}}
                />
                <InputInfoLabel title="Transaction fee" value="$0.000" />
              </div>
            </section>
            <FormButton
              label="Send"
              extraClass="mt-[80px]"
              onClick={() => {}}
            />
            <SelectBox
              state={openMomo}
              title="Select Momo"
              placeholder="Search Momo"
              childList={momoData}
              onPickChild={(list) => setMomo(list?.name)}
              onClose={() => setOpenMomo(false)}
            />
            <SelectBox
              state={openCurrency}
              title="Select Currency"
              placeholder="Search Currency"
              childList={currencyData}
              onPickChild={(list) => setCurrency(list?.name)}
              onClose={() => setOpenCurrency(false)}
            />
          </form>
        </div>
      }
    />
  );
};

export default SendMobile;
