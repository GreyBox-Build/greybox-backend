import { useEffect } from "react";

export const useScrollToTop = () => {
  useEffect(() => {
    // Scroll to the top of the component with smooth behavior
    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  }, []);
};
