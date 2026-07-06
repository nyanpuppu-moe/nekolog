import { defineConfig } from "@lingui/conf";
import { formatter } from "@lingui/format-po";
import babelExtractor from "@lingui/cli/api/extractors/babel";
import { astroExtractor } from "lingui-for-astro/extractor";

export default defineConfig({
  locales: ["en", "ko"],
  sourceLocale: "en",
  catalogs: [
    {
      path: "<rootDir>/src/i18n/locales/{locale}",
      include: ["src"],
      exclude: ["**/node_modules/**"],
    },
  ],
  format: formatter({ lineNumbers: false }),
  extractors: [astroExtractor, babelExtractor],
});
