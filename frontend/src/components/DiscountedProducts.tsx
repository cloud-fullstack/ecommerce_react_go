import { useState, useEffect } from "react";
import { Carousel } from "react-responsive-carousel";
import "react-responsive-carousel/lib/styles/carousel.min.css";
import ProductPreview from "./ProductPreview";
import "animate.css";
import apiClient from '../utils/api';
import Title from "./Title";

interface Product {
  product_id: string;
  product_name: string;
  store_id: string;
  picture_link: string;
  price: number;
  discounted_price: number;
  discounted: boolean;
  demo?: boolean;
  pricing?: boolean;
}

interface DiscountedProductsProps {
  title: string;
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

      // Map the API response to match the `Product` interface
      const formattedProducts = data.map((product: any) => ({
        product_id: product.product_id,
        product_name: product.product_name,
        store_id: product.store_id,
        picture_link: product.picture_link,
        price: product.price,
        discounted_price: product.discounted_price,
        discounted: product.discounted,
        demo: product.demo,
        pricing: product.pricing,
      }));

      setDiscountedProducts(formattedProducts);
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
          {title}
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
          additionalTransitions={0} // Add this line to reduce transition effect
        >
          {discountedProducts.map((product, i) => (
            <div key={product.product_id} className="product">
              <ProductPreview
                product_id={product.product_id}
                product_name={product.product_name}
                store_id={product.store_id}
                picture_link={product.picture_link}
                price={product.price}
                discounted_price={product.discounted_price}
                discounted={product.discounted}
                demo={product.demo}
                pricing={product.pricing}
                index={i}
                className="w-full h-48 sm:h-48 object-cover" // Ensure images cover the area
              />
            </div>
          ))}
        </Carousel>
      ) : (
        <p>No products available!</p>
      )}

      <div className="showArrow">
        <button onClick={() => window.scrollTo(0, 0)}>
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
        </button>
      </div>
    </div>
  );
};

export default DiscountedProducts;
