import React from 'react';

interface CookiePopupProps {
  onAccept: () => void;
}

const CookiePopup: React.FC<CookiePopupProps> = ({ onAccept }) => {
  return (
    <div className="fixed bottom-4 left-4 bg-white p-4 rounded-lg shadow-lg">
      <p>We use cookies to provide a better user experience.</p>
      <button onClick={onAccept} className="mt-2 bg-blue-500 text-white px-4 py-2 rounded-lg">
        Accept
      </button>
    </div>
  );
};

export default CookiePopup;