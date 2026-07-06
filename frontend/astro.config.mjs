// @ts-check
import { defineConfig, envField } from "astro/config";
import node from "@astrojs/node";
import react from "@astrojs/react";
import linguiForAstro from "lingui-for-astro/integration";
import linguiMacro from "unplugin-lingui-macro/vite";

// @ts-ignore We dont need full node types. vite env runs on nodejs
//            So we can just ignore it
try { process.loadEnvFile(); } catch {}
// @ts-ignore
const allowedHost = process.env.ASTRO_ALLOWED_HOST || 'localhost';
// @ts-ignore
const devHost = process.env.ASTRO_HOST === 'true' ? true : process.env.ASTRO_HOST;
// @ts-ignore
const devPort = process.env.ASTRO_PORT ? parseInt(process.env.ASTRO_PORT, 10) : undefined;

// https://astro.build/config
export default defineConfig({
  output: "server",
  integrations: [react(), linguiForAstro()],
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
  vite: {
    plugins: [linguiMacro()],
  },
  server: {
    host: devHost, 
    port: devPort,
    allowedHosts: [allowedHost]
  },
});
