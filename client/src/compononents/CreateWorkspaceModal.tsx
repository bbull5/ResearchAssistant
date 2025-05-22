import { useState } from 'react';


type Props = {
    userId:     number,
    onClose:    () => void;
    onSuccess:  () => void;
};


export default function CreateWorkspaceModal({ userId, onClose, onSuccess}: Props) {
    const [title, setTitle] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);

    const handleSubmit = async () => {
        if (!title.trim()) {
            setError('Title is required');
            return;
        }

        setLoading(true);
        setError(null);

        try {
            const res = await fetch('http://localhost:8080/workspace/create', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({user_id: userId, title}),
            });

            if (!res.ok) {
                const msg = await res.text();
                throw new Error(msg || 'Failed to create workspace')
            }

            onSuccess();
            onClose();
        } catch (err: any) {
            setError(err.message || 'Something went wrong');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white p-6 rounded shadow w-full max-w-md space-y-4">
                <h2 className="text-xl font-bold text-gray-700">Create New Workspace</h2>

                <input
                    type="text"
                    placeholder="Workspace title"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    className="w-full border px-3 py-2 rounded"
                    disabled={loading}
                />

                {error && <p className="text-red-600 text-sm">{error}</p>}

                <div className="flex justify-end space-x-3">
                    <button
                        onClick={onClose}
                        className="px-4 py-2 border rounded hover:bg-gray-100"
                        disabled={loading}
                    >Cancel</button>
                    <button
                        onClick={handleSubmit}
                        className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                        disabled={loading}
                    >
                        {loading ? 'Creating...' : 'Create'}
                    </button>
                </div>
            </div>
        </div>
    );
}