import React from "react";

interface Product {
  imageUrl: string;
  altText: string;
  name: string;
  price: number;
}

interface ProductButtonProps {
  product: Product;
}

const ProductButton: React.FC<ProductButtonProps> = ({ product }) => {
  return (
    <button type="button" className="group">
      <div className="aspect-w-1 aspect-h-1 w-full overflow-hidden rounded-lg bg-gray-200 xl:aspect-w-7 xl:aspect-h-8">
        <img
          src={product.imageUrl}
          alt={product.altText}
          className="h-full w-full object-cover object-center group-hover:opacity-75"
        />
      </div>
      <h3 className="mt-4 text-sm text-gray-700">{product.name}</h3>
      <p className="mt-1 text-lg font-medium text-gray-900">${product.price}</p>
    </button>
  );
};

export default ProductButton;