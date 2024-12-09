import React from 'react';

function FAQ() {
  const faqs = [
    { question: "What's the best thing about Switzerland?", answer: "I don't know, but the flag is a big plus. Lorem ipsum dolor sit amet consectetur adipisicing elit. Quas cupiditate laboriosam fugiat." },
    { question: "What's the best thing about Switzerland?", answer: "I don't know, but the flag is a big plus. Lorem ipsum dolor sit amet consectetur adipisicing elit. Quas cupiditate laboriosam fugiat." },
    { question: "What's the best thing about Switzerland?", answer: "I don't know, but the flag is a big plus. Lorem ipsum dolor sit amet consectetur adipisicing elit. Quas cupiditate laboriosam fugiat." },
    { question: "What's the best thing about Switzerland?", answer: "I don't know, but the flag is a big plus. Lorem ipsum dolor sit amet consectetur adipisicing elit. Quas cupiditate laboriosam fugiat." },
  ];

  return (
    <section className="py-16">
      <h2 className="text-3xl font-bold mb-8 px-8">Frequently asked questions</h2>
      <div className="max-w-6xl mx-auto grid grid-cols-2 gap-8">
        {faqs.map((faq, index) => (
          <div key={index}>
            <h3 className="font-bold mb-2">{faq.question}</h3>
            <p className="text-gray-600">{faq.answer}</p>
          </div>
        ))}
      </div>
    </section>
  );
}

export default FAQ;