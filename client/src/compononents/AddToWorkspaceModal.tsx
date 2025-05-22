import { useState } from 'react';


type Document = {
    id:     number;
    title:  string;
};

type Workspace = {
    id:     number;
    title:  string;
};

type Props = {
    documents:  Document[];
    workspaces: Workspace[];
    onClose:    () => void;
    onSuccess:  () => void;
};


export default function AddToWorkspaceModal({ documents, workspaces, onClose, onSuccess }: Props) {
    const [selectedDocs, setSelectedDocs] = useState<number[]>([]);
    const [selectedWorkspace, setSelectedWorkspace] =  useState<number | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async () => {
        if (!selectedWorkspace || selectedDocs.length == 0) {
            setError("Please select at least one document and a workspace");
            return;
        }

        setLoading(true);
        setError(null);

        try {
            for (const docId of selectedDocs) {
                const res = await fetch('http://localhost:8080/workspace/add-document', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ document_id: docId, workspace_id: selectedWorkspace }),
                });

                if (!res.ok) throw new Error (`Failed to add document ID ${docId}`);
            }

            onSuccess();
            onClose();
        } catch (err: any) {
            setError(err.message || 'Something went wrong');
        } finally {
            setLoading(false);
        }
    };

    const toggleDocSelection = (id: number) => {
        setSelectedDocs((prev) => 
            prev.includes(id) ? prev.filter((d) => d !== id) : [...prev, id]
        );
    };

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white p-6 rounded shadow w-full max-w-lg space-y-4">
                <h2 className="text-xl font-bold">Add Documents to Workspace</h2>

                <div>
                    <label className="block text-sm font-medium mb-1">Select Workspace</label>
                    <select
                        value={selectedWorkspace ?? ''}
                        onChange={(e) => setSelectedWorkspace(Number(e.target.value))}
                        className="w-full border rounded px-3 py-2"
                        disabled={loading}
                    >
                        <option value="">--Select a Workspace --</option>
                        {workspaces.map((ws) => (
                            <option key={ws.id} value={ws.id}>
                                {ws.title}
                            </option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium mb-1">Select Documents</label>
                    <div className="max-h-40 overflow-y-auto border rounded p-2 space-y-1">
                        {documents.map((doc) => (
                            <label key={doc.id} className="flex items-center space-x-2">
                                <input
                                    type="checkbox"
                                    checked={selectedDocs.includes(doc.id)}
                                    onChange={() => toggleDocSelection(doc.id)}
                                    disabled={loading}
                                />
                                <span>{doc.title}</span>
                            </label>
                        ))}
                    </div>
                </div>

                {error && <p className="text-red-600 text-sm">{error}</p>}

                <div className="flex justify-end space-x-3">
                    <button
                        onClick={onClose}
                        className="px-4 py-2 border rounded hover:bg-gray-100"
                        disabled={loading}
                    >
                        Cancel
                    </button>
                    <button 
                        onClick={handleSubmit}
                        className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                        disabled={loading}
                    >
                        {loading ? 'Adding...' : 'Add'}
                    </button>
                </div>
            </div>
        </div>
    );
}