import { routes } from "./Routes";
import { RouterProvider } from "react-router-dom";
import "./i18n";
function App() {
  return <RouterProvider router={routes}></RouterProvider>;
}

export default App;
