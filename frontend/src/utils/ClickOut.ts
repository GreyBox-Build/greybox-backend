import { useEffect } from "react";

type CheckOutType = {
  onState: Boolean | string | null;
  mainRef: React.RefObject<HTMLDivElement | HTMLUListElement | HTMLFormElement>;
  subRef: React.RefObject<HTMLImageElement | HTMLDivElement>;
  dispatch?: ({ type }: { type: string }) => {};
  dispatchType?: string;
  setFunc?: React.Dispatch<React.SetStateAction<Boolean>>;
};

export const useClickOut = ({
  onState,
  mainRef,
  subRef,
  dispatch,
  dispatchType,
  setFunc,
}: CheckOutType) => {
  useEffect(() => {
    const checkIfClickedOutside = (e: MouseEvent) => {
      if (
        onState &&
        mainRef?.current &&
        !mainRef?.current?.contains(e.target as Node) &&
        !subRef?.current?.contains(e.target as Node)
      ) {
        dispatch
          ? dispatch({ type: dispatchType ? dispatchType : "" })
          : setFunc && setFunc(!onState);
      }
    };
    document.addEventListener("mousedown", checkIfClickedOutside);
    return () => {
      document.removeEventListener("mousedown", checkIfClickedOutside);
    };
  }, [onState, mainRef, dispatchType, subRef, dispatch, setFunc]);
};
