import Navigation from "./home/Navigation";
import Footer from "./home/Footer";
import { useScrollToTop } from "../../utils/ScrollToTop";

const PrivacyPolicy = () => {
  useScrollToTop();
  return (
    <section>
      <section className=" bg-pink-1 pb-[69px] flex flex-col">
        <Navigation />
        <div className="flex flex-col px-[25px] md:px-[10%] m-[0_auto]">
          <h2 className="text-center text-[2rem] text-black-2 font-[700] leading-[40.63px] mt-[64px] mb-[16px]">
            Privacy Policy
          </h2>
          <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
            Greybox Technologies, Inc. is committed to protecting the privacy
            and security of our customers&#39; personal information. This
            Privacy Policy outlines how we collect, use, disclose, and safeguard
            your information when you use our products or services.
          </p>

          <section className="mb-6">
            <h2 className=" text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              1. Information We Collect:
            </h2>
            <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              We may collect the following types of information when you use our
              products or services:
            </p>
            <div>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                Personal Information: This may include your name, email address,
                contact information, and any other information you provide to us
                voluntarily.
              </p>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                Transaction Information: We may collect details of transactions
                you carry out through our platform, including blockchain
                addresses and transaction history.
              </p>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                Usage Information: We gather information about how you interact
                with our platform, such as your browsing activity, device
                information, and IP address.
              </p>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                Cookies and Similar Technologies: We may use cookies and similar
                technologies to collect information about your browsing behavior
                and preferences.
              </p>
            </div>
          </section>

          <section className="mb-6">
            <h2 className=" text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              2. How We Use Your Information:
            </h2>
            <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              We may use the information we collect for the following purposes:
            </p>
            <div>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                To provide and improve our products and services.
              </p>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                To personalize your experience and customize content.
              </p>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                To communicate with you about our products, services, and
                promotions.
              </p>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                To detect and prevent fraud, abuse, and unauthorized access.
              </p>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                To comply with legal and regulatory requirements.
              </p>
            </div>
          </section>

          <section className="mb-6">
            <h2 className=" text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              3. Information Sharing and Disclosure:
            </h2>
            <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              We may share your information with third parties in the following
              circumstances:
            </p>
            <div>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                With service providers who assist us in operating our platform
                and providing our services.
              </p>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                With business partners and affiliates for marketing and
                promotional purposes, with your consent.
              </p>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                In response to legal requests, court orders, or government
                inquiries.
              </p>
              <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
                In connection with a merger, acquisition, or sale of assets.
              </p>
            </div>
          </section>

          <section className="mb-6">
            <h2 className=" text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              4. Data Security:
            </h2>
            <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
              We implement appropriate technical and organizational measures to
              protect your information from unauthorized access, disclosure,
              alteration, or destruction.
            </p>
          </section>

          <section className="mb-6">
            <h2 className=" text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              5. Your Choices:
            </h2>
            <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
              You may choose not to provide certain information or opt-out of
              certain communications. However, this may limit your ability to
              access certain features of our platform.
            </p>
          </section>

          <section className="mb-6">
            <h2 className=" text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              6. Children&#39;s Privacy:
            </h2>
            <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
              Our products and services are not intended for individuals under
              the age of 18. We do not knowingly collect personal information
              from children without parental consent.
            </p>
          </section>
          <section className="mb-6">
            <h2 className=" text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              7. International Transfers:
            </h2>
            <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
              Your information may be transferred to, and processed in,
              countries other than your own. By using our platform, you consent
              to the transfer of your information to these countries.
            </p>
          </section>
          <section className="mb-6">
            <h2 className=" text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              8. Changes to this Privacy Policy:
            </h2>
            <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
              We may update this Privacy Policy from time to time. We will
              notify you of any changes by posting the new Privacy Policy on our
              website.
            </p>
          </section>
          <section className="mb-6">
            <h2 className=" text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              9. Contact Us:
            </h2>
            <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
              If you have any questions or concerns about our Privacy Policy or
              practices, please contact us at info@greyboxpay.com.
            </p>
          </section>

          <section className="mb-6">
            <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px] mb-6">
              This Privacy Policy was last updated on 01/01/2024.
            </p>
            <p className="text-justify text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
              Greybox Technologies, Inc. strives to maintain the highest
              standards of privacy protection and compliance with applicable
              laws and regulations. Thank you for trusting us with your
              information
            </p>
          </section>
        </div>
      </section>

      <Footer />
    </section>
  );
};

export default PrivacyPolicy;
