import React, { useState, useEffect } from "react";
import "animate.css";
import apiClient from '../utils/api';

// Define the Product interface
interface Product {
  product_id: string;
  product_name: string;
  store_id: string;
  picture_link: string;
  price: number;
  discounted_price: number;
  discounted: boolean;
}

const Products: React.FC = () => {
  const [discountedProducts, setDiscountedProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showMore, setShowMore] = useState(false);

  // Fetch discounted products
  const fetchDiscountedProducts = async () => {
    setLoading(true);
    try {
      const res = await apiClient.get("/api/discounted-products-frontpage/", {
        method: "GET",
      });
      const data = res.data;
      if (data.error) throw new Error(data.message);
      setDiscountedProducts(data);
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unknown error occurred");
      }
    } finally {
      setLoading(false);
    }
  };
  // Fetch data on mount
  useEffect(() => {
    fetchDiscountedProducts();
  }, []);

  // Function to show more products
  const showMoreDiscounted = () => {
    setShowMore(true);
  };

  if (loading) return <p>Loading products...</p>;
  if (error) return <p style={{ textAlign: "center" }}>{error}</p>;

  return (
    <div className="DiscountDiv">
      <div className="containerBorder">
        <div className="mx-auto max-w-2xl py-10 px-4 sm:px-6 lg:max-w-7xl lg:px-8">
          <h2 className="titleColor pb-3 font-bold tracking-tight animate__animated animate__fadeInDown animate__delay-0.8s">
            Discounted items
          </h2>

          <div className="mt-6 grid grid-cols-1 gap-y-10 gap-x-6 sm:grid-cols-2 lg:grid-cols-4 xl:gap-x-8">
            {discountedProducts.length > 0 ? (
              discountedProducts.map((item, i) => (
                (showMore || i < 8) && (
                  <div key={item.product_id} className="group relative">
                    <div className="min-h-80 aspect-w-1 aspect-h-1 w-full overflow-hidden rounded-md bg-gray-200 group-hover:opacity-75 lg:aspect-none lg:h-80 shadow-lg">
                      <img
                        src={item.picture_link}
                        alt={item.product_name}
                        className="h-full w-full object-cover object-center lg:h-full lg:w-full"
                      />
                    </div>
                    <div className="mt-4 flex justify-between">
                      <div>
                        <h3 className="text-sm text-gray-100">
                          <a href={`/store/${item.store_id}/${item.product_id}`}>
                            <span aria-hidden="true" className="absolute inset-0"></span>
                            {item.product_name}
                          </a>
                        </h3>
                        <p className="mt-1 text-sm text-gray-300 categoryText">Category</p>
                      </div>
                      <p className="text-sm font-medium priceDiscounted">{item.price} $</p>
                    </div>
                  </div>
                )
              ))
            ) : (
              <p>There are no discounted items right now!</p>
            )}
          </div>

          {discountedProducts.length > 8 && !showMore && (
            <div className="showButton">
              <button
                id="showMoreButtonHide"
                type="button"
                onClick={showMoreDiscounted}
                className="showMoreBtn items-center rounded-md border border-gray-900 px-3 py-2 text-sm font-medium leading-4 shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
              >
                Show more
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Products;