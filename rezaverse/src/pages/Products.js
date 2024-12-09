import React from 'react';
import { FaChevronLeft, FaChevronRight } from 'react-icons/fa';

function Products() {
  const products = [
    { id: 1, name: 'Product name', price: 'L5.99', image: '/images/product-1.jpg' },
    { id: 2, name: 'Product name', price: 'L5.99', image: '/images/product-2.jpg' },
    { id: 3, name: 'Product name', price: 'L5.99', image: '/images/product-3.jpg' },
    { id: 4, name: 'Product name', price: 'L5.99', image: '/images/product-4.jpg' },
    { id: 5, name: 'Product name', price: 'L5.99', image: '/images/product-5.jpg' },
    { id: 6, name: 'Product name', price: 'L5.99', image: '/images/product-6.jpg' },
    // Add more products as needed
  ];

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6 text-center">Our Products</h1>
      
      <div className="mb-8">
        <h2 className="text-2xl font-semibold mb-4">Featured Products</h2>
        <div className="relative">
          <button className="absolute left-0 top-1/2 transform -translate-y-1/2 bg-white rounded-full p-2 shadow-md">
            <FaChevronLeft className="text-gray-600" />
          </button>
          <div className="flex overflow-x-auto space-x-4 py-4">
            {products.map((product) => (
              <div key={product.id} className="flex-none w-48">
                <img src={product.image} alt={product.name} className="w-full h-48 object-cover rounded-lg mb-2" />
                <h3 className="font-semibold">{product.name}</h3>
                <p className="text-gray-600">{product.price}</p>
              </div>
            ))}
          </div>
          <button className="absolute right-0 top-1/2 transform -translate-y-1/2 bg-white rounded-full p-2 shadow-md">
            <FaChevronRight className="text-gray-600" />
          </button>
        </div>
      </div>
      
      <div>
        <h2 className="text-2xl font-semibold mb-4">All Products</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {products.map((product) => (
            <div key={product.id} className="border rounded-lg p-4">
              <img src={product.image} alt={product.name} className="w-full h-48 object-cover rounded-lg mb-2" />
              <h3 className="font-semibold">{product.name}</h3>
              <p className="text-gray-600">{product.price}</p>
              <button className="mt-2 bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors">
                Add to Cart
              </button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default Products;