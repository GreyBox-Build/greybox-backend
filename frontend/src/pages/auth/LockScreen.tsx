import AuthLayout from "./AuthLayout";
import { Link } from "react-router-dom";

const LockScreen = () => {
  return (
    <AuthLayout
      child={
        <div className="w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1 p-[51px_25px]">
          <section className="w-full flex items-center justify-center gap-x-[16px] mt-[143px]">
            <div className="w-[18px] h-[18px] bg-orange-2 rounded-[50%]">
              <input type="number" className="w-full opacity-0" />
            </div>

            <div className="w-[18px] h-[18px] bg-orange-2 rounded-[50%]">
              <input type="number" className="w-full opacity-0" />
            </div>
            <div className="w-[18px] h-[18px] bg-orange-2 rounded-[50%]">
              <input type="number" className="w-full opacity-0" />
            </div>
            <div className="w-[18px] h-[18px] bg-orange-2 rounded-[50%]">
              <input type="number" className="w-full opacity-0" />
            </div>
          </section>
          <p className="text-[0.875rem] text-black-3 leading-[18px] text-center mt-[175px]">
            Enter Your Greybox Pin
          </p>
          <section className="flex flex-col gap-y-[8px] mt-[55px]">
            <Link
              to={""}
              className=" text-[0.875rem] text-center text-orange-1 leading-[18px] font-[700]"
            >
              Forgot Pin
            </Link>
          </section>
        </div>
      }
    />
  );
};

export default LockScreen;
