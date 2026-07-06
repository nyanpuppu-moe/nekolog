import { defineMiddleware } from "astro:middleware";
import { setupI18n } from "@lingui/core";
import { setLinguiContext } from "lingui-for-astro";
import { catalog } from "./i18n/catalog";

const supportedLocales = ["ko", "en"];

function resolveLocale(request: Request): string {
  const acceptLanguage = request.headers.get("accept-language");

  if (acceptLanguage) {
    // Extract lang prefix from header (ex: ko-KR,ko;q=0.9)
    const preferredLang = acceptLanguage
      .split(",")[0]
      .substring(0, 2)
      .toLowerCase();
    if (supportedLocales.includes(preferredLang)) {
      return preferredLang;
    }
  }
  return "en";
}

export const onRequest = defineMiddleware(async (context, next) => {
  // Get page lang
  const locale = resolveLocale(context.request);

  // Create i18n conext
  const i18n = setupI18n({ locale, messages: catalog });
  setLinguiContext(context.locals, i18n);

  // Get rendered page and set Vary header
  const response = await next();
  response.headers.set("Vary", "Accept-Language");

  return response;
});
