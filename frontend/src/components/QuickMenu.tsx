import { useState, useEffect, useRef } from "react";
import useStore from "../stores/useStore";
import SignInInstructions from "./SignInInstructions"; // Assuming you have this component
import apiClient from '../utils/api';

const QuickMenu = () => {
  const [showMenu, setShowMenu] = useState(false);
  const [hudStatus, setHudStatus] = useState("");
  const menuRef = useRef<HTMLDivElement | null>(null); // Add type annotation for menuRef

  const {
    authToken,
    aviKey,
    aviLegacyName,
    aviLegacyNamePretty,
    profilePicture,
    updateProfilePicture,
    setAuthToken,
    setAviKey,
    setAviLegacyName,
    setProfilePicture,
  } = useStore();

  // Handle outside click and escape key press
  useEffect(() => {
    const handleOutsideClick = (event: MouseEvent) => {
      if (showMenu && !menuRef.current?.contains(event.target as Node)) {
        setShowMenu(false);
      }
    };

    const handleEscape = (event: KeyboardEvent) => {
      if (showMenu && event.key === "Escape") {
        setShowMenu(false);
      }
    };

    document.addEventListener("click", handleOutsideClick, false);
    document.addEventListener("keyup", handleEscape, false);

    return () => {
      document.removeEventListener("click", handleOutsideClick, false);
      document.removeEventListener("keyup", handleEscape, false);
    };
  }, [showMenu]); // Re-run when `showMenu` changes

  // Fetch HUD status when menu is shown
  useEffect(() => {
    if (showMenu) {
      setHudStatus("Fetching HUD status ...");

      // Use apiClient.post for POST requests
      apiClient
        .post(
          "/api/heartbeat-hud/", // Endpoint URL
          { avatar_key: aviKey }, // Request body
          {
            headers: {
              Authorization: `${authToken}.${aviKey}`, // Headers
            },
          }
        )
        .then((res) => {
          const data = res.data;
          if (data.error) {
            setHudStatus(
              data.message === "failed to heartbeat HUD"
                ? "HUD not worn. Get HUD at: xxx"
                : `HUD not worn: ${data.message}`
            );
          } else {
            setHudStatus("HUD is worn.");
          }
        })
        .catch((err) => {
          console.error("Error fetching HUD status:", err);
          setHudStatus("Error fetching HUD status.");
        });
    }
  }, [showMenu, authToken, aviKey]); // Re-run when `showMenu`, `authToken`, or `aviKey` changes

  // Sign out function
  const signOut = () => {
    setAuthToken("");
    setAviKey("");
    setAviLegacyName("");
    setProfilePicture("https://i.imgur.com/JlGo0fY.png");
    window.location.href = "/";
  };

  return (
    <div className="relative" ref={menuRef}>
      {aviKey ? (
        <button
          onClick={() => setShowMenu(!showMenu)}
          className="focus:outline-none focus:shadow-solid"
        >
          <img
            className="w-10 h-10 rounded-full"
            alt="profilePicture"
            src={profilePicture || "https://res.cloudinary.com/dutkkgjm5/image/upload/v1667497058/gtpj4fwtjlegsrfftkf6.png"}
          />
        </button>
      ) : (
        <a href="/signIn">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth="1.5"
            stroke="currentColor"
            className="w-8 h-8 text-white"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M17.982 18.725A7.488 7.488 0 0012 15.75a7.488 7.488 0 00-5.982 2.975m11.963 0a9 9 0 10-11.963 0m11.963 0A8.966 8.966 0 0112 21a8.966 8.966 0 01-5.982-2.275M15 9.75a3 3 0 11-6 0 3 3 0 016 0z"
            />
          </svg>
        </a>
      )}

      {showMenu && (
        <div className="absolute right-0 mt-2 w-48 bg-black rounded-md shadow-lg transform transition-transform duration-100 scale-100 origin-top-right">
          <div className="py-2">
            <span className="block px-4 py-2 text-amber-200">Hello {aviLegacyNamePretty}!</span>
            <a
              href="/order-history"
              className="block px-4 py-2 text-white hover:bg-amber-200 hover:text-gray-900"
            >
              Order History
            </a>
            {aviKey && (
              <a
                href={`/blogs/${aviKey}`}
                className="block px-4 py-2 text-white hover:bg-amber-200 hover:text-gray-900"
              >
                Your Blogs
              </a>
            )}
            <a
              href="/merchant"
              className="block px-4 py-2 text-white hover:bg-amber-200 hover:text-gray-900"
            >
              Creator Dashboard
            </a>
            {aviLegacyName ? (
              <a
                onClick={signOut}
                href="/"
                className="block px-4 py-2 text-white hover:bg-amber-200 hover:text-gray-900 cursor-pointer"
              >
                Logout
              </a>
            ) : (
              <SignInInstructions />
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default QuickMenu;