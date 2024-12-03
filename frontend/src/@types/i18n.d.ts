import "react-i18next";
import { resources } from "../i18n";

declare module "react-i18next" {
  type DefaultResources = (typeof resources)["en"]["translation"];
  interface CustomTypeOptions {
    resources: DefaultResources;
  }
}
// import "react-i18next";
// import { resources } from "../i18n";
// import Resources from "./resources";

// declare module "i18next" {
//   // DefaultResources = (typeof resources)["en"]["translation"];
//   interface CustomTypeOptions {
//     defaultNS: "ns1";
//     resources: Resources;
//   }
// }
