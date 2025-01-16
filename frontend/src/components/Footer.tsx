import React, { useEffect, useState } from "react";

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
      className="footerGradGray"
      id="footerId"
      aria-labelledby="footer-heading"
      style={{ display: "inline" }}
    >
      <h2 id="footer-heading" className="sr-only">
        Footer
      </h2>
      <div className="mx-auto max-w-7xl py-12 px-4 sm:px-6 lg:py-16 lg:px-8">
        <div className="xl:grid xl:grid-cols-3 xl:gap-8">
          {/* Footer content goes here */}
        </div>
      </div>
    </footer>
  );
};

export default Footer;