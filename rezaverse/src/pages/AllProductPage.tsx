import React, { useEffect, useState } from 'react';
import { Carousel } from 'react-responsive-carousel';
import 'react-responsive-carousel/lib/styles/carousel.min.css';
import ProductPreview from '../components/ProductPreview';
import { FunnelIcon } from '@heroicons/react/24/outline'; // Use the correct icon name
import apiClient from '../utils/api';

interface Product {
  product_id: string;
  store_id: string;
  product_name: string;
  picture_link: string;
  price: number;
  discounted_price?: number;
  discounted?: boolean;
  category: string; // Add this line
}

const Products = () => {
  const [current, setCurrent] = useState("Products");
  const [selected, setSelected] = useState("All");
  const [selectedFilter, setSelectedFilter] = useState("");
  const [searchProductName, setSearchProductName] = useState("");
  const [listOfProducts, setListOfProducts] = useState<Product[]>([]);
  const [showCategory, setShowCategory] = useState(false);
  const [showFilters, setShowFilters] = useState(false);

  const listOfCategory = [
    "All", "Animals", "Animated Objects", "Animations", "Apparel", "Audio and Video", 
    "Avatar Accessories", "Avatar Apparel", "Avatar Components", "Avatar Roleplay", 
    "Breedables", "Building Components", "Buildings and Various Structures", "Business", 
    "Celebrations", "Complete Avatars", "Furry", "Gachas", "Gadgets", "Home and Garden", 
    "Miscellaneous", "Real Estate", "Recreation and Entertainment", "Scripts", "Services", 
    "Vehicles", "Virtual Art", "Weapons"
  ];

  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setShowCategory(false);
        setShowFilters(false);
      }
    };
    document.addEventListener('keyup', handleEscape);
    return () => document.removeEventListener('keyup', handleEscape);
  }, []);

  const fetchProducts = async () => {
    const res = await apiClient.get("/api/frontpage-product-previews/");
    const json = await res.json();
    if (json.error) throw new Error(json.message);
  
    // Ensure `discounted` is always a boolean
    const products = json.map((product: any) => ({
      ...product,
      discounted: product.discounted || false, // Default to `false` if `discounted` is missing
    }));
  
    setListOfProducts(products);
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  const applyFilter = (arr: Product[], filter: string) => {
    switch (filter) {
      case "bestSell":
        return arr.filter(prod => prod.discounted && prod.discounted_price && prod.discounted_price < 10);
      case "lowerToHigh":
        return [...arr].sort((a, b) => (a.discounted_price || a.price) - (b.discounted_price || b.price));
      case "highToLower":
        return [...arr].sort((a, b) => (b.discounted_price || b.price) - (a.discounted_price || a.price));
      case "bestDiscount":
        return arr.filter(prod => prod.discounted).sort((a, b) => (a.discounted_price || 0) - (b.discounted_price || 0));
      default:
        return arr;
    }
  };

  const searchFilterCategory = (arr: Product[], searchTerm: string, filter: string) => {
    const filtered = applyFilter(arr, filter);
    return filtered.filter(order => 
      order.product_name.toLowerCase().includes(searchTerm.toLowerCase()) || 
      order.category === selected // Ensure `category` is used here
    );
  };

  return (
    <div className="bg-white">
      {/* Main Content */}
      <main className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="border-b border-gray-200 pt-16 pb-10">
          <h1 className="text-4xl font-bold tracking-tight text-gray-900">All Products</h1>
          <p className="mt-4 text-base text-gray-500">Explore our wide range of products.</p>
          <div className="mt-6">
            <input
              type="text"
              placeholder="Search products..."
              value={searchProductName}
              onChange={(e) => setSearchProductName(e.target.value)}
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            />
          </div>
        </div>

        <div className="pt-12 pb-24 lg:grid lg:grid-cols-3 lg:gap-x-8 xl:grid-cols-4">
          <aside>
            <h2 className="sr-only">Filters</h2>
            <div className="space-y-10 divide-y divide-gray-200">
              <div>
                <label className="block text-sm font-medium text-gray-700">Filter by</label>
                <select
                  value={selectedFilter}
                  onChange={(e) => setSelectedFilter(e.target.value)}
                  className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
                >
                  <option value="">Select a filter</option>
                  <option value="bestSell">Best selling</option>
                  <option value="lowerToHigh">Price: Low to High</option>
                  <option value="highToLower">Price: High to Low</option>
                  <option value="bestDiscount">Best Discount</option>
                </select>
              </div>
            </div>
          </aside>

          <section className="mt-6 lg:col-span-2 lg:mt-0 xl:col-span-3">
            <div className="grid grid-cols-1 gap-y-4 sm:grid-cols-2 sm:gap-x-6 sm:gap-y-10 lg:gap-x-8 xl:grid-cols-3">
              {listOfProducts.length > 0 ? (
                searchFilterCategory(listOfProducts, searchProductName, selectedFilter).map((product, index) => (
                  <ProductPreview
                    key={product.product_id}
                    productID={product.product_id}
                    storeID={product.store_id}
                    name={product.product_name} // Corrected prop name
                    pictureLink={product.picture_link}
                    price={product.price}
                    discount_price={product.discounted_price || product.price} // Fallback to `price` if `discounted_price` is undefined
                    discounted={product.discounted}
                    index={index} // Add `index` prop
                  />
                ))
              ) : (
                <p>No products available!</p>
              )}
            </div>
          </section>
        </div>
      </main>
    </div>
  );
};

export default Products;