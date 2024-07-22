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
import { banksData, currencyData } from "../../utils/Dummies";
import { useState } from "react";
import { useForm } from "react-hook-form";

const SendBank = () => {
  const navigate = useNavigate();

  const [openBank, setOpenBank] = useState<boolean>(false);
  const [openCurrency, setOpenCurrency] = useState<boolean>(false);

  const { control, setValue } = useForm();

  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative ">
            <span className="absolute left-[24px]" onClick={() => navigate(-1)}>
              <CancelIcon />
            </span>{" "}
            <h2 className=" text-black text-[1.5rem] font-[600]">
              Send Via Bank
            </h2>
          </div>
          <p className="text-black-3 text-[0.875rem] text-center">
            (Bal $0.00)
          </p>
          <form className="mt-[29px] px-[24px] pb-[80px]">
            <section className="flex flex-col gap-y-[32px]">
              <div>
                <InputLabel text="Select Bank Account" />
                <TextInput
                  name="bank"
                  control={control}
                  placeholder="Select bank"
                  readOnly
                  type="text"
                  onClick={() => {
                    setOpenBank(true);
                  }}
                  img={<DropDown />}
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
                  onClick={() => {
                    setOpenCurrency(true);
                  }}
                  img={<DropDown />}
                />
              </div>
              <div>
                <InputLabel text="Amount to send" />
                <TextInput
                  name="amount_to_send"
                  control={control}
                  placeholder="0"
                  type="text"
                />
                <InputInfoLabel title="Exchange Rate" value="1cUSD = USD" />
              </div>
              <div>
                <InputLabel text="Recipient will get" />
                <TextInput
                  name="recipient_will_get"
                  control={control}
                  placeholder="₦ Amount to send"
                  type="text"
                />
                <InputInfoLabel title="Transaction fee" value="$0.000" />
              </div>
            </section>
            <FormButton label="Send" extraClass="mt-[80px]" />
            <SelectBox
              state={openBank}
              title="Select Bank"
              placeholder="Search Bank"
              childList={banksData}
              onPickChild={(list) => {
                setValue("bank", list?.name);
              }}
              onClose={() => setOpenBank(false)}
            />
            <SelectBox
              state={openCurrency}
              title="Select Currency"
              placeholder="Search Currency"
              childList={currencyData}
              onPickChild={(list) => {
                setValue("currency", list?.name);
              }}
              onClose={() => setOpenCurrency(false)}
            />
          </form>
        </div>
      }
    />
  );
};

export default SendBank;
