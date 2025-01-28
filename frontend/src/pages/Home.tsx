import React from 'react';
import DiscountedProducts from '../components/DiscountedProducts';
import PersonalSelectionSection from '../components/PersonalSelection';
import MostLovedBlogs from '../components/MostLovedBlogs'; // Import the MostLovedBlogs component
import CookiePopup from '../components/CookiePopup';

interface HomeProps {
  showModal: boolean;
  openModal: () => void;
  closeModal: () => void;
  cookieAccepted: boolean;
  handleAcceptCookies: () => void;
}

const Home: React.FC<HomeProps> = ({
   cookieAccepted,
  handleAcceptCookies,
}) => {
  return (
    <div className="min-h-screen bg-white">
    

      {/* Personal Selection Section */}
      <PersonalSelectionSection />

      {/* Hottest Sales & Discounts Section */}
      <DiscountedProducts title="Hottest Sales & Discounts" />

      {/* Most Loved Blogs Section */}
      <MostLovedBlogs /> {/* Use the MostLovedBlogs component here */}

    

      {/* Modal Button */}
    

      {/* Cookie Popup */}
      {!cookieAccepted && (
        <CookiePopup onAccept={handleAcceptCookies} />
      )}
    </div>
  );
};

export default Home;