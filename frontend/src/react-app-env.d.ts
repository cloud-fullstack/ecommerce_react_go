declare namespace NodeJS {
  interface ProcessEnv {
    readonly REACT_APP_API_URL: string;
    readonly PUBLIC_URL: string;
    readonly INLINE_RUNTIME_CHUNK: string;
  }
}