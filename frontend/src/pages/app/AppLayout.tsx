import { useNavigate } from "react-router-dom";
import { Cards, Home, More, Send, Wallet } from "../../components/icons/Icons";
import { Menu } from "../../components/Cards";
import { useGetAuthUserQuery } from "../../appSlices/apiSlice";
import { useEffect } from "react";

const AppLayout = ({ child }: { child: React.ReactNode }) => {
  const navigate = useNavigate();
  const { error }: any = useGetAuthUserQuery({});
  useEffect(() => {
    setTimeout(() => {
      if (error?.status === 401) {
        localStorage.removeItem("access_token");
        navigate("/sign-in");
      }
    }, 500);
  });

  return (
    <div className=" bg-grey-box-bg bg-cover bg-no-repeat w-full min-h-[100vh] flex justify-center">
      {child}
      <section className="w-full h-[54px] flex justify-between fixed bottom-0 bg-grey-1 p-[5px_24px] md:p-[5px_65px]">
        <Menu
          icon={<Home />}
          label="Home"
          onClick={() => navigate("/dashboard")}
        />{" "}
        <Menu
          icon={<Send />}
          label="Send"
          onClick={() => navigate("/send-options")}
        />
        <Menu icon={<Cards />} label="Cards" onClick={() => navigate("")} />
        <Menu icon={<Wallet />} label="Wallet" onClick={() => navigate("")} />
        <Menu icon={<More />} label="More" onClick={() => navigate("")} />
      </section>
    </div>
  );
};

export default AppLayout;
