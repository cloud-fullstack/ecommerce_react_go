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