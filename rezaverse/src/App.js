import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Navigation from './components/Navigation';
import Footer from './components/Footer';
import Home from './pages/Home';
import Products from './pages/Products';
import FAQ from './pages/FAQ';
import CreatorDashboard from './pages/CreatorDashboard';

function App() {
  return (
    <div className="App">
      <Navigation />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/products" element={<Products />} />
        <Route path="/faq" element={<FAQ />} />
        <Route path="/creator-dashboard" element={<CreatorDashboard />} />
      </Routes>
      <Footer />
    </div>
  );
}

export default App;