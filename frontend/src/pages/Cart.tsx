import { useState } from "react";
import  useStore  from "../stores/useStore"; // Zustand store for global state

import "animate.css";

const Cart = () => {
  const { cart, setCart } = useStore(); // Access cart state and setter from Zustand
  const [shown, setShown] = useState(false); // State to control cart visibility

  const handleCheckout = () => {
    setShown(false);
    // Open BuyDialog (you can use a modal library like `react-modal` or `@headlessui/react`)
    // Example: open(BuyDialog);
  };

  const removeItem = (index: number) => {
    const newCart = [...cart];
    newCart.splice(index, 1);
    setCart(newCart);
  };

  const truePrice = (
    price: number,
    discountedPrice: number,
    discountActive: boolean,
    demo: boolean
  ): number => {
    if (demo) return 0;
    if (discountActive) return discountedPrice;
    return price;
  };

  const totalPrice = cart.reduce(
    (total, item) =>
      total +
      truePrice(item.price, item.discountedPrice, item.discountActive, item.demo),
    0
  );

  return (
    <div className="relative">
      {/* Cart Button */}
      <div className="flex items-center gap-2 cursor-pointer" onClick={() => setShown(!shown)}>
        <span className="text-white bg-amber-500 rounded-full w-6 h-6 flex items-center justify-center">
          {cart.length}
        </span>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth="1.5"
          stroke="currentColor"
          className="w-7 h-7 text-white"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M15.75 10.5V6a3.75 3.75 0 10-7.5 0v4.5m11.356-1.993l1.263 12c.07.665-.45 1.243-1.119 1.243H4.25a1.125 1.125 0 01-1.12-1.243l1.264-12A1.125 1.125 0 015.513 7.5h12.974c.576 0 1.059.435 1.119 1.007zM8.625 10.5a.375.375 0 11-.75 0 .375.375 0 01.75 0zm7.5 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z"
          />
        </svg>
      </div>

      {/* Cart Panel */}
      {shown && (
        <div
          className={`absolute right-0 top-16 bg-white border border-gray-200 rounded-lg shadow-lg p-4 w-96 animate__animated animate__fadeIn animate__fast`}
        >
          {cart.length > 0 ? (
            <>
              <div className="flex items-center mb-4">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth="1.5"
                  stroke="currentColor"
                  className="w-7 h-7 text-black"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M15.75 10.5V6a3.75 3.75 0 10-7.5 0v4.5m11.356-1.993l1.263 12c.07.665-.45 1.243-1.119 1.243H4.25a1.125 1.125 0 01-1.12-1.243l1.264-12A1.125 1.125 0 015.513 7.5h12.974c.576 0 1.059.435 1.119 1.007zM8.625 10.5a.375.375 0 11-.75 0 .375.375 0 01.75 0zm7.5 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z"
                  />
                </svg>
                <span className="ml-2 font-medium">Checkout your cart</span>
              </div>
              <hr className="mb-4" />
              <div className="space-y-4">
                {cart.map((product, index) => (
                  <div key={index} className="grid grid-cols-5 gap-4 items-center">
                    <img
                      src={product.picture_link}
                      alt="product"
                      className="rounded w-12 h-12"
                    />
                    <span className="text-sm">
                      <a
                        href={`/store/${product.storeID}/${product.productID}`}
                        className="text-blue-500 hover:underline"
                      >
                        {product.name}
                      </a>
                    </span>
                    <span className="text-sm">
                      L${truePrice(product.price, product.discountedPrice, product.discountActive, product.demo)}
                    </span>
                    {product.demo ? (
                      <span className="text-sm bg-black text-white px-2 py-1 rounded-full text-center">
                        DEMO
                      </span>
                    ) : product.discountActive ? (
                      <span className="text-sm bg-black text-white px-2 py-1 rounded-full text-center">
                        DISCOUNTED
                      </span>
                    ) : (
                      <span></span>
                    )}
                    <button
                      className="text-sm bg-amber-100 text-amber-700 px-3 py-1 rounded-md hover:bg-amber-200 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2"
                      onClick={() => removeItem(index)}
                    >
                      Remove
                    </button>
                  </div>
                ))}
              </div>
              <div className="mt-4 grid grid-cols-5 gap-4 items-center">
                <span></span>
                <span className="font-bold">Total</span>
                <span className="font-bold">L${totalPrice}</span>
                <span></span>
                <button
                  className="bg-[#4f2236] text-white px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-[#4f2236] focus:ring-offset-2"
                  onClick={handleCheckout}
                >
                  Checkout
                </button>
              </div>
            </>
          ) : (
            <p className="text-center">No items in cart!</p>
          )}
        </div>
      )}
    </div>
  );
};

export default Cart;