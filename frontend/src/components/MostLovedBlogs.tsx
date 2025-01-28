import { useState, useEffect } from "react";
import Slider from "react-slick";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import useStore from "../stores/useStore";
import { BlogPost } from "../types/types";
import apiClient from '../utils/api';
import Title from "./Title";

const MostLovedBlogs = () => {
  const [mostLovedPictures, setMostLovedPictures] = useState<BlogPost[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { profilePicture } = useStore();

  // Fetch most loved blog posts
  const fetchMostLovedPictures = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await apiClient.get("/api/most-loved-recent-blogs/");
      const data = res.data;

      // Ensure the data matches the BlogPost type
      if (data.error) throw new Error(data.message);

      // Transform the data if necessary
      const transformedData: BlogPost[] = data.map((item: any) => ({
        blog_post_id: item.blog_post_id,
        store_id: item.store_id,
        product_id: item.product_id,
        product_name: item.product_name,
        picture_link: item.picture_link,
        category_id: item.category_id,
        category_name: item.category_name,
        author_id: item.author_id,
        author_name: item.author_name,
      }));

      setMostLovedPictures(transformedData);
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Failed to load most loved blogs. Please try again later.");
      }
    } finally {
      setLoading(false);
    }
  };

  // Fetch data on mount
  useEffect(() => {
    fetchMostLovedPictures();
  }, []);

  // Slider settings
  const settings = {
    dots: false,
    infinite: true,
    speed: 500,
    slidesToShow: 3,
    slidesToScroll: 3,
    autoplay: true,
    autoplaySpeed: 5000,
    pauseOnHover: true,
    responsive: [
      {
        breakpoint: 1024,
        settings: {
          slidesToShow: 2,
          slidesToScroll: 2,
        },
      },
      {
        breakpoint: 768,
        settings: {
          slidesToShow: 1,
          slidesToScroll: 1,
        },
      },
    ],
  };

  // Render loading state
  if (loading) {
    return <p>Loading most recent loved pictures...</p>;
  }

  // Render error state
  if (error) {
    return <p style={{ textAlign: "center", color: "red" }}>{error}</p>;
  }

  return (
    <div className="most-loved-blogs relative px-4 pt-16 pb-20 sm:px-6 lg:px-8 lg:pt-24 lg:pb-28">
      <div className="absolute inset-0 h-1/3 sm:h-2/3 bg-white"></div>
      <div className="relative mx-auto max-w-7xl">
        <div className="text-center">
          <Title>The Most Loved Blogs</Title>
          <p className="mx-auto mt-3 max-w-2xl text-xl text-gray-500 sm:mt-4 animate__animated animate__fadeInRight animate__slow">
            A list of the most popular blogs by users
          </p>
        </div>
        <div className="mt-12">
          {mostLovedPictures.length === 0 ? (
            <p>No blogs available!</p>
          ) : (
            <Slider {...settings}>
              {mostLovedPictures.map((pic: BlogPost, i: number) => (
                <div key={i} className="px-2">
                  <a
                    href={`/store/${pic.store_id}/${pic.product_id}#blog=${pic.blog_post_id}`}
                    aria-label={`View blog post for ${pic.product_name}`}
                  >
                    <div className="flex flex-col overflow-hidden rounded-lg shadow-lg cursor-pointer">
                      <div className="flex-shrink-0">
                        <img
                          className="carouselPic w-full h-48 object-cover" // Use the same class as ProductPreview
                          src={pic.picture_link}
                          alt={`Banner for ${pic.product_name}`}
                        />
                      </div>
                      <div className="flex flex-1 flex-col justify-between bg-white p-6">
                        <div className="flex-1">
                          <p className="text-sm font-medium text-indigo-600">
                            <a
                              href={`/category/${pic.category_id}`}
                              className="hover:underline"
                              aria-label={`View category ${pic.category_name}`}
                            >
                              {pic.category_name}
                            </a>
                          </p>
                          <a
                            href={`/store/${pic.store_id}/${pic.product_id}#blog=${pic.blog_post_id}`}
                            className="mt-2 block"
                            aria-label={`View product ${pic.product_name}`}
                          >
                            <p className="text-xl font-semibold text-gray-900">
                              {pic.product_name}
                            </p>
                          </a>
                        </div>
                        <div className="mt-6 flex items-center">
                          <div className="flex-shrink-0">
                            <a
                              href={`/blog/${pic.blog_post_id}/`}
                              aria-label={`View blog post by ${pic.author_name}`}
                            >
                              <img
                                className="h-10 w-10 rounded-full"
                                src={profilePicture || ""}
                                alt={`Profile picture of ${pic.author_name}`}
                              />
                            </a>
                          </div>
                          <div className="ml-3">
                            <p className="text-sm font-medium text-gray-900">
                              <a
                                href={`/blogs/${pic.author_id}`}
                                className="hover:underline"
                                aria-label={`View profile of ${pic.author_name}`}
                              >
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