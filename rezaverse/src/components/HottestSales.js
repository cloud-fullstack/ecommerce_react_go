import React from 'react';

function HottestSales() {
  const products = [
    { id: 1, name: 'Product name', price: 'L5.99', image: '/images/product-1.jpg' },
    { id: 2, name: 'Product name', price: 'L5.99', image: '/images/product-2.jpg' },
    { id: 3, name: 'Product name', price: 'L5.99', image: '/images/product-3.jpg' },
    { id: 4, name: 'Product name', price: 'L5.99', image: '/images/product-4.jpg' },
    { id: 5, name: 'Product name', price: 'L5.99', image: '/images/product-5.jpg' },
    { id: 6, name: 'Product name', price: 'L5.99', image: '/images/product-6.jpg' },
  ];

  return (
    <section className="py-16">
      <div className="text-center mb-8">
        <h2 className="text-3xl font-bold mb-2">Hottest Sales & Discounts</h2>
        <h3 className="text-xl mb-2">Of the Moment</h3>
        <p className="text-gray-600 mb-4">The Must-Haves. Trending now.</p>
        <a href="#" className="text-blue-600">View All</a>
      </div>
      <div className="relative max-w-6xl mx-auto">
        <button className="absolute left-0 top-1/2 transform -translate-y-1/2 text-3xl">&lt;</button>
        <div className="flex justify-center space-x-4">
          {products.map((product) => (
            <div key={product.id} className="text-center">
              <img src={product.image} alt={product.name} className="rounded-lg mb-2" width="200" height="200" />
              <p className="text-sm">{product.name}</p>
              <p className="text-sm">{product.price}</p>
            </div>
          ))}
        </div>
        <button className="absolute right-0 top-1/2 transform -translate-y-1/2 text-3xl">&gt;</button>
      </div>
    </section>
  );
}

export default HottestSales;