import { useEffect } from "react";
import { useLocation } from "react-router-dom";

const Footer = () => {
  const location = useLocation();

  useEffect(() => {
    const footerElement = document.getElementById("footerId");
    if (!footerElement) return;

    if (location.pathname === "/products/") {
      footerElement.classList.remove("footerGradGray");
      footerElement.classList.add("footerGradWhite");
    }

    if (location.pathname === "/signIn") {
      footerElement.style.display = "none";
    } else {
      footerElement.style.display = "block"; // Ensure footer is visible on other routes
    }
  }, [location.pathname]);

  return (
    <footer
      className="footerGradGray py-8 px-4"
      id="footerId"
      aria-labelledby="footer-heading"
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