import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";

const Footer = () => {
  const location = useLocation();
  const [email, setEmail] = useState("");

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

  const handleSubscribe = (e: React.FormEvent) => {
    e.preventDefault();
    console.log("Subscribed with email:", email);
    setEmail(""); // Clear the input after submission
  };

  return (
    <footer className="footerGradGray" id="footerId" aria-labelledby="footer-heading">
      <h2 id="footer-heading" className="sr-only">
        Footer
      </h2>
      <div className="footer-container">
        <p>&copy; 2023 Everything SRL. All rights reserved.</p>
        <div className="footer-links">
          <a href="#">Privacy Policy</a>
          <a href="#">Terms of Service</a>
          <a href="#">Contact Us</a>
        </div>
        <form onSubmit={handleSubscribe} className="subscribe-section">
          <input
            type="email"
            placeholder="Enter your email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <button type="submit">Subscribe</button>
        </form>
      </div>
    </footer>
  );
};

export default Footer;