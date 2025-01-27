import React from 'react';

const PersonalSelectionSection = () => {
  return (
    <section className="py-8 px-4">
      <h1 className="text-3xl font-bold text-center mb-2">Our Personal Selection</h1>
      <p className="text-center text-gray-600 mb-4">Rezaverse. Reserves you The Best.</p>
      <div className="relative">
        {/* Left Scroll Button */}
        <button className="absolute left-0 top-1/2 transform -translate-y-1/2 bg-white rounded-full p-2 shadow-lg z-10">
          <i className="fas fa-chevron-left"></i>
        </button>

        {/* Scrollable Container */}
        <div className="scroll-container">
          <img
            src={`${process.env.PUBLIC_URL}/images/personal-selection-left.png`}
            alt="Left Selection"
            className="image-item"
          />
          <img
            src={`${process.env.PUBLIC_URL}/images/personal-selection-1.png`}
            alt="Selection 1"
            className="image-item"
          />
          <img
            src={`${process.env.PUBLIC_URL}/images/personal-selection-2.png`}
            alt="Selection 2"
            className="image-item"
          />
          <img
            src={`${process.env.PUBLIC_URL}/images/personal-selection-3.png`}
            alt="Selection 3"
            className="image-item"
          />
          <img
            src={`${process.env.PUBLIC_URL}/images/personal-selection-4.png`}
            alt="Selection 4"
            className="image-item"
          />
          <img
            src={`${process.env.PUBLIC_URL}/images/personal-selection-5.png`}
            alt="Selection 5"
            className="image-item"
          />
          <img
            src={`${process.env.PUBLIC_URL}/images/personal-selection-6.png`}
            alt="Selection 6"
            className="image-item"
          />
          <img
            src={`${process.env.PUBLIC_URL}/images/personal-selection-right.png`}
            alt="Right Selection"
            className="image-item"
          />
        </div>

        {/* Right Scroll Button */}
        <button className="absolute right-0 top-1/2 transform -translate-y-1/2 bg-white rounded-full p-2 shadow-lg z-10">
          <i className="fas fa-chevron-right"></i>
        </button>
      </div>
    </section>
  );
};

export default PersonalSelectionSection;