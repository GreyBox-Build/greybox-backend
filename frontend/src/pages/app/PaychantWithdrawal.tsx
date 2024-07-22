import { useNavigate } from "react-router-dom";
import { CancelIcon } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import { InputLabel, TextInput } from "../../components/inputs/TextInput";
import { FormButton } from "../../components/buttons/FormButton";
import { useForm } from "react-hook-form";

const PaychantWithdrawal = () => {
  const navigate = useNavigate();
  const { control } = useForm();
  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative ">
            <span className="absolute left-[24px]" onClick={() => navigate(-1)}>
              <CancelIcon />
            </span>
            <h2 className=" text-black text-[1.5rem] font-[600]">Withdraw</h2>
          </div>
          <form className="mt-[29px] px-[24px] pb-[80px]">
            <section className="flex flex-col gap-y-[32px]">
              <div>
                <InputLabel text="To: Your Email" />
                <TextInput
                  name="email"
                  control={control}
                  placeholder="Enter official email"
                  type="text"
                  onChange={() => {}}
                />
              </div>
              <div>
                <InputLabel text="Enter amount you want to withdraw" />
                <TextInput
                  name="amount"
                  control={control}
                  placeholder="0"
                  type="text"
                  onChange={() => {}}
                />
              </div>
            </section>
            <FormButton
              label="Submit Request"
              extraClass="mt-[80px]"
              onClick={() => {}}
            />
          </form>
        </div>
      }
    />
  );
};

export default PaychantWithdrawal;
