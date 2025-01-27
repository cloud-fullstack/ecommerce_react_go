import { useEffect, useState } from "react";

const Footer = () => {
  const [currentPath, setCurrentPath] = useState("");

  // Update the footer based on the current path
  useEffect(() => {
    const path = window.location.pathname;
    setCurrentPath(path);

    const footerElement = document.getElementById("footerId");
    if (!footerElement) return;

    // Update footer styles based on the current path
    if (path === "/products/") {
      footerElement.classList.remove("footerGradGray");
      footerElement.classList.add("footerGradWhite");
    }

    if (path === "/signIn") {
      footerElement.style.display = "none";
    }
  }, []); // Empty dependency array ensures this runs only once on mount

  return (
    <footer
      className="footerGradGray py-8 px-4" // Use the CSS class from index.css
      id="footerId"
      aria-labelledby="footer-heading"
      style={{ display: "inline" }}
    >
      <h2 id="footer-heading" className="sr-only">
        Footer
      </h2>
      <div className="flex justify-between">
        <p>&copy; 2023 Everything SRL. All rights reserved.</p>
        <div className="flex space-x-4">
          <a href="#" className="text-white hover:underline">
            Privacy Policy
          </a>
          <a href="#" className="text-white hover:underline">
            Terms of Service
          </a>
          <a href="#" className="text-white hover:underline">
            Contact Us
          </a>
        </div>
      </div>
    </footer>
  );
};

export default Footer;