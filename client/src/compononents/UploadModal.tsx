import { useCallback, useState } from 'react';
import { useDropzone } from 'react-dropzone';


type Props = {
    onClose: () => void;
    onSuccess: () => void;
    userId: number;
    workspaceId?: number;
}


export default function UploadModal({onClose, onSuccess, userId, workspaceId}: Props) {
    const [title, setTitle] = useState('');
    const [file, setFile] = useState<File | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [uploading, setUploading] = useState(false);

    const onDrop = useCallback((acceptedFiles: File[]) => {
        if (acceptedFiles.length > 0) {
            setFile(acceptedFiles[0]);
        }
    }, []);

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        accept: { 'application/pdf': [] },
        maxFiles: 1,
        onDrop,
    });

    const handleSubmit = async () => {
        if (!file || !title) {
            setError('Title and PDF file are required.');
            return;
        }

        const formData = new FormData();
        formData.append('title', title);
        formData.append('pdf', file);
        formData.append('user_id', String(userId));
        if (workspaceId) formData.append('workspace_id', String(workspaceId));

        setUploading(true);
        setError(null);

        try {
            const res = await fetch('http://localhost:8080/documents/upload', {
                method: 'POST',
                body: formData,
            });

            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || 'Upload failed');
            }

            onSuccess();
            onClose();
        } catch (err: any) {
            setError(err.message || 'Something went wrong');
        } finally {
            setUploading(false);
        }
    };

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white p-6 rounded w-full max-w-md shadow space-y-4">

                <input
                    type="text"
                    placeholder="Title"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    required
                />

                <div
                    {...getRootProps()}
                    className="border-2 border-dashed border-gray-400 rounded p-6 text-center cursor-pointer"
                >
                    <input {...getInputProps()} />
                    {file ? (
                        <p>{file.name}</p>
                    ) : isDragActive ? (
                        <p>Drop the PDF here...</p>
                    ) : (
                        <p>Drag & drop a PDF, or click to select one</p>
                    )}
                </div>

                {error && <p className="text-red-600 text-sm">{error}</p>}

                <div className="flex justify-end space-x-3">
                    <button onClick={onClose} className="px-4 py-2 border rounded hover:bgf-gray-100">
                        Cancel
                    </button>
                    <button
                        onClick={handleSubmit}
                        disabled={uploading}
                        className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                    >
                        {uploading ? 'Uploading...' : 'Submit'}
                    </button>
                </div>
            </div>
        </div>
    );
}