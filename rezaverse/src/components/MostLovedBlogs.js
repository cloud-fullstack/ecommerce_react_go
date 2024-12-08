import React from 'react';

function MostLovedBlogs() {
  const blogs = [
    { id: 1, title: 'Producto', author: 'skillsthegamer', image: '/images/blog-1.jpg' },
    { id: 2, title: 'LTB - Daphne Costume (Freya)', author: 'animation', image: '/images/blog-2.jpg' },
    { id: 3, title: 'Test product', author: 'animation', image: '/images/blog-3.jpg' },
  ];

  return (
    <section className="bg-gray-100 py-16">
      <div className="text-center mb-8">
        <h2 className="text-3xl font-bold mb-2">The Most Loved Blogs</h2>
        <p className="text-gray-600">A list of the most popular blogs by users</p>
      </div>
      <div className="relative max-w-6xl mx-auto">
        <button className="absolute left-0 top-1/2 transform -translate-y-1/2 text-3xl">&lt;</button>
        <div className="flex justify-center space-x-8">
          {blogs.map((blog) => (
            <div key={blog.id} className="bg-white rounded-lg p-4 w-1/3">
              <img src={blog.image} alt={blog.title} className="rounded-lg mb-4" width="300" height="200" />
              <p className="text-gray-500 text-sm">Category</p>
              <h3 className="font-bold mb-2">{blog.title}</h3>
              <div className="flex items-center">
                <img src="/images/avatar.jpg" alt="User avatar" className="rounded-full mr-2" width="24" height="24" />
                <span className="text-sm text-gray-600">{blog.author}</span>
              </div>
            </div>
          ))}
        </div>
        <button className="absolute right-0 top-1/2 transform -translate-y-1/2 text-3xl">&gt;</button>
      </div>
    </section>
  );
}

export default MostLovedBlogs;