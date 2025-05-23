type Props = {
  pdfUrl: string;
  title: string;
  onClose: () => void;
};

export default function ViewDocumentModal({ pdfUrl, title, onClose }: Props) {
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
      <div className="bg-white rounded shadow-lg p-4 max-w-4xl w-full h-[90vh] relative">
        <button onClick={onClose} className="absolute top-2 right-2 text-gray-600 text-lg">✕</button>
        <h2 className="text-lg font-bold mb-3">{title}</h2>
        <iframe
          src={pdfUrl}
          title="PDF Viewer"
          className="w-full h-full border"
        />
      </div>
    </div>
  );
}
