import React from "react";

interface ProductPreviewProps {
  name?: string;
  storeID?: string;
  productID?: string;
  pictureLink?: string;
  price?: number;
  discountedPrice?: number;
  discountActive?: boolean;
  demo?: boolean;
  pricing?: boolean;
  index?: number;
  className?: string; // Add className prop
}

const ProductPreview: React.FC<ProductPreviewProps> = ({
  name = "Unnamed Product",
  storeID = "",
  productID = "",
  pictureLink = "https://via.placeholder.com/150",
  price = 0,
  discountedPrice = 0,
  discountActive = false,
  demo = false,
  pricing = true,
  index,
  className = "", // Default to empty string
}) => {
  const truncatedName = name.length > 14 ? `${name.slice(0, 14)}..` : name;
  const realPrice = discountActive ? discountedPrice : price;
  const formattedPrice = realPrice.toFixed(2);

  return (
    <div className="productDisplay" key={index}>
      <a href={storeID && productID ? `/store/${storeID}/${productID}` : "#"}>
        <img
          src={pictureLink}
          alt={name}
          className={`carouselPic ${className} object-cover`} // Apply className prop
          onError={(e) => {
            e.currentTarget.src = "https://via.placeholder.com/150";
          }}
        />
        <div className="nameTitle mt-2 p-2 bg-white shadow-lg">
          <span className="productN block text-sm font-semibold">{truncatedName}</span>
          {pricing && <span className="productP block text-sm text-gray-600">${formattedPrice}</span>}
          {demo && <span className="demoLabel block text-xs text-blue-600">Demo Available</span>}
        </div>
      </a>
    </div>
  );
};

export default ProductPreview;