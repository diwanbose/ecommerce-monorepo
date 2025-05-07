export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  image: string;
  stock: number;
}

export interface CartItem {
  productId: number;
  quantity: number;
  product: Product;
}

export interface Order {
  id: number;
  items: CartItem[];
  total: number;
  status: 'pending' | 'processing' | 'shipped' | 'delivered';
  paymentMethod: 'credit_card' | 'netbanking' | 'cod';
  createdAt: string;
}

export interface FeatureFlag {
  name: string;
  enabled: boolean;
}

export interface User {
  id: number;
  email: string;
  name: string;
} 