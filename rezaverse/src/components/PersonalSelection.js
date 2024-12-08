import React from 'react';

function PersonalSelection() {
  return (
    <section className="text-center py-8">
      <h2 className="text-3xl font-bold mb-2">Our Personal Selection</h2>
      <p className="text-gray-600 mb-4">Rezaverse. Reserves you The Best.</p>
      <div className="flex justify-center space-x-2 mb-8">
        <a href="#" className="text-blue-600">Show Products</a>
        <span className="text-gray-400">|</span>
        <a href="#" className="text-blue-600">All Products</a>
      </div>
      <div className="relative max-w-6xl mx-auto">
        <button className="absolute left-0 top-1/2 transform -translate-y-1/2 text-3xl">&lt;</button>
        <div className="flex justify-between px-16">
          <img src="/images/personal-selection-left.png" alt="Stylish shoppers" className="w-1/4" />
          <div className="grid grid-cols-3 gap-4 w-2/4">
            <img src="/images/personal-selection-1.png" alt="Product 1" className="rounded-lg" />
            <img src="/images/personal-selection-2.png" alt="Product 2" className="rounded-lg" />
            <img src="/images/personal-selection-3.png" alt="Product 3" className="rounded-lg" />
            <img src="/images/personal-selection-4.png" alt="Product 4" className="rounded-lg" />
            <img src="/images/personal-selection-5.png" alt="Product 5" className="rounded-lg" />
            <img src="/images/personal-selection-6.png" alt="Product 6" className="rounded-lg" />
          </div>
          <img src="/images/personal-selection-right.png" alt="Fashion models" className="w-1/4" />
        </div>
        <button className="absolute right-0 top-1/2 transform -translate-y-1/2 text-3xl">&gt;</button>
      </div>
    </section>
  );
}

export default PersonalSelection;