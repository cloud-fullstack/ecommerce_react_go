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
    <div className={`productDisplay ${className}`} key={index}>
      <a href={store_id && product_id ? `/store/${store_id}/${product_id}` : "#"}>
        {/* Product Image */}
        <img
          src={picture_link}
          alt={product_name}
          className="carouselPic aspect-h-1 w-full overflow-hidden bg-gray-200 group-hover:opacity-75 lg:aspect-none"
          onError={(e) => {
            e.currentTarget.src = "https://via.placeholder.com/150"; // Fallback to a placeholder image
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