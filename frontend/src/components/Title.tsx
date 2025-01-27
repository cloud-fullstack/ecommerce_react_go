// components/Title.tsx
import React, { ReactNode } from "react";

interface TitleProps {
  children: ReactNode; // Type for children (can be any valid React node)
  className?: string;  // Optional className prop
}

const Title: React.FC<TitleProps> = ({ children, className = "" }) => {
  return (
    <div className={`text-center titleSponsorised ${className}`}>
      <h2 className="animate__animated animate__fadeInDown pt-5 font-bold tracking-tight text-3xl sm:text-4xl text-gray-900">
        {children}
      </h2>
    </div>
  );
};

export default Title;