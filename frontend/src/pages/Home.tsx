import React from 'react';
import Layout from '../components/Layout'; // Import the Layout component
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
    <Layout> {/* Wrap the Home content with the Layout component */}
      {/* Main Content */}
      <main className="flex-grow">
        {/* Personal Selection Section */}
        <PersonalSelectionSection />

        {/* Hottest Sales & Discounts Section */}
        <DiscountedProducts title="Hottest Sales & Discounts" />

        {/* Most Loved Blogs Section */}
        <MostLovedBlogs />
      </main>

      {/* Cookie Popup */}
      {!cookieAccepted && <CookiePopup onAccept={handleAcceptCookies} />}

      {/* Modal Button */}
      <button
        onClick={openModal}
        className="fixed bottom-4 right-4 bg-blue-500 text-white px-4 py-2 rounded-lg"
      >
        Open Modal
      </button>
    </Layout>
  );
};

export default Home;