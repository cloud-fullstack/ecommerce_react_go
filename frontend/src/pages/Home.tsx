import React from 'react';
import NavBar from '../components/NavBar';
import Footer from '../components/Footer';
import DiscountedProducts from '../components/DiscountedProducts';
import PersonalSelectionSection from '../components/PersonalSelection';
import MostLovedBlogs from '../components/MostLovedBlogs';
import CookiePopup from '../components/CookiePopup';

interface HomeProps {
  showModal: boolean;
  openModal: () => void;
  closeModal: () => void;
  cookieAccepted: boolean;
  handleAcceptCookies: () => void;
}

const Home: React.FC<HomeProps> = ({
  showModal,
  openModal,
  closeModal,
  cookieAccepted,
  handleAcceptCookies,
}) => {
  return (
    <div className="min-h-screen bg-white flex flex-col">
      {/* Navigation */}
      <NavBar />

      {/* Main Content */}
      <main className="flex-grow">
        {/* Personal Selection Section */}
        <PersonalSelectionSection />

        {/* Hottest Sales & Discounts Section */}
        <DiscountedProducts title="Hottest Sales & Discounts" />

        {/* Most Loved Blogs Section */}
        <MostLovedBlogs />
      </main>

      {/* Footer Section */}
      <Footer />

      {/* Cookie Popup */}
      {!cookieAccepted && <CookiePopup onAccept={handleAcceptCookies} />}

      {/* Modal Button */}
      <button
        onClick={openModal}
        className="fixed bottom-4 right-4 bg-blue-500 text-white px-4 py-2 rounded-lg"
      >
        Open Modal
      </button>
    </div>
  );
};

export default Home;