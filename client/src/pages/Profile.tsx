import { useState } from 'react';
import { Link } from 'react-router-dom';


export default function ProfilePage() {
    const [activeTab, setActiveTab] = useState<'account' | 'api' | 'none'>('account');

    const renderContent = () => {
        switch (activeTab) {
            case 'account':
                return (
                    <div className="space-y-4">
                        <h2 className="text-xl font-semibold">Account Information</h2>
                        <div>
                            <p className="text-gray-700">Username: <span className="font-medium">[username]</span></p>
                            <p className="text-gray-700">Email: <span className="font-medium">[email@example.com]</span></p>
                        </div>
                        <button className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 text-sm">
                            Change Password
                        </button>
                    </div>
                );
            case 'api':
                return (
                    <div className="space-y-4">
                        <h2 className="text-xl font-semibold">API Keys</h2>
                        {['OpenAI', 'Mistral', 'Claude', 'VoyageAI'].map((service) => (
                            <div key={service} className="flex items-center space-x-4">
                                <label className="w-24 text-gray-700">{service}:</label>
                                <input
                                    type="text"
                                    placeholder="API Key"
                                    className="border border-gray-300 rounded px-3 py-1 flex-1"
                                />
                                <button className="bg-green-600 text-white px-3 py-1 rounded hover:bg-green-700 text-sm">
                                    Register
                                </button>
                            </div>
                        ))}
                    </div>
                );
            default:
                return null;
        }
    };

    return (
        <div className="min-h-screen bg-gray-50 text-gray-800">
            {/* Navbar */}
            <nav className="bg-white shadow p-4 flex justify-between items-center">
                <div className="text-xl font-bold text-blue-600">Research Assistant</div>
                <div className="space-x-6">
                    <Link to="/" className="hover:underline">Home</Link>
                    <Link to="/" className="hover:underline">Workspaces</Link>
                    <Link to="/profile" className="hover:underline">Profile</Link>
                    <Link to="/" className="hover:underline">Logout</Link>
                </div>
            </nav>

            {/* Layout */}
            <div className="flex p-6 space-x-6">
                {/* Sidebar Menu */}
                <aside className="w-64 bg-white rounded shadow p-4 flex flex-col justify-between">
                    <div className="space-y-2">
                        <button
                            onClick={() => setActiveTab('account')}
                            className={`block w-full text-left px-3 py-2 rounded text-sm hover:bg-gray-100 ${
                                activeTab === 'account' ? 'bg-blue-100 font-semibold text-blue-700' : ''
                            }`}
                        >
                            Account
                        </button>
                        <button
                            onClick={() => setActiveTab('api')}
                            className={`block w-full text-left px-3 py-2 rounded text-sm hover:bg-gray-100 ${
                                activeTab === 'api' ? 'bg-blue-100 font-semibold text-blue-700' : ''
                            }`}
                        >
                            API Keys
                        </button>
                    </div>
                    <div>
                        <button className="text-red-600 text-sm hover:underline">Log Out</button>
                    </div>
                </aside>

                {/* Main Display Panel */}
                <main className="flex-1 bg-white p-6 rounded shadow">
                    {renderContent()}
                </main>
            </div>
        </div>
    );
}