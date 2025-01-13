import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import 'animate.css';
import apiClient from '../utils/api'; // Import the centralized API client

const SignIn = () => {
  const [authToken, setAuthToken] = useState("");
  const [screenWidth, setScreenWidth] = useState(window.innerWidth);
  const [screenType, setScreenType] = useState("computer");
  const navigate = useNavigate();

  // Handle screen size changes
  useEffect(() => {
    const handleResize = () => {
      setScreenWidth(window.innerWidth);
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  // Determine screen type based on width
  useEffect(() => {
    if (screenWidth < 768) {
      setScreenType("phone");
    } else {
      setScreenType("computer");
    }
  }, [screenWidth]);

  // Handle "Enter the website logged out" click
  const handleEnter = () => {
    if (!authToken) {
      setAuthToken("fakeToken");
    }
    navigate("/"); // Navigate to the home page
  };

  // Handle Second Life Sign In
  const handleSecondLifeSignIn = async () => {
    try {
      const response = await apiClient.post("/api/gen-token/", {
        hash: "shopaTMAC#3", // Replace with actual hash logic
        legacy_name: "user", // Replace with actual user name
        avatar_key: "user-key" // Replace with actual avatar key
      });

      const data = res.data;
      if (data.token) {
        setAuthToken(data.token);
        navigate("/"); // Navigate to the home page after successful sign-in
      }
    } catch (error) {
      console.error("Error signing in:", error);
    }
  };

  return (
    <div className="flex min-h-full flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <img
          className="mx-auto h-12 w-auto"
          src="https://res.cloudinary.com/dutkkgjm5/image/upload/v1667398956/kwfvofm5ewxpoox8lhgb.png"
          alt="Your Company"
        />
        <h2 className="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900">
          Sign in to your account
        </h2>
        <p className="mt-2 text-center text-sm text-gray-600 animate__animated animate__pulse animate__infinite animate__slow">
          Or
          <a
            href="/"
            onClick={handleEnter}
            className="font-medium text-indigo-600 hover:text-indigo-500"
          >
            enter the website logged out
          </a>
        </p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
          <div className="mt-6 grid grid-cols-1 gap-3">
            <div>
              <a
                href="#"
                className="inline-flex w-full justify-center rounded-md border border-gray-300 bg-white py-2 px-4 text-md font-medium text-gray-500 shadow-sm hover:bg-gray-50"
              >
                <span className="sr-only">Sign in with Second Life</span>
                <img
                  className="mt-3 h-10"
                  src="https://res.cloudinary.com/dutkkgjm5/image/upload/v1667493466/n3ihgzy9iars9beazspv.svg"
                  alt="Second Life"
                />
                <span className="mt-3">Get a HUD on the marketplace.</span>
              </a>
            </div>
          </div>

          <div className="mt-6">
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-gray-300"></div>
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="bg-white px-2 text-gray-500">Or continue with</span>
              </div>
            </div>

            <div className="mt-6 grid grid-cols-1 gap-3">
              <div>
                <button
                  onClick={handleSecondLifeSignIn} // Ensure the function is connected here
                  className="inline-flex w-full justify-center rounded-md border border-gray-300 bg-white py-2 px-4 text-md font-medium text-gray-500 shadow-sm hover:bg-gray-50"
                >
                  <span className="sr-only">Sign in with Second Life</span>
                  <img
                    className="mt-1 h-10"
                    src="https://res.cloudinary.com/dutkkgjm5/image/upload/v1667493466/n3ihgzy9iars9beazspv.svg"
                    alt="Second Life"
                  />
                  <span className="pt-1">
                    {screenType === "computer"
                      ? "Safely Sign up through Second Life to verify your avatar identity with the HUD."
                      : "If you want to log in from your phone, please share the login URL provided by your Second Life HUD to your mobile device. Remember not to share the URL with other people."}
                  </span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SignIn;