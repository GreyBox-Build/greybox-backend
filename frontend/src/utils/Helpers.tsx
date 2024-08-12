import moment from "moment";
import { ZodIssue } from "zod";

export const assignLocalError = (path: string, state: ZodIssue[]) => {
  const targetError = state?.find((err) => err?.path[0] === path);
  if (targetError) {
    return targetError?.message;
  }
  return null;
};

export const removeLocalError = (
  path: string,
  state: ZodIssue[],
  setState: React.Dispatch<React.SetStateAction<ZodIssue[]>>
) => {
  const filteredErrors = state?.filter((err) => err?.path[0] !== path);
  if (filteredErrors) {
    setState(filteredErrors);
  }
};

export const returnAsset = (chain: string) => {
  if (chain?.toLocaleLowerCase() === "celo") {
    return "cUSD";
  } else {
    return "USDC";
  }
};

export const groupByDate = (array: any[]) => {
  const grouped: any = {};

  array?.forEach((transaction) => {
    const date = moment(transaction?.timestamp).format("D/MM/YY");
    if (!grouped[date]) {
      grouped[date] = [];
    }
    grouped[date].push(transaction);
  });
  return [
    {
      date: Object.keys(grouped),
      transactions: Object.values(grouped),
    },
  ];
};

export const findSubArray = (arr: any[], param: string) => {
  for (let i = 0; i < arr.length; i++) {
    if (moment(arr[i][0]?.timestamp).format("D/MM/YY") === param) {
      return arr[i];
    }
  }
};
