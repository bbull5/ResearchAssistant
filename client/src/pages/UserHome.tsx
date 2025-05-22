import { Link } from 'react-router-dom';
import { useEffect, useState } from 'react';
import UploadModal from '../compononents/UploadModal';
import CreateWorkspaceModal from '../compononents/CreateWorkspaceModal';

type Document = {
  id: number;
  title: string;
  workspace_id: number;
  uploaded_at: string;
};

type Workspace = {
  id: number;
  title: string;
  created_at: string;
};

export default function HomePage() {
  const [showUploadModal, setShowUploadModal] = useState(false);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [documents, setDocuments] = useState<Document[]>([]);
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [loadingDocs, setLoadingDocs] = useState(true);
  const [loadingWorkspaces, setLoadingWorkspaces] = useState(true);

  const userId = 1; // TEMP — this should eventually come from session/auth context

  const fetchDocuments = async () => {
    setLoadingDocs(true);
    try {
      const res = await fetch(`http://localhost:8080/documents/get?user_id=${userId}`);
      const data = await res.json();
      setDocuments(data);
    } catch (err) {
      console.error('Error fetching documents:', err);
    } finally {
      setLoadingDocs(false);
    }
  };

  const fetchWorkspaces = async () => {
    setLoadingWorkspaces(true);
    try {
      const res = await fetch(`http://localhost:8080/workspace/get?user_id=${userId}`);
      const data = await res.json();
      setWorkspaces(data);
    } catch (err) {
      console.error('Error fetching workspaces:', err);
    } finally {
      setLoadingWorkspaces(false);
    }
  };

  const handleDeleteWorkspace = async (id: number) => {
    const confirmed = window.confirm("Are you sure you want to delete this workspace?");
    if (!confirmed) return;

    try {
      const res = await fetch('http://localhost:8080/workspace/delete', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id }),
      });

      if (!res.ok) {
        const msg = await res.text();
        throw new Error(msg || 'Failed to delete workspace');
      }

      // Refresh workspace list
      fetchWorkspaces();
    } catch (err) {
      console.error('Error deleting workspace:', err);
    }
  };

  useEffect(() => {
    fetchDocuments();
    fetchWorkspaces();
  }, []);

  return (
    <div className="min-h-screen bg-gray-50 text-gray-800">
      {/* NavBar */}
      <nav className="bg-white shadow p-4 flex justify-between items-center">
        <div className="text-xl font-bold text-blue-600">Research Assistant</div>
        <div className="space-x-6">
          <Link to="/" className="hover:underline">Home</Link>
          <Link to="/" className="hover:underline">Workspaces</Link>
          <button
            onClick={() => setShowUploadModal(true)}
            className="hover:underline"
          >
            Upload
          </button>
          <Link to="/" className="hover:underline">Profile</Link>
          <Link to="/" className="hover:underline">Logout</Link>
        </div>
      </nav>

      {/* Main Layout */}
      <div className="p-6">
        <h2 className="text-2xl font-semibold mb-6">Welcome, [username]!</h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Workspaces Panel */}
          <div className="bg-white p-4 rounded shadow">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-xl font-bold">Workspaces</h3>
              <button 
                onClick={() => setShowCreateModal(true)}
                className="bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700 text-sm"
              >
                Create Workspace
              </button>
            </div>
            {loadingWorkspaces ? (
              <p>Loading workspaces...</p>
            ) : workspaces.length === 0 ? (
              <p className="text-gray-500 text-sm">No workspaces yet.</p>
            ) : (
              <ul className="space-y-2">
                {workspaces.map((ws) => (
                  <li key={ws.id} className="flex justify-between items-center border-b py-2">
                    <span>{ws.title}</span>
                    <div className="space-x-2">
                      <button className="text-blue-600 text-sm hover:underline">Open</button>
                      <button 
                        className="text-red-600 text-sm hover:underline"
                        onClick={() => handleDeleteWorkspace(ws.id)}
                      >
                        Delete
                      </button>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </div>

          {/* Documents Panel */}
          <div className="bg-white p-4 rounded shadow">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-xl font-bold">Uploaded Documents</h3>
              <button
                onClick={() => setShowUploadModal(true)}
                className="bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700 text-sm"
              >
                Upload PDF
              </button>
            </div>
            {loadingDocs ? (
              <p>Loading documents...</p>
            ) : documents.length === 0 ? (
              <p className="text-gray-500 text-sm">No documents uploaded yet.</p>
            ) : (
              <ul className="space-y-2">
                {documents.map((doc) => (
                  <li key={doc.id} className="flex justify-between items-center border-b py-2">
                    <div>
                      <p className="font-medium">{doc.title}</p>
                      <p className="text-sm text-gray-500">
                        Uploaded: {new Date(doc.uploaded_at).toLocaleDateString()}
                      </p>
                    </div>
                    <div className="space-x-2">
                      <button className="text-blue-600 text-sm hover:underline">Add to Workspace</button>
                      <button className="text-red-600 text-sm hover:underline">Delete</button>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      </div>

      {showUploadModal && (
        <UploadModal
          onClose={() => setShowUploadModal(false)}
          onSuccess={() => {
            fetchDocuments();
            fetchWorkspaces();
          }}
          userId={userId}
        />
      )}
      {showCreateModal && (
        <CreateWorkspaceModal
          userId={userId}
          onClose={() => setShowCreateModal(false)}
          onSuccess={fetchWorkspaces}
        />
      )}
    </div>
  );
}
