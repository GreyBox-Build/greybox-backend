import {
  FiveStarsSmall,
  User1Picture,
  User2Picture,
  User3Picture,
} from "../../../components/icons/Icons";

const Testimonials = () => {
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
            TESTIMONIALS
          </h2>
          <p className="max-w-[471px] text-[2rem] text-black-2 font-[700] leading-[40.63px] mt-[12px] ">
            What Our Clients Say About Us
          </p>
          <p className="text-[1rem] text-justify text-black-3 leading-[22px] max-w-[491px] mt-[20px]">
            At Greybox Technologies, Inc., we're passionate about transforming
            the way businesses and individuals in Africa manage their finances.
            Our mission is to provide a reliable, low-cost solution for
            cross-border payments and remittances.
          </p>
        </section>
        <section className=" bg-round-grey-bg bg-center bg-contain bg-no-repeat w-full md:w-[50%]">
          <div className="w-full min-h-[400px] flex flex-col items-center justify-center gap-y-[27px]">
            <TestimonialsCard
              userImg={<User1Picture />}
              userName="Koshie, Tema"
              comment="“Greybox is easy to use and saves us the trouble of waiting days before our suppliers receive funds”"
            />
            <TestimonialsCard
              userImg={<User2Picture />}
              userName="Fred, Accra"
              comment="“With Greybox, I don't have to worry about cedi depreciation against the dollar. That is a huge relief”"
              right
            />
            <TestimonialsCard
              userImg={<User3Picture />}
              userName="Akua, Tema"
              comment="“Greybox gave me the convenience of serving my clients outside my country without any delays.”"
            />
          </div>
        </section>
      </div>
    </section>
  );
};

export default Testimonials;
