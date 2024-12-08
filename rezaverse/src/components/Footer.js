import React from 'react';

function Footer() {
  return (
    <footer className="py-16 px-8">
      <div className="max-w-6xl mx-auto grid grid-cols-5 gap-8">
        <div>
          <img src="/images/logo.png" alt="The Rezaverse logo" className="mb-4" width="120" height="40" />
        </div>
        {['Support', 'Company', 'Legal'].map((title, index) => (
          <div key={index}>
            <h4 className="font-bold mb-4">{title}</h4>
            <ul className="space-y-2">
              {['Status', 'Documentation', 'Guides', 'API Status', 'Support Ticket'].map((item, i) => (
                <li key={i}><a href="#" className="text-gray-600">{item}</a></li>
              ))}
            </ul>
          </div>
        ))}
        <div>
          <h4 className="font-bold mb-4">Subscribe to our newsletter</h4>
          <p className="text-gray-600 mb-4">The latest news, articles, and resources, sent to your inbox weekly.</p>
          <div className="flex">
            <input type="email" placeholder="Enter your email" className="flex-grow p-2 border rounded-l" />
            <button className="bg-yellow-400 px-4 py-2 rounded-r">Subscribe</button>
          </div>
        </div>
      </div>
      <div className="mt-8 text-center text-gray-600">
        <p>© 2023 Emptytute SRL. All rights reserved.</p>
      </div>
    </footer>
  );
}

export default Footer;