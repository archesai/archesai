import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

const faqs = [
  {
    category: "General",
    questions: [
      {
        answer:
          "Arches AI leverages advanced machine learning techniques to transform your documents into 'word embeddings'. These embeddings allow users to search by semantic meaning, rather than the exact language. This is especially valuable when trying to decipher unstructured text information, like textbooks, articles, and documentation.",
        question: "How does Arches AI work?",
      },
      {
        answer:
          "Arches AI offers a multitude of features: semantic search through documents, support for multiple files searching, an API for seamless integration with other platforms, and the ability to embed interactive widgets on your website or application for easy access and message by users.",
        question: "What are the main features of Arches AI?",
      },
      {
        answer:
          "Absolutely! Arches AI provides API tokens, allowing developers to easily integrate its powerful search capabilities into other platforms, applications, or websites.",
        question: "Can I integrate Arches AI with my application?",
      },
      {
        answer:
          "Embedding Arches AI widgets allows users to access and interact with the platform directly from your website or application. It provides a seamless user experience, ensuring that users can easily search and retrieve information without navigating away from your platform.",
        question: "What are the benefits of embedding Arches AI widgets?",
      },
    ],
  },
  {
    category: "Features",
    questions: [
      {
        answer:
          "Yes, Arches AI supports searching across multiple files. When creating an chatbot, you have the option to either search through all of your files or a specific subset.",
        question: "Can I search multiple files simultaneously?",
      },
      {
        answer:
          "API tokens are unique keys provided by Arches AI to developers. They enable the integration of Arches AI's features into other platforms. By using these tokens, developers can make authenticated requests to Arches AI's services.",
        question: "How do API tokens work?",
      },
      {
        answer:
          "Arches AI's embedded widgets are designed with flexibility in mind. Developers can customize the appearance and behavior to match the design and functionality of their platform, ensuring a consistent and seamless user experience.",
        question: "How customizable are the embedded widgets?",
      },
    ],
  },
  {
    category: "Security",
    questions: [
      {
        answer:
          "Arches AI prioritizes data security. Documents are stored in encrypted cloud storage with strict security protocols in place to protect against unauthorized access. Your information is safeguarded from potential threats, and all documents can be easily deleted through the 'Files' page if needed.",
        question: "How secure is my data with Arches AI?",
      },
      {
        answer: `Arches AI is committed to protecting your privacy. We do not sell or share your information with third parties. We only use your data to provide you with the best possible experience. For more information, please refer to our Privacy Policy.`,
        question: "How does Arches AI protect my privacy?",
      },
      {
        answer: `Arches AI does not claim ownership of your documents. We only use your data to provide you with the best possible experience. For more information, please refer to our Privacy Policy.`,
        question: "How does Arches AI protect my intellectual property?",
      },
    ],
  },
  {
    category: "Pricing",
    questions: [
      {
        answer: `Credits are based on token usage. Here's a breakdown:

          - Upload: 1 credit = 50 tokens. If you index a 1000-word document, that's 20 credits.
          - Chat: 1 credit = 1 token. A query with a context length of 1000 words uses 1000 credits.
          - Image: Each image processing costs 500 credits.
          - Animation: 85 credits per frame. A 10-second video at 12fps will cost 10,200 credits.`,
        question: "How are credits calculated in Arches AI?",
      },
      {
        answer:
          "Yes, unused credits from one month will roll over to the next. These accrued credits don't expire and can be used at any time, regardless of the status of your subscription.",
        question: "Do unused credits carry over to the next billing cycle?",
      },
    ],
  },
];

export const FAQ = () => {
  return (
    <section className="container py-24 sm:py-32" id="faq">
      <h2 className="text-3xl md:text-4xl font-bold mb-4">
        Frequently Asked{" "}
        <span className="bg-gradient-to-b from-primary/60 to-primary text-transparent bg-clip-text">
          Questions
        </span>
      </h2>

      {faqs.map(({ category, questions }, i) => (
        <>
          <h3 className="font-bold text-lg mt-8">{category}</h3>
          <Accordion
            className="w-full AccordionRoot"
            collapsible
            key={i}
            type="single"
          >
            {questions.map(({ answer, question }, i) => (
              <AccordionItem key={i} value={i.toString()}>
                <AccordionTrigger className="text-left">
                  {question}
                </AccordionTrigger>

                <AccordionContent>{answer}</AccordionContent>
              </AccordionItem>
            ))}
          </Accordion>
        </>
      ))}

      <h3 className="font-medium mt-4">
        Still have questions?{" "}
        <a
          className="text-primary transition-all border-primary hover:border-b-2"
          href="#"
          rel="noreferrer noopener"
        >
          Contact us
        </a>
      </h3>
    </section>
  );
};
