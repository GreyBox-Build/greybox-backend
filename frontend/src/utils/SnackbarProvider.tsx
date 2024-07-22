import { MaterialDesignContent } from "notistack";
import { IconButton, styled } from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import { useSnackbar } from "notistack";

export const StyledMaterialDesignContent = styled(MaterialDesignContent)(
  () => ({
    "&.notistack-MuiContent-success": {
      backgroundColor: "#CD5928",
    },
    "&.notistack-MuiContent-error": {
      backgroundColor: "#970C0C",
    },
  })
);

export const SnackbarCloseButton = ({
  snackbarId,
}: {
  snackbarId: number | string;
}) => {
  const { closeSnackbar } = useSnackbar();

  const style = {
    color: "#fff",
    height: "15px",
    width: "15px",
  };

  const ibStyle: any = {
    position: "absolute",
    top: "5px",
    right: "5px",
    height: "20px",
    width: "20px",
  };

  return (
    <IconButton onClick={() => closeSnackbar(snackbarId)} style={ibStyle}>
      <CloseIcon style={style} />
    </IconButton>
  );
};
