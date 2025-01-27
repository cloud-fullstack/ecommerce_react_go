import React, { useState, useEffect } from "react";
import { Carousel } from "react-responsive-carousel";
import "react-responsive-carousel/lib/styles/carousel.min.css";
import ProductPreview from "./ProductPreview";
import "animate.css";
import apiClient from '../utils/api';

interface Product {
  product_id: string;
  name: string;
  storeID: string;
  pictureLink: string;
  price: number;
  discountedPrice: number;
  discountActive: boolean;
  demo?: boolean;
  pricing?: boolean;
  index?: number;
}

interface DiscountedProductsProps {
  title: string; // Add a title prop
}

const DiscountedProducts: React.FC<DiscountedProductsProps> = ({ title }) => {
  const [discountedProducts, setDiscountedProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [screenWidth, setScreenWidth] = useState(window.innerWidth);
  const [numberOfCart, setNumberOfCart] = useState(4);

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

  // Handle screen size changes
  const handleResize = () => {
    const newScreenWidth = window.innerWidth;
    setScreenWidth(newScreenWidth);

    // Update the number of carousel items based on screen width
    if (newScreenWidth >= 767 && newScreenWidth <= 1280) setNumberOfCart(3);
    else if (newScreenWidth <= 767 && newScreenWidth >= 640) setNumberOfCart(3);
    else if (newScreenWidth <= 639 && newScreenWidth >= 481) setNumberOfCart(2);
    else if (newScreenWidth <= 480) setNumberOfCart(1);
  };

  // Fetch data on mount
  useEffect(() => {
    fetchDiscountedProducts();
  }, []);

  // Add event listener for screen size changes
  useEffect(() => {
    window.addEventListener("resize", handleResize);

    // Cleanup function to remove the event listener
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  // Update number of carousel items when screenWidth changes
  useEffect(() => {
    handleResize();
  }, [screenWidth]);

  if (loading) return <p>Loading Products...</p>;
  if (error) return <p style={{ textAlign: "center" }}>{error}</p>;

  return (
    <div className="carouselDiv">
      <div className="text-center titleSponsorised">
        <h2 className="animate__animated animate__fadeInDown pt-5 font-bold tracking-tight carouTitle">
          {title} {/* Use the title prop here */}
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
                product_id={product.product_id}
                product_name={product.name}
                store_id={product.storeID}
                picture_link={product.pictureLink}
                price={product.price}
                discounted_price={product.discountedPrice}
                discounted={product.discountActive}
                demo={product.demo}
                pricing={product.pricing}
                index={i}
                className="custom-class-name" // Pass the className prop here
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