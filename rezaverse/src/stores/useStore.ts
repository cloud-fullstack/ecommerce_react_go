import { create } from "zustand";

interface StoreState {
  authToken: string | null;
  aviKey: string | null;
  aviLegacyName: string | null;
  aviLegacyNamePretty: string | null;
  profilePicture: string | null;
  animationLoaded: boolean; // Add animationLoaded
  setAuthToken: (token: string | null) => void;
  setAviKey: (key: string | null) => void;
  setAviLegacyName: (name: string | null) => void;
  setAviLegacyNamePretty: (name: string | null) => void;
  setProfilePicture: (picture: string | null) => void;
  setAnimationLoaded: (loaded: boolean) => void; // Add setter for animationLoaded
  updateProfilePicture: (token: string, key: string) => void;
  cart: Array<{
    productID: string;
    name: string;
    price: number;
    discountedPrice: number;
    discountActive: boolean;
    demo: boolean;
    picture_link: string;
    storeID: string;
  }>;
  
  setCart: (cart: StoreState["cart"]) => void;
}

const useStore = create<StoreState>((set) => ({
  authToken: null,
  aviKey: null,
  aviLegacyName: null,
  aviLegacyNamePretty: null,
  profilePicture: null,
  animationLoaded: false, // Initialize animationLoaded
  setAuthToken: (token) => set({ authToken: token }),
  setAviKey: (key) => set({ aviKey: key }),
  setAviLegacyName: (name) => set({ aviLegacyName: name }),
  setAviLegacyNamePretty: (name) => set({ aviLegacyNamePretty: name }),
  setProfilePicture: (picture) => set({ profilePicture: picture }),
  setAnimationLoaded: (loaded) => set({ animationLoaded: loaded }), // Setter for animationLoaded
  cart: [], // Initialize cart as an empty array
  setCart: (cart) => set({ cart }), // Setter for cart
  updateProfilePicture: (token, key) => {
    // Implement your logic to update the profile picture
    console.log("Updating profile picture with token:", token, "and key:", key);
  },
}));

export default useStore;