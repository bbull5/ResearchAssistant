import { Link } from 'react-router-dom';

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-gray-100 text-gray-800">
      <header className="flex justify-between items-center p-6 bg-white shadow">
        <h1 className="text-2xl font-bold text-blue-600">Research Assistant</h1>
        <div className="space-x-4">
          <Link to="/login" className="text-blue-600 font-medium hover:underline">Login</Link>
          <Link to="/register" className="text-blue-600 font-medium hover:underline">Register</Link>
        </div>
      </header>

      <main className="flex flex-col items-center justify-center text-center p-8 mt-20">
        <h2 className="text-4xl font-extrabold mb-4">Your AI-powered research companion</h2>
        <p className="text-lg max-w-2xl mb-8 text-gray-700">
          Upload research papers, organize them into workspaces by topic, and harness AI to generate summaries, insights, and even have conversations with your documents.
        </p>
        <Link to="/register">
          <button className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
            Get Started
          </button>
        </Link>
      </main>
    </div>
  );
}
