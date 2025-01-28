import React from 'react';

interface ProductCardProps {
  category: string;
  discounted: boolean;
  discounted_price: number;
  picture_link: string;
  price: number;
  product_id: string;
  product_name: string;
  store_id: string;
  store_name: string;
}

const ProductCard: React.FC<ProductCardProps> = ({
  category,
  discounted,
  discounted_price,
  picture_link,
  price,
  product_id,
  product_name,
  store_id,
  store_name,
}) => {
  return (
    <div className="product-card">
      <img src={picture_link} alt={product_name} className="product-image" />
      <div className="product-info">
        <h3 className="product-name">{product_name}</h3>
        <p className="product-price">
          ${price.toFixed(2)}
          {discounted && <span>Discounted</span>}
        </p>
        <p>
          <a href={`/product-details/${product_id}`} className="view-details">View Details</a>
        </p>
      </div>
    </div>
  );
};

export default ProductCard;