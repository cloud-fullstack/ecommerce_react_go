import React from 'react';
import PersonalSelection from '../components/PersonalSelection';
import HottestSales from '../components/HottestSales';
import MostLovedBlogs from '../components/MostLovedBlogs';
import FAQ from '../components/FAQ';

function Home() {
  return (
    <div>
      <PersonalSelection />
      <HottestSales />
      <MostLovedBlogs />
      <FAQ />
    </div>
  );
}