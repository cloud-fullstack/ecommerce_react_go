import React, { useState, useEffect } from "react";
import { Carousel } from "react-responsive-carousel";
import "react-responsive-carousel/lib/styles/carousel.min.css";
import ProductPreview from "./ProductPreview";
import "animate.css";
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

const DiscountedProducts = () => {
  const [discountedProducts, setDiscountedProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [screenWidth, setScreenWidth] = useState(window.innerWidth);
  const [numberOfCart, setNumberOfCart] = useState(4);

  // Fetch discounted products
  const fetchDiscountedProducts = async () => {
    try {
      const res = await apiClient.get("/api/discounted-products-frontpage/", {
        method: "GET",
      });
      const data = res.data;
      if (data.error) throw new Error(data.message);
      setDiscountedProducts(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  // Handle screen size changes
  const handleResize = () => {
    setScreenWidth(window.innerWidth);
    if (screenWidth >= 767 && screenWidth <= 1280) setNumberOfCart(3);
    else if (screenWidth <= 767 && screenWidth >= 640) setNumberOfCart(3);
    else if (screenWidth <= 639 && screenWidth >= 481) setNumberOfCart(2);
    else if (screenWidth <= 480) setNumberOfCart(1);
  };

  // Fetch data on mount
  useEffect(() => {
    fetchDiscountedProducts();
  }, []);

  // Add event listener for screen size changes
  useEffect(() => {
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, [screenWidth]);

  if (loading) return <p>Loading Products...</p>;
  if (error) return <p style={{ textAlign: "center" }}>{error}</p>;

  return (
    <div className="carouselDiv">
      <div className="text-center titleSponsorised">
        <h2 className="animate__animated animate__fadeInDown pt-5 font-bold tracking-tight carouTitle">
          Hot sales & Discount
        </h2>
      </div>

      {discountedProducts.length > 0 ? (
        <Carousel
          showThumbs={false}
          infiniteLoop
          autoPlay
          interval={5000}
          showArrows={true}
          showStatus={false}
          showIndicators={false}
          centerMode
          centerSlidePercentage={100 / numberOfCart}
        >
          {discountedProducts.map((product, i) => (
            <div key={product.product_id} className="product">
              <ProductPreview
                name={product.product_name}
                storeID={product.store_id}
                productID={product.product_id}
                pictureLink={product.picture_link}
                price={product.price}
                discount_price={product.discounted_price}
                discounted={product.discounted}
                index={i}
              />
            </div>
          ))}
        </Carousel>
      ) : (
        <p>No products available!</p>
      )}

      <div className="showArrow">
        <a href="/#blogsTarget">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth="1.5"
            stroke="currentColor"
            className="w-12 h-12 animate__animated animate__infinite animate__bounce animate__slow text-white"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M19.5 5.25l-7.5 7.5-7.5-7.5m15 6l-7.5 7.5-7.5-7.5"
            />
          </svg>
        </a>
      </div>
    </div>
  );
};

export default DiscountedProducts;