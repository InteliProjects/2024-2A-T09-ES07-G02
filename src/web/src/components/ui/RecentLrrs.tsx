import { Bell } from "lucide-react";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Badge } from "./badge";

const recentLrrs = [
  { name: 'Relatório 1', url: 'https://s3.amazonaws.com/mock-bucket/report1.pdf', tags: ['TAG 1', 'TAG 2'] },
  { name: 'Relatório 2', url: 'https://s3.amazonaws.com/mock-bucket/report2.pdf', tags: ['TAG 3', 'TAG 4'] },
  { name: 'Relatório 3', url: 'https://s3.amazonaws.com/mock-bucket/report3.pdf', tags: ['TAG 5', 'TAG 1'] },
  { name: 'Relatório 4', url: 'https://s3.amazonaws.com/mock-bucket/report4.pdf', tags: ['TAG 2', 'TAG 5', 'TAG 1', 'TAG 3', 'TAG 4', 'asdasda'] },
  { name: 'Relatório 5', url: 'https://s3.amazonaws.com/mock-bucket/report5.pdf', tags: ['TAG 3', 'TAG 1'] },
  { name: 'Relatório 6', url: 'https://s3.amazonaws.com/mock-bucket/report6.pdf', tags: ['TAG 4', 'TAG 2'] },
];

const tagColors = {
  "TAG 1": "bg-blue-700 text-blue-200",
  "TAG 2": "bg-green-700 text-green-200",
  "TAG 3": "bg-red-700 text-red-200",
  "TAG 4": "bg-yellow-700 text-yellow-200",
  "TAG 5": "bg-purple-700 text-purple-200",
};

export function RecentLrrsPopover() {
  return (
    <Popover>
      <PopoverTrigger asChild>
        <button className="relative p-2 rounded-full hover:bg-gray-700 focus:outline-none">
          <Bell className="h-6 w-6 text-black dark:text-white" />
        </button>
      </PopoverTrigger>
      <PopoverContent className="w-[28rem] p-4 border border-gray-700 rounded-lg shadow-sm bg-gray-50 dark:bg-gray-800">
        <h4 className="font-semibold mb-2 border-b border-gray-600 pb-2 text-gray-900 dark:text-gray-300">Últimas LRRs</h4>
        <div className="max-h-72 overflow-y-auto">
          <ul className="space-y-3">
            {recentLrrs.map((lrr, index) => (
              <li
                key={index}
                className="flex flex-col p-3 border border-gray-600 rounded-lg bg-gray-100 dark:bg-gray-700 shadow-sm"
              >
                <a
                  href={lrr.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-lg font-medium underline text-blue-600 dark:text-blue-400 mb-2"
                >
                  {lrr.name}
                </a>
                <div className="flex flex-wrap gap-2">
                  {lrr.tags.map((tag, idx) => (
                    <Badge
                      key={idx}
                      className={`${tagColors[tag] || "bg-gray-600 text-gray-300"} whitespace-nowrap`}
                    >
                      {tag}
                    </Badge>
                  ))}
                </div>
              </li>
            ))}
          </ul>
        </div>
      </PopoverContent>
    </Popover>
  );
}
