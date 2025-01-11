import React, { useState, useEffect } from "react";
import Slider from "react-slick";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import ProductPreview from "../components/ProductPreview";
import DiscountedProducts from "../components/DiscountedProducts";
import MostLovedBlogs from "../components/MostLovedBlogs";
import Faq from "../components/Faq";
import useStore from "../stores/useStore"; // Zustand store for global state
import "animate.css";

const Home = () => {
  const [products, setProducts] = useState([]);
  const [animationLoaded, setAnimationLoaded] = useState(false);
  const [screenWidth, setScreenWidth] = useState(window.innerWidth);
  const [slidesToShow, setSlidesToShow] = useState(3);
  const { authToken } = useStore();

  const fetchProductPreviews = async () => {
    try {
      const res = await fetch("__API_URL__/api/frontpage-product-previews/", {
        method: "GET",
      });
      const data = await res.json();
      if (data.error) throw new Error(data.message);
      setProducts(data);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchProductPreviews();
  }, []);

  useEffect(() => {
    const handleResize = () => {
      setScreenWidth(window.innerWidth);
      if (screenWidth >= 1280) setSlidesToShow(3);
      else if (screenWidth >= 1024) setSlidesToShow(2);
      else if (screenWidth >= 768) setSlidesToShow(2);
      else if (screenWidth >= 640) setSlidesToShow(1);
      else setSlidesToShow(1);
    };

    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, [screenWidth]);

  useEffect(() => {
    const currentPath = window.location.pathname;
    const homeDiv = document.getElementById("homeDiv");

    if (currentPath === "/" && !animationLoaded) {
      setTimeout(() => {
        if (homeDiv) {
          homeDiv.classList.add("animate__fadeIn");
          homeDiv.style.display = "block";
        }
        setAnimationLoaded(true);
      }, 4000);
    } else if (homeDiv) {
      homeDiv.classList.add("animate__fadeIn");
      homeDiv.style.display = "block";
    }
  }, [animationLoaded]);

  const settings = {
    dots: false,
    infinite: true,
    speed: 500,
    slidesToShow: slidesToShow,
    slidesToScroll: slidesToShow,
    autoplay: true,
    autoplaySpeed: 5000,
    pauseOnHover: true,
  };

  return (
    <div>
      <div className="homePage" style={{ display: "none" }} id="homeDiv">
        <div className="carouselDiv">
          <div className="backgroundDiv">
            <div className="tl">
              <img
                className="imgCurtainL"
                alt="leftCurtain"
                src="https://res.cloudinary.com/dutkkgjm5/image/upload/v1666010618/d7bursci4kn3j438c2wj.png"
              />
            </div>
            <div className="tr">
              <img
                className="imgCurtainR"
                alt="rightCurtain"
                src="https://res.cloudinary.com/dutkkgjm5/image/upload/v1666010618/qr6uyye70x86cl5cpjos.png"
              />
            </div>
          </div>

          <div className="text-center titleSponsorised">
            <h2 className="animate__animated animate__fadeInDown pt-5 font-bold tracking-tight sm:text-4xl carouTitle">
              Our personal selection
            </h2>
          </div>

          {products.length > 0 ? (
            <Slider {...settings}>
              {products.map((product, i) => (
                <div key={i} className="product">
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
            </Slider>
          ) : (
            <p>No products available!</p>
          )}

          <div className="showButton">
            <button
              type="button"
              className="animate__animated animate__fadeInRight shadow-xl showMoreBtnBig items-center rounded-md border border-gray-900 mx-4 px-3 py-2 text-sm leading-4 text-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
            >
              <a href="/products" className="colorLinkBtn">
                See all products
              </a>
            </button>
          </div>

          <div className="showArrow">
            <button
              type="button"
              className="animate__animated animate__fadeInRight showMoreBtn items-center rounded-md border-none lg:mr-20 md:mr-10 px-3 py-2 text-sm leading-4 text-700"
            >
              <a href="#discountTarget" className="whiteLink">
                BEST DISCOUNTS
              </a>
            </button>
            <a href="#discountTarget">
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
            <button
              type="button"
              className="animate__animated animate__fadeInRight showMoreBtn items-center rounded-md border-none lg:ml-20 md:ml-5 px-3 py-2 text-sm text-700"
            >
              <a href="#blogsTarget" className="whiteLink">
                BEST BLOGS
              </a>
            </button>
          </div>
        </div>

        <div id="discountTarget">
          <DiscountedProducts />
        </div>

        <div id="blogsTarget">
          <MostLovedBlogs />
        </div>

        <div id="faqTarget">
          <Faq />
        </div>
      </div>
    </div>
  );
};

export default Home;