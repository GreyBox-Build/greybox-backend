import AppLayout from "./AppLayout";
import { useNavigate } from "react-router-dom";
import { FlatButton } from "../../components/buttons/FlatButton";
import { FormButton } from "../../components/buttons/FormButton";

const SendOption = () => {
  const navigate = useNavigate();
  return (
    <AppLayout
      child={
        <div className="p-[25px] m-auto w-full h-fit md:w-[50.33%] lg:w-[45.33%] bg-grey-5 rounded-[32px_32px_0px_0px]">
          <FlatButton
            label="Send Via Mobile Money"
            extraClass="bg-grey-1 mb-[8px]"
            onClick={() => navigate("/send-via-mobile-money")}
          />
          <FlatButton
            label="Send Via Bank"
            extraClass="bg-[transparent] mb-[34px]"
            onClick={() => navigate("/send-via-bank")}
          />

          <FormButton label="Cancel" onClick={() => navigate(-1)} />
        </div>
      }
    />
  );
};

export default SendOption;
