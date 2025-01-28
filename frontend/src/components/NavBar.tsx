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
  const menuRef = useRef<HTMLDivElement>(null);

  // Fetch product previews
  const fetchProductPreviews = async () => {
    setIsSearchLoading(true);
    try {
      const res = await apiClient.get("/api/frontpage-product-previews/");
      const data = res.data;
      if (data.error) throw new Error(data.message);
      setSearchedProducts(data);
    } catch (err) {
      console.error(err);
      alert("Failed to fetch product previews. Please try again later.");
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

  // Handle outside click for mobile menu
  useEffect(() => {
    const handleOutsideClick = (event: MouseEvent) => {
      if (showMobileMenu && !menuRef.current?.contains(event.target as Node)) {
        setShowMobileMenu(false);
      }
    };

    const handleEscape = (event: KeyboardEvent) => {
      if (showMobileMenu && event.key === "Escape") {
        setShowMobileMenu(false);
      }
    };

    document.addEventListener("click", handleOutsideClick);
    document.addEventListener("keyup", handleEscape);

    return () => {
      document.removeEventListener("click", handleOutsideClick);
      document.removeEventListener("keyup", handleEscape);
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
        className="h-10 rounded-full"
      />

      {/* Desktop Links */}
      <div className="hidden md:flex space-x-6">
        {["Home", "Products", "FAQ"].map((link) => (
          <a
            key={link}
            href={link === "FAQ" ? "/#faqTarget" : `/${link.toLowerCase()}`}
            onClick={() => setCurrent(link)}
            className={`text-gray-600 hover:text-gray-900 ${
              current === link ? "font-bold" : ""
            }`}
            aria-label={link}
          >
            {link === "Home" ? "What's Trending" : link}
          </a>
        ))}
        {aviLegacyName && (
          <a
            href="/merchant"
            className="text-gray-600 hover:text-gray-900"
            aria-label="Creator Dashboard"
          >
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
            className="pl-3 pr-10 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            aria-label="Search products"
          />
          <i className="fas fa-search absolute right-3 top-3 text-gray-400 pointer-events-none"></i>

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
                        product_name={product.product_name}
                        store_id={product.store_id}
                        product_id={product.product_id}
                        picture_link={product.picture_link}
                        price={product.price}
                        discounted_price={product.discounted_price}
                        discounted={product.discounted}
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
          className="md:hidden text-gray-600 hover:text-gray-900"
          aria-label="Toggle mobile menu"
        >
          <i className="fas fa-bars"></i>
        </button>
      </div>

      {/* Mobile Menu */}
      {showMobileMenu && (
        <div
          ref={menuRef}
          className="md:hidden absolute top-16 right-4 bg-white border rounded-lg shadow-lg z-50"
          role="dialog"
          aria-modal="true"
        >
          <div className="space-y-2 p-4">
            {["Home", "Products", "FAQ"].map((link) => (
              <a
                key={link}
                href={link === "FAQ" ? "/#faqTarget" : `/${link.toLowerCase()}`}
                onClick={() => setShowMobileMenu(false)}
                className="block text-gray-600 hover:text-gray-900"
                aria-label={link}
              >
                {link === "Home" ? "What's Trending" : link}
              </a>
            ))}
            {aviLegacyName && (
              <a
                href="/merchant"
                onClick={() => setShowMobileMenu(false)}
                className="block text-gray-600 hover:text-gray-900"
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