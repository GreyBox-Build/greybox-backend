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
import "./i18n";

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement
);
root.render(
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
);
