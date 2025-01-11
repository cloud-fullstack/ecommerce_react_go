import React, { useState, useEffect } from 'react';
import { Routes, Route, useLocation } from 'react-router-dom';
import NavBar from './components/NavBar';
import Home from './pages/Home';
import ProductPreview from './components/ProductPreview';
import Footer from './components/Footer';
import AllProductPage from './pages/AllProductPage';
import SignIn from './pages/SignIn';
import Cart from './pages/Cart';


function App() {
  const [showModal, setShowModal] = useState(false);
  const location = useLocation();

  const openModal = () => {
    setShowModal(true);
  };

  const closeModal = () => {
    setShowModal(false);
  };

  useEffect(() => {
    console.log("Modal state changed:", showModal);
  }, [showModal]);

  return (
    <main>
      <div className="content">
        {showModal && (
          <div className="modal">
            <button onClick={closeModal}>Close Modal</button>
            <Routes>
           
            <Route path="/" element={<Home />} />
            <Route path="/signIn" element={<SignIn />} />
            <Route path="/products" element={<AllProductPage />} />
            <Route path="/store/:storeID/:productID" element={<Product />} />
            <Route path="/cart" element={<Cart />} />
          
            </Routes>
          </div>
        )}
      </div>

      <button onClick={openModal}>Open Modal</button>
      <CookiePopUp />
      <footer>
        <Footer />
      </footer>
    </main>
  );
}

export default App;