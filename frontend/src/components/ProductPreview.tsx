import React from "react";

interface ProductPreviewProps {
  name?: string;
  storeID?: string;
  productID?: string;
  pictureLink?: string;
  price?: number;
  discountedPrice?: number; // Expected prop name
  discountActive?: boolean; // Expected prop name
  demo?: boolean;
  pricing?: boolean;
  index?: number;
}

const ProductPreview: React.FC<ProductPreviewProps> = ({
  name = "Unnamed Product", // Default product name
  storeID = "", // Default store ID
  productID = "", // Default product ID
  pictureLink = "https://via.placeholder.com/150", // Default placeholder image
  price = 0, // Default price
  discountedPrice = 0, // Default discounted price
  discountActive = false, // Default discount status
  demo = false, // Default demo label visibility
  pricing = true, // Default to showing pricing
  index, // Optional index
}) => {
  // Truncate the product name if it's too long
  const truncatedName = name.length > 14 ? `${name.slice(0, 14)}..` : name;

  // Calculate the real price based on whether the discount is active
  const realPrice = discountActive ? discountedPrice : price;

  // Format the price to 2 decimal places
  const formattedPrice = realPrice.toFixed(2);

  return (
    <div className="productDisplay" key={index}>
      <a href={storeID && productID ? `/store/${storeID}/${productID}` : "#"}>
        {/* Product Image */}
        <img
          src={pictureLink}
          alt={name}
          className="carouselPic aspect-h-1 w-full overflow-hidden bg-gray-200 group-hover:opacity-75 lg:aspect-none"
          onError={(e) => {
            // Fallback to a placeholder image if the provided image fails to load
            e.currentTarget.src = "https://via.placeholder.com/150";
          }}
        />

        {/* Product Details */}
        <div className="nameTitle aspect-h-1 w-full overflow-hidden bg-gray-200 group-hover:opacity-75 lg:aspect-none shadow-lg">
          {/* Product Name */}
          <span className="productN">{truncatedName}</span>

          {/* Pricing Information */}
          {pricing && <span className="productP">${formattedPrice}</span>}

          {/* Demo Label */}
          {demo && <span className="demoLabel">Demo Available</span>}
        </div>
      </a>
    </div>
  );
};

export default ProductPreview;