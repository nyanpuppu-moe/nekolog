// Extract frame-light/dark and atom one hjs theme
import ThemeBase from "@/styles/theme-base.css?raw";
import Print from "@/styles/print.css?raw";

import FrameCss from "@milkdown/crepe/theme/frame.css?raw";
import FrameDarkCss from "@milkdown/crepe/theme/frame-dark.css?raw";

import HJSAtomOneLight from "highlight.js/styles/atom-one-light.min.css?raw";
import HJSAtomOneDark from "highlight.js/styles/atom-one-dark.min.css?raw";

export const FrameCssRoot = FrameCss.replace(/^ *\.milkdown *{/, "{");
export const FrameDarkCssRoot = FrameDarkCss.replace(/^ *\.milkdown *{/, "{");

export const Light = `
&${FrameCssRoot}
& .milkdown.view .ProseMirror {${HJSAtomOneLight}}
`;

export const Dark = `
&${FrameDarkCssRoot}
& .milkdown.view .ProseMirror {${HJSAtomOneDark}}
&{ --crepe-color-selected: #525252; }
`;

export const ThemeCss = `
@layer theme {
  @media (prefers-color-scheme: light) {
    :root {${Light}}
  }
  @media (prefers-color-scheme: dark) {
    :root {${Dark}}
  }
  :root.light {${Light}}
  :root.dark {${Dark}}
}
@layer core {
  :root {${ThemeBase}}
  html {
    background: var(--crepe-color-background);
  }
}
@layer print {
  @media print {
    :root {${Light}}
    ${Print}
  }
}
`;
