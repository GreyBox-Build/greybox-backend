import Navigation from "./home/Navigation";
import Footer from "./home/Footer";
import {
  LocationIcon,
  MailFIcon,
  PhoneIcon,
} from "../../components/icons/Icons";
import { HomeButton } from "../../components/buttons/HomeButton";
import { ContactTextInput } from "../../components/inputs/ContactTextInput";
import { ContactTextArea } from "../../components/inputs/ContactTextArea";
import { useScrollToTop } from "../../utils/ScrollToTop";

const Contact = () => {
  const ContactCard = ({
    icon,
    title,
    desc,
  }: {
    icon: React.ReactNode;
    title: string;
    desc: string;
  }) => {
    return (
      <div className="flex items-center gap-x-[17px]">
        <div className="flex items-center justify-center p-[9px] bg-pink-1 rounded-[4px]">
          {icon}
        </div>
        <div className="flex flex-col gap-y-[6px] text-white">
          <p className="text-[1.125rem] font-[600] leading-[18px]">{title}</p>
          <p className="text-[0.875rem] text-grey-5 leading-[18px]">{desc}</p>
        </div>
      </div>
    );
  };
  useScrollToTop();
  return (
    <section>
      <section className=" bg-pink-1 pb-[69px] flex flex-col">
        <Navigation />
        <div className="flex flex-col items-center px-[25px]">
          <h2 className="text-center text-[2rem] text-black-2 font-[700] leading-[40.63px] mt-[64px] mb-[16px]">
            Contact Us
          </h2>
          <p className="text-center text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
            At Greybox Technologies, Inc., we're passionate about transforming
            the way businesses and individuals in Africa manage their finances.
            Our mission is to provide a reliable, low-cost solution for
            cross-border payments and remittances.
          </p>
        </div>
      </section>
      <section className="flex flex-col md:flex-row px-[25px] md:px-[5%] lg:px-[10%] py-[81px] justify-center gap-y-[74px]">
        <section className="rounded-[16px] md:rounded-[16px_0px_0px_16px] bg-orange-1 p-[61px_45px] text-white shadow-md">
          <h2 className="text-[1.5rem] font-[600] leading-[24px]">
            Contact Information
          </h2>
          <p className="text-[0.875rem] text-grey-5 leading-[18px] max-w-[297px]">
            If you have any questions or concerns, please contact us
          </p>

          <div className="flex flex-col gap-y-[40px] mt-[51px]">
            <ContactCard
              title="Email"
              desc="info@greyboxpay.com"
              icon={<MailFIcon />}
            />
            <ContactCard
              title="Phone"
              desc="+233 2022680388"
              icon={<PhoneIcon />}
            />
            <ContactCard
              title="Address"
              desc="Denovo Plaza, Community 10. Tema"
              icon={<LocationIcon />}
            />
          </div>
        </section>
        <form className="bg-white rounded-[16px] md:rounded-[0px_16px_16px_0px] px-[34px] pt-[52px] pb-[26px]  w-full md:w-[70%] shadow-md">
          <section className="flex flex-col gap-y-[26px] mb-[55px]">
            <ContactTextInput placeholder="Full name" />
            <ContactTextInput placeholder="Email Address" />
            <ContactTextArea placeholder="Message" />
          </section>

          <HomeButton
            label="Send Message"
            onClick={() => {}}
            extraClass="bg-orange-1 text-white w-[225px]"
          />
        </form>
      </section>
      <Footer />
    </section>
  );
};

export default Contact;
