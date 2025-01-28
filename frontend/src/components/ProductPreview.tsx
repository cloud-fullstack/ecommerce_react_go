import React from "react";

interface Product {
  product_id: string;
  product_name: string;
  store_id: string;
  picture_link: string;
  price: number;
  discounted_price: number;
  discounted: boolean;
  category?: string; // Optional, for `ProductCard`
  store_name?: string; // Optional, for `ProductCard`
}

interface ProductPreviewProps extends Product {
  demo?: boolean;
  pricing?: boolean;
  index?: number;
  className?: string; // Add className prop
}

const ProductPreview: React.FC<ProductPreviewProps> = ({
  product_name = "Unnamed Product",
  store_id = "",
  product_id = "",
  picture_link = "",
  price = 0,
  discounted_price = 0,
  discounted = false,
  demo = false,
  pricing = true,
  index,
  className = "", // Default to empty string
}) => {
  // Truncate the product name if it's too long
  const truncatedName = product_name.length > 14 ? `${product_name.slice(0, 14)}..` : product_name;

  // Calculate the real price based on whether the discount is active
  const realPrice = discounted ? discounted_price : price;

  // Format the price to 2 decimal places
  const formattedPrice = realPrice.toFixed(2);

  return (
    <div
      className={`productDisplay ${className} flex flex-col items-center justify-center p-2`}
      key={index}
    >
      <a
        href={store_id && product_id ? `/store/${store_id}/${product_id}` : "#"}
        className="flex flex-col items-center"
      >
        {/* Product Image */}
        <img
          src={picture_link}
          alt={product_name}
          className="carouselPic w-full h-48 object-cover rounded-lg bg-gray-200"
          onError={(e) => {
            e.currentTarget.src = ""; // Fallback to a placeholder image
          }}
        />

        {/* Product Details */}
        <div className="nameTitle w-full text-center mt-2">
          {/* Product Name */}
          <span className="productN block font-bold text-sm">{truncatedName}</span>

          {/* Pricing Information */}
          {pricing && (
            <span className="productP block text-green-600 text-sm">
              ${formattedPrice}
            </span>
          )}

          {/* Demo Label */}
          {demo && (
            <span className="demoLabel block text-xs text-blue-500 mt-1">
              Demo Available
            </span>
          )}
        </div>
      </a>
    </div>
  );
};

export default ProductPreview;
