import React, { useState, useEffect } from "react";
import useStore from "../../stores/useStore"; // Zustand store for global state
import "animate.css";
import apiClient from '../../utils/api'; // Import the apiClient

interface ProductPreviewProps {
  name?: string;
  storeID?: string;
  productID?: string;
  pictureLink?: string;
  price?: number;
  discountedPrice?: number;
  discountActive?: boolean;
  demo?: boolean;
  pricing?: boolean;
}

interface BuyDialogProps {
  demoProduct?: ProductPreviewProps;
  resendOrderID?: string;
  onClose: () => void;
}

const BuyDialog: React.FC<BuyDialogProps> = ({ demoProduct, resendOrderID, onClose }) => {
  const { authToken, aviKey, cart, setCart } = useStore();
  const [hudStatus, setHudStatus] = useState(null);
  const [orderID, setOrderID] = useState("");
  const [orderSent, setOrderSent] = useState(false);
  const [revealOrderDetails, setRevealOrderDetails] = useState(false);
  const [error, setError] = useState(null);

  const getHUDStatus = async () => {
    if (!authToken || !aviKey) {
      throw new Error("not signed in");
    }
    try {
      const response = await apiClient.post("/api/heartbeat-hud/", {
        avatar_key: aviKey,
      });
      if (response.data.error) {
        throw new Error(response.data.message);
      }
      return response.data;
    } catch (err) {
      throw new Error(err.response?.data?.message || "Failed to fetch HUD status");
    }
  };

  const sendOrder = async (resendOrderID: string | null) => {
    let orderLines = cart.map((item) => ({
      product_id: item.productID, // Ensure item.productID is always a string
      demo: item.demo,
    }));

    if (demoProduct) {
      // Ensure demoProduct.productID is defined
      if (!demoProduct.productID) {
        throw new Error("Demo product ID is missing.");
      }
      orderLines = [{ product_id: demoProduct.productID, demo: true }];
    }

    try {
      let response;
      if (resendOrderID) {
        // For resend-order endpoint
        response = await apiClient.post("/api/resend-order/", {
          avatar_buyer: aviKey,
          order_id: resendOrderID,
        });
      } else {
        // For create-order endpoint
        response = await apiClient.post("/api/create-order/", {
          avatar_buyer: aviKey,
          order_lines: orderLines,
          avatar_key: aviKey,
        });
      }

      if (response.data.error) {
        throw new Error(response.data.message);
      }

      setCart([]); // Clear the cart
      return response.data.order_id;
    } catch (err) {
      throw new Error(err.response?.data?.message || "Failed to send order");
    }
  };

  const handleSendOrder = async () => {
    setRevealOrderDetails(true);
    try {
      const orderID = await sendOrder(resendOrderID || null);
      setOrderID(orderID);
      setOrderSent(true);
    } catch (err) {
      setError(err.message);
    }
  };

  const truePrice = (price: number, discountedPrice: number, discountActive: boolean, demo: boolean) => {
    if (demo) return 0;
    if (discountActive) return discountedPrice;
    return price;
  };

  const totalPrice = cart.reduce(
    (total, item) =>
      total + truePrice(item.price, item.discountedPrice, item.discountActive, item.demo),
    0
  );

  useEffect(() => {
    getHUDStatus()
      .then((hud) => setHudStatus(hud))
      .catch((err) => setError(err.message));
  }, []);

  return (
    <div className="p-4">
      {error ? (
        <>
          {error === "not signed in" ? (
            <p className="text-center">Please Sign In before purchasing.</p>
          ) : (
            <>
              <p className="text-center">Error: Could not find HUD in-world.</p>
              <p className="text-center">Are you logged in?</p>
              <p className="text-center">Are you wearing the HUD and clicked the Register button?</p>
              <p className="text-center mb-2">
                Is your HUD currently accepting an order? If so, cancel the order so you can accept
                this one.
              </p>
            </>
          )}
          <button
            className="block mx-auto rounded-md bg-amber-100 text-amber-700 px-3 py-2 hover:bg-amber-200 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2"
            onClick={onClose}
          >
            Close
          </button>
        </>
      ) : (
        <>
          {!revealOrderDetails ? (
            <>
              {resendOrderID ? (
                <p className="text-center">Resend order to your HUD?</p>
              ) : demoProduct ? (
                <p className="text-center">Send this free demo to your avatar?</p>
              ) : (
                <p className="text-center">
                  Send this order to HUD for L${totalPrice}?
                </p>
              )}
              <button
                className="block mx-auto rounded-md bg-[#4f2236] text-white px-3 py-2 hover:bg-[#3a1a2a] focus:outline-none focus:ring-2 focus:ring-[#4f2236] focus:ring-offset-2"
                onClick={handleSendOrder}
              >
                {demoProduct ? "Send Free Demo" : resendOrderID ? "Resend Order" : "Send Order"}
              </button>
            </>
          ) : (
            <>
              {orderSent ? (
                <>
                  <p className="text-center">
                    Order <code className="bg-blue-100 text-blue-700 px-2 py-1 rounded">{orderID}</code> submitted!
                  </p>
                  {!demoProduct && (
                    <p className="text-center">Pay your HUD in-world to finish purchasing.</p>
                  )}
                  <button
                    className="block mx-auto rounded-md bg-amber-100 text-amber-700 px-3 py-2 hover:bg-amber-200 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2"
                    onClick={onClose}
                  >
                    Close
                  </button>
                </>
              ) : (
                <p className="text-center">Submitting Order ...</p>
              )}
            </>
          )}
        </>
      )}
    </div>
  );
};

export default BuyDialog;