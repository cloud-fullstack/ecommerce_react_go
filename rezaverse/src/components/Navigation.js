import React from 'react';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSearch, faShoppingCart, faUserCircle } from '@fortawesome/free-solid-svg-icons';

function Navigation() {
  return (
    <nav className="flex items-center justify-between px-4 py-4">
      <div className="flex items-center">
        <img src="/images/logo.png" alt="The Rezaverse logo" className="mr-8" width="120" height="40" />
        <div className="space-x-6">
          <Link to="/" className="text-gray-700">What's Trending</Link>
          <Link to="/products" className="text-gray-700">Products</Link>
          <Link to="/faq" className="text-gray-700">FAQ</Link>
          <Link to="/creator-dashboard" className="text-gray-700">Creator Dashboard</Link>
        </div>
      </div>
      <div className="flex items-center space-x-4">
        <div className="relative">
          <input type="search" placeholder="Search" className="pl-8 pr-4 py-1 border rounded-lg" />
          <FontAwesomeIcon icon={faSearch} className="absolute left-2 top-2 text-gray-400" />
        </div>
        <FontAwesomeIcon icon={faShoppingCart} />
        <FontAwesomeIcon icon={faUserCircle} />
      </div>
    </nav>
  );
}

export default Navigation;