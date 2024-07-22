import { useSnackbar } from "notistack";

export const useCopyTextToClipboard = () => {
  const { enqueueSnackbar } = useSnackbar();
  const copyText = async (text: string, message: string) => {
    enqueueSnackbar(message, { variant: "success" });
    if ("clipboard" in navigator) {
      return await navigator.clipboard.writeText(text);
    }
  };

  return copyText;
};
