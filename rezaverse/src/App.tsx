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
  const [cookieAccepted, setCookieAccepted] = useState(
    localStorage.getItem('cookieAccepted') === 'true'
  );
  const location = useLocation();

  const openModal = () => {
    setShowModal(true);
  };

  const closeModal = () => {
    setShowModal(false);
  };

  const handleAcceptCookies = () => {
    setCookieAccepted(true);
    localStorage.setItem('cookieAccepted', 'true');
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
              <Route path="/store/:storeID/:productID" element={<ProductPreview />} />
              <Route path="/cart" element={<Cart />} />
            </Routes>
          </div>
          
        )}
      </div>

      <button onClick={openModal}>Open Modal</button>

      {/* Render Cookie Popup if cookies are not accepted */}
      {!cookieAccepted && <CookiePopUp onAccept={handleAcceptCookies} />}

      <footer>
        <Footer />
      </footer>
    </main>
  );
}

// CookiePopUp Component
const CookiePopUp = ({ onAccept }: { onAccept: () => void }) => {
  return (
    <div className="container mx-auto opacity-95" style={styles.container}>
      <div className="rounded-lg" style={{ backgroundColor: 'rgb(255, 255, 255)' }}>
        <div className="w-72 bg-white rounded-lg shadow-md p-6" style={{ cursor: 'auto' }}>
          <div className="w-16 mx-auto relative -mt-10 mb-3">
            <img
              className="-mt-1"
              src="https://www.svgrepo.com/show/30963/cookie.svg"
              alt="Cookie Icon SVG"
            />
          </div>
          <span className="w-full sm:w-48 block leading-normal text-gray-800 text-md mb-3">
            We use cookies to provide a better user experience.
          </span>
          <div className="flex items-center justify-between">
            <a className="text-xs text-gray-400 mr-1 hover:text-gray-800" href="#">
              Privacy Policy
            </a>
            <div className="w-1/2">
              <button
                type="button"
                onClick={onAccept}
                className="py-2 px-4 bg-indigo-600 hover:bg-indigo-700 focus:ring-indigo-500 focus:ring-offset-indigo-200 text-white w-full transition ease-in duration-200 text-center text-base font-semibold shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2 rounded-lg"
              >
                Accept
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

// Styles for the container
const styles = {
  container: {
    marginTop: '70vh',
    maxWidth: '20%',
    position: 'fixed',
    marginLeft: '40vw',
  },
  '@media (min-width: 640px) and (max-width: 767px)': {
    container: {
      marginTop: '70vh',
      maxWidth: '20%',
      position: 'fixed',
      marginLeft: '30%',
    },
  },
};

export default App;