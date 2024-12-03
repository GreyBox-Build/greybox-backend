import Example from "../../../components/language";
import { useScrollToTop } from "../../../utils/ScrollToTop";
import About from "./About";
import Footer from "./Footer";
import GetStarted from "./GetStarted";
import Intro from "./Intro";
import Navigation from "./Navigation";
import Process from "./Process";
import Services from "./Services";
import Testimonials from "./Testimonials";
import WhyChoose from "./WhyChoose";

const Home = () => {
  useScrollToTop();
  return (
    <>
      <Example />
      <Navigation />
      <Intro />
      <Services />
      <About />
      <WhyChoose />
      <Process />
      <GetStarted />
      <Testimonials />
      <Footer />
    </>
  );
};

export default Home;
