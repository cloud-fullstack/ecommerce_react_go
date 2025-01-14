import React from 'react';

const Faq = () => {
  // Define styles as template literals
  const styles = `
    .faq-container {
      max-width: 800px;
      margin: 0 auto;
      padding: 20px;
      font-family: Arial, sans-serif;
    }
    .faq-title {
      text-align: center;
      color: #333;
      margin-bottom: 20px;
    }
    .faq-list {
      display: flex;
      flex-direction: column;
      gap: 15px;
    }
    .faq-item {
      background-color: #f9f9f9;
      border: 1px solid #ddd;
      border-radius: 8px;
      padding: 15px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }
    .faq-item h3 {
      margin: 0 0 10px;
      color: #444;
    }
    .faq-item p {
      margin: 0;
      color: #666;
    }
  `;

  return (
    <>
      <style>{styles}</style>
      <div className="faq-container">
        <h2 className="faq-title">Frequently Asked Questions</h2>
        <div className="faq-list">
          <div className="faq-item">
            <h3>What is this project about?</h3>
            <p>This project is a React-based application designed to showcase various features and components.</p>
          </div>
          <div className="faq-item">
            <h3>How do I get started?</h3>
            <p>To get started, simply navigate through the application and explore the available features.</p>
          </div>
          <div className="faq-item">
            <h3>Where can I find more information?</h3>
            <p>You can find more information in the documentation or by contacting support.</p>
          </div>
        </div>
      </div>
    </>
  );
};

export default Faq;