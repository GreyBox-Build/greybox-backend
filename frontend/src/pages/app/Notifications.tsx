import { useNavigate } from "react-router-dom";
import { NotificationsCard } from "../../components/Cards";
import { CancelIcon } from "../../components/icons/Icons";
import { notificationsData } from "../../utils/Dummies";
import AppLayout from "./AppLayout";

const Notifications = () => {
  const navigate = useNavigate();
  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative ">
            <span className="absolute left-[24px]" onClick={() => navigate(-1)}>
              <CancelIcon />
            </span>{" "}
            <h2 className=" text-black text-[1.5rem] font-[600]">
              Notifications
            </h2>
          </div>
          <section className="flex flex-col gap-y-[32px] mt-[47px]">
            {notificationsData.map((note, index) => (
              <NotificationsCard
                key={index}
                text={note.text}
                subText={note.subText}
                date={note.date}
              />
            ))}
          </section>
        </div>
      }
    />
  );
};

export default Notifications;
