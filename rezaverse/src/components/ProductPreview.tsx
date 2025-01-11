import React from "react";

interface ProductPreviewProps {
  name: string;
  storeID: string;
  productID: string;
  pictureLink: string;
  price: number;
  discount_price: number;
  discounted?: boolean; // `discounted` is optional
  index: number;
  pricing?: boolean;
}

const ProductPreview: React.FC<ProductPreviewProps> = ({
  name,
  storeID,
  productID,
  pictureLink,
  price,
  discount_price,
  discounted = false, // Default to `false` if `discounted` is undefined
  index,
  pricing = true,
}) => {
  // Truncate name if it's too long
  const truncatedName = name.length > 14 ? `${name.slice(0, 14)}..` : name;

  // Calculate real price
  const realPrice = discounted ? discount_price : price;

  return (
    <div className="productDisplay">
      <a href={`/store/${storeID}/${productID}`}>
        <img
          src={pictureLink}
          alt={name}
          className="carouselPic aspect-h-1 w-full overflow-hidden bg-gray-200 group-hover:opacity-75 lg:aspect-none"
        />
        <div className="nameTitle aspect-h-1 w-full overflow-hidden bg-gray-200 group-hover:opacity-75 lg:aspect-none shadow-lg">
          <span className="productN">{truncatedName}</span>
          {pricing && <span className="productP">${realPrice}</span>}
        </div>
      </a>
    </div>
  );
};

export default ProductPreview;