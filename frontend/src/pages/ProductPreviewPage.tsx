import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import ProductPreview from "../components/ProductPreview";
import apiClient from "../utils/api";

const ProductPreviewPage: React.FC = () => {
  const { storeID, productID } = useParams<{ storeID: string; productID: string }>();
  const [product, setProduct] = useState<{
    product_id: string;
    product_name: string;
    store_id: string;
    picture_link: string;
    price: number;
    discounted_price: number;
    discounted: boolean;
    demo?: boolean;
    pricing?: boolean;
  } | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Fetch product data based on storeID and productID
  useEffect(() => {
    const fetchProduct = async () => {
      setLoading(true);
      try {
        const res = await apiClient.get(`/api/store/${storeID}/${productID}`);
        const data = res.data;

        if (data.error) throw new Error(data.message);

        // Map the API response to match the `Product` interface
        const formattedProduct = {
          product_id: data.product_id,
          product_name: data.product_name || data.name, // Map `name` to `product_name`
          store_id: data.store_id || data.storeID, // Map `storeID` to `store_id`
          picture_link: data.picture_link || data.pictureLink, // Map `pictureLink` to `picture_link`
          price: data.price,
          discounted_price: data.discounted_price || data.discountedPrice, // Map `discountedPrice` to `discounted_price`
          discounted: data.discounted || data.discountActive, // Map `discountActive` to `discounted`
          demo: data.demo,
          pricing: data.pricing,
        };

        setProduct(formattedProduct);
      } catch (err) {
        if (err instanceof Error) {
          setError(err.message);
        } else {
          setError("An unknown error occurred");
        }
      } finally {
        setLoading(false);
      }
    };

    if (storeID && productID) {
      fetchProduct();
    }
  }, [storeID, productID]);

  if (loading) return <p>Loading product details...</p>;
  if (error) return <p style={{ textAlign: "center" }}>{error}</p>;
  if (!product) return <p>Product not found!</p>;

  return (
    <ProductPreview
      product_id={product.product_id}
      product_name={product.product_name} // Correct prop name
      store_id={product.store_id}
      picture_link={product.picture_link}
      price={product.price}
      discounted_price={product.discounted_price}
      discounted={product.discounted}
      demo={product.demo}
      pricing={product.pricing}
    />
  );
};

export default ProductPreviewPage;