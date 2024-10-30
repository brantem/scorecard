/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string;
  readonly VITE_GENERATOR_DELAY: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
