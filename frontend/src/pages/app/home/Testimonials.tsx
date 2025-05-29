import { useTranslation } from "react-i18next";
import {
  FiveStarsSmall,
  User1Picture,
  User2Picture,
  User3Picture,
} from "../../../components/icons/Icons";

const Testimonials = () => {
  const { t }: { t: any } = useTranslation();
  const TestimonialsCard = ({
    userImg,
    userName,
    comment,
    right,
  }: {
    userImg: React.ReactNode;
    userName: string;
    comment: string;
    right?: boolean;
  }) => (
    <div
      className={`w-[85%] flex items-center p-[14px_20px] gap-x-[28px] rounded-[8px] ${
        right ? "bg-grey-5 self-end" : "bg-pink-1 self-start"
      }`}
    >
      {!right && <div> {userImg}</div>}
      <div>
        <div className="flex items-center justify-between">
          <p className="text-[1.125rem] text-black-2 font-[600] leading-[20px]">
            {userName}
          </p>
          <FiveStarsSmall />
        </div>
        <p className="text-[0.875rem] text-black-3 leading-[18px]">{comment}</p>
      </div>
      {right && <div> {userImg}</div>}
    </div>
  );
  return (
    <section className="pb-[100px] pt-[100px] bg-white px-[25px] md:px-[5%] lg:px-[10%]">
      <div className="w-full flex flex-col md:flex-row items-center justify-center gap-x-[2%] gap-y-[20px]">
        <section className="w-full md:w-[50%] ">
          <h2 className="text-[1.125rem] text-orange-1 font-[700] leading-[24.75px]">
            {t("testimonial")}
          </h2>
          <p className="max-w-[471px] text-[2rem] text-black-2 font-[700] leading-[40.63px] mt-[12px] ">
            {t("testimonialTitle")}
          </p>
          <p className="text-[1rem] text-justify text-black-3 leading-[22px] max-w-[491px] mt-[20px]">
            {t("testimonialSubTitle")}
          </p>
        </section>
        <section className=" bg-round-grey-bg bg-center bg-contain bg-no-repeat w-full md:w-[50%]">
          <div className="w-full min-h-[400px] flex flex-col items-center justify-center gap-y-[27px]">
            <TestimonialsCard
              userImg={<User1Picture />}
              userName="Koshie, Tema"
              comment={t("test1")}
            />
            <TestimonialsCard
              userImg={<User2Picture />}
              userName="Fred, Accra"
              comment={t("test2")}
              right
            />
            <TestimonialsCard
              userImg={<User3Picture />}
              userName="Akua, Tema"
              comment={t("test3")}
            />
          </div>
        </section>
      </div>
    </section>
  );
};

export default Testimonials;
