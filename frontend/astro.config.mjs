// @ts-check
import { defineConfig, envField } from "astro/config";
import node from "@astrojs/node";
import react from "@astrojs/react";

// https://astro.build/config
export default defineConfig({
  output: "server",
  integrations: [react()],
  adapter: node({
    mode: "standalone",
  }),
  env: {
    schema: {
      PUBLIC_SITE_NAME: envField.string({
        context: "client",
        access: "public",
        default: "NekoLog",
      }),
    },
  },
});
