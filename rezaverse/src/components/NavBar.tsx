import React, { useState, useEffect, useRef } from "react";
import QuickMenu from "./QuickMenu";
import Cart from "../pages/Cart";
import { BlogPost } from "../types/types"; // Import the BlogPost type
import useStore from "../stores/useStore"; // Zustand store for global state
import apiClient from '../utils/api';

// Define the Product type
interface Product {
  product_id: string;
  product_name: string;
  store_id: string;
  picture_link: string;
  store_name: string;
}

const NavBar = () => {
  const [current, setCurrent] = useState("");
  const [show, setShow] = useState(false); // Mobile menu state
  const [searchProductName, setSearchProductName] = useState("");
  const [searchedProducts, setSearchedProducts] = useState<Product[]>([]); // Use the Product type
  const { animationLoaded, aviLegacyName, setAnimationLoaded } = useStore(); // Destructure setAnimationLoaded
  const menuRef = useRef<HTMLElement | null>(null); // Ref for mobile menu with explicit type

  // Fetch product previews
  const fetchProductPreviews = async () => {
    try {
      const res = await apiClient.get("/api/frontpage-product-previews/", {
        method: "GET",
      });
      const data = res.data;
      if (data.error) throw new Error(data.message);
      setSearchedProducts(data);
    } catch (err) {
      console.error(err);
    }
  };

  // Handle search filter
  const searchNavbarFilter = (arr: Product[], searchTerm: string) => {
    searchTerm = searchTerm.toLowerCase();
    return arr.filter((product) => {
      return (
        product.product_name.toLowerCase() === searchTerm ||
        product.product_name.includes(searchTerm)
      );
    });
  };

  // Handle mobile menu outside click
  const handleOutsideClick = (event: MouseEvent) => {
    if (show && !menuRef.current?.contains(event.target as Node)) {
      setShow(false);
    }
  };

  // Handle escape key press
  const handleEscape = (event: KeyboardEvent) => {
    if (show && event.key === "Escape") {
      setShow(false);
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
  }, [show]);

  // Fetch data on mount
  useEffect(() => {
    fetchProductPreviews();
  }, []);

  // Handle animation on mount
  useEffect(() => {
    const currentPath = window.location.pathname;
    const navElement = document.getElementById("navv");

    if (currentPath !== "/" || animationLoaded) {
      if (navElement) {
        navElement.style.display = "";
      }
    } else {
      setTimeout(() => {
        if (navElement) {
          navElement.classList.add("animate__fadeIn");
          navElement.style.display = "";
        }
        setAnimationLoaded(true); // Use setAnimationLoaded from the store
      }, 4000);
    }

    if (currentPath === "/signIn" && navElement) {
      navElement.style.display = "none";
    }

    // Highlight active link
    if (currentPath === "/") {
      document.getElementById("HL")?.classList.add("navLink");
    } else if (currentPath === "/products") {
      document.getElementById("PL")?.classList.add("navLink");
    }
  }, [animationLoaded, setAnimationLoaded]); // Add setAnimationLoaded to dependencies

  return (
    <nav className="bgLight shadow sticky top-0 z-50" style={{ display: "none" }} id="navv">
      <div className="mx-auto max-w-7xl px-4 md:px-6 lg:px-8">
        <div className="flex h-16 justify-between">
          <div className="flex">
            <div className="flex flex-shrink-0 items-center">
              <a href="/">
                <img
                  className="block h-8 w-auto lg:hidden"
                  src="https://res.cloudinary.com/dutkkgjm5/image/upload/v1667932486/u05sslpwv1g621uqocxu.png"
                  alt="Your Company"
                />
                <img
                  className="hidden h-12 w-auto lg:block"
                  src="https://res.cloudinary.com/dutkkgjm5/image/upload/v1667932486/u05sslpwv1g621uqocxu.png"
                  alt="Your Company"
                />
              </a>

              {/* Search for mobile */}
              <div className="flex flex-1 items-center justify-center px-2 md:hidden lg:ml-6 lg:justify-end searchBarMobile">
                <div className="w-full max-w-lg lg:max-w-xs">
                  <label htmlFor="search" className="sr-only">
                    Search
                  </label>
                  <div className="relative inputSearchMobile">
                    <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                      <svg
                        className="h-5 w-5 text-gray-400"
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                        aria-hidden="true"
                      >
                        <path
                          fillRule="evenodd"
                          d="M9 3.5a5.5 5.5 0 100 11 5.5 5.5 0 000-11zM2 9a7 7 0 1112.452 4.391l3.328 3.329a.75.75 0 11-1.06 1.06l-3.329-3.328A7 7 0 012 9z"
                          clipRule="evenodd"
                        />
                      </svg>
                    </div>
                    <input
                      id="search"
                      value={searchProductName}
                      onChange={(e) => setSearchProductName(e.target.value)}
                      className="block w-full rounded-md border border-gray-300 bg-white py-2 pl-10 pr-3 leading-5 placeholder-gray-500 focus:border-indigo-500 focus:placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-indigo-500 sm:text-sm"
                      placeholder="Search"
                      type="search"
                    />
                  </div>

                  {searchProductName && (
                    <div style={{ overflowY: "scroll" }}>
                      <ul
                        role="list"
                        className="scrolledSearch w-50 h-96 divide-y divide-gray-200 absolute overflow-y-scroll overflow-y-hidden"
                      >
                        {searchNavbarFilter(searchedProducts, searchProductName).map((product) => (
                          <a key={product.product_id} href={`/store/${product.store_id}/${product.product_id}`}>
                            <li className="flex py-4 navbarList bg-white">
                              <img
                                className="h-10 w-10 ml-2 rounded"
                                src={product.picture_link}
                                alt={product.product_name}
                              />
                              <div className="ml-3">
                                <p className="text-sm font-medium text-gray-900">{product.product_name}</p>
                                <p className="text-xs italic text-gray-500">{product.store_name}</p>
                              </div>
                            </li>
                          </a>
                        ))}
                      </ul>
                    </div>
                  )}
                </div>
              </div>

              <div className="md:hidden openMobileCartBtn">
                <Cart />
              </div>
              <div className="md:hidden openMobileProfileBtn">
                <QuickMenu />
              </div>
              <div>
                <button
                  onClick={() => setShow(!show)}
                  className="md:hidden text-white openMobileNavBtn"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth="1.5"
                    stroke="currentColor"
                    className="w-9 h-9 text-white"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
                    />
                  </svg>
                </button>
              </div>
            </div>

            {/* Desktop Links */}
            <div className="hidden md:ml-6 md:flex md:space-x-8">
              <a
                id="HL"
                href="/"
                onClick={() => setCurrent("Home")}
                className={`inline-flex items-center border-b-2 px-1 pt-1 text-sm font-medium fontLinkNav ${
                  current === "Home" ? "border-500 navLink" : "border-transparent navLinkColor"
                }`}
              >
                What's Trending
              </a>
              <a
                id="PL"
                href="/products"
                onClick={() => setCurrent("Products")}
                className={`inline-flex items-center border-b-2 px-1 pt-1 text-sm font-medium fontLinkNav ${
                  current === "Products" ? "border-500 navLink" : "border-transparent navLinkColor"
                }`}
              >
                Products
              </a>
              <a
                href="/#faqTarget"
                className="inline-flex items-center border-b-2 border-transparent px-1 pt-1 text-sm font-medium navLinkColor fontLinkNav"
              >
                FAQ
              </a>
              {aviLegacyName && (
                <a
                  href="/merchant"
                  className="inline-flex items-center border-b-2 border-transparent px-1 pt-1 text-sm font-medium navLinkColor fontLinkNav"
                >
                  Creator dashboard
                </a>
              )}
            </div>
          </div>

          {/* Desktop Search and Cart */}
          <div className="hidden md:ml-6 md:flex md:items-center">
            <div className="flex flex-1 items-center justify-center px-2 lg:ml-6 lg:justify-end">
              <div className="w-full max-w-lg lg:max-w-xs">
                <label htmlFor="search" className="sr-only">
                  Search
                </label>
                <div className="relative">
                  <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                    <svg
                      className="h-5 w-5 text-gray-400"
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                      aria-hidden="true"
                    >
                      <path
                        fillRule="evenodd"
                        d="M9 3.5a5.5 5.5 0 100 11 5.5 5.5 0 000-11zM2 9a7 7 0 1112.452 4.391l3.328 3.329a.75.75 0 11-1.06 1.06l-3.329-3.328A7 7 0 012 9z"
                        clipRule="evenodd"
                      />
                    </svg>
                  </div>
                  <input
                    id="search"
                    value={searchProductName}
                    onChange={(e) => setSearchProductName(e.target.value)}
                    className="block w-full rounded-md border border-gray-300 bg-white py-2 pl-10 pr-3 leading-5 placeholder-gray-500 focus:border-indigo-500 focus:placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-indigo-500 sm:text-sm"
                    placeholder="Search"
                    type="search"
                  />
                </div>

                {searchProductName && (
                  <div style={{ overflowY: "scroll" }}>
                    <ul
                      role="list"
                      className="scrolledSearch w-50 h-96 divide-y divide-gray-200 absolute overflow-y-scroll overflow-y-hidden"
                    >
                      {searchNavbarFilter(searchedProducts, searchProductName).map((product) => (
                        <a key={product.product_id} href={`/store/${product.store_id}/${product.product_id}`}>
                          <li className="flex py-4 navbarList bg-white">
                            <img
                              className="h-10 w-10 ml-2 rounded"
                              src={product.picture_link}
                              alt={product.product_name}
                            />
                            <div className="ml-3">
                              <p className="text-sm font-medium text-gray-900">{product.product_name}</p>
                              <p className="text-xs italic text-gray-500">{product.store_name}</p>
                            </div>
                          </li>
                        </a>
                      ))}
                    </ul>
                  </div>
                )}
              </div>
            </div>

            <Cart />
            <QuickMenu />
          </div>
        </div>
      </div>

      {/* Mobile Menu */}
      {show && (
        <div ref={menuRef as React.RefObject<HTMLDivElement>} className="md:hidden absolute bg-black mobileMenuDiv" id="mobile-menu">
          <div className="space-y-1 pt-2 pb-3">
            <a
              href="/"
              className="block border-l-4 py-2 pl-3 pr-4 text-base font-medium mobileLink"
            >
              What's Trending
            </a>
            <a
              href="/products"
              className="block border-l-4 border-transparent py-2 pl-3 pr-4 text-base font-medium mobileLink hover:border-gray-300 hover:bg-gray-50 hover:text-gray-700"
            >
              Products
            </a>
            <a
              href="/#faqTarget"
              className="block border-l-4 border-transparent py-2 pl-3 pr-4 text-base font-medium mobileLink hover:border-gray-300 hover:bg-gray-50 hover:text-gray-700"
            >
              FAQ
            </a>
            {aviLegacyName && (
              <a
                href="/merchant"
                className="block border-l-4 border-transparent py-2 pl-3 pr-4 text-base font-medium mobileLink hover:border-gray-300 hover:bg-gray-50 hover:text-gray-700"
              >
                Creator dashboard
              </a>
            )}
          </div>
        </div>
      )}
    </nav>
  );
};

export default NavBar;