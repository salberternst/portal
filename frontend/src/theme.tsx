import { nanoLightTheme, nanoDarkTheme } from "react-admin";

nanoLightTheme.palette.background.default = "#ff";
nanoLightTheme.palette.background.paper = "#ff";
nanoLightTheme.palette.common.white = "#ff";
nanoLightTheme.palette.primary.main = "#ff";
nanoLightTheme.palette.primary.light = "#ff";
nanoLightTheme.palette.primary.dark = "#ff";
nanoLightTheme.palette.secondary.main = "#ff";
nanoLightTheme.palette.secondary.light = "#ff";
nanoLightTheme.palette.secondary.dark = "#ff";
nanoLightTheme.palette.text.primary = "#33";
nanoLightTheme.palette.text.secondary = "#33";
// nanoLightTheme.components = {
//     ...nanoLightTheme.components,
//     ...{
//         "RaMenuItemLink": {
//             "styleOverrides": {
//                 "root": {
//                     "paddingLeft": "8px",
//                     "paddingRight": "8px",
//                     "&.RaMenuItemLink-active": {
//                         "color": "rgb(27, 67, 206)",
//                         "fontWeight": 700,
//                         "& .MuiSvgIcon-root": {
//                             "fill": "rgb(27, 67, 206)"
//                         }
//                     }
//                 }
//             }
//         }
//     }
// }

export const theme = nanoLightTheme;
export const darkTheme = nanoDarkTheme;
