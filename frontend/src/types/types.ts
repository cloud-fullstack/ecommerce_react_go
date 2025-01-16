// src/types.ts
export interface BlogPost {
  id: number;
  title: string;
  content: string;
  author: string;
  createdAt: string;
  picture_link: string;
  store_id: number;
  product_id: number;
  blog_post_id: number;
  author_name: string;
  author_id: number;
  product_name: string;
  category_id: number; // Add this line
}