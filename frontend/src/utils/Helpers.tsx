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
