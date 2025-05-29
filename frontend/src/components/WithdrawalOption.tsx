import AppLayout from "../pages/app/AppLayout";
import { useNavigate } from "react-router-dom";
import { FlatButton } from "./buttons/FlatButton";
import { FormButton } from "./buttons/FormButton";

const WithdrawalOption = () => {
  const navigate = useNavigate();
  return (
    <AppLayout
      child={
        <div className="p-[25px] m-auto w-full h-fit md:w-[50.33%] lg:w-[45.33%] bg-grey-5 rounded-[32px_32px_0px_0px]">
          <FlatButton
            label="Withdraw Via Bank "
            extraClass="bg-grey-1 mb-[8px]"
            onClick={() => navigate("/withdraw-via-bank")}
          />
          <FlatButton
            label="Withdraw Via Mobile Money"
            extraClass="bg-grey-1 mb-[8px]"
            onClick={() => navigate("/send-via-mobile-money")}
          />

          <FormButton label="Cancel" onClick={() => navigate(-1)} />
        </div>
      }
    />
  );
};

export default WithdrawalOption;
