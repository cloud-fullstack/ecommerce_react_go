import React from "react";

interface ProductPreviewProps {
  name?: string; // Make optional
  storeID?: string; // Make optional
  productID?: string; // Make optional
  pictureLink?: string; // Make optional
  price?: number; // Make optional
  discountedPrice?: number; // Make optional
  discountActive?: boolean; // Make optional
  demo?: boolean; // Make optional
  pricing?: boolean; // Optional prop
}

const ProductPreview: React.FC<ProductPreviewProps> = ({
  name = "Unnamed Product", // Default value
  storeID = "", // Default value
  productID = "", // Default value
  pictureLink = "https://via.placeholder.com/150", // Default placeholder image
  price = 0, // Default value
  discountedPrice = 0, // Default value
  discountActive = false, // Default value
  demo = false, // Default value
  pricing = true, // Default to `true` if `pricing` is undefined
}) => {
  // Truncate name if it's too long
  const truncatedName = name.length > 14 ? `${name.slice(0, 14)}..` : name;

  // Calculate real price
  const realPrice = discountActive ? discountedPrice : price;

  return (
    <div className="productDisplay">
      <a href={storeID && productID ? `/store/${storeID}/${productID}` : "#"}>
        <img
          src={pictureLink}
          alt={name}
          className="carouselPic aspect-h-1 w-full overflow-hidden bg-gray-200 group-hover:opacity-75 lg:aspect-none"
        />
        <div className="nameTitle aspect-h-1 w-full overflow-hidden bg-gray-200 group-hover:opacity-75 lg:aspect-none shadow-lg">
          <span className="productN">{truncatedName}</span>
          {pricing && <span className="productP">${realPrice}</span>}
          {demo && <span className="demoLabel">Demo Available</span>}
        </div>
      </a>
    </div>
  );
};

export default ProductPreview;