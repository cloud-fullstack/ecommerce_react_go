/// <reference types="react-scripts" />

declare namespace NodeJS {
    interface ProcessEnv {
      readonly REACT_APP_API_SERVER_DEV: string;
      readonly REACT_APP_FRONTEND_DEV: string;
      readonly REACT_APP_API_SERVER: string;
      readonly REACT_APP_FRONTEND: string;
      // Add other environment variables here as needed
    }
  }