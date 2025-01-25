import { useState, useEffect, useRef } from "react";
import QuickMenu from "./QuickMenu";
import Cart from "../pages/Cart";
import ProductPreview from "./ProductPreview";
import useStore from "../stores/useStore";
import apiClient from '../utils/api';

interface Product {
  product_id: string;
  product_name: string;
  store_id: string;
  picture_link: string;
  price: number;
  discounted_price: number;
  discounted: boolean;
}

const NavBar = () => {
  const [current, setCurrent] = useState("");
  const [showMobileMenu, setShowMobileMenu] = useState(false);
  const [searchProductName, setSearchProductName] = useState("");
  const [searchedProducts, setSearchedProducts] = useState<Product[]>([]);
  const [isSearchLoading, setIsSearchLoading] = useState(false);
  const { aviLegacyName } = useStore();
  const menuRef = useRef<HTMLElement | null>(null);

  // Fetch product previews
  const fetchProductPreviews = async () => {
    setIsSearchLoading(true);
    try {
      const res = await apiClient.get("/frontpage-product-previews/");
      const data = res.data;
      if (data.error) throw new Error(data.message);
      setSearchedProducts(data);
    } catch (err) {
      console.error(err);
      // Display a user-friendly error message
    } finally {
      setIsSearchLoading(false);
    }
  };

  // Handle search filter
  const searchNavbarFilter = (arr: Product[], searchTerm: string) => {
    searchTerm = searchTerm.toLowerCase();
    return arr.filter((product) =>
      product.product_name.toLowerCase().includes(searchTerm)
    );
  };

  // Handle mobile menu outside click
  const handleOutsideClick = (event: MouseEvent) => {
    if (showMobileMenu && !menuRef.current?.contains(event.target as Node)) {
      setShowMobileMenu(false);
    }
  };

  // Handle escape key press
  const handleEscape = (event: KeyboardEvent) => {
    if (showMobileMenu && event.key === "Escape") {
      setShowMobileMenu(false);
    }
  };

  // Add event listeners for mobile menu
  useEffect(() => {
    document.addEventListener("click", handleOutsideClick, false);
    document.addEventListener("keyup", handleEscape, false);

    return () => {
      document.removeEventListener("click", handleOutsideClick, false);
      document.removeEventListener("keyup", handleEscape, false);
    };
  }, [showMobileMenu]);

  // Fetch data on mount
  useEffect(() => {
    fetchProductPreviews();
  }, []);

  return (
    <nav className="flex items-center justify-between px-4 py-3 bg-white shadow sticky top-0 z-50">
      {/* Logo */}
      <img
        src="https://placehold.co/120x40"
        alt="The Rezaverse logo"
        className="h-10"
      />

      {/* Desktop Links */}
      <div className="hidden md:flex space-x-6">
        <a
          href="/"
          onClick={() => setCurrent("Home")}
          className={`text-gray-600 ${current === "Home" ? "font-bold" : ""}`}
          aria-label="What's Trending"
        >
          What's Trending
        </a>
        <a
          href="/products"
          onClick={() => setCurrent("Products")}
          className={`text-gray-600 ${current === "Products" ? "font-bold" : ""}`}
          aria-label="Products"
        >
          Products
        </a>
        <a href="/#faqTarget" className="text-gray-600" aria-label="FAQ">
          FAQ
        </a>
        {aviLegacyName && (
          <a href="/merchant" className="text-gray-600" aria-label="Creator Dashboard">
            Creator Dashboard
          </a>
        )}
      </div>

      {/* Search, Cart, and QuickMenu */}
      <div className="flex items-center space-x-4">
        {/* Search Bar */}
        <div className="relative">
          <input
            type="search"
            placeholder="Search"
            value={searchProductName}
            onChange={(e) => setSearchProductName(e.target.value)}
            className="pl-3 pr-10 py-2 border rounded-lg"
            aria-label="Search products"
          />
          <i className="fas fa-search absolute right-3 top-3 text-gray-400"></i>

          {/* Search Results Dropdown */}
          {searchProductName && (
            <div className="absolute mt-2 w-64 bg-white border rounded-lg shadow-lg z-50">
              {isSearchLoading ? (
                <p className="p-2 text-gray-600">Loading...</p>
              ) : (
                <ul className="divide-y divide-gray-200">
                  {searchNavbarFilter(searchedProducts, searchProductName).map((product) => (
                    <li key={product.product_id} className="p-2 hover:bg-gray-100">
                      <ProductPreview
                        name={product.product_name}
                        storeID={product.store_id}
                        productID={product.product_id}
                        pictureLink={product.picture_link}
                        price={product.price}
                        discountedPrice={product.discounted_price}
                        discountActive={product.discounted}
                      />
                    </li>
                  ))}
                </ul>
              )}
            </div>
          )}
        </div>

        {/* Cart */}
        <Cart />

        {/* QuickMenu */}
        <QuickMenu />

        {/* Mobile Menu Toggle */}
        <button
          onClick={() => setShowMobileMenu(!showMobileMenu)}
          className="md:hidden text-gray-600"
          aria-label="Toggle mobile menu"
        >
          <i className="fas fa-bars"></i>
        </button>
      </div>

      {/* Mobile Menu */}
      {showMobileMenu && (
        <div
          ref={menuRef as React.RefObject<HTMLDivElement>}
          className="md:hidden absolute top-16 right-4 bg-white border rounded-lg shadow-lg z-50"
          role="dialog"
          aria-modal="true"
        >
          <div className="space-y-2 p-4">
            <a
              href="/"
              onClick={() => setShowMobileMenu(false)}
              className="block text-gray-600"
              aria-label="What's Trending"
            >
              What's Trending
            </a>
            <a
              href="/products"
              onClick={() => setShowMobileMenu(false)}
              className="block text-gray-600"
              aria-label="Products"
            >
              Products
            </a>
            <a
              href="/#faqTarget"
              onClick={() => setShowMobileMenu(false)}
              className="block text-gray-600"
              aria-label="FAQ"
            >
              FAQ
            </a>
            {aviLegacyName && (
              <a
                href="/merchant"
                onClick={() => setShowMobileMenu(false)}
                className="block text-gray-600"
                aria-label="Creator Dashboard"
              >
                Creator Dashboard
              </a>
            )}
          </div>
        </div>
      )}
    </nav>
  );
};

export default NavBar;