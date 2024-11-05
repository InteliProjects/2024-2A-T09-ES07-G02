import React from "react";
import { Badge } from "./badge";
import PdfViewerDialog from "./PdfViewer";

// Card properties
interface LrrCardProps {
  rule: string;         // The rule field from API response
  agency: string;       // The regulatory agency
  s3_url: string;       // The S3 URL for the PDF
  tags: string[] | null; // Tags, which might be null
}

// Temporary tag mock
const tagColors: Record<string, string> = {
  "TAG 1": "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300",
  "TAG 2": "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
  "TAG 3": "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
  "TAG 4": "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
  "TAG 5": "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300",
};

const LrrCard: React.FC<LrrCardProps> = ({ rule, agency, s3_url, tags }) => {
  // Use an empty array if tags is null
  const tagsToDisplay = tags || [];

  return (
    <div className="flex items-center border border-gray-300 rounded-2xl bg-white shadow-sm dark:bg-gray-800 dark:border-gray-600 w-11/12 p-4">
      <PdfViewerDialog pdfUrl={s3_url} />
      <div className="flex flex-col flex-grow gap-4">
        <span className="font-semibold text-lg dark:text-gray-100">{rule}</span>
        <span className="font-bold dark:text-gray-300">{agency}</span>
        <div className="flex space-x-2 mt-2">
          {tagsToDisplay.length > 0 ? (
            tagsToDisplay.map((tag) => 
              <Badge
                key={tag} // Use the tag itself as key if unique
                className={tagColors[tag] || "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300"}
              >
                {tag}
              </Badge>
            )
          ) : (
            <span className="text-gray-400 dark:text-gray-500">No tags</span>
          )}
        </div>
      </div>
    </div>
  );
};

export default LrrCard;
