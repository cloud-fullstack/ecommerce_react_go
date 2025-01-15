import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import { BrowserRouter } from 'react-router-dom';


// Find the root element
const rootElement = document.getElementById('root');

// Check if the root element exists
if (!rootElement) {
  throw new Error("Root element with ID p'root' not found in the DOM.");
}

// Create the root and render the app
const root = ReactDOM.createRoot(rootElement);
root.render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>
);