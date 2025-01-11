import React, { useState, useEffect } from "react";
import Slider from "react-slick";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import useStore from "../stores/useStore"; // Import the Zustand store
import { BlogPost } from "../types/types"; // Import the BlogPost type

const MostLovedBlogs = () => {
  const [mostLovedPictures, setMostLovedPictures] = useState<BlogPost[]>([]);
  const { profilePicture } = useStore(); // Access the Zustand store

  const fetchMostLovedPictures = async () => {
    try {
      const res = await fetch("__API_URL__/api/most-loved-recent-blogs/", {
        method: "GET",
      });
      const data = await res.json();
      if (data.error) throw new Error(data.message);
      setMostLovedPictures(data);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchMostLovedPictures();
  }, []);

  const settings = {
    dots: false,
    infinite: true,
    speed: 500,
    slidesToShow: 3,
    slidesToScroll: 3,
    autoplay: true,
    autoplaySpeed: 5000,
    pauseOnHover: true,
  };

  return (
    <div className="relative px-4 pt-16 pb-20 sm:px-6 lg:px-8 lg:pt-24 lg:pb-28 bg-[#4f2236]">
      <div className="absolute inset-0 h-1/3 sm:h-2/3 bg-white"></div>
      <div className="relative mx-auto max-w-7xl">
        <div className="text-center">
          <h2 className="text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl animate__animated animate__fadeInDown">
            The most loved blogs
          </h2>
          <p className="mx-auto mt-3 max-w-2xl text-xl text-gray-500 sm:mt-4 animate__animated animate__fadeInRight animate__slow">
            A list of the most popular blogs by users
          </p>
        </div>
        <div className="mt-12">
          {mostLovedPictures.length === 0 ? (
            <p>Loading most recent loved pictures...</p>
          ) : (
            <Slider {...settings}>
              {mostLovedPictures.map((pic, i) => (
                <div key={i} className="px-2">
                  <a href={`/store/${pic.store_id}/${pic.product_id}#blog=${pic.blog_post_id}`}>
                    <div className="flex flex-col overflow-hidden rounded-lg shadow-lg cursor-pointer">
                      <div className="flex-shrink-0">
                        <img
                          className="h-48 w-full object-cover"
                          src={pic.picture_link}
                          alt="bannerLovedBlog"
                        />
                      </div>
                      <div className="flex flex-1 flex-col justify-between bg-white p-6">
                        <div className="flex-1">
                          <p className="text-sm font-medium text-indigo-600">
                            <a href="#" className="hover:underline">
                              Category
                            </a>
                          </p>
                          <a href="#" className="mt-2 block">
                            <p className="text-xl font-semibold text-gray-900">
                              {pic.product_name}
                            </p>
                          </a>
                        </div>
                        <div className="mt-6 flex items-center">
                          <div className="flex-shrink-0">
                            <a href={`/blog/${pic.blog_post_id}/`}>
                              <img
                                className="h-10 w-10 rounded-full"
                                src={profilePicture}
                                alt={pic.author_name}
                              />
                            </a>
                          </div>
                          <div className="ml-3">
                            <p className="text-sm font-medium text-gray-900">
                              <a href={`/blogs/${pic.author_id}`} className="hover:underline">
                                {pic.author_name}
                              </a>
                            </p>
                          </div>
                        </div>
                      </div>
                    </div>
                  </a>
                </div>
              ))}
            </Slider>
          )}
        </div>
      </div>
    </div>
  );
};

export default MostLovedBlogs;