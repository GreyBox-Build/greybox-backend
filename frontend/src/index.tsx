import ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App";
import { Provider } from "react-redux";
import { store } from "./app/store";
import { SnackbarProvider } from "notistack";
import {
  SnackbarCloseButton,
  StyledMaterialDesignContent,
} from "./utils/SnackbarProvider";
import { MoonPayProvider } from "@moonpay/moonpay-react";

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement
);
root.render(
  <MoonPayProvider apiKey={process.env.REACT_APP_MOONPAY_KEY!} debug>
    <Provider store={store}>
      <SnackbarProvider
        Components={{
          success: StyledMaterialDesignContent,
          error: StyledMaterialDesignContent,
        }}
        anchorOrigin={{
          vertical: "top",
          horizontal: "center",
        }}
        hideIconVariant
        autoHideDuration={5000}
        action={(snackbarId) => <SnackbarCloseButton snackbarId={snackbarId} />}
      >
        <App />
      </SnackbarProvider>
    </Provider>
  </MoonPayProvider>
);
