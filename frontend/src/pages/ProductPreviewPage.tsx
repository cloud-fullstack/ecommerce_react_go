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
        setProduct(data);
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