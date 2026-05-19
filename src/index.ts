import { Container, getContainer } from "@cloudflare/containers";
import { env } from "cloudflare:workers";

export class ApiContainer extends Container {
  defaultPort = 8080;
  sleepAfter = "30m";

  envVars = {
    PORT: "8080",
    GIN_MODE: env.GIN_MODE,
    JWT_SECRET: env.JWT_SECRET,
    SUPABASE_PROJECT_ID: env.SUPABASE_PROJECT_ID,
    SUPABASE_URL: env.SUPABASE_URL,
    SUPABASE_ANON_KEY: env.SUPABASE_ANON_KEY,
    SUPABASE_DB_PASSWORD: env.SUPABASE_DB_PASSWORD,
    SUPABASE_REGION: env.SUPABASE_REGION,
    SUPABASE_POOLER: env.SUPABASE_POOLER,
    SUPABASE_DB_MODE: env.SUPABASE_DB_MODE,
  };
}

export default {
  async fetch(request: Request): Promise<Response> {
    return getContainer(env.API_CONTAINER, "primary").fetch(request);
  },
};
